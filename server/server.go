// Package server implements a Gemini server that serves Ghost content.

package server

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"git.sr.ht/~adnano/go-gemini"
	"git.sr.ht/~adnano/go-gemini/certificate"
	"github.com/mplewis/ghostini/cache"
	"github.com/mplewis/ghostini/ghost"
	"github.com/mplewis/ghostini/parse"
	"github.com/mplewis/ghostini/render"
	"github.com/mplewis/ghostini/types"
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

// NewGemini creates a new Gemini server.
func NewGemini(cfg types.Config) (*gemini.Server, error) {
	if !(strings.HasPrefix(cfg.GhostSite, "http://") || strings.HasPrefix(cfg.GhostSite, "https://")) {
		log.Fatalf("Ghost site URL must start with http:// or https://")
	}

	// load certs
	certs := &certificate.Store{}
	for _, domain := range strings.Split(cfg.Domains, ",") {
		certs.Register(domain)
		fmt.Printf("Registered domain %s\n", domain)
	}
	if cfg.GeminiCertsPath != "" {
		err := certs.Load(cfg.GeminiCertsPath)
		if err != nil {
			return nil, err
		}
		fmt.Printf("Loaded certificates from %s\n", cfg.GeminiCertsPath)
	}

	host := ghost.Host{SiteURL: cfg.GhostSite, ContentKey: cfg.ContentKey}
	host.SiteURL = strings.TrimSuffix(host.SiteURL, "/")

	s := Server{cache.New(), host}
	// warm cache and verify connectivity
	_, err := ghost.GetPosts(s.cache, s.host, 1)
	if err != nil {
		return nil, err
	}

	mux := &gemini.Mux{}
	mux.Handle("/", s)

	server := &gemini.Server{
		Handler:        gemini.LoggingMiddleware(mux),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   1 * time.Minute,
		GetCertificate: certs.Get,
		Addr:           fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	}
	return server, nil
}
