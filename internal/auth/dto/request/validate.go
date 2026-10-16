package request

import (
	"regexp"
)

func CheckEmail(email string) error {
	if len(email) < 1 {
		return ErrorEmailIsNotEmpty
	}
	_, err := regexp.MatchString(`^[0-9a-zA-Z.\-+]+@[0-9a-zA-Z]+\.[0-9a-zA-Z]{2,}$`, email)
	if err != nil {
		return ErrorEmailNotValid
	}
	return nil
}

func CheckPasword(password string) error {
	if len(password) < 2 {
		return ErrorPasswordIsNotEmpty
	}
	if len(password) < 2 || len(password) > 20 {
		return ErrorPasswordNotValid
	}
	return nil
}

func CheckName(name string) error {
	if len(name) < 2 {
		return ErrorNameIsNotEmpty
	}
	if len(name) < 2 || len(name) > 20 {
		return ErrorNameNotValid
	}
	return nil
}
