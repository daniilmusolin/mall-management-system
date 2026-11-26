package handler

import (
	"encoding/json"
	"net/http"

	"mall-api/internal/domain"
	"mall-api/internal/service"
)

type PaymentHandler struct {
	svc *service.PaymentService
}

func NewPaymentHandler(s *service.PaymentService) *PaymentHandler {
	return &PaymentHandler{svc: s}
}

func (h *PaymentHandler) RegisterRoutes(mux *http.ServeMux, auth func(http.Handler) http.Handler) {
	mux.Handle("GET /api/payments", auth(RequireRole(domain.RoleAdmin, domain.RoleAccountant)(http.HandlerFunc(h.GetAll))))
	mux.Handle("POST /api/payments", auth(RequireRole(domain.RoleAdmin, domain.RoleAccountant)(http.HandlerFunc(h.Create))))
	mux.Handle("PUT /api/payments/{id}", auth(RequireRole(domain.RoleAdmin, domain.RoleAccountant)(http.HandlerFunc(h.Update))))
	mux.Handle("DELETE /api/payments/{id}", auth(RequireRole(domain.RoleAdmin)(http.HandlerFunc(h.Delete))))
}

func (h *PaymentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.GetAll(r.Context(), ParsePagination(r))
	if err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, items)
}

func (h *PaymentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p domain.Payment
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

func (h *PaymentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "некорректный id")
		return
	}
	var p domain.Payment
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

func (h *PaymentHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
