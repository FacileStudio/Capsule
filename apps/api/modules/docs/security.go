package docs

import (
	"github.com/FacileStudio/tronc/middleware"
	"github.com/go-chi/chi/v5"
)

const scalarRegistryOrigin = "https://api.scalar.com"

// referenceRouter returns router with the reference page's own CSP additions,
// on top of the ones tronc's apiref already makes for the Scalar bundle and
// its fonts. Scalar evaluates code at start-up and queries its public API
// registry for the search box. Both are allowed on /docs and /docs/openapi.json
// only: the app's own policy keeps script-src without unsafe-eval.
func referenceRouter(router chi.Router) chi.Router {
	return router.With(middleware.SecurityOverride(middleware.SecurityOverrideConfig{
		CSPSources: map[string][]string{
			"script-src":  {"'unsafe-eval'"},
			"connect-src": {scalarRegistryOrigin},
		},
	}))
}
