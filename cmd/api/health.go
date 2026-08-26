package main

import (
	"net/http"
)

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Returns service status, environment, and version
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	error
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "ok", "env": app.config.env, "version": version}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		writeErrorJSON(w, http.StatusInternalServerError, "error writing JSON response")
	}
}
