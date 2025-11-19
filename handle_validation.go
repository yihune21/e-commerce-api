package main

import (
	"errors"
	"unicode"
)

func PasswordChecker(password string) error {
	var (
		hasUpperCase bool
		hasLowerCase bool
		hasDigit     bool
		hasSpecial   bool
	)

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpperCase = true
		} else if unicode.IsLower(char) {
			hasLowerCase = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else {
			hasSpecial = true
		}
	}
    if len(password) < 8 {
		return errors.New("Password is too short. It should be at least 8 characters long")
	} else if !hasUpperCase {
		return errors.New("Password should contain at least one uppercase letter")
	} else if !hasLowerCase {
		return errors.New("Password should contain at least one lowercase letter")
	} else if !hasDigit {
		return errors.New("Password should contain at least one digit")
	} else if !hasSpecial {
		return errors.New("Password should contain at least one special character")
	}
	
	return nil
}

