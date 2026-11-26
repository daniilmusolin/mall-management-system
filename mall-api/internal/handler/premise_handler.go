package handler

import (
	"encoding/json"
	"net/http"

	"mall-api/internal/domain"
	"mall-api/internal/service"
)

type PremiseHandler struct {
	svc *service.PremiseService
}

func NewPremiseHandler(s *service.PremiseService) *PremiseHandler {
	return &PremiseHandler{svc: s}
}

func (h *PremiseHandler) RegisterRoutes(mux *http.ServeMux, auth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/premises", auth(RequireRole(domain.RoleAdmin, domain.RoleManager, domain.RoleAccountant, domain.RoleSupport, domain.RoleTechnician)(http.HandlerFunc(h.GetAll))))
	mux.Handle("POST /api/premises", auth(RequireRole(domain.RoleAdmin)(http.HandlerFunc(h.Create))))
	mux.Handle("PUT /api/premises/{id}", auth(RequireRole(domain.RoleAdmin)(http.HandlerFunc(h.Update))))
	mux.Handle("DELETE /api/premises/{id}", auth(RequireRole(domain.RoleAdmin)(http.HandlerFunc(h.Delete))))
}

func (h *PremiseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.GetAll(r.Context(), ParsePagination(r))
	if err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, items)
}

func (h *PremiseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p domain.Premise
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	if err := h.svc.Create(r.Context(), &p); err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusCreated, p)
}

func (h *PremiseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "некорректный id")
		return
	}
	var p domain.Premise
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	p.ID = id
	if err := h.svc.Update(r.Context(), &p); err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, p)
}

func (h *PremiseHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
