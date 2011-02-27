package main

import (
    "strconv"
    "github.com/garyburd/twister/server"
    "github.com/garyburd/twister/web"
)

type ViewCtx struct {
    left, right interface{}
}

// Render main page
func show(req *web.Request) {
    id, _ := strconv.Atoi(req.Param.Get("artnum"))
    main_view.Exec(
        req.Respond(web.StatusOK),
        ViewCtx{getArticleList(), getArticle(id)},
    )
}

// Render edit page
func edit(req *web.Request) {
    id, _ := strconv.Atoi(req.Param.Get("artnum"))
    edit_view.Exec(
            req.Respond(web.StatusOK),
            ViewCtx{getArticleList(), getArticle(id)},
    )
}

// Update database and render main page
func update(req *web.Request) {
    id, _ := strconv.Atoi(req.Param.Get("artnum"))
    if req.Param.Get("submit") == "Save" {
        id = updateArticle(
            id, req.Param.Get("title"), req.Param.Get("body"),
        )
    }
    // Redirect to the main page which will show the specified article
    req.Redirect("/" + strconv.Itoa(id), false)
    // We could show this article directly using:
    //     req.Param.Set("artnum", strconv.Itoa(id))
    //     show(req)
    // but see: http://en.wikipedia.org/wiki/Post/Redirect/Get
}

func main() {
    viewInit()
    mysqlInit()

    router := web.NewRouter().
        Register("/edit/<artnum:.*>", "GET", edit).
        Register("/style.css", "GET", web.FileHandler("static/style.css")).
        Register("/<artnum:.*>", "GET", show, "POST", update)

    handler := web.ProcessForm(10000, false, router)

    server.Run(":1111", handler)
}
