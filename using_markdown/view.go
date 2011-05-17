package main

import (
    "bytes"
    "github.com/ziutek/kview"
    "github.com/knieriem/markdown"
)

// Our Wiki pages
var main_view, edit_view kview.View

func viewInit() {
    // Load layout template
    layout := kview.New("layout.kt")

    // Load template which shows list of articles
    article_list := kview.New("list.kt")

    // Create main page
    main_view = layout.Copy()
    main_view.Div("left", article_list)
    main_view.Div("right", kview.New("show.kt", utils))

    // Create edit page
    edit_view = layout.Copy()
    edit_view.Div("left", article_list)
    edit_view.Div("right", kview.New("edit.kt"))
}

var (
    mde = markdown.Extensions{
        Smart:        true,
        Dlists:       true,
        FilterHTML:   true,
        FilterStyles: true,
    }
    utils = map[string]interface{} {
        "markdown": func(txt string) []byte {
            var buf bytes.Buffer
            doc := markdown.Parse(txt, mde)
            doc.WriteHtml(&buf)
            return buf.Bytes()
        },
    }
)
