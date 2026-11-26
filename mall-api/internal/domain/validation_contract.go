package domain

import "time"

func (c *LeaseContract) ValidateForCreate() error {
	v := NewValidator()
	c.PremiseID = RequireID(v, "premise_id", c.PremiseID)
	c.TenantID = RequireID(v, "tenant_id", c.TenantID)
	start := RequireDate(v, "start_date", c.StartDate)
	end := RequireDate(v, "end_date", c.EndDate)
	if !start.IsZero() && !end.IsZero() && !end.After(start) {
		v.Add("end_date", "дата окончания должна быть позже даты начала")
	}
	c.RatePerSqm = RequirePositiveFloat(v, "rate_per_sqm", c.RatePerSqm)
	c.DepositAmount = RequireNonNegativeFloat(v, "deposit_amount", c.DepositAmount)
	return v.Err()
}

func (c *LeaseContract) ValidateForUpdate() error {
	v := NewValidator()
	if c.StartDate != "" {
		RequireDate(v, "start_date", c.StartDate)
	}
	if c.EndDate != "" {
		RequireDate(v, "end_date", c.EndDate)
	}
	if c.StartDate != "" && c.EndDate != "" {
		s, _ := time.Parse("2006-01-02", c.StartDate)
		e, _ := time.Parse("2006-01-02", c.EndDate)
		if !s.IsZero() && !e.IsZero() && !e.After(s) {
			v.Add("end_date", "дата окончания должна быть позже даты начала")
		}
	}
	if c.RatePerSqm != 0 {
		RequirePositiveFloat(v, "rate_per_sqm", c.RatePerSqm)
	}
	if c.DepositAmount < 0 {
		v.Add("deposit_amount", "не может быть отрицательной")
	}
	if c.Status != "" {
		RequireEnum(v, "status", c.Status, "active", "terminated", "expired")
	}
	return v.Err()
}
