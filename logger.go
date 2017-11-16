package main

import (
	"fmt"
	"log"
	"os"

	"comail.io/go/colog"
)

// DefaultStackDepth is the default level to plumb to when logging
const DefaultStackDepth = 1

// NewCoLogLogger creates a new CoLog based logger
func NewCoLogLogger(domain string) (logger *log.Logger) {
	cl := colog.NewCoLog(os.Stdout, fmt.Sprintf("%s ", domain), log.LstdFlags|log.Lshortfile)

	return cl.NewLogger()
}
