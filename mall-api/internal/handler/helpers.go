package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"mall-api/internal/domain"
)

func SendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func SendError(w http.ResponseWriter, status int, msg string) {
	SendJSON(w, status, map[string]string{"error": msg})
}

func ParsePagination(r *http.Request) domain.Pagination {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	return domain.Pagination{Limit: limit, Offset: offset}
}

func ParseID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, domain.ErrBadRequest
	}
	return id, nil
}
