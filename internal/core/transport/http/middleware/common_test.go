package core_http_middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORS(t *testing.T) {
	t.Setenv(corsAllowedOriginsEnv, "https://app.example.com, http://localhost:3000")

	calls := 0
	handler := CORS()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusOK)
	}))

	t.Run("allowed origin receives CORS headers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
		req.Header.Set("Origin", "https://app.example.com")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		if res.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", res.Code, http.StatusOK)
		}
		if got := res.Header().Get("Access-Control-Allow-Origin"); got != "https://app.example.com" {
			t.Fatalf("Access-Control-Allow-Origin = %q", got)
		}
		if got := res.Header().Get("Vary"); got != "Origin" {
			t.Fatalf("Vary = %q, want Origin", got)
		}
	})

	t.Run("disallowed origin receives no CORS permission", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
		req.Header.Set("Origin", "https://attacker.example")
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		if got := res.Header().Get("Access-Control-Allow-Origin"); got != "" {
			t.Fatalf("unexpected Access-Control-Allow-Origin: %q", got)
		}
	})

	t.Run("valid preflight is accepted without calling handler", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/v1/tasks", nil)
		req.Header.Set("Origin", "http://localhost:3000")
		req.Header.Set("Access-Control-Request-Method", http.MethodPost)
		req.Header.Set("Access-Control-Request-Headers", "content-type, authorization")
		res := httptest.NewRecorder()
		callsBefore := calls

		handler.ServeHTTP(res, req)

		if res.Code != http.StatusNoContent {
			t.Fatalf("status = %d, want %d", res.Code, http.StatusNoContent)
		}
		if calls != callsBefore {
			t.Fatal("preflight unexpectedly reached the next handler")
		}
	})

	t.Run("disallowed preflight is rejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/api/v1/tasks", nil)
		req.Header.Set("Origin", "https://attacker.example")
		req.Header.Set("Access-Control-Request-Method", http.MethodPost)
		res := httptest.NewRecorder()

		handler.ServeHTTP(res, req)

		if res.Code != http.StatusForbidden {
			t.Fatalf("status = %d, want %d", res.Code, http.StatusForbidden)
		}
	})

	t.Run("unsupported method or header is rejected", func(t *testing.T) {
		for _, tc := range []struct {
			name   string
			method string
			header string
		}{
			{name: "method", method: http.MethodPut},
			{name: "header", method: http.MethodPost, header: "x-custom-token"},
		} {
			t.Run(tc.name, func(t *testing.T) {
				req := httptest.NewRequest(http.MethodOptions, "/api/v1/tasks", nil)
				req.Header.Set("Origin", "https://app.example.com")
				req.Header.Set("Access-Control-Request-Method", tc.method)
				if tc.header != "" {
					req.Header.Set("Access-Control-Request-Headers", tc.header)
				}
				res := httptest.NewRecorder()

				handler.ServeHTTP(res, req)

				if res.Code != http.StatusForbidden {
					t.Fatalf("status = %d, want %d", res.Code, http.StatusForbidden)
				}
			})
		}
	})
}

func TestParseAllowedOrigins(t *testing.T) {
	origins := parseAllowedOrigins(" https://a.example, ,https://b.example, * ")
	if len(origins) != 2 {
		t.Fatalf("got %d origins, want 2", len(origins))
	}
	for _, origin := range []string{"https://a.example", "https://b.example"} {
		if _, ok := origins[origin]; !ok {
			t.Errorf("origin %q missing", origin)
		}
	}
}
