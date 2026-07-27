// Package backoffice serves the static admin UI. Infra, no model/service —
// same exception as health/docs/apierror. It holds no business logic: every
// page talks to the regular /api endpoints from the browser.
package backoffice

import (
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts the UI at /backoffice. Unlike the API domains it is
// registered on the engine root rather than the /api group: it is a UI, not
// part of the API surface, and its assets should not appear under /api.
func RegisterRoutes(r *gin.Engine) error {
	sub, err := fs.Sub(assets, "assets")
	if err != nil {
		return err
	}

	r.StaticFS("/backoffice", http.FS(sub))
	return nil
}
