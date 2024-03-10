package service

import (
	"errors"
	"fmt"
	"forum/internal/models"
	s "forum/internal/storage"
	"time"

	"github.com/gofrs/uuid/v5"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	storage s.Users
}

func NewUserService(storage s.Users) *UserService {
	return &UserService{storage}
}

func (u *UserService) SignUp(input UserSignUpInput) error {
	hashPass, err := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
	if err != nil {
		return err
	}
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashPass),
	}
	if err = u.storage.Create(user); err != nil {
		return err
	}
	return nil
}

func (u *UserService) SignIn(input UserSignInInput) (string, time.Time, error) {
	user, err := u.storage.GetByCredentials(input.Email)
	if err != nil {
		return "", time.Time{}, err
	}
	token, hasSession, err := u.CheckSession(user.Id)
	if hasSession {
		err = u.storage.DeleteSession(token)
		if err != nil {
			return "", time.Time{}, err
		}
	}
	if err != nil {
		return "", time.Time{}, fmt.Errorf("UserService.SignIn session check error: %w", err)
	}
	if err != nil {
		return "", time.Time{}, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", time.Time{}, models.ErrMissMatchEmailOrPass
		} else {
			return "", time.Time{}, fmt.Errorf("hash and password comparing error: %w", err)
		}
	}

	return u.createSession(user.Id)
}

func (u *UserService) Logout(token string) error {
	return u.storage.DeleteSession(token)
}

func (u *UserService) GetByToken(token string) (*models.User, error) {
	return u.storage.GetByToken(token)
}

func (u *UserService) createSession(userID int) (string, time.Time, error) {
	u2, err := uuid.NewV4()
	if err != nil {
		return "", time.Time{}, fmt.Errorf("create session error: %w", err)
	}

	currentTime := time.Now()
	localTimeZone := currentTime.Location()
	expiredDate := time.Now().Add(time.Minute * 60).In(localTimeZone)
	session := models.Session{
		UserID:      userID,
		Token:       u2.String(),
		ExpiredDate: expiredDate,
	}

	err = u.storage.SetSession(session)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("create session error: %w", err)
	}
	return u2.String(), expiredDate, nil
}

func (u *UserService) CheckSession(userID int) (string, bool, error) {
	token, ok, err := u.storage.CheckSession(userID)
	if err != nil {
		return "", false, err
	}
	return token, ok, nil
}
