package docs

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FacileStudio/tronc/apiref"
	"github.com/FacileStudio/tronc/httpjson"
	"github.com/go-chi/chi/v5"
	"gopkg.in/yaml.v3"
)

//go:embed openapi.yaml
var specYAML []byte

// spec returns the embedded document as JSON.
//
// The YAML stays the source of truth: it is hand-written and carries schemas,
// examples and prose that the generated route registries used elsewhere in the
// suite cannot express. It is converted rather than served as-is so that
// /docs/openapi.json means the same thing in every Facile app.
func spec() (json.RawMessage, error) {
	var document any
	if err := yaml.Unmarshal(specYAML, &document); err != nil {
		return nil, fmt.Errorf("openapi.yaml is not valid YAML: %w", err)
	}
	encoded, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("openapi.yaml does not convert to JSON: %w", err)
	}
	return encoded, nil
}

// RegisterRoutes mounts the reference page and the spec. Mount it on the root
// router, beside /api and before the SPA catch-all.
func RegisterRoutes(router chi.Router) error {
	document, err := spec()
	if err != nil {
		return err
	}

	router.Get(apiref.BasePath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "public, max-age=3600")
		_, _ = w.Write([]byte(page))
	})

	router.Get(apiref.SpecPath, func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		httpjson.WriteJSON(w, http.StatusOK, document)
	})

	return nil
}

// page mirrors the reference page tronc/apiref serves, down to the pinned
// Scalar bundle, so Capsule's hand-written spec renders like the generated ones.
const page = `<!doctype html>
<html>
<head>
  <title>Capsule API</title>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <style>body { margin: 0; }</style>
</head>
<body>
  <script id="api-reference" data-url="` + apiref.SpecPath + `"></script>
  <script src="` + apiref.ScalarScriptURL + `"></script>
</body>
</html>`
