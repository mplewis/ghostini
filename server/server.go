// Package server implements a Gemini server that serves Ghost content.

package server

import (
	"context"
	_ "embed"
	"regexp"
	"time"

	"git.sr.ht/~adnano/go-gemini"
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
func (s Server) ServeGemini(ctx context.Context, w gemini.ResponseWriter, r *gemini.Request) {
	if r.URL.Path == "/" {
		page := parse.Int(r.URL.Query().Get("page"), 1)
		resp, err := ghost.GetPosts(s.cache, s.host, page)
		if err != nil {
			w.WriteHeader(gemini.StatusTemporaryFailure, "")
			return
		}
		w.WriteHeader(gemini.StatusSuccess, "")
		render.Index(w, s.host, resp)
		return
	}

	if matches := slugMatcher.FindStringSubmatch(r.URL.Path); len(matches) > 0 {
		slug := matches[1]
		resp, found, err := ghost.GetPost(s.cache, s.host, slug)
		if err != nil {
			w.WriteHeader(gemini.StatusTemporaryFailure, "")
			return
		}
		if !found {
			w.WriteHeader(gemini.StatusNotFound, "")
			return
		}
		w.WriteHeader(gemini.StatusSuccess, "")
		render.Post(w, resp.Posts[0])
		return
	}

	w.WriteHeader(gemini.StatusNotFound, "")
	w.Write([]byte("not found"))
}

// New creates a new Gemini server.
func New(host ghost.Host) (*gemini.Server, error) {
	s := Server{cache.New(), host}
	// warm cache and verify connectivity
	_, err := ghost.GetPosts(s.cache, s.host, 1)
	if err != nil {
		return nil, err
	}

	mux := &gemini.Mux{}
	mux.Handle("/", s)

	server := &gemini.Server{
		Handler:      gemini.LoggingMiddleware(mux),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 1 * time.Minute,
	}
	return server, nil
}
