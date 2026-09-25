package unit

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_CreateAndGet(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)

	body := `{"code":"T1-502","kind":"apartment","floor":5,"area_m2":80,"coefficient":0.015}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/units", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/units/no-existe", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestHandler_RejectsUnknownFields(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)

	body := `{"code":"T1-502","kind":"apartment","floor":5,"area_m2":80,"coefficient":0.015,"is_admin":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/units", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func TestHandler_RejectsInvalidInput(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)

	body := `{"code":"","kind":"apartment","floor":5,"area_m2":80,"coefficient":0.015}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/units", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422, body = %s", rec.Code, rec.Body)
	}
}

func TestHandler_RejectsDuplicateCode(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)

	body := `{"code":"T1-502","kind":"apartment","floor":5,"area_m2":80,"coefficient":0.015}`

	req := httptest.NewRequest(http.MethodPost, "/api/v1/units", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}

	req = httptest.NewRequest(http.MethodPost, "/api/v1/units", strings.NewReader(body))
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409, body = %s", rec.Code, rec.Body)
	}
}

func TestHandler_List(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)

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
