package main

import (
	"net/http"

	"github.com/FacileStudio/tronc/httpjson"
)

var version = "dev"

// versionHandler answers the release the binary was built from. The Dockerfile
// stamps it from apps/client/package.json through scripts/version.sh.
func versionHandler(w http.ResponseWriter, _ *http.Request) {
	httpjson.WriteJSON(w, http.StatusOK, map[string]string{"version": version})
}
