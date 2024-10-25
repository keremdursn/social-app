package models


type Comment struct {
	CommentID int    `json:"commentid" gorm:"primaryKey;autoIncrement"`
	PostID    int    `json:"postid"`
	Post      Post   `gorm:"foreignKey:PostID"`
	UserID    int    `json:"userid"`
	User      User   `gorm:"foreignKey:UserID"`
	CommentText   string `json:"commenttext"`
	LikeCount int    `gorm:"default:0"`
}

type LikeCommand struct {
	// LikeCommandID int     `json:"like_id" gorm:"primaryKey;autoIncrement"`
	CommentID     int     `json:"comment_id" gorm:"index:idx_user_post,unique"`
	Comment       Comment `gorm:"foreignKey:CommentID"`
	UserID        int     `json:"user_id" gorm:"index:idx_user_post,unique"`
	User          User    `gorm:"foreignKey:UserID"`
}

type AnswerCommand struct {
	// AnswerCommandID int     `json:"answer_id" gorm:"primaryKey;autoIncrement"`
	CommentID       int     `gorm:"index" json:"comment_id"` 
	Comment         Comment `gorm:"foreignKey:CommentID"`
	UserID          int     `json:"user_id"`
	User            User    `gorm:"foreignKey:UserID"`
	Answer          string  `json:"answer"`
}
