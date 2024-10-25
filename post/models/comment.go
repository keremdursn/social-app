package models

type Comment struct {
	CommentID   int    `json:"commentid" gorm:"primaryKey;autoIncrement"`
	PostID      int    `json:"postid"`
	Post        Post   `json:"foreignKey:PostId"`
	UserID      int    `json:"userid"`
	User        User   `gorm:"foreignKey:UserID"`
	CommentText string `json:"commenttext"`
	LikeCount   int    `gorm:"default:0"`
}

type LikeComment struct {
	User      User    `gorm:"foreigenKey:UserID"`
	Comment   Comment `json:"foreignKey:CommentID"`
	CommentID int     `json:"commentid"`
	UserID    int     `json:"userid"`
}

type AnswerComment struct {
	CommentID int     `gorm:"index" json:"comment_id"`
	Comment   Comment `gorm:"foreignKey:CommentID"`
	UserID    int     `json:"user_id"`
	User      User    `gorm:"foreignKey:UserID"`
	Answer    string  `json:"answer"`
}
