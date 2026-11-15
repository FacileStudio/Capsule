package main

import (
	"log/slog"
	"net/http"

	"github.com/FacileStudio/tronc/middleware"
	"github.com/FacileStudio/tronc/spa"
)

const catalogOrigin = "https://api-ardoise.facile.studio"

// securityConfig extends tronc's default policy for Capsule's client: the
// inline scripts adapter-static writes, hashed from the build the binary
// serves, the Ardoise plan catalog the landing's pricing block fetches, and
// blob: media, which the reveal page's audio and video players load from an
// object URL of the decrypted file. The Vision tracker is a pinned copy the client serves itself, so script-src
// names no analytics origin. Capsule needs nothing else.
// Web Crypto, the Shiki regex engine and the lazy chunks are
// same-origin and need no unsafe-eval. The full policy is enforced.
func securityConfig(clientDir string, logger *slog.Logger) middleware.SecurityConfig {
	var sources []string
	if spa.Available(clientDir) {
		hashes, err := spa.InlineScriptHashes(clientDir)
		if err != nil {
			logger.Error("failed to hash the client's inline scripts", slog.Any("error", err))
		}
		sources = append(sources, hashes...)
	}
	return middleware.SecurityConfig{
		CSPEnforce: true,
		CSPSources: map[string][]string{
			"script-src":  sources,
			"connect-src": {catalogOrigin},
			"media-src":   {"blob:"},
		},
	}
}

// noReferrer stops the browser sending a Referer from the client's pages. A
// capsule link carries its decryption key in the fragment, which browsers never
// put in a Referer, but the path is a capsule id and no other site needs it.
// SecurityHeaders keeps a value the handler set.
func noReferrer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(middleware.HeaderReferrerPolicy, "no-referrer")
		next.ServeHTTP(w, r)
	})
}
