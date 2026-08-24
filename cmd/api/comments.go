package main

import (
	"net/http"

	"github.com/peterintech/psocial/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required"`
}

func (app *application) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post, err := getPostFromContext(r.Context())
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	comment := &store.Comment{
		Content: payload.Content,
		PostID:  post.ID,
		UserID:  post.UserID,
	}

	if err := app.store.Comments.Create(r.Context(), comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
