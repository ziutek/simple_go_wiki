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
