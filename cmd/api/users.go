package main

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/peterintech/psocial/internal/store"
)

type userKey string

const userContextKey userKey = "user"

type publicUserProfile struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

type relationshipResponse struct {
	IsFollowing bool `json:"is_following"`
}

// GetUser godoc
//
//	@Summary		Fetches a user profile
//	@Description	Fetches a user profile by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		string	true	"User ID"
//	@Success		200		{object}	store.User
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	if userID == "" {
		app.badRequestError(w, r, errors.New("userID is required"))
		return
	}

	user, err := app.getUser(r.Context(), userID)
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

	profile := publicUserProfile{}
	if user != nil {
		profile = publicUserProfile{ID: user.ID, Username: user.Username, CreatedAt: user.CreatedAt}
	}
	if err := app.jsonResponse(w, http.StatusOK, profile); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// getCurrentUserHandler returns the authenticated user's private account details.
func (app *application) getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	if err := app.jsonResponse(w, http.StatusOK, app.getUserFromContext(r)); err != nil {
		app.internalServerError(w, r, err)
	}
}

// getRelationshipHandler returns whether the authenticated user follows the profile.
func (app *application) getRelationshipHandler(w http.ResponseWriter, r *http.Request) {
	viewer := app.getUserFromContext(r)
	targetID := chi.URLParam(r, "userID")
	if viewer.ID == targetID {
		_ = app.jsonResponse(w, http.StatusOK, relationshipResponse{IsFollowing: false})
		return
	}
	if _, err := app.getUser(r.Context(), targetID); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			app.notFoundError(w, r, err)
		} else {
			app.internalServerError(w, r, err)
		}
		return
	}
	following, err := app.store.Followers.IsFollowing(r.Context(), viewer.ID, targetID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}
	if err := app.jsonResponse(w, http.StatusOK, relationshipResponse{IsFollowing: following}); err != nil {
		app.internalServerError(w, r, err)
	}
}

type FollowUser struct {
	UserID string `json:"user_id" validate:"required"`
}

// FollowUser godoc
//
//	@Summary		Follow a user
//	@Description	Follows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	string		true	"User ID"
//	@Param			payload	body	FollowUser	true	"Follow payload"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		409	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := app.getUserFromContext(r)
	followedID := chi.URLParam(r, "userID")

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followerUser.ID, followedID); err != nil { // todo: change after auth
		switch err {
		case store.ErrConflict:
			app.conflictError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// UnfollowUser godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollows a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	string		true	"User ID"
//	@Param			payload	body	FollowUser	true	"Unfollow payload"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	unFollowedUser := app.getUserFromContext(r)
	unfollowedID := chi.URLParam(r, "userID")

	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, unFollowedUser.ID, unfollowedID); err != nil { // todo: change after auth
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// ActivateUser godoc
//
//	@Summary		Activate a user account
//	@Description	Activates a user account using an activation token
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			token	path	string	true	"Activation Token"
//	@Success		204
//	@Failure		400	{object}	error
//	@Failure		404	{object}	error
//	@Failure		500	{object}	error
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")
	if token == "" {
		app.badRequestError(w, r, errors.New("token is required"))
		return
	}

	ctx := r.Context()

	if err := app.store.Users.Activate(ctx, token); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getUserFromContext(r *http.Request) *store.User {
	user := r.Context().Value(userContextKey)
	if user == nil {
		return nil
	}
	return user.(*store.User)
}
