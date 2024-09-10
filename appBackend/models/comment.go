package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	CommentID    int   `json:"commentid" gorm:"primaryKey"`
	PostID       int    `json:"postid"`
	Post         Post   `gorm:"foreignKey:PostID"`
	UserID       int    `json:"userid"`
	User         User   `gorm:"foreigenKey:UserID"`
	Comment      string `json:"comment"`
	LikeCount    int    `gorm:"default:0"`
	CommentCount int    `gorm:"default:0"`
}

type LikeCommand struct {
	gorm.Model
	LikeCommandID int    `json:"like_id" gorm:"primaryKey;autoIncrement"`
	CommentID     int     `gorm:"primarykey" json:"comment_id"`
	Comment       Comment `gorm:"foreigenKey:CommentID"`
	UserID        int     `json:"user_id"`
	User          User    `gorm:"foreigenKey:UserID"`
}

type AnswerCommand struct {
	gorm.Model
	AnswerCommandID int    `json:"like_id" gorm:"primaryKey;autoIncrement"`
	CommentID       int     `gorm:"primarykey" json:"comment_id"`
	Comment         Comment `gorm:"foreigenKey:CommentID"`
	UserID          int     `json:"user_id"`
	User            User    `gorm:"foreigenKey:UserID"`
	Answer          string  `json:"answer"`
}
