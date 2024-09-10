package models

import (
	"gorm.io/gorm"
)

type Like struct {
	gorm.Model
	LikeID uint `json:"like_id" gorm:"primaryKey;autoIncrement"`
	UserID uint `json:"userid" gorm:"index:,unique"`
	User   User `gorm:"foreignKey:UserID"`
	PostID int  `json:"post_id" gorm:"index:,unique"`
	Post   Post `gorm:"foreignKey:PostID"`
}
