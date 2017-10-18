package main

import (
	"strconv"

	"github.com/garyburd/twister/server"
	"github.com/garyburd/twister/web"
)

//ViewCtx Struct view context
type ViewCtx struct {
	Left, Right interface{}
}

// Render main page
func show(req *web.Request) {
	id, _ := strconv.Atoi(req.URLParam["artnum"])
	mainView.Exec(
		req.Respond(web.StatusOK),
		ViewCtx{getArticleList(), getArticle(id)},
	)
}

// Render edit page
func edit(req *web.Request) {
	id, _ := strconv.Atoi(req.URLParam["artnum"])
	editView.Exec(
		req.Respond(web.StatusOK),
		ViewCtx{getArticleList(), getArticle(id)},
	)
}

// Update database and render main page
func update(req *web.Request) {
	id, _ := strconv.Atoi(req.URLParam["artnum"])
	if req.Param.Get("submit") == "Save" {
		id = updateArticle(
			id, req.Param.Get("title"), req.Param.Get("body"),
		)
	}
	// Redirect to the main page which will show the specified article
	req.Redirect("/"+strconv.Itoa(id), false)
}

func main() {
	router := web.NewRouter().
		Register("/style.css", "GET", web.FileHandler("static/style.css", nil)).
		Register("/favicon.ico", "GET", web.FileHandler("static/favicon.ico", nil)).
		Register("/edit/<artnum:[0-9]*>", "GET", edit).
		Register("/<artnum:[0-9]*>", "GET", show, "POST", update)

	handler := web.ProcessForm(10e3, false, router)

	server.Run(":2222", handler)
}
