package models

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Mail     string `json:"mail"`
	PhotoUrl string `json:"image"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ChangePassword struct {
	OldPassword  string `json:"oldpassword"`
	NewPassword1 string `json:"newpassword1"`
	NewPassword2 string `json:"newpassword2"`
}
