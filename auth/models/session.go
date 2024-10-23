package models

type Session struct {
	UserID   int `gorm:"primaryKey;autoIncrement"`
	Token    string
}
