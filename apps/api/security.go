package main

import (
	"log/slog"
	"net/http"

	"github.com/FacileStudio/tronc/middleware"
	"github.com/FacileStudio/tronc/spa"
)

const analyticsOrigin = "https://vision.facile.studio"

// securityConfig extends tronc's default policy for Capsule's client: the
// inline scripts adapter-static writes, hashed from the build the binary
// serves, and the Vision analytics script app.html loads. Capsule needs
// nothing else. Web Crypto, the Shiki regex engine and the lazy chunks are
// same-origin and need no unsafe-eval. The policy stays report-only.
func securityConfig(clientDir string, logger *slog.Logger) middleware.SecurityConfig {
	sources := []string{analyticsOrigin}
	if spa.Available(clientDir) {
		hashes, err := spa.InlineScriptHashes(clientDir)
		if err != nil {
			logger.Error("failed to hash the client's inline scripts", slog.Any("error", err))
		}
		sources = append(sources, hashes...)
	}
	return middleware.SecurityConfig{CSPSources: map[string][]string{"script-src": sources}}
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
