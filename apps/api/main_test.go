package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	mux := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Se esperaba status 200 para /healthz, se obtuvo %d", rec.Code)
	}

	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("error leyendo body: %v", err)
	}

	if string(body) != "ok" {
		t.Errorf("Se esperaba 'ok', se obtuvo '%s'", string(body))
	}
}

func TestReadyz(t *testing.T) {
	mux := setupRouter()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	// Sin DB conectada da 503, con DB da 200
	if rec.Code != http.StatusOK && rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("Se esperaba status 200 o 503 para /readyz, se obtuvo %d", rec.Code)
	}
}

func TestRouterRoutes(t *testing.T) {
	mux := setupRouter()

	// Probar OPTIONS en /api/notificaciones y /api/v1/notificaciones
	for _, path := range []string{"/api/notificaciones", "/api/v1/notificaciones"} {
		req := httptest.NewRequest(http.MethodOptions, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("Se esperaba status 200 para OPTIONS %s, se obtuvo %d", path, rec.Code)
		}
	}
}
