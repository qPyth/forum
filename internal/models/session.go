package models

import (
	"errors"
	"time"
)

type Session struct {
	ID          int
	UserID      int
	Token       string
	ExpiredDate time.Time
}

var (
	ErrSetSession = errors.New("set session error")
)
