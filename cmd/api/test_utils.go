package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/peterintech/psocial/internal/auth"
	"github.com/peterintech/psocial/internal/ratelimiter"
	"github.com/peterintech/psocial/internal/store"
	"github.com/peterintech/psocial/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T, cfg config) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()
	rateLimiter := ratelimiter.NewFixedWindowRateLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	testAuth := &auth.TestAuthenticator{}

	app := &application{
		logger:        logger,
		config:        cfg,
		store:         mockStore,
		cacheStorage:  mockCacheStore,
		rateLimiter:   rateLimiter,
		authenticator: testAuth,
	}

	return app
}

func executeRequest(req *http.Request, mux *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
