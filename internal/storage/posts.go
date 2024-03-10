package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"log"
	"strings"
)

type PostsStorage struct {
	db *sql.DB
}

func NewPostsStorage(db *sql.DB) *PostsStorage {
	return &PostsStorage{
		db,
	}
}

func (s *PostsStorage) Create(post models.Post) (int, error) {
	const op = "PostStorage.Create"
	stmt, err := s.db.Prepare(`INSERT INTO posts (title, content, image_path, user_id, category_id, like_count, dislike_count)
	values (?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, fmt.Errorf("%s: db.Prepare error: %w", op, err)
	}
	result, err := stmt.Exec(post.Title, post.Content, post.Image, post.UserID, post.CategoryID, post.LikeCount, post.DislikeCount)
	if err != nil {
		return 0, fmt.Errorf("%s: Exec error: %w", op, err)
	}
	lastRowID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s-lastrow error: %w", op, err)
	}

	for _, tag := range post.Tags {
		stmt, err = s.db.Prepare(`INSERT INTO posts_tags (post_id, tags_id) VALUES (?,?)`)
		if err != nil {
			return 0, fmt.Errorf("%s - prepare error: %w", op, err)
		}
		_, err = stmt.Exec(lastRowID, tag)
		if err != nil {
			return 0, fmt.Errorf("%s - posttags exec error: %w", op, err)
		}
	}
	return int(lastRowID), nil
}

func (s *PostsStorage) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	rows, err := s.db.Query("SELECT * FROM categories")
	if err != nil {
		return nil, fmt.Errorf("PostStorage.GetAllCategories: scan error:%w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Url); err != nil {
			return nil, fmt.Errorf("PostStorage.GetAllCategories: scan error: %w", err)
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (s *PostsStorage) GetPostByID(id int) (*models.Post, error) {
	var post models.Post
	err := s.db.QueryRow(`SELECT * FROM posts WHERE id=?`, id).Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.UserID, &post.CategoryID, &post.CreatedDate, &post.LikeCount, &post.DislikeCount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoPost
		} else {
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostsStorage) GetPostsByUserID(userID int) ([]models.Post, error) {
	op := "PostsStorage.GetPostByUserID"
	var posts []models.Post

	rows, err := s.db.Query(`SELECT * FROM posts where user_id = ?`, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("%s - query error: %w", op, err)
		}
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		if err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.UserID, &post.CategoryID, &post.CreatedDate, &post.LikeCount, &post.DislikeCount); err != nil {
			return nil, fmt.Errorf("%s: scan error: %w", op, err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostsStorage) GetPostsByCategory(category string) ([]models.Post, error) {
	const op = "GetPostByCategory"
	var posts []models.Post
	rows, err := s.db.Query(`
		select p.*
		from posts p
		join categories c on c.id = p.category_id
		where c.name = ?
		`, category)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrPostsByCatNotFound
		}
		return nil, fmt.Errorf("%s query error: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.UserID, &post.CategoryID, &post.CreatedDate, &post.LikeCount, &post.DislikeCount); err != nil {
			return nil, fmt.Errorf("PostStorage.PostByCategory: scan error: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (s *PostsStorage) PostsCountByCategory(id int) (int, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(*) FROM posts WHERE category_id =?`, id).Scan(&count)
	if err != nil {
		return -1, fmt.Errorf("PostsCountByCategory - scan error: %w", err)
	}

	return count, nil
}

func (s *PostsStorage) GetLastPostByCategory(id int) (models.Post, error) {
	var post models.Post
	err := s.db.QueryRow(`SELECT * FROM posts WHERE category_id =?`, id).Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.UserID, &post.CategoryID, &post.CreatedDate, &post.LikeCount, &post.DislikeCount)
	if err != nil {
		if errors.As(err, &sql.ErrNoRows) {
			return models.Post{}, models.ErrPostsByCatNotFound
		}
		return models.Post{}, fmt.Errorf("GetLastPostByCategory - scan error: %w", err)
	}
	return post, nil
}

func (s *PostsStorage) GetCategoryByUrl(url string) (*models.Category, error) {
	var category models.Category

	err := s.db.QueryRow(`SELECT * from categories WHERE url=?`, url).Scan(&category.ID, &category.Name, &category.Description, &category.Url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrWrongCategory
		}
		return nil, fmt.Errorf("PostStorage.GetCategoryByUrl - scan error: %w", err)
	}

	return &category, nil
}

