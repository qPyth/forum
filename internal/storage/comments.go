package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

type CommentsStorage struct {
	db *sql.DB
}

func (s *CommentsStorage) Create(comment models.Comment) error {
	op := "CommentStorage.Create"
	stmt, err := s.db.Prepare("INSERT INTO comments (content, user_id, post_id) values (?,?,?)")
	if err != nil {
		return fmt.Errorf("%s-query prepare error: %w", op, err)
	}
	_, err = stmt.Exec(comment.Content, comment.UserID, comment.PostID)
	if err != nil {
		return fmt.Errorf("%s-query exec error: %w", op, err)
	}
	return nil
}

func NewCommentsStorage(db *sql.DB) *CommentsStorage {
	return &CommentsStorage{db: db}
}

func (s *CommentsStorage) GetCountOnPost(categoryID int) (map[int]int, error) {
	rows, err := s.db.Query(`
	SELECT posts.id, COUNT(comments.id) as comments_count
	FROM posts
	LEFT JOIN comments ON comments.post_id = posts.id
	WHERE category_id=?
	GROUP BY posts.id
`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("CommentStorage.GetCountOnPosts - query error: %w", err)
	}
	defer rows.Close()

	commentsCountMap := make(map[int]int)
	for rows.Next() {
		var id, commentsCount int
		err = rows.Scan(&id, &commentsCount)
		if err != nil {
			return nil, fmt.Errorf("CommentStorage.GetCountOnPosts - scan error: %w", err)
		}
		commentsCountMap[id] = commentsCount
	}
	return commentsCountMap, nil
}

func (s *CommentsStorage) GetByPostID(postID int) ([]models.Comment, error) {
	op := "CommentStorage.GetByPostID"
	var comments []models.Comment

	rows, err := s.db.Query(`SELECT * FROM comments WHERE post_id=?`, postID)
	if err != nil {
		return nil, fmt.Errorf("%s - query error: %w", op, err)
		
	}
	defer rows.Close()

	for rows.Next() {
		var comment models.Comment

		err = rows.Scan(&comment.ID, &comment.Content, &comment.UserID, &comment.PostID, &comment.CreatedDate, &comment.LikeCount, &comment.DislikeCount)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, models.ErrCommentsNotFound
			} else {
				return nil, fmt.Errorf("%s - scan error:%w", op, err)
			}
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (s *CommentsStorage) GetAuthors(postID int) (map[int]string, error) {
	rows, err := s.db.Query(`
			SELECT comments.id, users.username
			FROM comments
			INNER JOIN users ON users.id = comments.user_id
			WHERE comments.post_id = ?
		`, postID)
	if err != nil {
		return nil, fmt.Errorf("UsersStorage.GetAuthorsOfPost - query error: %w", err)
	}
	defer rows.Close()

	authorsMap := make(map[int]string)
	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("UsersStorage.GetAuthorsOfPost - scan error: %w", err)
		}
		authorsMap[id] = name
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("UsersStorage.GetAuthorsOfPost - rows error: %w", err)
	}
	return authorsMap, nil
}

func (s *CommentsStorage) UpdateVotes(postID int, action string, delta int) error {
	query := fmt.Sprintf("UPDATE comments SET %s = %s + ? WHERE id = ?", action, action)
	_, err := s.db.Exec(query, delta, postID)
	if err != nil {
		return fmt.Errorf("CommentStorage.UpdateVotes - exec error:%w", err)
	}

	return nil
}

func (s *CommentsStorage) GetByID(ID int) (*models.Comment, error) {
	var c models.Comment
	err := s.db.QueryRow(`SELECT * FROM comments WHERE id=?`, ID).Scan(&c.ID, &c.Content, &c.UserID, &c.PostID, &c.CreatedDate, &c.LikeCount, &c.DislikeCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoPost
		} else {
			return nil, err
		}
	}
	return &c, nil
}

func (s *CommentsStorage) CommentExist(commentID int) (bool, error) {
	var commentExists bool

	err := s.db.QueryRow("SELECT EXISTS (SELECT 1 FROM comments WHERE id = ?)", commentID).Scan(&commentExists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, models.ErrWrongVoteItem
		}
		return false, fmt.Errorf("CommentStorage.CommentExist - scan error: %w", err)
	}
	return commentExists, nil
}
