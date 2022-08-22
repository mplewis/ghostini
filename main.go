// Package main runs a Gemini server that serves Ghost content.

package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"

	"github.com/mplewis/figyr"
	"github.com/mplewis/ghostini/server"
	"github.com/mplewis/ghostini/types"
)

// main starts the server.
func main() {
	var cfg types.Config
	figyr.MustParse(&cfg)

	server, err := server.NewGemini(cfg)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Serving Ghost content from %s on %s\n", cfg.GhostSite, server.Addr)
	log.Fatal(server.ListenAndServe(context.Background()))
}
