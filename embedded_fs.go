//go:build !dev
// +build !dev

/*
embedded_fs.go - Production filesystem handling for the Everato application

This file provides embedded filesystem functionality for serving static files and
migrations in production builds. It embeds the files directly into the binary,
eliminating the need for external file access during deployment.

Only used in production builds (!dev). For development, see live_fs.go.
*/
package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"os"
)

// PublicFS serves public files (CSS, JS, images) from the embedded filesystem in production.
// It returns an HTTP handler that serves static content directly from the binary.
//
// The //go:embed directive embeds the entire public directory into the compiled binary,
// allowing the application to serve static assets without external file dependencies.
//
// Returns:
//   - An http.Handler that serves embedded static files

//go:embed public
var publicFS embed.FS

func PublicFS() http.Handler {
	fs, err := fs.Sub(publicFS, "public")
	if err != nil {
		panic("Failed to create sub filesystem: " + err.Error())
	}

	// Return a file server that serves files from the embedded filesystem
	return http.FileServer(http.FS(fs))
}

//go:embed internal/db/migrations/*.sql
var migrationsFS embed.FS

// MigrationsFS provides access to embedded database migration SQL scripts.
// In production, migrations are embedded directly into the binary for reliable deployment.
//
// This function:
// 1. Verifies that the embedded migrations directory exists
// 2. Returns the embedded filesystem containing all migration scripts
//
// Returns:
//   - An fs.FS filesystem interface to access migration scripts
//   - An error if the embedded migrations directory is not found
func MigrationsFS() (fs.FS, error) {
	// Check if the migrations directory exists in the embedded filesystem
	_, err := fs.Stat(migrationsFS, "internal/db/migrations")
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("Migrations directory does not exist in the embedded filesystem")
	}

	return migrationsFS, nil
}
