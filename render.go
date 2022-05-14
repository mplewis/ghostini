package main

import (
	_ "embed"
	"fmt"
	"html"
	"io"
	"strings"
	"text/template"
)

//go:embed templates/index.gmi
var tmplIndexRaw string
var tmplIndex = template.Must(template.New("index").Parse(tmplIndexRaw))

type indexViewModel struct {
	Host            string
	Posts           []Post
	Link            string
	PrevPagePresent bool
	PrevPage        string
	NextPagePresent bool
	NextPage        string
}

func renderIndex(w io.Writer, h host, p postsResp) error {
	indexViewModel := indexViewModel{Host: h.apiUrl, Posts: p.Posts}
	fmt.Printf("%+v\n", p.Meta.Pagination)
	if p.Meta.Pagination.Prev != 0 {
		indexViewModel.PrevPagePresent = true
		indexViewModel.PrevPage = fmt.Sprintf("/?page=%d", p.Meta.Pagination.Prev)
	}
	if p.Meta.Pagination.Next != 0 {
		indexViewModel.NextPagePresent = true
		indexViewModel.NextPage = fmt.Sprintf("/?page=%d", p.Meta.Pagination.Next)
	}
	for i := range p.Posts {
		p.Posts[i].Title = html.UnescapeString(p.Posts[i].Title)
		ex := p.Posts[i].Excerpt
		ex = strings.ReplaceAll(ex, "\n", " ") + "…"
		ex = html.UnescapeString(ex)
		p.Posts[i].Excerpt = ex
	}
	return tmplIndex.Execute(w, indexViewModel)
}
