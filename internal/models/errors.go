package models

import "errors"

var (
	ErrUsernameExist        = errors.New("user with this username already exists")
	ErrEmailExist           = errors.New("user with this email already exists")
	ErrUserNotFound         = errors.New("user with these credentials not found")
	ErrMissMatchEmailOrPass = errors.New("incorrect email or password")
	ErrPostsByCatNotFound   = errors.New("posts for this category not found")
	ErrNoPost               = errors.New("will be first")
	ErrWrongCategory        = errors.New("choose right category")
	ErrImageSize            = errors.New("image size over 20mb")
	ErrImageExtension       = errors.New("image extension is wrong, must jpeg, svg, gif or png")
	ErrCommentsNotFound     = errors.New("comments for this post not found")
	ErrDirExists            = errors.New("mkdir ui/static/images/post_images/: file exists")
	ErrPostNotExists        = errors.New("post doesn't exist")
	ErrCommentLength        = errors.New("the comment should not be more than 300 characters")
	ErrWrongAction          = errors.New("action must be like or dislike")
	ErrWrongVoteItem        = errors.New("vote item must be post or comment")
	ErrPostTitleLength      = errors.New("the post title should not be more than 50 characters")
	ErrPostContentLength    = errors.New("the post content should not be more than 500 characters")
	ErrNotValidTagID        = errors.New("not valid tag id")
)