func (s *PostsStorage) PostExist(postID int) bool {
	var postExists bool

	err := s.db.QueryRow("SELECT EXISTS (SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&postExists)
	if err != nil {
		return false
	}
	return postExists
}

func (s *PostsStorage) GetVotesCount(postID int) (int, int, error) {
	var likeCount, dislikeCount int

	query := `SELECT like_count, dislike_count FROM posts WHERE id = ?`
	err := s.db.QueryRow(query, postID).Scan(&likeCount, &dislikeCount)
	if err != nil {
		return 0, 0, fmt.Errorf("PostsStorage.GetVotesCount - scan error:%w", err)
	}

	return likeCount, dislikeCount, nil
}

func (s *PostsStorage) UpdateVotes(postID int, action string, delta int) error {
	query := fmt.Sprintf("UPDATE posts SET %s = %s + ? WHERE id = ?", action, action)
	_, err := s.db.Exec(query, delta, postID)
	if err != nil {
		return fmt.Errorf("PostsStorage.UpdateVotes - exec error:%w", err)
	}

	return nil
}

func (s *PostsStorage) GetMapCategories() (map[int]string, error) {
	categories := make(map[int]string)

	rows, err := s.db.Query(`SELECT id, url FROM categories`)
	if err != nil {
		return nil, fmt.Errorf("PostStorage.GetMapCategories - query error: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var url string
		err = rows.Scan(&id, &url)
		if err != nil {
			return nil, fmt.Errorf("PostStorage.GetMapCategories - scan error: %w", err)
		}
		categories[id] = url
	}

	return categories, nil
}

func (s *PostsStorage) GetVotesPosts(userID int) ([]models.Post, error) {
	var posts []models.Post

	query := `
		SELECT posts.*
		FROM posts
		INNER JOIN votes ON posts.id = votes.item_id
		WHERE votes.user_id = ? AND votes.item = 'post'
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("PostsStorage.GetVotesPosts - query error:%w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.UserID, &post.CategoryID, &post.CreatedDate, &post.LikeCount, &post.DislikeCount)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			} else {
				return nil, fmt.Errorf("PostsStorage.GetVotesPosts - scan error:%w", err)
			}
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return posts, nil
}

func (s *PostsStorage) TagExists(tagID int) (bool, error) {
	var tagExists bool

	err := s.db.QueryRow("SELECT EXISTS (SELECT 1 FROM tags WHERE id = ?)", tagID).Scan(&tagExists)
	if err != nil {
		return false, fmt.Errorf("PostStorage.TagExists - scan error: %w", err)
	}
	return tagExists, nil
}

func (s *PostsStorage) GetTags(categoryID int) ([]models.Tag, error) {
	var op = "PostStorage.GetTags"
	var tags []models.Tag

	rows, err := s.db.Query(`SELECT * FROM tags where category_id = ?`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s - query error: %w", op, err)
	}
	defer rows.Close()
	for rows.Next() {
		var tag models.Tag
		err = rows.Scan(&tag.ID, &tag.CategoryID, &tag.Name)
		if err != nil {
			return nil, fmt.Errorf("%s - scan error: %w", op, err)
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (s *PostsStorage) GetByTags(tagIDs []int) ([]models.Post, error) {
	op := "PostStorage.GetByTags"
	query := `
        SELECT DISTINCT p.*
        FROM posts p
        JOIN posts_tags pt ON p.ID = pt.post_id
        JOIN tags t ON pt.tags_id = t.id
        WHERE t.id IN (?)
    `

	args := make([]interface{}, len(tagIDs))
	for i, id := range tagIDs {
		args[i] = id
	}

	placeholders := make([]string, len(tagIDs))
	for i := range tagIDs {
		placeholders[i] = "?"
	}
	placeholderStr := strings.Join(placeholders, ", ")

	query = strings.Replace(query, "?", placeholderStr, 1)

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s - query error:%w", op, err)
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var post models.Post
		if err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.Image, &post.UserID, &post.CategoryID, &post.CreatedDate, &post.LikeCount, &post.DislikeCount); err != nil {
			return nil, fmt.Errorf("%s - scan error:%w", op, err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNotValidTagID
		}
		return nil, fmt.Errorf("%s - rows error:%w", op, err)
	}
	return posts, nil
}

func (s *PostsStorage) GetTagsByPostID(postID int) ([]string, error) {
	var tags []string
	rows, err := s.db.Query(`
	SELECT t.name
	FROM tags t
	JOIN posts_tags on t.id = posts_tags.tags_id
	WHERE posts_tags.post_id = ?`, postID)
	if err != nil {
		return nil, fmt.Errorf("PostStorage.GetTagsByPostID - query error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, fmt.Errorf("PostStorage.GetTagsByPostID - scan error: %w", err)
		}
		tags = append(tags, name)
	}
	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoPost
		}
		return nil, fmt.Errorf("PostStorage.GetTagsByPostID - rows error:%w", err)
	}
	return tags, nil
}

func (s *PostsStorage) GetTagsByCategoryID(categoryID int) (map[int][]string, error) {
	m := make(map[int][]string)
	op := "PostStoratge.GetTagsByCategoryID"
	rows, err := s.db.Query(`
		SELECT p.id, t.name
		FROM posts p
		JOIN posts_tags pt ON p.id = pt.post_id
		JOIN tags t ON pt.tags_id = t.id
		WHERE p.category_id =?`, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%s - query error: %w", op, err)
	}
	for rows.Next() {
		var id int
		var tag string
		err = rows.Scan(&id, &tag)
		if err != nil {
			return nil, fmt.Errorf("%s - scan error: %w", op, err)
		}
		m[id] = append(m[id], tag)
	}
	if err = rows.Err(); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoPost
		}
		return nil, fmt.Errorf("%s - rows error:%w", op, err)
	}
	return m, nil
}
