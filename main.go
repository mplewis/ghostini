package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/a-h/gemini"
)

type Server struct{}

var allPostsMatcher = regexp.MustCompile(`^/posts/?$`)
var slugMatcher = regexp.MustCompile(`^/posts/([^/]+)/?$`)

// /posts => get page 1 of all posts
// /posts?page=2 => get page 2 of all posts
// /posts/some-blog-post-slug => get post with slug "some-blog-post-slug"
func (s Server) ServeGemini(w gemini.ResponseWriter, r *gemini.Request) {
	fmt.Println(r.URL.Path)

	if r.URL.Path == "/" {
		fmt.Println("home")
		w.SetHeader(gemini.CodeRedirect, "/posts")
		return
	}

	if allPostsMatcher.MatchString(r.URL.Path) {
		fmt.Println("all posts")
		w.SetHeader(gemini.CodeSuccess, "")
		w.Write([]byte("all posts"))
		return
	}

	if matches := slugMatcher.FindStringSubmatch(r.URL.Path); len(matches) > 0 {
		slug := matches[1]
		fmt.Printf("one post: %s\n", slug)
		w.SetHeader(gemini.CodeSuccess, "")
		w.Write([]byte(fmt.Sprintf("one post: %s", slug)))
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

var c = newCache()
var h = host{
	apiUrl:     mustEnv("GHOST_SITE"),
	contentKey: mustEnv("CONTENT_KEY"),
}

func main() {
	if !(strings.HasPrefix(h.apiUrl, "http://") || strings.HasPrefix(h.apiUrl, "https://")) {
		log.Fatalf("GHOST_SITE must start with http:// or https://")
	}
	h.apiUrl = strings.TrimSuffix(h.apiUrl, "/")
	fmt.Printf("Starting server for Ghost site at %s\n", h.apiUrl)

	// resp, err := getPosts(c, h, 1)
	// check(err)
	// fmt.Println(resp)
	// resp, err = getPosts(c, h, 1)
	// check(err)
	// fmt.Println(resp)

	cert, err := tls.LoadX509KeyPair("tmp/localhost.crt", "tmp/localhost.key")
	check(err)

	domain := gemini.NewDomainHandler("localhost", cert, Server{})
	err = gemini.ListenAndServe(context.Background(), ":1965", domain)
	if err != nil {
		log.Fatal("error:", err)
	}
}
