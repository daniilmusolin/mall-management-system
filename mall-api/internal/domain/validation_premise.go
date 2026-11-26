package domain

func (p *Premise) ValidateForCreate() error {
	v := NewValidator()
	p.RoomNumber = RequireRoomNumber(v, "room_number", p.RoomNumber)
	if p.Floor < 0 {
		v.Add("floor", "этаж не может быть отрицательным")
	}
	p.Area = RequirePositiveFloat(v, "area", p.Area)
	p.Type = RequireString(v, "type", p.Type, 1, 50)
	p.Zone = RequireString(v, "zone", p.Zone, 0, 20)
	return v.Err()
}

func (p *Premise) ValidateForUpdate() error {
	v := NewValidator()
	p.RoomNumber = RequireRoomNumber(v, "room_number", p.RoomNumber)
	if p.Floor < 0 {
		v.Add("floor", "этаж не может быть отрицательным")
	}
	p.Area = RequirePositiveFloat(v, "area", p.Area)
	p.Type = RequireString(v, "type", p.Type, 1, 50)
	p.Zone = RequireString(v, "zone", p.Zone, 0, 20)
	p.Status = RequireEnum(v, "status", p.Status, "free", "rented", "maintenance")
	return v.Err()
}
