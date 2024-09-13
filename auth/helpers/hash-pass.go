package helpers

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPass(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

// CheckPass compares a bcrypt hashed password with its plain-text version
func CheckPass(hashedPassword, password string) error {
	// Verilen şifre ile hashlenmiş şifreyi karşılaştır
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err // Şifre uyuşmuyorsa hata döner
	}
	return nil
}
