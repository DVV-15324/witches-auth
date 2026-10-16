package repository

func (r *RefreshTokenRepository) cacheKeyToken(token string) string {
	return "refresh:token:" + token
}
func (r *RefreshTokenRepository) cacheKeyJKT(jkt string) string {
	return "refresh:jkt:" + jkt
}
