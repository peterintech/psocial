package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/peterintech/psocial/internal/store"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		app.badRequestError(w, r, errors.New("userID is required"))
		return
	}

	ctx := r.Context()
	user, err := app.store.Users.GetByID(ctx, userID)
	if err != nil {
		switch err {
		case store.ErrNotFound:
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
