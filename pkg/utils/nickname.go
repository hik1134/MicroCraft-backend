package utils

import "strings"

func NicknameFromEmail(email string) string {
	email = strings.TrimSpace(email)
	at := strings.Index(email, "@")
	if at <= 0 {
		return "user"
	}
	return email[:at]
}