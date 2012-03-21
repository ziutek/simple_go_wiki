package main

import (
    "strconv"
    "github.com/hoisie/web.go"
)

type ViewCtx struct {
    Left, Right interface{}
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
        // If we insert new article, we change art_num to its id. This
        // allows to show the article immediately after its creation.
        art_num = strconv.Itoa(id)
    }
    // Redirect to the main page which will show the specified article
    wr.Redirect(303, "/" + art_num)
    // We could show this article directly using show(wr, art_num)
    // but see: http://en.wikipedia.org/wiki/Post/Redirect/Get
}

func main() {
    web.Get("/edit/(.*)", edit)
    web.Get("/(.*)", show)
    web.Post("/(.*)", update)
    web.Run("0.0.0.0:2222")
}
