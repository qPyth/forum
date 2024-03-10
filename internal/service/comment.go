package service

import (
	"forum/internal/models"
	s "forum/internal/storage"
)

type CommentService struct {
	commentStorage s.Comments
	postStorage    s.Posts
}

func NewCommentService(commentStorage s.Comments, postStorage s.Posts) *CommentService {
	return &CommentService{commentStorage: commentStorage, postStorage: postStorage}
}

func (c CommentService) Create(comment models.Comment) error {
	postExist := c.postStorage.PostExist(comment.PostID)
	if !postExist {
		return models.ErrPostNotExists
	}
	if len([]rune(comment.Content)) > 300 || len([]rune(comment.Content)) < 1 {
		return models.ErrCommentLength
	}
	return c.commentStorage.Create(comment)
}
