package models

type Session struct {
	UserID   int `gorm:"primaryKey;autoIncrement"`
	Token    string
	Expiry   int64
	IP       string
	IsActive bool `gorm:"default:true"`
}
