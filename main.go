// Package main runs a Gemini server that serves Ghost content.

package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"strings"

	"git.sr.ht/~adnano/go-gemini/certificate"
	"github.com/mplewis/figyr"
	"github.com/mplewis/ghostini/ghost"
	"github.com/mplewis/ghostini/server"
)

type Config struct {
	GhostSite  string `figyr:"required,description=The base URL of your Ghost website"`
	ContentKey string `figyr:"required,description=The Content API key for your Ghost website"`
	// GeminiCertsPath string `figyr:"required,description=The path to the certificates and keys for your Gemini domains"`
}

// main starts the server.
func main() {
	var cfg Config
	figyr.MustParse(&cfg)

	if !(strings.HasPrefix(cfg.GhostSite, "http://") || strings.HasPrefix(cfg.GhostSite, "https://")) {
		log.Fatalf("Ghost site URL must start with http:// or https://")
	}

	// TODO: Load certs dynamically from directory
	certs := &certificate.Store{}
	domains := []string{"localhost", "kesdev.com"}
	for _, domain := range domains {
		certs.Register(domain)
	}

	host := ghost.Host{SiteURL: cfg.GhostSite, ContentKey: cfg.ContentKey}
	host.SiteURL = strings.TrimSuffix(host.SiteURL, "/")

	server, err := server.New(host)
	if err != nil {
		log.Fatalf("Error connecting to %s: %s\n", host.SiteURL, err)
	}
	// TODO: untangle initialization
	server.GetCertificate = certs.Get
	server.Addr = ":1965"

	fmt.Printf("Serving Ghost content from %s on %s\n", host.SiteURL, server.Addr)
	log.Fatal(server.ListenAndServe(context.Background()))
}
