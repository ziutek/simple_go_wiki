# How to write database-driven Web application using Go

In this tutorial I tried to explain how you can use *web.go*,
*kview*/*kasia.go* and *MyMySQL* together to write a simple database driven web application. As usual, the example application will be a simple Wiki.

## Prerequisites

* Programming experience.
* Basic knowledge about HTML and HTTP.
* Knowledge about MySQL and mysql command-line tool.
* MySQL account with permissions to create tables.
* Last version of Go compiler - see [http://golang.org/doc/install.html]

## Database

Let's start by creating a definition of the article in the database.

If you have got your own MySQL server installation, you have full privileges to
it. In this case, you can create a separate database for this example:

    $ mysql -u root -p
    Enter password: xxxxxxxx
    Welcome to the MySQL monitor.  Commands end with ; or \g.
    Your MySQL connection id is 10136
    Server version: 5.1.49-3~bpo50+1 (Debian)

    mysql> create database test;
    Query OK, 1 row affected (0.03 sec)
    
    mysql> use test
    Database changed
    
If you have only simple MySQL account, with privileges to one database, you will
need to modify further examples, using the name of your database and your user
name. You must also make sure that your database doesn't contain a table called
*articles*.

Next we will create *articles* table:

    mysql> CREATE TABLE articles (
                id    INT AUTO_INCREMENT PRIMARY KEY,
                title VARCHAR(80) NOT NULL,
                body  TEXT NOT NULL
            );
    Query OK, 0 rows affected (0.02 sec)

Next we may create separate user for our application and grant him access to
*articles* table:

    mysql> grant insert,update,delete on articles to testuser@localhost;
    Query OK, 0 rows affected (0.00 sec)

    mysql> set password for testuser@localhost = password('TestPasswd9')
    Query OK, 0 rows affected (0.00 sec)
    
## View

Lets write some code in Go. To define view we use *kview* and *kasia.go*
packages. You may install them this way:

    $ git clone git://github.com/ziutek/kasia.go
    $ cd kasia.go && make install
    $ cd .. 
    $ git clone git://github.com/ziutek/kview
    $ cd kview && make install
    $ cd ..
 
Next we will create directory for our project:

    $ mkdir simple_go_wiki
    $ cd simple_go_wiki
    $ mkdir templates

The *templates* will be used for our Kasia templates.

In this directory we will create our *view.go* file:

    package main

    import "kview"

    // Our Wiki pages
    var main_view, edit_view kview.View

    func viewInit() {
        // Load layout template
        layout := kview.New("layout.kt")

        // Load template which shows list of articles
        article_list := kview.New("list.kt")

        // Create main page
        main_view = layout.Copy()
        main_view.Div("Left", article_list)
        main_view.Div("Right", kview.New("show.kt"))

        // Create edit page
        edit_view = layout.Copy()
        edit_view.Div("Left", article_list)
        edit_view.Div("Right", kview.New("edit.kt"))
    }


You can see, that our service will consist of two pages:

* *main_view* - using which, we will be able to read articles,
* *edit_view* - using which, we will be able to create and edit articles.

Both pages will consists two columns:

* Left - list of articles,
* Right - column specific to the page. 

Lets create our first Kasia template. It will define layout of oure site. We
must create *layout.kt* file in *templates* directory:

    <!DOCTYPE HTML PUBLIC '-//W3C//DTD HTML 4.01 Transitional//EN'>
    <html>
    <head>
        <meta http-equiv='Content-type' content='text/html; charset=utf-8'>
        <link href='/style.css' type='text/css' rel='stylesheet'>
        <title>Simple Wiki</title>
    </head>
        <body>
            <div id='Container0'>
                <div id='Container1'>
                    <div id='Left'>$Left.Render(left)</div>
                    <div id='Right'>$Right.Render(right)</div>
                </div>
            </div>
        </body>
    </html>

This simple layout is responsible for:

* create proper HTML document with the appropriate *doctype*, *head* and *body*
  sections,
* render *Left* and *Right* divs (subtemplates) using the data available in
  *left* and *right* variables.

The *Render* method is defined in *kview* package.

Next we will create *list.kt* which will be rendered in *Left* div.

    <a href='/edit/'>New article</a>
    <hr>
    <ul id='List'>
    $for _, art in articles:
        <li><a href='$art.Data[id]'>$art.Data[title]</a></li>
    $end
    </ul>

This simple template will be rendered to *New article* URL and to the list of
URLs to articles stored in the database.

As you can see it uses *for* statement to iterate over *articles* list. For each
item, it uses *art.Data[id]* variable to create relative URL, and
*art.Data[title]* variable to render the title of the article. *id* and *title*
are also variables. They will contain indexes to the appropriate item in *Data*
slice. *art.Data* will contain the raw row fetched from the MySQL database.

Lets create *show.kt* which will be template for rendering articles:

    $if [0]:
        <div>
            <h4>$title</h4>
            $body
        </div>
        <div id='Actions'><a href='/edit/$id'>Edit</a></div>
    $else:
        <h4>Simple Wiki</h4>
        <p>This application was written entirely in Go language, using the
        following external packages:</p>
        <ul>
            <li><a href='https://github.com/hoisie/web.go'>web.go</a></li>
            <li><a href='https://github.com/ziutek/kasia.go'>kasia.go</a></li>
            <li><a href='https://github.com/ziutek/kview'>kview</a></li>
            <li><a href='https://github.com/ziutek/mymysql'>MyMySQL</a></li>
        </ul>
    $end

As you can see it uses *if*/*else* statement to check that 0 item of the context
stack is *nil* or not *nil*. This item is the *right* variable which we pass to
the *Right.Render* method. If it isn't *nil* there is article selected and we
can render *title* and *body* variables. Otherwise we print an alternative text.

