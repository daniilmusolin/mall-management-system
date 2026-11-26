package handler

import (
	"encoding/json"
	"net/http"

	"mall-api/internal/domain"
	"mall-api/internal/service"
)

type ContractHandler struct {
	svc *service.ContractService
}

func NewContractHandler(s *service.ContractService) *ContractHandler {
	return &ContractHandler{svc: s}
}

func (h *ContractHandler) RegisterRoutes(mux *http.ServeMux, auth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/contracts", auth(RequireRole(domain.RoleAdmin, domain.RoleManager, domain.RoleAccountant)(http.HandlerFunc(h.GetAll))))
	mux.Handle("POST /api/contracts", auth(RequireRole(domain.RoleAdmin, domain.RoleManager)(http.HandlerFunc(h.Create))))
	mux.Handle("PUT /api/contracts/{id}", auth(RequireRole(domain.RoleAdmin, domain.RoleManager)(http.HandlerFunc(h.Update))))
	mux.Handle("DELETE /api/contracts/{id}", auth(RequireRole(domain.RoleAdmin)(http.HandlerFunc(h.Delete))))
}

func (h *ContractHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.GetAll(r.Context(), ParsePagination(r))
	if err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, items)
}

func (h *ContractHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c domain.LeaseContract
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	if err := h.svc.Create(r.Context(), &c); err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusCreated, c)
}

func (h *ContractHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "некорректный id")
		return
	}
	var c domain.LeaseContract
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	c.ID = id
	if err := h.svc.Update(r.Context(), &c); err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, c)
}

func (h *ContractHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
