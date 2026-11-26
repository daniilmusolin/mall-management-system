package handler

import (
	"encoding/json"
	"net/http"

	"mall-api/internal/domain"
	"mall-api/internal/service"
)

type TenantHandler struct {
	svc *service.TenantService
}

func NewTenantHandler(s *service.TenantService) *TenantHandler {
	return &TenantHandler{svc: s}
}

func (h *TenantHandler) RegisterRoutes(mux *http.ServeMux, auth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/tenants", auth(RequireRole(domain.RoleAdmin, domain.RoleManager)(http.HandlerFunc(h.GetAll))))
	mux.Handle("POST /api/tenants", auth(RequireRole(domain.RoleAdmin, domain.RoleManager)(http.HandlerFunc(h.Create))))
	mux.Handle("PUT /api/tenants/{id}", auth(RequireRole(domain.RoleAdmin, domain.RoleManager)(http.HandlerFunc(h.Update))))
	mux.Handle("DELETE /api/tenants/{id}", auth(RequireRole(domain.RoleAdmin)(http.HandlerFunc(h.Delete))))
}

func (h *TenantHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.GetAll(r.Context(), ParsePagination(r))
	if err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, items)
}

func (h *TenantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var t domain.Tenant
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	if err := h.svc.Create(r.Context(), &t); err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusCreated, t)
}

func (h *TenantHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "некорректный id")
		return
	}
	var t domain.Tenant
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	t.ID = id
	if err := h.svc.Update(r.Context(), &t); err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, t)
}

func (h *TenantHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "некорректный id")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		SendMappedError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
