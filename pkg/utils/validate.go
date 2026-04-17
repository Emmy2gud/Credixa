package utils

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateRegister(fullname, email, password string) error {
	if strings.TrimSpace(fullname) == "" {
		return errors.New("fullname is required")
	}

	if !isValidEmail(email) {
		return errors.New("invalid email")
	}

	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}

	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func validatePassword(password string) (bool, string) {
	// Regex pattern explanation:
	// ^                 : start of string
	// (?=.*[a-z])       : at least one lowercase letter
	// (?=.*[A-Z])       : at least one uppercase letter
	// (?=.*[0-9])       : at least one digit
	// (?=.*[!@#$%^&*])  : at least one special character
	// .{8,}             : at least 8 characters long
	// $                 : end of string
	pattern := `^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#$%^&*]).{8,}$`

	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, "Invalid regex pattern"
	}

	if re.MatchString(password) {
		return true, "Password is valid"
	}
	return false, "Password must be at least 8 characters long, include uppercase, lowercase, number, and special character"
}
