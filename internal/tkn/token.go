package tkn

// заглушки
func IsAdminToken(token string) bool {
	return token == "admin_token"
}

func IsUserToken(token string) bool {
	return token == "user_token"
}
