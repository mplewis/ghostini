// Package main runs a Gemini server that serves Ghost content.

package main

import (
	"context"
	"crypto/tls"
	_ "embed"
	"fmt"
	"log"
	"strings"

	"github.com/a-h/gemini"
	"github.com/mplewis/figyr"
	"github.com/mplewis/ghostini/ghost"
	"github.com/mplewis/ghostini/server"
)

type Config struct {
	GhostSite      string `figyr:"required,description=The base URL of your Ghost website"`
	ContentKey     string `figyr:"required,description=The Content API key for your Ghost website"`
	GeminiCertPath string `figyr:"required,description=The path to the certificate for your Gemini domain"`
	GeminiKeyPath  string `figyr:"required,description=The path to the key for your Gemini domain"`
}

// main starts the server.
func main() {
	var cfg Config
	figyr.MustParse(&cfg)

	if !(strings.HasPrefix(cfg.GhostSite, "http://") || strings.HasPrefix(cfg.GhostSite, "https://")) {
		log.Fatalf("Ghost site URL must start with http:// or https://")
	}

	cert, err := tls.LoadX509KeyPair(cfg.GeminiCertPath, cfg.GeminiKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	host := ghost.Host{SiteURL: cfg.GhostSite, ContentKey: cfg.ContentKey}
	host.SiteURL = strings.TrimSuffix(host.SiteURL, "/")

	server, err := server.New(host)
	if err != nil {
		log.Fatalf("Error connecting to %s: %s\n", host.SiteURL, err)
	}

	domain := gemini.NewDomainHandler("localhost", cert, server)
	fmt.Printf("Serving Ghost content from %s\n", host.SiteURL)
	log.Fatal(gemini.ListenAndServe(context.Background(), ":1965", domain))
}
