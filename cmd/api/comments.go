package main

import (
	"net/http"

	"github.com/peterintech/psocial/internal/store"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required"`
}

// CreateComment godoc
//
//	@Summary		Create a comment
//	@Description	Creates a new comment on a post
//	@Tags			comments
//	@Accept			json
//	@Produce		json
//	@Param			postID	path		string					true	"Post ID"
//	@Param			payload	body		CreateCommentPayload	true	"Comment payload"
//	@Success		201		{object}	store.Comment
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/posts/{postID}/comments [post]
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
	user := app.getUserFromContext(r)

	comment := &store.Comment{
		Content: payload.Content,
		PostID:  post.ID,
		UserID:  user.ID,
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
