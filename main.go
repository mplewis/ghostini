// Package main runs a Gemini server that serves Ghost content.

package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/a-h/gemini"
	"github.com/mplewis/ghostini/ghost"
	"github.com/mplewis/ghostini/server"
)

// parseInt parses a string into an integer, with a default fallback for any empty/error cases.
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

// check crashes if an error is present.
func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

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
	check(err)

	host := ghost.Host{APIURL: mustEnv("GHOST_SITE"), ContentKey: mustEnv("CONTENT_KEY")}
	host.APIURL = strings.TrimSuffix(host.APIURL, "/")
	if !(strings.HasPrefix(host.APIURL, "http://") || strings.HasPrefix(host.APIURL, "https://")) {
		log.Fatalf("GHOST_SITE must start with http:// or https://")
	}

	server := server.New(host)
	domain := gemini.NewDomainHandler("localhost", cert, server)
	fmt.Printf("Starting server for Ghost site at %s\n", host.APIURL)
	log.Fatal(gemini.ListenAndServe(context.Background(), ":1965", domain))
}
