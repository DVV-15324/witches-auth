package repository

import (
	"fmt"
)

func (r *AuthRepository) cacheKeyAuthByEmail(email string) string {
	return fmt.Sprintf("auth:email:%s", email)
}

func (r *AuthRepository) cacheKeyAuthByUserID(userID int64) string {
	return fmt.Sprintf("auth:user:%d", userID)
}
