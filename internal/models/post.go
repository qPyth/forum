package models

import (
	"time"
)

type Post struct {
	ID           int
	Title        string
	Content      string
	Image        string
	UserID       int
	CategoryID   int
	Tags         []int
	CreatedDate  time.Time
	LikeCount    int
	DislikeCount int
}
