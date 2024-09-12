package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	CommentID int    `json:"commentid" gorm:"primaryKey"`
	PostID    int    `json:"postid"`
	Post      Post   `gorm:"foreignKey:PostID"`
	UserID    int    `json:"userid"`
	User      User   `gorm:"foreignKey:UserID"`
	Comment   string `json:"comment"`
	LikeCount int    `gorm:"default:0"`
}

type LikeCommand struct {
	gorm.Model
	LikeCommandID int     `json:"like_id" gorm:"primaryKey;autoIncrement"`
	CommentID     int     `json:"comment_id" gorm:"index:idx_user_post,unique"`
	Comment       Comment `gorm:"foreignKey:CommentID"`
	UserID        int     `json:"user_id" gorm:"index:idx_user_post,unique"`
	User          User    `gorm:"foreignKey:UserID"`
}

type AnswerCommand struct {
	gorm.Model
	AnswerCommandID int     `json:"answer_id" gorm:"primaryKey;autoIncrement"`
	CommentID       int     `gorm:"index" json:"comment_id"` 
	Comment         Comment `gorm:"foreignKey:CommentID"`
	UserID          int     `json:"user_id"`
	User            User    `gorm:"foreignKey:UserID"`
	Answer          string  `json:"answer"`
}
