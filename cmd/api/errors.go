package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeErrorJSON(w, http.StatusInternalServerError, err.Error())
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeErrorJSON(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeErrorJSON(w, http.StatusNotFound, err.Error())
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeErrorJSON(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedError(w http.ResponseWriter, r *http.Request, err error) {

	app.logger.Warnw("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeErrorJSON(w, http.StatusUnauthorized, err.Error())
}

func (app *application) unauthorizedBasicError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	w.Header().Set("WWW-Authenticate", `Basic realm="Restricted", charset="UTF-8"`)

	writeErrorJSON(w, http.StatusUnauthorized, err.Error())
}
