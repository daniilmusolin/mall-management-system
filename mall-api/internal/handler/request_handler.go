package handler

import (
	"encoding/json"
	"net/http"

	"mall-api/internal/domain"
	"mall-api/internal/service"
)

type RequestHandler struct {
	svc *service.RequestService
}

func NewRequestHandler(s *service.RequestService) *RequestHandler {
	return &RequestHandler{svc: s}
}

func (h *RequestHandler) RegisterRoutes(mux *http.ServeMux, auth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/requests", auth(RequireRole(domain.RoleAdmin, domain.RoleTechnician, domain.RoleSupport)(http.HandlerFunc(h.GetAll))))
	mux.Handle("POST /api/requests", auth(RequireRole(domain.RoleAdmin, domain.RoleSupport)(http.HandlerFunc(h.Create))))
	mux.Handle("PUT /api/requests/{id}", auth(RequireRole(domain.RoleAdmin, domain.RoleTechnician)(http.HandlerFunc(h.Update))))
	mux.Handle("DELETE /api/requests/{id}", auth(RequireRole(domain.RoleAdmin)(http.HandlerFunc(h.Delete))))
}

func (h *RequestHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.GetAll(r.Context(), ParsePagination(r))
	if err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, items)
}

func (h *RequestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var s domain.ServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	if err := h.svc.Create(r.Context(), &s); err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusCreated, s)
}

func (h *RequestHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "некорректный id")
		return
	}
	var s domain.ServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	s.ID = id
	if err := h.svc.Update(r.Context(), &s); err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, s)
}

func (h *RequestHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
