package tenant

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tu-usuario/habitat/internal/platform/httpx"
)

// Handler expone los casos de uso de inquilinos por HTTP.
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Routes registra las rutas usando el router de la stdlib (Go 1.22+).
func (h *Handler) Routes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/tenants", h.create)
	mux.HandleFunc("GET /api/v1/tenants", h.list)
	mux.HandleFunc("GET /api/v1/tenants/{id}", h.get)
	mux.HandleFunc("PATCH /api/v1/tenants/{id}", h.update)
	mux.HandleFunc("DELETE /api/v1/tenants/{id}", h.delete)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var in CreateInput
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	t, err := h.svc.Register(r.Context(), in)
	if err != nil {
		writeErr(w, err)
		return
	}
	httpx.JSON(w, http.StatusCreated, t)
}

// optString devuelve un puntero al valor si el param vino, o nil si no vino.
func optString(q url.Values, key string) *string {
	v := strings.TrimSpace(q.Get(key)) // quita espacios: " Ana " → "Ana"
	if v == "" {
		return nil // no vino → no filtrar por este campo
	}
	return &v
}

func parseFilter(q url.Values) (FilterTenants, error) {
	f := FilterTenants{
		FullName: optString(q, "full_name"),
		Email:    optString(q, "email"),
		Phone:    optString(q, "phone"),
		UnitID:   optString(q, "unit_id"),
	}

	// Los emails se guardan en minúsculas, así que el filtro también
	if f.Email != nil {
		lower := strings.ToLower(*f.Email)
		f.Email = &lower
	}

	// Active: texto → bool → puntero
	if v := q.Get("active"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return FilterTenants{}, errors.New("active must be true or false")
		}
		f.Active = &b
	}
	return f, nil
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f, err := parseFilter(q)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	var page uint64 // si no viene, queda en 0 (primera página)
	if v := q.Get("page"); v != "" {
		p, err := strconv.ParseUint(v, 10, 64) // base 10, 64 bits
		if err != nil || p == 0 {              // no es un número positivo
			httpx.Error(w, http.StatusBadRequest, "page must be a positive number")
			return
		}
		page = p
	}

	var limit uint64 = 2 // valor por defecto
	if v := q.Get("limit"); v != "" {
		l, err := strconv.ParseUint(v, 10, 64) // base 10, 64 bits
		if err != nil || l == 0 || l > 100 {   // no es un número positivo
			httpx.Error(w, http.StatusBadRequest, "limit must be a valid number")
			return
		}
		limit = l
	}

	ts, err := h.svc.List(r.Context(), limit, page, f)
	if err != nil {
		writeErr(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"data": ts})
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	t, err := h.svc.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, t)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	var in UpdateInput
	if err := httpx.Decode(r, &in); err != nil {
		httpx.Error(w, http.StatusBadRequest, err.Error())
		return
	}
	t, err := h.svc.Update(r.Context(), r.PathValue("id"), in)
	if err != nil {
		writeErr(w, err)
		return
	}
	httpx.JSON(w, http.StatusOK, t)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	err := h.svc.Delete(r.Context(), r.PathValue("id"))
	if err != nil {
		writeErr(w, err)
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}

// writeErr traduce errores de dominio a códigos HTTP.
func writeErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, ErrInvalidName), errors.Is(err, ErrInvalidEmail):
		httpx.Error(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, ErrDuplicateMail):
		httpx.Error(w, http.StatusConflict, err.Error())
	default:
		httpx.Error(w, http.StatusInternalServerError, "internal server error")
	}
}
