package main

import (
	"net/http"

	"github.com/peterintech/psocial/internal/store"
)

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
	}

	ctx := r.Context()

	feed, err := app.store.Posts.GetUserFeeds(ctx, "23", fq)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
