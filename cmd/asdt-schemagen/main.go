// Package main is the composition root for the asdt-schemagen binary.
// It regenerates the machine-generated inline schema blocks in the producing
// specialist step .md files from their canonical schemas/*.schema.yaml sources.
// Run from the repository root via `make generate`.
package main

import (
	"fmt"
	"os"

	"github.com/vitualizz/asdt/internal/installer"
)

func main() {
	root := "."
	if err := installer.GenerateSchemaRegions(root); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
