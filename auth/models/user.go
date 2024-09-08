package models

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Mail     string `json:"mail"`
}
