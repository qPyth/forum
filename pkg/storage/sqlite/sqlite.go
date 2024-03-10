package sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	Conn *sql.DB
}

func NewStorage(pathToDb string) (*Storage, error) {
	conn, err := sql.Open("sqlite3", pathToDb)
	if err != nil {
		return nil, fmt.Errorf("datebase open: %w", err)
	}

	return &Storage{conn}, nil
}

func (s *Storage) ApplyMigrations() error {
	pathToMigrations := "./pkg/storage/migrations"
	dir, err := os.Open(pathToMigrations)
	if err != nil {
		return fmt.Errorf("open migrations dir: %w", err)
	}

	files, err := dir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".sql") {
			filePath := fmt.Sprintf("%s/%s", pathToMigrations, fileName)
			file, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("open migration file: %w", err)
			}
			_, err = s.Conn.Exec(string(file))
			if err != nil {
				return fmt.Errorf("exec migration: %w", err)
			}
		}
	}
	return nil
}
