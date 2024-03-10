package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/internal/models"
)

type VoteStorage struct {
	db *sql.DB
}

func NewVoteStorage(db *sql.DB) *VoteStorage {
	return &VoteStorage{db: db}
}

func (v VoteStorage) Create(vote models.Vote) error {
	op := "VoteStorage.Create"
	stmt, err := v.db.Prepare(`INSERT INTO votes (user_id, action, item_id, item) values (?,?,?,?)`)
	if err != nil {
		return fmt.Errorf("%s - prepare error: %w", op, err)
	}

	_, err = stmt.Exec(vote.UserID, vote.Action, vote.ItemID, vote.Item)
	if err != nil {
		return fmt.Errorf("%s - exec error: %w", op, err)
	}
	return nil
}

func (v VoteStorage) IsExists(vote models.Vote) (*models.Vote, error) {
	var res models.Vote

	err := v.db.QueryRow("SELECT * FROM votes WHERE user_id = ? and item_id = ? and item = ?", vote.UserID, vote.ItemID, vote.Item).Scan(&res.ID, &res.UserID, &res.Action, &res.ItemID, &res.Item)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("VotesStorage.IsExists - scan error: %w", err)
		}
	}
	return &res, nil
}

func (v VoteStorage) Delete(voteID int) error {
	stmt, err := v.db.Prepare(`DELETE from votes where id = ?`)
	if err != nil {
		return fmt.Errorf("VoteStorage.Delete - prepare error:%w", err)
	}
	_, err = stmt.Exec(voteID)
	if err != nil {
		return fmt.Errorf("VoteStorage.Delete - exec error:%w", err)
	}
	return nil
}

func (v VoteStorage) Update(voteID int, action string) error {
	stmt, err := v.db.Prepare(`UPDATE votes SET action = ? where id = ?`)
	if err != nil {
		return fmt.Errorf("VoteStorage.Update - prepare error:%w", err)
	}
	_, err = stmt.Exec(action, voteID)
	if err != nil {
		return fmt.Errorf("VoteStorage.Update - exec error:%w", err)
	}
	return nil
}
