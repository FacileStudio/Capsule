// Package migrations embeds the ordered SQL that owns Capsule's schema.
//
// It is a package rather than a bare directory so that the tests can reach the
// same migrations the binary runs. Embedding them in package main would leave
// every test package building its own schema from the GORM struct tags, which
// is exactly the drift this replaced.
package migrations

import "embed"

// FS holds the migration files at its root, which is the layout
// tronc/migrate expects — no fs.Sub needed.
//
//go:embed *.sql
var FS embed.FS