At the end, lets create *edit.kt*:

<form action='/$id' method='post'>
  <div>
    <input name='title' value='$title'>
    <textarea name='body'>$body</textarea>
  </div>
  <div id='Actions'>
    <input type='submit' value='Cancel'>
    <input type='submit' name='submit' value='Save'>
  </div>
</form>

I think that there is nothing to comment here.

We need a style sheet to set the appearance of our website. You can find it in
*static/style.css* file. 

## Communication with MySQL server

For communication with the MySQL server we use *MyMySQL* package. Lets install
it:

    $ cd ..
    $ git clone git://github.com/ziutek/mymysql
    $ cd mymysql && make install
    $ cd ../simple_go_wiki

Now, we can write oure MySQL connector. Lets create the *mysql.go* file. In
first part of this file we import necessary packages, define const and declare
global variables:

    package main

    import (
        "os"
        "log"
        "mymy"
    )

    const (
        db_proto = "tcp"
        db_addr  = "127.0.0.1:3306"
        db_user  = "testuser"
        db_pass  = "TestPasswd9"
        db_name  = "test"
    )

    var (
        // MySQL connection handler
        db = mymy.New(db_proto, "", db_addr, db_user, db_pass, db_name)

        // Prepared statements
        artlist_stmt, article_stmt, update_stmt *mymy.Statement
    )

Next we define some utility functions for MySQL errors handling:

    func mysqlError(err os.Error) (ret bool) {
        ret = (err != nil)
        if ret {
            log.Println("MySQL error:", err)
        }
        return
    }

    func mysqlErrExit(err os.Error) {
        if mysqlError(err) {
            os.Exit(1)
        }
    }

In our application we will use the MyMySQL *autorecon* interface. This is a set
of functions that do not require connect to the database before use them. More
importantly, they don't need to manually reconnect in case of network error or
MySQL server reboot. Lets define the initialisation function:

    func mysqlInit() {
        var err os.Error

        // Initialisation command
        db.Register("SET NAMES utf8")

        // Prepare server-side statements
        artlist_stmt, err = db.PrepareAC("SELECT id, title FROM articles")
        mysqlErrExit(err)

        article_stmt, err = db.PrepareAC(
            "SELECT title, body FROM articles WHERE id = ?",
        )
        mysqlErrExit(err)

        update_stmt, err = db.Prepare(
            "INSERT articles (id, title, body) VALUES (?, ?, ?)" +
            " ON DUPLICATE KEY UPDATE title=VALUES(title), body=VALUES(body)",
        )
        mysqlErrExit(err)
    }

The *Register* method registers commands for executing immediately after
establishing the connection to the database. The *PrepareAC* prepare the
server-side prepared statement. We sould use prepared statements mainly for
security reasons.  With use of prepared statements we don't need any
escape function becouse SQL logic and data are completely separated. Without use
of prepared statements there is always a risk of the SQL injection attack.

