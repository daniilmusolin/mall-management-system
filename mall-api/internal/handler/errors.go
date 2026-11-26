package handler

import (
	"errors"
	"net/http"

	"mall-api/internal/domain"
)

type ErrorPayload struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

func MapError(err error) (int, ErrorPayload) {
	if err == nil {
		return http.StatusOK, ErrorPayload{}
	}

	var ve *domain.ValidationErrors
	if errors.As(err, &ve) {
		fields := make(map[string]string, len(ve.List))
		for _, e := range ve.List {
			fields[e.Field] = e.Message
		}
		return http.StatusUnprocessableEntity, ErrorPayload{
			Error:  "проверьте правильность заполнения полей",
			Fields: fields,
		}
	}

	var single *domain.ValidationError
	if errors.As(err, &single) {
		return http.StatusUnprocessableEntity, ErrorPayload{
			Error:  single.Message,
			Fields: map[string]string{single.Field: single.Message},
		}
	}

	switch {
	case errors.Is(err, domain.ErrNotFound):
		return http.StatusNotFound, ErrorPayload{Error: "запись не найдена"}
	case errors.Is(err, domain.ErrConflict):
		return http.StatusConflict, ErrorPayload{Error: "конфликт данных"}
	case errors.Is(err, domain.ErrUnauthorized):
		return http.StatusUnauthorized, ErrorPayload{Error: "требуется авторизация"}
	case errors.Is(err, domain.ErrForbidden):
		return http.StatusForbidden, ErrorPayload{Error: "недостаточно прав"}
	case errors.Is(err, domain.ErrValidation), errors.Is(err, domain.ErrBadRequest):
		return http.StatusBadRequest, ErrorPayload{Error: err.Error()}
	default:
		return http.StatusInternalServerError, ErrorPayload{Error: "внутренняя ошибка сервера"}
	}
}

func SendMappedError(w http.ResponseWriter, err error) {
	status, payload := MapError(err)
	SendJSON(w, status, payload)
}
