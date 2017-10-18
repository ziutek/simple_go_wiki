package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

//ViewCtx struct
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
		// If we insert new article, we change art_num to its id. This allows
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
