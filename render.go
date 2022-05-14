package main

import (
	_ "embed"
	"fmt"
	"html"
	"io"
	"strings"
	"text/template"

	"github.com/LukeEmmet/html2gemini"
)

//go:embed templates/index.gmi
var tmplIndexRaw string
var tmplIndex = template.Must(template.New("index").Parse(tmplIndexRaw))

//go:embed templates/post.gmi
var tmplPostRaw string
var tmplPost = template.Must(template.New("post").Parse(tmplPostRaw))

type indexViewModel struct {
	Host            string
	Posts           []PostMeta
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

func renderPost(w io.Writer, p Post) error {
	p.Title = html.UnescapeString(p.Title)
	c := html2gemini.NewTraverseContext(*html2gemini.NewOptions())
	text, err := html2gemini.FromString(p.HTML, *c)
	if err != nil {
		return err
	}
	p.HTML = text
	return tmplPost.Execute(w, p)
}
