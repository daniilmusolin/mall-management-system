package domain

func (p *Payment) ValidateForCreate() error {
	v := NewValidator()
	p.ContractID = RequireID(v, "contract_id", p.ContractID)
	p.Amount = RequirePositiveFloat(v, "amount", p.Amount)
	p.PaymentDate = RequireString(v, "payment_date", p.PaymentDate, 10, 10)
	if p.PaymentDate != "" {
		RequireDate(v, "payment_date", p.PaymentDate)
	}
	p.PaymentType = RequireEnum(v, "payment_type", p.PaymentType, "rent", "deposit", "utilities", "penalty")
	if len(p.ReferenceNumber) > 100 {
		v.Add("reference_number", "максимум 100 символов")
	}
	return v.Err()
}

func (p *Payment) ValidateForUpdate() error {
	v := NewValidator()
	if p.Amount != 0 {
		RequirePositiveFloat(v, "amount", p.Amount)
	}
	if p.PaymentDate != "" {
		RequireDate(v, "payment_date", p.PaymentDate)
	}
	if p.PaymentType != "" {
		RequireEnum(v, "payment_type", p.PaymentType, "rent", "deposit", "utilities", "penalty")
	}
	if p.Status != "" {
		RequireEnum(v, "status", p.Status, "pending", "completed", "cancelled")
	}
	if len(p.ReferenceNumber) > 100 {
		v.Add("reference_number", "максимум 100 символов")
	}
	return v.Err()
}
