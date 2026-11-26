package domain

func (in *RegisterInput) Validate() error {
	v := NewValidator()
	in.Email = RequireEmail(v, "email", in.Email)
	in.Password = RequirePassword(v, "password", in.Password)
	in.FullName = RequireString(v, "full_name", in.FullName, 2, 255)
	return v.Err()
}

func (in *LoginInput) Validate() error {
	v := NewValidator()
	in.Email = RequireEmail(v, "email", in.Email)
	if in.Password == "" {
		v.Add("password", "обязательное поле")
	}
	return v.Err()
}
