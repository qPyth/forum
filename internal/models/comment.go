package models

import "time"

type Comment struct {
	ID           int
	Content      string
	UserID       int
	PostID       int
	CreatedDate  time.Time
	LikeCount    int
	DislikeCount int
}
