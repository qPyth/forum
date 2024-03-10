package service

import (
	"forum/internal/models"
	s "forum/internal/storage"
	"mime/multipart"
	"time"
)

type UsersService interface {
	SignUp(input UserSignUpInput) error
	SignIn(input UserSignInInput) (string, time.Time, error)
	Logout(token string) error
	GetByToken(token string) (*models.User, error)
	CheckSession(userID int) (string, bool, error)
}

type UserSignUpInput struct {
	Username string
	Email    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type PostsService interface {
	Create(input PostCreateInput) (int, error)
	GetCategoriesWithInfo() (Categories, error)
	GetPostsByCategory(category string, tagsIDs []string) (*PostsByCategory, error)
	GetByID(postID int) (*PostViewData, error)
	GetCreatedPosts(userID int) ([]models.Post, map[int]string, error)
	GetVotedPosts(userID int) ([]models.Post, map[int]string, error)
	GetTagsByCategory(categoryUrl string) ([]models.Tag, error)
}

type PostViewData struct {
	Post     *models.Post
	Author   string
	Comments []models.Comment
	Authors  map[int]string
	Tags     []string
}

type Image struct {
	File       multipart.File
	FileHeader *multipart.FileHeader
}

type PostCreateInput struct {
	Title    string
	Content  string
	Category string
	Tags     []string
	Image    Image
	UserId   int
	HasImage bool
}

type Categories struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	PostsCount  int    `json:"posts_count"`
	LatestPost  Post   `json:"latest_post"`
}

type Post struct {
	Author string `json:"author"`
	models.Post
}

type PostsByCategory struct {
	Category      *models.Category `json:"category"`
	CommentsCount map[int]int      `json:"comments_count"`
	Authors       map[int]string   `json:"authors"`
	Posts         []models.Post    `json:"posts"`
	Tags          map[int][]string `json:"tags"`
}

type CommentsService interface {
	Create(comment models.Comment) error
}

type VotesService interface {
	MakeVote(input VoteInput) (*VoteOutput, error)
}

type VoteInput struct {
	ID        int    `json:"id"`
	Action    string `json:"action"`
	IsPost    bool   `json:"is_post"`
	IsComment bool   `json:"is_comment"`
	UserID    int    `json:"user_id"`
}

type VoteOutput struct {
	Action       string `json:"action"`
	LikeCount    int    `json:"like_count"`
	DislikeCount int    `json:"dislike_count"`
}

type Services struct {
	Users    UsersService
	Posts    PostsService
	Comments CommentsService
	Votes    VotesService
}

func NewServices(db *s.Storages) *Services {
	return &Services{
		Users:    NewUserService(db.Users),
		Posts:    NewPostService(db),
		Comments: NewCommentService(db.Comments, db.Posts),
		Votes:    NewVoteService(db.Posts, db.Comments, db.Votes),
	}
}
