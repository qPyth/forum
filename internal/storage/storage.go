package storage

import (
	"database/sql"
	"forum/internal/models"
)

type Users interface {
	Create(user models.User) error
	GetByCredentials(email string) (models.User, error)
	SetSession(session models.Session) error
	UpdateSession(session models.Session) error
	DeleteSession(token string) error
	GetByID(id int) (models.User, error)
	GetByToken(token string) (*models.User, error)
	GetAuthorsOfPost(category int) (map[int]string, error)
	CheckSession(userID int) (string, bool, error)
}

type Posts interface {
	Create(post models.Post) (int, error)
	GetAllCategories() ([]models.Category, error)
	GetMapCategories() (map[int]string, error)
	GetPostByID(id int) (*models.Post, error)
	GetPostsByUserID(userID int) ([]models.Post, error)
	GetPostsByCategory(category string) ([]models.Post, error)
	PostsCountByCategory(id int) (int, error)
	GetLastPostByCategory(id int) (models.Post, error)
	GetCategoryByUrl(url string) (*models.Category, error)
	PostExist(postID int) bool
	GetVotesCount(postID int) (int, int, error)
	UpdateVotes(postID int, action string, delta int) error
	GetVotesPosts(userID int) ([]models.Post, error)
	TagExists(tagID int) (bool, error)
	GetTags(categoryID int) ([]models.Tag, error)
	GetByTags(tagIDs []int) ([]models.Post, error)
	GetTagsByPostID(postID int) ([]string, error)
	GetTagsByCategoryID(categoryID int) (map[int][]string, error)
}

type Comments interface {
	Create(comment models.Comment) error
	GetCountOnPost(categoryID int) (map[int]int, error)
	GetByPostID(postID int) ([]models.Comment, error)
	GetAuthors(postID int) (map[int]string, error)
	UpdateVotes(postID int, action string, delta int) error
	GetByID(ID int) (*models.Comment, error)
	CommentExist(commentID int) (bool, error)
}

type Votes interface {
	Create(vote models.Vote) error
	Delete(voteID int) error
	Update(voteID int, action string) error
	IsExists(vote models.Vote) (*models.Vote, error)
}

const tokenTable = "sessions"

type Storages struct {
	Users    Users
	Posts    Posts
	Comments Comments
	Votes    Votes
}

func NewStorages(db *sql.DB) *Storages {
	return &Storages{
		Users:    NewUsersStorage(db),
		Posts:    NewPostsStorage(db),
		Comments: NewCommentsStorage(db),
		Votes:    NewVoteStorage(db),
	}
}
