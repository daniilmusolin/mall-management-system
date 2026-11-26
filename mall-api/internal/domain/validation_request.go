package domain

func (s *ServiceRequest) ValidateForCreate() error {
	v := NewValidator()
	s.PremiseID = RequireID(v, "premise_id", s.PremiseID)
	s.RequestType = RequireString(v, "request_type", s.RequestType, 2, 100)
	s.Description = RequireString(v, "description", s.Description, 5, 2000)
	s.Priority = RequireEnum(v, "priority", s.Priority, "low", "medium", "high", "critical")
	return v.Err()
}

func (s *ServiceRequest) ValidateForUpdate() error {
	v := NewValidator()
	if s.RequestType != "" {
		RequireString(v, "request_type", s.RequestType, 2, 100)
	}
	if s.Description != "" {
		RequireString(v, "description", s.Description, 5, 2000)
	}
	if s.Priority != "" {
		RequireEnum(v, "priority", s.Priority, "low", "medium", "high", "critical")
	}
	if s.Status != "" {
		RequireEnum(v, "status", s.Status, "open", "in_progress", "closed", "cancelled")
	}
	return v.Err()
}
