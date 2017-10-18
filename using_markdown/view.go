package main

import (
	"bytes"

	"github.com/knieriem/markdown"
	"github.com/ziutek/kview"
)

// Our Wiki pages
var mainView, editView kview.View

func viewInit() {
	// Load layout template
	layout := kview.New("layout.kt")

	// Load template which shows list of articles
	articleList := kview.New("list.kt")

	// Create main page
	mainView = layout.Copy()
	mainView.Div("left", articleList)
	mainView.Div("right", kview.New("show.kt", utils))

	// Create edit page
	editView = layout.Copy()
	editView.Div("left", articleList)
	editView.Div("right", kview.New("edit.kt"))
}

var (
	mde = markdown.Extensions{
		Smart:        true,
		Dlists:       true,
		FilterHTML:   true,
		FilterStyles: true,
	}
	utils = map[string]interface{}{
		"markdown": func(txt string) []byte {
			var buf bytes.Buffer
			doc := markdown.Parse(txt, mde)
			doc.WriteHtml(&buf)
			return buf.Bytes()
		},
	}
)
