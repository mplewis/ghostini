// Package server implements a Gemini server that serves Ghost content.

package server

import (
	_ "embed"
	"fmt"
	"regexp"

	"github.com/a-h/gemini"
	"github.com/mplewis/ghostini/cache"
	"github.com/mplewis/ghostini/ghost"
	"github.com/mplewis/ghostini/parse"
	"github.com/mplewis/ghostini/render"
)

// slugMatcher matches URL paths for slugs with optional trailing slashes, such as /my-slug or /my-slug/.
var slugMatcher = regexp.MustCompile(`^/([^/]+)/?$`)

// Server implements a Gemini server that serves Ghost content.
type Server struct {
	cache *cache.Cache
	host  ghost.Host
}

// ServeGemini handles routing and rendering.
func (s Server) ServeGemini(w gemini.ResponseWriter, r *gemini.Request) {
	if r.URL.Path == "/" {
		page := parse.Int(r.URL.Query().Get("page"), 1)
		resp, err := ghost.GetPosts(s.cache, s.host, page)
		if err != nil {
			w.SetHeader(gemini.CodeTemporaryFailure, "")
			return
		}
		w.SetHeader(gemini.CodeSuccess, "")
		render.Index(w, s.host, resp)
		return
	}

	if matches := slugMatcher.FindStringSubmatch(r.URL.Path); len(matches) > 0 {
		slug := matches[1]
		resp, found, err := ghost.GetPost(s.cache, s.host, slug)
		if err != nil {
			w.SetHeader(gemini.CodeTemporaryFailure, "")
			return
		}
		if !found {
			w.SetHeader(gemini.CodeNotFound, "")
			return
		}

		w.SetHeader(gemini.CodeSuccess, "")
		render.Post(w, resp.Posts[0])
		return
	}

	fmt.Println("invalid path")
	w.SetHeader(gemini.CodeNotFound, "")
	w.Write([]byte("not found"))
}

func New(host ghost.Host) Server {
	return Server{cache.New(), host}
}
