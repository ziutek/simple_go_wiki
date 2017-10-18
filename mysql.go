package main

import (
	"log"
	"os"

	"github.com/ziutek/mymysql/autorc"
	"github.com/ziutek/mymysql/mysql"
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

//ArticleList struct
type ArticleList struct {
	ID, Title int
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

//Article struct
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
