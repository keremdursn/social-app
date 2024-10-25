package models


type Post struct {
	PostID       int      `json:"postid" gorm:"primaryKey;autoIncrement"`
	User         User     `gorm:"foreignKey:UserID"`
	UserID       int      `json:"userid"`
	PostDesc     string   `json:"postdesc"`
	ImageIDs     []string `json:"image_ids" gorm:"type:jsonb"`
	ImageURLs    []string `json:"image_urls" gorm:"type:jsonb"`
	LikeCount    int      `json:"likecount" gorm:"default:0"`
	CommentCount int      `json:"commentcount" gorm:"default:0"`
	IsActive     bool     `gorm:"default:true"`
}
