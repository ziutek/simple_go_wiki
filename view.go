package main

import "kview"

// Our Wiki pages
var main_view, edit_view kview.View

func viewInit() {
    // Load layout template
    layout := kview.New("layout.kt")

    // Load template which shows list of articles
    article_list := kview.New("list.kt")

    // Create main page
    main_view = layout.Copy()
    main_view.Div("Left", article_list)
    main_view.Div("Right", kview.New("show.kt"))

    // Create edit page
    edit_view = layout.Copy()
    edit_view.Div("Left", article_list)
    edit_view.Div("Right", kview.New("edit.kt"))
}
