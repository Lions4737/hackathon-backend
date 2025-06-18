package model

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	FirebaseUID  string    `gorm:"unique;type:varchar(128)" json:"firebase_uid"`
	Username     string    `gorm:"type:varchar(50)" json:"username"`
	Description  string    `gorm:"type:varchar(140)" json:"description"`
	ProfileImage string    `gorm:"type:varchar(255)" json:"profile_image"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Post struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `json:"user_id"`
	Content      string    `gorm:"type:text" json:"content"`
	IsReply      bool      `gorm:"default:false" json:"is_reply"`
	ParentPostID *uint     `json:"parent_post_id"`
	LikeCount    int       `gorm:"default:0" json:"like_count"`
	ReplyCount   int       `gorm:"default:0" json:"reply_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
}

type Like struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint      `json:"user_id"`
	PostID    uint      `json:"post_id"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Post      Post      `gorm:"foreignKey:PostID" json:"post"`
}