Lets write part of code that will be used to get data for left column of our
web pages. 

    type ArticleList struct {
        id, title int
        articles  []*mymy.Row
    }

    // Returns list of articles for article_list.kt template. We don't create
    // map because it is expensive work. Instead, we provide indexes to id and
    // title fields, and raw query result.
    func getArticleList() *ArticleList {
        rows, res, err := artlist_stmt.ExecAC()
        if mysqlError(err) {
            return nil
        }
        return &ArticleList{
            id:       res.Map["id"],
            title:    res.Map["title"],
            articles: rows,
        }
    }

Then define functions for getting and updating articles:

    type Article struct {
        id          int
        title, body string
    }

    // Get an article
    func getArticle(id int) (article *Article) {
        rows, res, err := article_stmt.ExecAC(id)
        if mysqlError(err) {
            return
        }
        if len(rows) != 0 {
            article = &Article{
                id:    id,
                title: rows[0].Str(res.Map["title"]),
                body:  rows[0].Str(res.Map["body"]),
            }
        }
        return
    }

    // Insert or update an article
    func updateArticle(id int, title, body string) {
        _, _, err := update_stmt.ExecAC(id, title, body)
        mysqlError(err)
    }

Last function uses MySQL *INSERT ... ON DUPLICATE KEY UPDATE* query. It inserts
or updates article depending on whether it exists or not exists.

## Controller

Next we should write the last part of our application witch is responsible for
interaction with the user. Lets create *controller.go* file:

    package main

    import (
        "web"
        "strconv"
    )

    type ViewCtx struct {
        left, right interface{}
    }

    // Render main page
    func show(wr *web.Context, art_num string) {
        id, _ := strconv.Atoi(art_num)
        main_view.Exec(wr, ViewCtx{getArticleList(), getArticle(id)})
    }

    // Render edit page
    func edit(wr *web.Context, art_num string) {
        id, _ := strconv.Atoi(art_num)
        edit_view.Exec(wr, ViewCtx{getArticleList(), getArticle(id)})
    }

    // Update database and render main page
    func update(wr *web.Context, art_num string) {
        if wr.Request.Params["submit"] == "Save" {
            id, _ := strconv.Atoi(art_num) // id == 0 means new article
            updateArticle(id, wr.Request.Params["title"], wr.Request.Params["body"])
        }
        show(wr, art_num)
    }

    func main() {
        viewInit()
        mysqlInit()

        web.Get("/edit/(.*)", edit)
        web.Get("/(.*)", show)
        web.Post("/(.*)", update)
        web.Run("0.0.0.0:1111")
    }

We use *web.go* framework for binding handlers to specified URLs and HTTP
methods. URLs are specified by regular expressions.

The *show* handler, binded to *GET* method and "/(.*)" URL scheme, is
responsible for render the main page witch allows the user to read selected
articles. The "/(.*)" regular expression match any URL and return it's path
part as article number. So if URL looks like:

    http://www.simple-go-wiki.org/19

it will return "19" as an article number. If URL looks like:

    http://www.simple-go-wiki.org/edit/19

it will return "edit/19" as an article number. Therefore, we must bind *edit*
handler before *show* handler. If article number isn't a number it will be
converted by *strconv.Atoi* to 0 which means unknown article.

The *edit* handler, binded to *GET* method and "/edit/(.*)" URL scheme, is
responsible for edit or create new article.

The *update* handler, binded to *POST* method and "/(.*)" URL scheme, is
responsible for update article in database. It modify article in database only
when user push the *Save* button on edit page. After updating the database this
handler calls *show* handler for render the main page.

## Building the application

Lets create *Makefile* for our project:

    include $(GOROOT)/src/Make.inc

    TARG=wiki
    GOFILES=\
        view.go\
        controller.go\
        mysql.go\

    include $(GOROOT)/src/Make.cmd

Next we can build our application:

    $ make
    8g -o _go_.8 view.go controller.go mysql.go 
    8l -o wiki _go_.8

and launch it:

    $ ./wiki 
    2011/01/26 19:44:55 web.go serving 0.0.0.0:1111

If may get this application from Github:

    $ git git@github
