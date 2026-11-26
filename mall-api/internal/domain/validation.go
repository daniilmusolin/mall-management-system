package domain

import (
	"net/mail"
	"regexp"
	"strings"
	"time"
)

var (
	innRe     = regexp.MustCompile(`^\d{10}$|^\d{12}$`)
	phoneRe   = regexp.MustCompile(`^\+?[0-9]{7,15}$`)
	roomNumRe = regexp.MustCompile(`^[A-Za-zА-Яа-я0-9\-/]{1,20}$`)
)

type Validator struct {
	Errors []*ValidationError
}

func NewValidator() *Validator { return &Validator{Errors: make([]*ValidationError, 0)} }

func (v *Validator) Add(field, msg string) {
	v.Errors = append(v.Errors, NewValidationError(field, msg))
}

func (v *Validator) HasErrors() bool { return len(v.Errors) > 0 }

func (v *Validator) Err() error {
	if !v.HasErrors() {
		return nil
	}
	if len(v.Errors) == 1 {
		return v.Errors[0]
	}
	return &ValidationErrors{List: v.Errors}
}

func RequireString(v *Validator, field, value string, minLen, maxLen int) string {
	value = strings.TrimSpace(value)
	if value == "" {
		if minLen > 0 {
			v.Add(field, "обязательное поле")
		}
		return value
	}
	l := len([]rune(value))
	if l < minLen {
		v.Add(field, "слишком короткое значение")
	}
	if maxLen > 0 && l > maxLen {
		v.Add(field, "слишком длинное значение")
	}
	return value
}

func RequireEmail(v *Validator, field, value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		v.Add(field, "обязательное поле")
		return value
	}
	if _, err := mail.ParseAddress(value); err != nil {
		v.Add(field, "некорректный email")
	}
	if len(value) > 254 {
		v.Add(field, "слишком длинный email")
	}
	return value
}

func RequirePassword(v *Validator, field, value string) string {
	if len(value) < 8 {
		v.Add(field, "минимум 8 символов")
	}
	if len(value) > 72 {
		v.Add(field, "максимум 72 символа")
	}
	return value
}

func RequirePositiveFloat(v *Validator, field string, value float64) float64 {
	if value <= 0 {
		v.Add(field, "должно быть больше нуля")
	}
	return value
}

func RequireNonNegativeFloat(v *Validator, field string, value float64) float64 {
	if value < 0 {
		v.Add(field, "не может быть отрицательным")
	}
	return value
}

func RequireID(v *Validator, field string, value int64) int64 {
	if value <= 0 {
		v.Add(field, "обязательное поле")
	}
	return value
}

func RequireEnum(v *Validator, field, value string, allowed ...string) string {
	if value == "" {
		v.Add(field, "обязательное поле")
		return value
	}
	for _, a := range allowed {
		if value == a {
			return value
		}
	}
	v.Add(field, "недопустимое значение: "+value)
	return value
}

func RequireDate(v *Validator, field, value string) time.Time {
	if value == "" {
		v.Add(field, "обязательное поле")
		return time.Time{}
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		v.Add(field, "формат даты должен быть YYYY-MM-DD")
		return time.Time{}
	}
	return t
}

func RequireINN(v *Validator, field, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		v.Add(field, "обязательное поле")
		return value
	}
	if !innRe.MatchString(value) {
		v.Add(field, "ИНН должен содержать 10 или 12 цифр")
	}
	return value
}

func RequirePhone(v *Validator, field, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	if !phoneRe.MatchString(value) {
		v.Add(field, "некорректный номер телефона")
	}
	return value
}

func RequireRoomNumber(v *Validator, field, value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		v.Add(field, "обязательное поле")
		return value
	}
	if !roomNumRe.MatchString(value) {
		v.Add(field, "допустимы буквы, цифры, дефис и слэш (1–20 символов)")
	}
	return value
}
