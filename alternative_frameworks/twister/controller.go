package main

import (
    "log"
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
    if req.Param.GetDef("submit", "") == "Save" {
        id, _ := strconv.Atoi(req.Param.GetDef("artnum", ""))
        id = updateArticle(
            id, req.Param.GetDef("title", ""), req.Param.GetDef("body", ""),
        )
        // If we insert new article, we change artnum to its id. This allows
        // show the article immediately after its creation.
        req.Param.Set("artnum", strconv.Itoa(id))
    }
    // Show modified/created article
    show(req)
}

func main() {
    viewInit()
    mysqlInit()

    router := web.NewRouter().
        Register("/edit/<artnum:.*>", "GET", edit).
        Register("/style.css", "GET", web.FileHandler("static/style.css")).
        Register("/<artnum:.*>", "GET", show, "POST", update)

    handler := web.ProcessForm(100, false, router)

    err := server.ListenAndServe(":1111", &server.Config{Handler: handler})
    if err != nil {
        log.Exitln("ListenAndServe:", err)
    }
}
