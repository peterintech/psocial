package main

import (
	"net/http"

	"github.com/peterintech/psocial/internal/store"
)

// GetFeed godoc
//
//	@Summary		Get user feed
//	@Description	Fetches paginated posts from followed users and own posts
//	@Tags			feed
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit"			default(10)
//	@Param			offset	query		int		false	"Offset"		default(0)
//	@Param			sort	query		string	false	"Sort order"	Enums(asc, desc)	default(desc)
//	@Param			tags	query		string	false	"Comma-separated tags"
//	@Param			search	query		string	false	"Search term"
//	@Success		200		{object}	[]store.PostWithMetadata
//	@Failure		400		{object}	error
//	@Failure		500		{object}	error
//	@Router			/users/feed [get]
func (app *application) getFeedHandler(w http.ResponseWriter, r *http.Request) {

	pagination := store.PaginatedFeedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := pagination.Parse(r)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeeds(ctx, "35", fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
