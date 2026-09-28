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

func seedThreeTenants(t *testing.T, mux *http.ServeMux) {
	t.Helper()
	bodies := []string{
		`{"full_name":"Ana Gómez","email":"ana@example.com","unit_id":"U1"}`,
		`{"full_name":"Carlos Ruiz","email":"carlos@example.com","unit_id":"U2"}`,
		`{"full_name":"Beatriz Soto","email":"beatriz@example.com","unit_id":"U1"}`,
	}
	for _, body := range bodies {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants", strings.NewReader(body))
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("seed status = %d, body = %s", rec.Code, rec.Body)
		}
	}
}

func TestHandler_List_NoParams(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)
	seedThreeTenants(t, mux)

	// Sin "page" en la query, el handler debe usar la primera página por
	// defecto en lugar de propagar page=0 y hacer panic en el repositorio.
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
}

func TestHandler_List_FilterByUnitID(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)
	seedThreeTenants(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants?unit_id=U1&limit=10", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), "ana@example.com") || !strings.Contains(rec.Body.String(), "beatriz@example.com") {
		t.Fatalf("body = %s, want ana y beatriz (unit_id=U1)", rec.Body)
	}
	if strings.Contains(rec.Body.String(), "carlos@example.com") {
		t.Fatalf("body = %s, no debería incluir a carlos (unit_id=U2)", rec.Body)
	}
}

func TestHandler_List_Pagination(t *testing.T) {
	mux := http.NewServeMux()
	NewHandler(NewService(NewMemoryRepository())).Routes(mux)
	seedThreeTenants(t, mux)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants?limit=2&page=1", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), "ana@example.com") || !strings.Contains(rec.Body.String(), "carlos@example.com") {
		t.Fatalf("página 1 = %s, want ana y carlos", rec.Body)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/v1/tenants?limit=2&page=2", nil)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	if !strings.Contains(rec.Body.String(), "beatriz@example.com") {
		t.Fatalf("página 2 = %s, want beatriz", rec.Body)
	}
}

func TestHandler_List_InvalidQueryParams(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{name: "active no booleano", query: "?active=quizas"},
		{name: "page en 0", query: "?page=0"},
		{name: "page no numérico", query: "?page=abc"},
		{name: "limit en 0", query: "?limit=0"},
		{name: "limit mayor a 100", query: "?limit=101"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			NewHandler(NewService(NewMemoryRepository())).Routes(mux)

			req := httptest.NewRequest(http.MethodGet, "/api/v1/tenants"+tt.query, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400, body = %s", rec.Code, rec.Body)
			}
		})
	}
}
