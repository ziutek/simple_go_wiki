package main

import "github.com/ziutek/kview"

// Our Wiki pages
var mainView, editView kview.View

func init() {
	// Load layout template
	layout := kview.New("layout.kt")

	// Load template which shows list of articles
	articleList := kview.New("list.kt")

	// Create main page
	mainView = layout.Copy()
	mainView.Div("left", articleList)
	mainView.Div("right", kview.New("show.kt"))

	// Create edit page
	editView = layout.Copy()
	editView.Div("left", articleList)
	editView.Div("right", kview.New("edit.kt"))
}
