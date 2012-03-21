package main

import "github.com/ziutek/kview"

// Our Wiki pages
var main_view, edit_view kview.View

func init() {
    // Load layout template
    layout := kview.New("layout.kt")

    // Load template which shows list of articles
    article_list := kview.New("list.kt")

    // Create main page
    main_view = layout.Copy()
    main_view.Div("left", article_list)
    main_view.Div("right", kview.New("show.kt"))

    // Create edit page
    edit_view = layout.Copy()
    edit_view.Div("left", article_list)
    edit_view.Div("right", kview.New("edit.kt"))
}
