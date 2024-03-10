package models

type Tag struct {
	ID         int
	CategoryID int
	Name       string
}

type PostTags struct {
	PostID int
	TagID  int
}
