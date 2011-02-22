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
    id, _ := strconv.Atoi(req.Param.GetDef("artnum", ""))
    main_view.Exec(
        req.Respond(web.StatusOK),
        ViewCtx{getArticleList(), getArticle(id)},
    )
}

// Render edit page
func edit(req *web.Request) {
    id, _ := strconv.Atoi(req.Param.GetDef("artnum", ""))
    edit_view.Exec(
            req.Respond(web.StatusOK),
            ViewCtx{getArticleList(), getArticle(id)},
    )
}

// Update database and render main page
func update(req *web.Request) {
    id, _ := strconv.Atoi(req.Param.GetDef("artnum", ""))
    if req.Param.GetDef("submit", "") == "Save" {
        id = updateArticle(
            id, req.Param.GetDef("title", ""), req.Param.GetDef("body", ""),
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
        Register("/style.css", "GET", web.FileHandler("static/style.css")).
        Register("/favicon.ico", "GET", web.FileHandler("static/favicon.ico")).
        Register("/edit/<artnum:.*>", "GET", edit).
        Register("/<artnum:.*>", "GET", show, "POST", update)

    handler := web.ProcessForm(10e3, false, router)

    server.Run(":1111", handler)
}
