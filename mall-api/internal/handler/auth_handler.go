package handler

import (
	"encoding/json"
	"net/http"

	"mall-api/internal/domain"
	"mall-api/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: s}
}

func (h *AuthHandler) RegisterRoutes(mux *http.ServeMux, auth func(http.Handler) http.Handler) {
	mux.HandleFunc("POST /api/auth/register", h.Register)
	mux.HandleFunc("POST /api/auth/login", h.Login)
	mux.HandleFunc("POST /api/auth/refresh", h.Refresh)

	mux.Handle("GET /api/users", auth(RequireRole(domain.RoleAdmin)(http.HandlerFunc(h.GetAllUsers))))
	mux.Handle("PUT /api/users/{id}", auth(RequireRole(domain.RoleAdmin)(http.HandlerFunc(h.UpdateUser))))
	mux.Handle("DELETE /api/users/{id}", auth(RequireRole(domain.RoleAdmin)(http.HandlerFunc(h.DeleteUser))))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var in domain.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	resp, err := h.svc.Register(r.Context(), in)
	if err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var in domain.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	resp, err := h.svc.Login(r.Context(), in)
	if err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var in domain.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}
	resp, err := h.svc.RefreshToken(r.Context(), in.RefreshToken)
	if err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.GetAllUsers(r.Context(), ParsePagination(r))
	if err != nil {
		SendMappedError(w, err)
		return
	}
	SendJSON(w, http.StatusOK, users)
}

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "некорректный id")
		return
	}

	var in domain.UpdateUserInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		SendError(w, http.StatusBadRequest, "некорректный JSON")
		return
	}

	if err := h.svc.UpdateUser(r.Context(), id, in); err != nil {
		SendMappedError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := ParseID(r)
	if err != nil {
		SendError(w, http.StatusBadRequest, "некорректный id")
		return
	}

	if err := h.svc.DeleteUser(r.Context(), id); err != nil {
		SendMappedError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
