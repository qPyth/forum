package http

import (
	"forum/pkg/validation"
)

func signUpValidation(username, email, password, confirmPassword string) error {
	if !validation.IsValidUname(username) {
		return usernameFormatError
	}
	if !validation.IsValidEmail(email) {
		return emailFormatError
	}
	if len([]byte(password)) < 8 {
		return passLenError
	}
	if !validation.IsValidPassword(password) {
		return passFormatError
	}
	if password != confirmPassword {
		return passMatchError
	}
	return nil
}
