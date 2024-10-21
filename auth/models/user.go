package models

type User struct {
	ID       int   `gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Password string `json:"password"`
	Mail     string `json:"mail"`
	PPId     string `json:"image"`
	PPUrl    string `json:"ppurl"`
}

type Login struct {
	Mail string `json:"mail"`
	Password string `json:"password"`
}

type ChangePassword struct {
	OldPassword  string `json:"oldpassword"`
	NewPassword1 string `json:"newpassword1"`
	NewPassword2 string `json:"newpassword2"`
}
