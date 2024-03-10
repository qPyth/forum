package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
	"github.com/mattn/go-sqlite3"
	"strings"
)

type UsersStorage struct {
	db *sql.DB
}

func NewUsersStorage(db *sql.DB) *UsersStorage {
	return &UsersStorage{db: db}
}

func (s *UsersStorage) Create(user models.User) error {
	stmt, err := s.db.Prepare("INSERT INTO users (username, email, password) values (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("user storage prepare error: %w", err)
	}

	if _, err = stmt.Exec(user.Username, user.Email, user.Password); err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			if strings.Contains(sqliteErr.Error(), "UNIQUE constraint failed: users.username") {
				return models.ErrUsernameExist
			} else if strings.Contains(sqliteErr.Error(), "UNIQUE constraint failed: users.email") {
				return models.ErrEmailExist
			}
		}
		return fmt.Errorf("user storage exec error: %w", err)
	}

	return nil
}

func (s *UsersStorage) GetByCredentials(email string) (models.User, error) {
	var user models.User
	err := s.db.QueryRow("SELECT id, username, email, password FROM users WHERE email = ?", email).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, models.ErrUserNotFound
		}
		return models.User{}, fmt.Errorf("scan error: %w", err)
	}
	return user, nil
}

func (s *UsersStorage) SetSession(session models.Session) error {

	stmt := fmt.Sprintf("INSERT INTO %s (user_id, token, expired_date) VALUES (?, ?, ?)", tokenTable)

	_, err := s.db.Exec(stmt, session.UserID, session.Token, session.ExpiredDate)
	if err != nil {
		return fmt.Errorf("set session error: %w", err)
	}
	return nil
}

func (s *UsersStorage) CheckSession(userID int) (string, bool, error) {
	query := "SELECT token FROM sessions WHERE user_id = ?"
	var token string
	err := s.db.QueryRow(query, userID).Scan(&token)
	if err == nil {
		return token, true, nil
	} else if errors.Is(err, sql.ErrNoRows) {
		return "", false, nil
	} else {
		return "", false, fmt.Errorf("UserStorage.CheckSession scan error: %w", err)
	}
}

func (s *UsersStorage) UpdateSession(session models.Session) error {
	stmt := fmt.Sprintf("UPDATE %s SET token=?, expired_date=? where user_id=?", tokenTable)
	_, err := s.db.Exec(stmt, session.Token, session.ExpiredDate, session.UserID)
	if err != nil {
		return fmt.Errorf("update session error:%w", err)
	}
	return nil
}

func (s *UsersStorage) DeleteSession(token string) error {
	stmt := fmt.Sprintf("DELETE FROM %s where token=?", tokenTable)

	_, err := s.db.Exec(stmt, token)
	if err != nil {
		return fmt.Errorf("delete session error: %w", err)
	}
	return nil
}

func (s *UsersStorage) GetByID(id int) (models.User, error) {
	var user models.User

	err := s.db.QueryRow(`SELECT * FROM users where id = ?`, id).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("UserRepo.GetByID - scan error:%w", err)
	}
	return user, nil
}

func (s *UsersStorage) GetByToken(token string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRow(`
		SELECT users.id, users.username, users.email, users.password
		FROM users
		JOIN sessions ON users.id = sessions.user_id
		WHERE sessions.token = ?`, token).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}
		return nil, fmt.Errorf("UsersStorage.GetByToken - scan error: %w", err)
	}
	return &user, nil
}

func (s *UsersStorage) GetAuthorsOfPost(categoryID int) (map[int]string, error) {

	rows, err := s.db.Query(`
	SELECT posts.id, users.username
	FROM posts
	INNER JOIN users ON users.id = posts.user_id
	WHERE posts.category_id = ?
`, categoryID)
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
