// Package main runs a Gemini server that serves Ghost content.

package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/a-h/gemini"
	"github.com/mplewis/ghostini/ghost"
	"github.com/mplewis/ghostini/server"
)

// mustEnv returns the value of an environment variable or crashes if it is not set.
func mustEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	log.Fatalf("Missing mandatory environment variable %s", key)
	return ""
}

// main starts the server.
func main() {
	cert, err := tls.LoadX509KeyPair("tmp/localhost.crt", "tmp/localhost.key")
	if err != nil {
		log.Fatal(err)
	}

	host := ghost.Host{APIURL: mustEnv("GHOST_SITE"), ContentKey: mustEnv("CONTENT_KEY")}
	host.APIURL = strings.TrimSuffix(host.APIURL, "/")
	if !(strings.HasPrefix(host.APIURL, "http://") || strings.HasPrefix(host.APIURL, "https://")) {
		log.Fatalf("GHOST_SITE must start with http:// or https://")
	}

	server, err := server.New(host)
	if err != nil {
		log.Fatalf("Error connecting to %s: %s\n", host.APIURL, err)
	}

	domain := gemini.NewDomainHandler("localhost", cert, server)
	fmt.Printf("Serving Ghost content from %s\n", host.APIURL)
	log.Fatal(gemini.ListenAndServe(context.Background(), ":1965", domain))
}
