# How to write database-driven Web application using Go

In this tutorial I tried to explain how you can use *web.go*,
*kview*/*kasia.go* and *MyMySQL* together to write a simple database driven web application. As usual, the example application will be a simple Wiki.

## Prerequisites

* Programming experience.
* Basic knowledge about HTML and HTTP.
* Knowledge about MySQL and mysql command-line tool.
* MySQL account with permissions to create tables.
* Last version of Go compiler - see [Go homepage](http://golang.org/doc/install.html)

## Database

Let's start by creating a definition of the article in the database.

If you have your own MySQL server installation, you have full privileges to
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
            ) DEFAULT CHARSET=utf8;
    Query OK, 0 rows affected (0.02 sec)

Next we may create separate user for our application and grant him access to
*articles* table:

    mysql> GRANT INSERT,UPDATE,DELETE ON articles TO testuser@localhost;
    Query OK, 0 rows affected (0.00 sec)

    mysql> SET PASSWORD FOR testuser@localhost = PASSWORD('TestPasswd9')
    Query OK, 0 rows affected (0.00 sec)
    
## View

Lets write some code in Go. To define the application view we will use *kview*
and *kasia.go* packages. You may install them this way:

    $ git clone git://github.com/ziutek/kasia.go
    $ cd kasia.go && make install
    $ cd .. 
    $ git clone git://github.com/ziutek/kview
    $ cd kview && make install
    $ cd ..
 
Next we will create the directory for our project:

    $ mkdir simple_go_wiki
    $ cd simple_go_wiki
    $ mkdir templates static

The *templates* directory will be used for our Kasia templates. The *static*
directory will be used for static files like *style.css*.  Static files are 
served by *web.go*.

In the *simple_go_wiki* directory we can create our *view.go* file:

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


As you can see, our service will consist of two pages:

* *main_view* - using which the user will be able to read articles,
* *edit_view* - using which the user will be able to create and edit articles.

Both pages will consists of two columns:

* Left - list of articles,
* Right - column specific to the page. 

Lets create our first Kasia template. It will define the layout of our site. We
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

The *Render* method is defined in *kview* package. It renders subtemplate with
specified data in the place of its occurrence.

Next we will create *list.kt* which will be rendered in *Left* div.

    <a href='/edit/'>New article</a>
    <hr>
    <ul id='List'>
    $for _, art in articles:
        <li><a href='$art.Data[id]'>$art.Data[title]</a></li>
    $end
    </ul>

This simple template prints the *New article* URL and the list of URLs to
articles stored in the database.

As you can see it uses a *for* statement to iterate over the *articles* list. For each
item, it uses *art.Data[id]* variable to create relative URL, and
*art.Data[title]* variable to print the title of the article. *id* and *title*
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

As you can see it uses the *if* / *else* statement to determine that element 0 of the
context stack array is *nil* or not *nil*. This item is the *right* variable which we
pass to the *Right.Render* method. If it isn't *nil* there is an article
selected and we can render *title* and *body* variables. Otherwise we print our
alternative text.

Finally, lets create *edit.kt*:

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

This is self explanatory.

We need a style sheet to set the appearance of our website. You can find it in
*static/style.css* file. 

## Communication with MySQL server

For communication with the MySQL server we use *MyMySQL* package. Lets install
it:

    $ cd ..
    $ git clone git://github.com/ziutek/mymysql
    $ cd mymysql && make install
    $ cd ../simple_go_wiki

Now, we can write the MySQL connector for our application. Lets create the
*mysql.go* file. In the first part of this file we import necessary packages,
define const and declare global variables:

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
After declaration, the MySQL connection handler is ready for connect to the
database. But we will not make this connection explicitly.

In our application we will use the MyMySQL *autorecon* interface. This is a set
of functions that do not require a connection to the database before using them. More
importantly, they don't need to manually reconnect in case of network error or
MySQL server reboot.

Next we will define some utility functions for MySQL errors handling:

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


Lets define the initialisation function. It is called once from *main* function
and initialises our MySQL connector.

    func mysqlInit() {
        var err os.Error

        // Initialisation command
        db.Register("SET NAMES utf8")

        // Prepare server-side statements

        artlist_stmt, err = db.PrepareAC("SELECT id, title FROM articles")
        mysqlErrExit(err)

        article_stmt, err = db.PrepareAC("SELECT title, body FROM articles WHERE id = ?")
        mysqlErrExit(err)

        update_stmt, err = db.Prepare(
            "INSERT articles (id, title, body) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE title=VALUES(title), body=VALUES(body)",
        )
        mysqlErrExit(err)
    }

The *Register* method registers commands for executing immediately after
establishing the connection to the database. The *PrepareAC* prepare the
server-side prepared statement. *AC* suffix means that it is a function from
MyMySQL *autorecon* interface. Therefore, during the first *PrepareAC*  call the
connection will be established.

Why do we use prepared statements instead of ordinary queries? We use them
mainly for security reasons. With prepared statements we don't need any escape
function for user input, because SQL logic and data are completely separated.
Without use of prepared statements there is always a risk of the SQL injection
attack.

Lets write the code that will be used to get data for left column of our
web pages. 

    type ArticleList struct {
        id, title int
        articles  []*mymy.Row
    }

    // Returns list of articles for list.kt template. We don't create
    // map because it is to expensive work. Instead, we provide indexes to id
    // and title fields, and raw query result.
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

    // Insert or update an article. It return id of updated/inserted article
    func updateArticle(id int, title, body string) int {
        _, res, err := update_stmt.ExecAC(id, title, body)
        if mysqlError(err) {
            return 0
        }
        return int(res.InsertId)
    }

The last function uses MySQL *INSERT ... ON DUPLICATE KEY UPDATE* query. It inserts
or updates article depending on whether it exists or not exists in the table.

## Controller

We need to install *web.go*:

    $ cd ..
    $ git clone git://github.com/hoisie/web.go
    $ cd web.go && make install
    $ cd ../simple_go_wiki

Next we can write the last part of our application which is responsible for
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
            id = updateArticle(
                id, wr.Request.Params["title"], wr.Request.Params["body"],
            )
            // If we insert new article, we change art_num to its id. This allows
            // show the article immediately after its creation.
            art_num = strconv.Itoa(id)
        }
        // Show modified/created article
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

The *show* handler, binded to *GET* method and "/(.\*)" URL scheme, is
responsible for render the main page witch allows the user to select and read
articles. The "/(.\*)" regular expression matches any URL and returns it's path
part as article number. So if URL looks like:

    http://www.simple-go-wiki.org/19

it will return "19" as an article number. If URL looks like:

    http://www.simple-go-wiki.org/edit/19

it will return "edit/19" as an article number. Therefore, we must bind *edit*
handler before *show* handler. The  article number will be converted by
*strconv.Atoi* to integer value. If it is empty string or it isn't a number it
will be converted to 0, which means unknown article.

The *edit* handler, bound to *GET* method and "/edit/(.\*)" URL scheme, is
responsible for edit or create new article.

The *update* handler, bound to *POST* method and "/(.\*)" URL scheme, is
responsible for update an article in database. It modify it in database
only when user push the *Save* button on edit page, which is checked using
*wr.Request.Params["submit"]* variable. After updating the database this
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

You may get this tutorial and example application from Github using the following
command:

    $ git clone git://github.com/ziutek/simple_go_wiki.git

## Exercices

Try to add the ability to delete an article.
