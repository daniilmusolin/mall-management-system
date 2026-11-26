package domain

func (t *Tenant) ValidateForCreate() error {
	v := NewValidator()
	t.Name = RequireString(v, "name", t.Name, 2, 255)
	t.INN = RequireINN(v, "inn", t.INN)
	t.LegalAddress = RequireString(v, "legal_address", t.LegalAddress, 5, 500)
	t.ContactPhone = RequirePhone(v, "contact_phone", t.ContactPhone)
	t.Email = RequireEmail(v, "email", t.Email)
	return v.Err()
}

func (t *Tenant) ValidateForUpdate() error {
	return t.ValidateForCreate()
}
