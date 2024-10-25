package models



type Like struct {

	UserID int  `json:"user_id" gorm:"index:idx_user_post,unique"`
	User   User `gorm:"foreignKey:UserID"`
	PostID int  `json:"post_id" gorm:"index:idx_user_post,unique"`
	Post   Post `gorm:"foreignKey:PostID"`
}
