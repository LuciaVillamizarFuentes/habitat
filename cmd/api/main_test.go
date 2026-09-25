package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tu-usuario/habitat/internal/tenant"
	"github.com/tu-usuario/habitat/internal/unit"
)

func newTestMux() http.Handler {
	tenantSvc := tenant.NewService(tenant.NewMemoryRepository())
	unitSvc := unit.NewService(unit.NewMemoryRepository())
	return newMux(tenantSvc, unitSvc)
}

func TestNewMux_Healthz(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), `"status":"ok"`) {
		t.Fatalf("body = %s, want it to contain status ok", rec.Body)
	}
}

func TestNewMux_TenantRoutesWired(t *testing.T) {
	mux := newTestMux()

	body := `{"full_name":"Carlos Ruiz","email":"carlos@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
}

func TestNewMux_UnitRoutesWired(t *testing.T) {
	mux := newTestMux()

	body := `{"code":"T1-502","kind":"apartment","floor":5,"area_m2":80,"coefficient":0.015}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/units", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/units", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), "T1-502") {
		t.Fatalf("body = %s, want it to contain T1-502", rec.Body)
	}
}

func TestNewMux_UnknownRoute(t *testing.T) {
	mux := newTestMux()

	req := httptest.NewRequest(http.MethodGet, "/no-existe", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}
