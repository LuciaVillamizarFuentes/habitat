package tenant

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_CreateAndGet(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)

	body := `{"full_name":"Carlos Ruiz","email":"carlos@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/tenants/no-existe", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestHandler_RejectsUnknownFields(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)

	body := `{"full_name":"X","email":"x@example.com","is_admin":true}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}

func createTenant(t *testing.T, mux *http.ServeMux) string {
	t.Helper()
	body := `{"full_name":"Carlos Ruiz","email":"carlos@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}

	var created Tenant
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return created.ID
}

func TestHandler_Update(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)
	id := createTenant(t, mux)

	body := `{"full_name":"Carlos Actualizado"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tenants/"+id, strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), "Carlos Actualizado") {
		t.Fatalf("body = %s, want it to contain el nombre actualizado", rec.Body)
	}
}

func TestHandler_Update_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)

	body := `{"full_name":"X"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tenants/no-existe", strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404, body = %s", rec.Code, rec.Body)
	}
}

func TestHandler_Update_InvalidEmail(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)
	id := createTenant(t, mux)

	body := `{"email":"no-es-email"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tenants/"+id, strings.NewReader(body))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422, body = %s", rec.Code, rec.Body)
	}
}

func TestHandler_Delete(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)
	id := createTenant(t, mux)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tenants/"+id, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204, body = %s", rec.Code, rec.Body)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/tenants/"+id, nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), `"active":false`) {
		t.Fatalf("body = %s, want active:false tras soft delete", rec.Body)
	}
}

func TestHandler_Delete_NotFound(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tenants/no-existe", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404, body = %s", rec.Code, rec.Body)
	}
}
