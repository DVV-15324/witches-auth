package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

func RandomStr(lenght int64) (string, error) {
	b := make([]byte, lenght)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (h *Hash) GenerateFromPassword(password string, salt string) (string, error) {
	spStr := fmt.Sprintf("%s.%s", salt, password)
	hashP, err := bcrypt.GenerateFromPassword([]byte(spStr), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashP), nil
}

func (h *Hash) CompareHashAndPassword(passwordStr string, password string, salt string) bool {
	spStr := fmt.Sprintf("%s.%s", salt, password)
	b := bcrypt.CompareHashAndPassword([]byte(passwordStr), []byte(spStr))
	return b == nil
}
