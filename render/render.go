// Package render implements a Gemini renderer for Ghost posts.

package render

import (
	_ "embed"
	"fmt"
	"html"
	"io"
	"strings"
	"text/template"

	"github.com/LukeEmmet/html2gemini"
	"github.com/mplewis/ghostini/ghost"
)

//go:embed templates/index.gmi
var tmplIndexRaw string
var tmplIndex = template.Must(template.New("index").Parse(tmplIndexRaw))

//go:embed templates/post.gmi
var tmplPostRaw string
var tmplPost = template.Must(template.New("post").Parse(tmplPostRaw))

// indexViewModel is a view model for the posts index page.
type indexViewModel struct {
	Host            string
	Posts           []ghost.PostMeta
	Link            string
	PrevPagePresent bool
	PrevPage        string
	NextPagePresent bool
	NextPage        string
}

// Index renders the index page, listing Ghost posts.
func Index(w io.Writer, h ghost.Host, p ghost.PostsResp) error {
	indexViewModel := indexViewModel{Host: h.APIURL, Posts: p.Posts}
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

// Post renders a single Ghost post.
func Post(w io.Writer, p ghost.Post) error {
	p.Title = html.UnescapeString(p.Title)
	opts := html2gemini.NewOptions()
	opts.LinkEmitFrequency = 1
	c := html2gemini.NewTraverseContext(*opts)
	text, err := html2gemini.FromString(p.HTML, *c)
	if err != nil {
		return err
	}
	p.HTML = text
	return tmplPost.Execute(w, p)
}
