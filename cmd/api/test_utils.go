package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/peterintech/psocial/internal/store"
	"github.com/peterintech/psocial/internal/store/cache"
	"go.uber.org/zap"
)

func newTestApplication(t *testing.T) *application {
	t.Helper()

	logger := zap.NewNop().Sugar()
	mockStore := store.NewMockStore()
	mockCacheStore := cache.NewMockStore()

	app := &application{
		logger: logger,
		config: config{
			redisCfg: redisConfig{
				enabled: false,
			},
		},
		store:        mockStore,
		cacheStorage: mockCacheStore,
	}

	return app
}

func executeRequest(req *http.Request, mux *chi.Mux) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}
