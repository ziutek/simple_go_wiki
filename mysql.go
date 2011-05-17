package main

import (
    "os"
    "log"
    mymy "github.com/ziutek/mymysql"
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

    update_stmt, err = db.PrepareAC(
        "INSERT articles (id, title, body) VALUES (?, ?, ?)" +
        " ON DUPLICATE KEY UPDATE title=VALUES(title), body=VALUES(body)",
    )
    mysqlErrExit(err)
}

type ArticleList struct {
    Id, Title int
    Articles  []*mymy.Row
}

// Returns list of articles for list.kt template. We don't create map
// because it is to expensive work. Instead, we provide raw query result
// and indexes to id and title fields.
func getArticleList() *ArticleList {
    rows, res, err := artlist_stmt.ExecAC()
    if mysqlError(err) {
        return nil
    }
    return &ArticleList{
        Id:       res.Map["id"],
        Title:    res.Map["title"],
        Articles: rows,
    }
}

type Article struct {
    Id          int
    Title, Body string
}

// Get an article
func getArticle(id int) (article *Article) {
    rows, res, err := article_stmt.ExecAC(id)
    if mysqlError(err) {
        return
    }
    if len(rows) != 0 {
        article = &Article{
            Id:    id,
            Title: rows[0].Str(res.Map["title"]),
            Body:  rows[0].Str(res.Map["body"]),
        }
    }
    return
}

// Insert or update an article. It returns id of updated/inserted article.
func updateArticle(id int, title, body string) int {
    _, res, err := update_stmt.ExecAC(id, title, body)
    if mysqlError(err) {
        return 0
    }
    return int(res.InsertId)
}
