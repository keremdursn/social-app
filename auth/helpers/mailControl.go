package helpers

import (
	"auth/database"
	"auth/models"
	"log"
)

// func IsValidEmail(email string) bool {
// 	// Basit bir e-posta regex deseni
// 	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
// 	re := regexp.MustCompile(emailRegex)
// 	return re.MatchString(email)
// }

func MailControl(mail string) error {
	// if !IsValidEmail(mail) {
	// 	return errors.New("invalid email format")
	// }

	db := database.DB.Db
	var user models.User

	result := db.Where("mail = ?", mail).First(&user).Error
	if result.Error != nil {
		log.Printf("record not found: %v", result.Error)
	}
	return nil
}
