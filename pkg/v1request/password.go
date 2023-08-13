package v1request

type PasswordchangeRequest struct {
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	NewPAsswordconfirm string `json:"new_password_confirm"`
}

func (p *PasswordchangeRequest) Validate() error {
	if p.OldPassword == "" {
		return ErrEmptyOldPassword
	}
	if p.NewPassword == "" {
		return ErrEmptyNewPassword
	}
	if p.NewPAsswordconfirm != p.NewPassword {
		return ErrPasswordNotMatch
	}
	return nil
}
