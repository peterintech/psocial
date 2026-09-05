package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/peterintech/psocial/internal/ratelimiter"
)

func TestRateLimiterMiddleware(t *testing.T) {
	cfg := config{
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: 20,
			TimeFrame:            time.Hour,
			Enabled:              true,
		},
		addr: ":8080",
	}

	app := newTestApplication(t, cfg)
	ts := httptest.NewServer(app.mount())
	defer ts.Close()

	client := &http.Client{}
	baseUrl := ts.URL + "/v1/auth/register" // Using the /token endpoint for testing
	mockIP := "192.168.1.1"
	marginOfError := 2
	body := []byte(`{"username":"testuser", "email":"testuser@gmail.com","password":"password"}`)

	for i := 0; i < cfg.rateLimiter.RequestsPerTimeFrame+marginOfError; i++ {
		req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("could not create request: %v", err)
		}

		req.Header.Set("X-Forwarded-For", mockIP)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("could not send request: %v", err)
		}
		resp.Body.Close()

		if i < cfg.rateLimiter.RequestsPerTimeFrame {
			if resp.StatusCode != http.StatusCreated {
				t.Errorf("expected status OK; got %v", resp.Status)
			}
		} else {
			if resp.StatusCode != http.StatusTooManyRequests {
				t.Errorf("expected status Too Many Requests; got %v", resp.Status)
			}
		}
	}
}
