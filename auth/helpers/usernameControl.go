package helpers

import (
	"auth/database"
	"auth/models"
)

func UsernameControl(username string) error {
	db := database.DB.Db
	var user models.User
	
	err := db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return err
	}
	return nil
}
