package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/peterintech/psocial/internal/store"
	"go.uber.org/zap"
)

type conflictingPostStore struct{}

func (*conflictingPostStore) Create(context.Context, *store.Post) error {
	return nil
}

func (*conflictingPostStore) GetByID(context.Context, string) (*store.Post, error) {
	return nil, store.ErrNotFound
}

func (*conflictingPostStore) Update(context.Context, *store.Post) error {
	return store.ErrVersionConflict
}

func (*conflictingPostStore) Delete(context.Context, string) error {
	return nil
}

func (*conflictingPostStore) GetUserFeeds(context.Context, string, store.PaginatedFeedQuery) ([]*store.PostWithMetadata, error) {
	return nil, nil
}

func TestUpdatePostReturnsConflictForStaleVersion(t *testing.T) {
	app := &application{
		store:  store.Storage{Posts: &conflictingPostStore{}},
		logger: zap.NewNop().Sugar(),
	}
	post := &store.Post{
		ID:      "42",
		Title:   "Original title",
		Content: "Original content",
		Version: 3,
	}

	req := httptest.NewRequest(http.MethodPatch, "/v1/posts/42", strings.NewReader(`{"title":"Updated title"}`))
	req = req.WithContext(context.WithValue(req.Context(), postContextKey, post))
	recorder := httptest.NewRecorder()

	app.updatePostHandler(recorder, req)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected status %d, got %d", http.StatusConflict, recorder.Code)
	}
	if !strings.Contains(recorder.Body.String(), store.ErrVersionConflict.Error()) {
		t.Fatalf("expected conflict response to contain %q, got %q", store.ErrVersionConflict, recorder.Body.String())
	}
}
