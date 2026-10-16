package request

import (
	"strings"
)

type Login struct {
	Email    string     `json:"email"`
	Password string     `json:"password"`
	JWK      JWKRequest `json:"jwk"`
}

type JWKRequest struct {
	Kty string `json:"kty" example:"EC"`
	Crv string `json:"crv" example:"P-256"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

func (r *Login) Validate() error {
	r.Email = strings.TrimSpace(r.Email)
	r.Password = strings.TrimSpace(r.Password)
	err_email := CheckEmail(r.Email)
	if err_email != nil {
		return err_email
	}
	err_password := CheckPasword(r.Password)
	if err_password != nil {
		return err_password
	}
	return nil
}
