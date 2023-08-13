package v1request

import "errors"

var (
	ErrEmptyOldPassword = errors.New("old password is empty")
	ErrEmptyNewPassword = errors.New("new password is empty")
	ErrPasswordNotMatch = errors.New("new password and new password confirm not match")
)
