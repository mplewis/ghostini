package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/a-h/gemini"
	"github.com/mplewis/ghostini/cache"
	"github.com/mplewis/ghostini/ghost"
)

type Server struct{}

var slugMatcher = regexp.MustCompile(`^/([^/]+)/?$`)

func parseInt(s string, dfault int) int {
	if s == "" {
		return dfault
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return dfault
	}
	return i
}

func (s Server) ServeGemini(w gemini.ResponseWriter, r *gemini.Request) {
	if r.URL.Path == "/" {
		page := parseInt(r.URL.Query().Get("page"), 1)
		resp, err := ghost.GetPosts(c, h, page)
		if err != nil {
			w.SetHeader(gemini.CodeTemporaryFailure, "")
			return
		}
		w.SetHeader(gemini.CodeSuccess, "")
		renderIndex(w, h, resp)
		return
	}

	if matches := slugMatcher.FindStringSubmatch(r.URL.Path); len(matches) > 0 {
		slug := matches[1]
		resp, found, err := ghost.GetPost(c, h, slug)
		if err != nil {
			w.SetHeader(gemini.CodeTemporaryFailure, "")
			return
		}
		if !found {
			w.SetHeader(gemini.CodeNotFound, "")
			return
		}

		w.SetHeader(gemini.CodeSuccess, "")
		renderPost(w, resp.Posts[0])
		return
	}

	fmt.Println("invalid path")
	w.SetHeader(gemini.CodeNotFound, "")
	w.Write([]byte("not found"))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func mustEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	log.Fatalf("Missing mandatory environment variable %s", key)
	return ""
}

var c = cache.New()
var h = ghost.Host{
	APIURL:     mustEnv("GHOST_SITE"),
	ContentKey: mustEnv("CONTENT_KEY"),
}

func main() {
	if !(strings.HasPrefix(h.APIURL, "http://") || strings.HasPrefix(h.APIURL, "https://")) {
		log.Fatalf("GHOST_SITE must start with http:// or https://")
	}
	h.APIURL = strings.TrimSuffix(h.APIURL, "/")
	fmt.Printf("Starting server for Ghost site at %s\n", h.APIURL)

	cert, err := tls.LoadX509KeyPair("tmp/localhost.crt", "tmp/localhost.key")
	check(err)

	domain := gemini.NewDomainHandler("localhost", cert, Server{})
	err = gemini.ListenAndServe(context.Background(), ":1965", domain)
	if err != nil {
		log.Fatal("error:", err)
	}
}
