# How to write database-driven Web application using Go

In this tutorial I tried to explain how you can use *kview*/*kasia.go* and
*MyMySQL* to write a simple database driven web application in Go. As usual,
the example application will be a simple Wiki.

## Prerequisites

* Programming experience.
* Basic knowledge about HTML and HTTP.
* Knowledge about MySQL and mysql command-line tool.
* MySQL account with permissions to create tables.
* Last weekly version of Go compiler - see
  [Go homepage](http://weekly.golang.org/doc/install.html)

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

    mysql> GRANT SELECT,INSERT,UPDATE,DELETE ON articles TO testuser@localhost;
    Query OK, 0 rows affected (0.00 sec)

    mysql> SET PASSWORD FOR testuser@localhost = PASSWORD('TestPasswd9');
    Query OK, 0 rows affected (0.00 sec)

## View

Lets write some code in Go. To define the application view we will use *kview*
and *kasia.go* packages. You may install them this way:

    $ go get github.com/ziutek/kview

It automatically downloads, builds and installs *kasia.go* and *kview*.

Next we will create the directory for our project:

    $ mkdir simple_go_wiki
    $ cd simple_go_wiki
    $ mkdir templates static

The *templates* directory will be used for our HTML templates. The *static*
directory will be used for static files like *style.css*.

In the *simple_go_wiki* directory we can create our *view.go* file:

    package main

    import "github.com/ziutek/kview"

    // Our Wiki pages
    var mainView, editView kview.View

    func init() {
        // Load layout template
        layout := kview.New("layout.kt")

        // Load template which shows list of articles
        articleList := kview.New("list.kt")

        // Create main page
        mainView = layout.Copy()
        mainView.Div("left", articleList)
        mainView.Div("right", kview.New("show.kt"))

        // Create edit page
        editView = layout.Copy()
        editView.Div("left", articleList)
        editView.Div("right", kview.New("edit.kt"))
    }

As you can see, our service will consist of two pages:

* *mainView* - using which the user will be able to read articles,
* *editView* - using which the user will be able to create and edit articles.

Both pages will consists of two columns:

* left - list of articles,
* right - column specific to the page. 

Lets create our first Kasia template. It will define the layout of our site. We
have to create *layout.kt* file in *templates* directory:

    <!doctype html>
    <html>
        <head>
            <meta charset='utf-8'>
            <link href='/style.css' type='text/css' rel='stylesheet'>
            <title>Simple Wiki</title>
        </head>
        <body>
            <div id='Container0'>
                <div id='Container1'>
                    <div id='Left'>$left.Render(Left)</div>
                    <div id='Right'>$right.Render(Right)</div>
                </div>
            </div>
        </body>
    </html>

This simple layout is responsible for:

* create proper HTML document with the appropriate *doctype*, *head* and *body*
  sections,
* render *left* and *right* divs (subtemplates) using the data available in
  *Left* and *Right* variables.

The *Render* method is defined in *kview* package. It renders subview with
specified data in the place of its occurrence. This subview can have its own
divs and can render another subviews (but we don't use this property in our
simple Wiki).

Next we will create *list.kt* which will be rendered in *Left* div.

    <a href='/edit/'>New article</a>
    <hr>
    <ul id='List'>
    $for _, art in Articles:
        <li><a href='$art[Id]'>$art[Title]</a></li>
    $end
    </ul>

This simple template prints the *New article* link and the list of links to
articles stored in the database.

As you can see it uses a *for* statement to iterate over the *Articles* list
(slice). For each item, it uses *art[Id]* to create relative URL, and
*art[Title]* to print the title of the article. *Articles*, *Id* and
*Title* are members of *ArticleList* struct (defined later in this tutorial). 
*Id* and *Title* will contain indexes to the appropriate item in *Row* slice,
*art* will contain the raw row fetched from the MySQL database.

*for* statement create two variables (*_* and *art*) in the local context.
First is the iteration number, second is the list element. We don't use
iteration number in our example but it may be useful:

    $for nn+, art in Articles:
        <div class='$if even(nn):Even$else:Odd$end'>
            $nn. <a href='$art[Id]'>$art[Title]</a>
        </div>
    $end

To make this work we need to provide *even(int)* function in context stack. '+'
after *nn* means that we counts from 1, not from 0.

Lets create *show.kt* which will be template for rendering articles:

    $if Id:
        <div>
            <h4>$Title</h4>
            $Body
        </div>
        <div id='Actions'><a href='/edit/$Id'>Edit</a></div>
    $else:
        <h4>Simple Wiki</h4>
        <p>This application was wrote entirely in Go language, using the
        following external packages:</p>
        <ul>
            <li><a href='https://github.com/ziutek/kasia.go'>kasia.go</a></li>
            <li><a href='https://github.com/ziutek/kview'>kview</a></li>
            <li><a href='https://github.com/ziutek/mymysql'>MyMySQL</a></li>
        </ul>
    $end

As you can see it uses the *if - else* statement to determine that is there
the article selected or not. If article is selected then *id* field of the
*Article* struct (defined later) has no zero value so we can render *title*
and *body* variables. Otherwise we print our alternative text.

#### An interlude about the context stack.

To render some template with *kview* package, you have to use *Exec* method in
Go code or *Render* method in template code. Usually you need to pass some
variables to render in template code. For example if you use *Exec* or *Render*
method like this:

    v.Exec(wr, a, b)
    v.Render(a, b)

the template associated with *v* view will be rendered with the following
context stack:

    []interface{}{globals, a, b} 

As you can see there is the *globals* variable at the bottom of the stack.
The *globals* is a map containing global symbols:

* subviews (subtemplates) added to *v* by *Div* method,
* *len* and *fmt* utility functions,
* yours symbols which you pass to *New* function as additional parameters.

For more information see [kview
documentation](https://github.com/ziutek/kview/blob/master/README.md).  

Your *b* variable is at the top of the stack. If you write template like this:

    $x  $@[1].y

then *Exec* or *Render* method will look for *x* and *y* attributes as follows:

1. *x* will be first searched in *b*, and if not found, it will be searched in
   *a*, and if *a* also doesn't contains field of name *x*, it will be searched
   in *gobals*.
2. *y* will be searched only in *a* because you specify directly element of
   context stack (*@[1]*) in which to look for it.

The *@* symbol means the context stack itself. So you can use the context stack
as usual or as a parameter list:

* Go code:
  `v.Exec(os.Stdout, "Hello", "world!", map[string]string{"kasia": "Katy"})`
* template: `$@[1] $@[2]  $@[1] $kasia!`
* output: `Hello world!  Hello Katy!`

At last, you can print full context stack to check their contents as follows:

    $for i, v in @:
        $i: $v<br>
    $end

Fro more information see [Kasia.go
documentation](https://github.com/ziutek/kasia.go/blob/master/README.md).

After this small interlude we should return to our work. Lets create last
template in *edit.kt* file:

    <form action='/$Id' method='post'>
        <div>
            <input name='title' value='$Title'>
            <textarea name='body'>$Body</textarea>
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

    $ go get github.com/ziutek/mymysql/autorc

Now, we can write the MySQL connector for our application. Lets create the
*mysql.go* file. In the first part of this file we import necessary packages,
define const and declare global variables:
 
	package main

    import (
        "os"
        "log"
        "github.com/ziutek/mymysql/mysql"
        "github.com/ziutek/mymysql/autorc"
        _ "github.com/ziutek/mymysql/thrsafe"
    )

    const (
        dbProto = "tcp"
        dbAddr  = "127.0.0.1:3306"
        dbUser  = "testuser"
        dbPass  = "TestPasswd9"
        dbName  = "test"
    )

    var (
        // MySQL connection handler
        db = autorc.New(dbProto, "", dbAddr, dbUser, dbPass, dbName)

        // Prepared statements
        artlistStmt, articleStmt, updateStmt *autorc.Stmt
    )

After declaration, the MySQL connection handler is ready for connect to the
database. But we will not make this connection explicitly.

In our application we will use the MyMySQL *autorecon* interface. This is a set
of functions that do not require a connection to the database before using them. More
importantly, they don't need to manually reconnect in case of network error or
MySQL server reboot.

Next we will define some utility functions for MySQL errors handling:

    func mysqlError(err error) (ret bool) {
        ret = (err != nil)
        if ret {
            log.Println("MySQL error:", err)
        }
        return
    }

    func mysqlErrExit(err error) {
        if mysqlError(err) {
            os.Exit(1)
        }
    }


Lets define the initialisation function. 

    func init() {
        var err error

        // Initialisation command
        db.Raw.Register("SET NAMES utf8")

        // Prepare server-side statements

        artlistStmt, err = db.Prepare("SELECT id, title FROM articles")
        mysqlErrExit(err)

        articleStmt, err = db.Prepare(
            "SELECT title, body FROM articles WHERE id = ?",
        )
        mysqlErrExit(err)

        updateStmt, err = db.Prepare(
            "INSERT articles (id, title, body) VALUES (?, ?, ?)" +
            " ON DUPLICATE KEY UPDATE title=VALUES(title), body=VALUES(body)",
        )
        mysqlErrExit(err)
    }

The *Register* method registers commands for executing immediately after
establishing the connection to the database. The *Prepare* prepares the
server-side prepared statement. We use functions from *mymysql/autorc* package
Therefore, during the first *Prepare* call the connection will be established.

Why do we use prepared statements instead of ordinary queries? We use them
mainly for security reasons. With prepared statements we don't need any escape
function for user input, because SQL logic and data are completely separated.
Without use of prepared statements there is always a risk of the SQL injection
attack.

Lets write the code that will be used to get data for left column of our
web pages. 

    type ArticleList struct {
        IS, Title int
        Articles  []mysql.Row
    }

    // Returns list of articles for list.kt template. We don't create map
    // because it is to expensive work. Instead, we provide raw query result
    // and indexes to id and title fields.
    func getArticleList() *ArticleList {
        rows, res, err := artlistStmt.Exec()
        if mysqlError(err) {
            return nil
        }
        return &ArticleList{
            ID:       res.Map("id"),
            Title:    res.Map("title"),
            Articles: rows,
        }
    }

Then define functions for getting and updating articles:

    type Article struct {
        ID          int
        Title, Body string
    }

    // Get an article
    func getArticle(id int) (article *Article) {
        rows, res, err := articleStmt.Exec(id)
        if mysqlError(err) {
            return
        }
        if len(rows) != 0 {
            article = &Article{
                ID:    id,
                Title: rows[0].Str(res.Map("title")),
                Body:  rows[0].Str(res.Map("body")),
            }
        }
        return
    }

    // Insert or update an article. It returns id of updated/inserted article.
    func updateArticle(id int, title, body string) int {
            _, res, err := updateStmt.Exec(id, title, body)
        if mysqlError(err) {
            return 0
        }
        return int(res.InsertId())
    }

The last function uses MySQL *INSERT ... ON DUPLICATE KEY UPDATE* query. It
inserts or updates article depending on whether it exists or not exists in the
table.

## Controller

The last part of our application which is responsible for
interaction with the user. Lets create *controller.go* file:

    package main

    import (
	    "log"
	    "net/http"
	    "strconv"
	    "strings"
    )

    type ViewCtx struct {
	    Left, Right interface{}
    }

    // Render main page
    func show(wr http.ResponseWriter, artNum string) {
	    id, _ := strconv.Atoi(artNum)
	    mainView.Exec(wr, ViewCtx{getArticleList(), getArticle(id)})
    }

    // Render edit page
    func edit(wr http.ResponseWriter, artNum string) {
	    id, _ := strconv.Atoi(artNum)
	    editView.Exec(wr, ViewCtx{getArticleList(), getArticle(id)})
    }

    // Update database and render main page
    func update(wr http.ResponseWriter, req *http.Request, artNum string) {
	    if req.FormValue("submit") == "Save" {
		    id, _ := strconv.Atoi(artNum) // id == 0 means new article
		    id = updateArticle(
			    id, req.FormValue("title"), req.FormValue("body"),
		    )
		    // If we insert new article, we change artNum to its id. This allows
		    // show the article immediately after its creation.
		    artNum = strconv.Itoa(id)
	    }
	    // Redirect to the main page which will show the specified article
	    http.Redirect(wr, req, "/"+artNum, 303)
	    // We could show this article directly using show(wr, art_num)
	    // but see: http://en.wikipedia.org/wiki/Post/Redirect/Get
    }

    // Decide which handler to use basis on the request method and URL path.
    func router(wr http.ResponseWriter, req *http.Request) {
	    rootPath := "/"
	    editPath := "/edit/"

	    switch req.Method {
	    case "GET":
		    switch {
		    case req.URL.Path == "/style.css" || req.URL.Path == "/favicon.ico":
			    http.ServeFile(wr, req, "static"+req.URL.Path)

		    case strings.HasPrefix(req.URL.Path, editPath):
			    edit(wr, req.URL.Path[len(editPath):])
    
		    case strings.HasPrefix(req.URL.Path, rootPath):
			    show(wr, req.URL.Path[len(rootPath):])
		    }

	    case "POST":
		    switch {
		    case strings.HasPrefix(req.URL.Path, rootPath):
			    update(wr, req, req.URL.Path[len(rootPath):])
		    }
	    }
    }
    
    func main() {
	    err := http.ListenAndServe(":2222", http.HandlerFunc(router))
	    if err != nil {
		    log.Fatalln("ListenAndServe:", err)
	    }
    }

The *show* handler, binded to *GET* method and "/(.\*)" URL scheme, is
responsible for render the main page witch allows the user to select and read
articles.

The *edit* handler, bound to *GET* method and "/edit/(.\*)" URL scheme, is
responsible for edit or create new article.

The *update* handler, bound to *POST* method and "/(.\*)" URL scheme, is
responsible for update an article in database. It modify it in database
only when user push the *Save* button on edit page, which is checked using
*wr.Request.Params["submit"]* variable. After updating the database this
handler sends the redirect response which redirects to the main page.

## Run the application

	$ go run *.go

You may get this tutorial and example application from Github using the following
command:

    $ git clone git://github.com/ziutek/simple_go_wiki.git

os install it in your *GOPATH* using *go* command:

	$ go get github.com/ziutek/simple_go_wiki

## Other frameworks

There are versions of *controller.go* file for
[web.go](http://www.getwebgo.com/) package and
[twister](https://github.com/garyburd/twister) package. You can find them in
*alternative_frameworks* directory. 

## Using Markdown to format the article body

To use the [Markdown](http://daringfireball.net/projects/markdown/syntax) syntax
in article body we can modify *getArticle* function to use
[markdown package](https://github.com/knieriem/markdown) to convert
article body fetched from database before pass it into *body* field of *Article*
struct. But we'll go another way: we define *markdown* utility function (using
*markdown* package) and we'll provide it for using inside the *show.kt*
template. It seems to be more general solution: we create additional tool which
we can use in any template code wherever we needed. To do this we should to
modify two files: *view.go* and *show.kt* template.

We define *utils* map with only one utility function:

    utils = map[string]interface{} {
        "markdown": func(txt string) []byte {
            var buf bytes.Buffer
            doc := markdown.Parse(txt, mde)
            doc.WriteHtml(&buf)
            return buf.Bytes()
        },
    }

and add its contents to the *globals* for the *show.kt* template:

    mainView.Div("Right", kview.New("show.kt", utils))

At last we should change `$body` to `$:markdown(body)` in the *show.kt* file.
For details see the *viwe.go* and *templates/show.kt* files in *using_markdown*
directory.

To build and run example wiki with *Markdown* support type:

    $ go get github.com/knieriem/markdown
    $ cd using_markdown
    $ ./run.sh

## Exercices

Try to add the ability to delete an article.
