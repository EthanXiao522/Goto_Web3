package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/xyd/web3-learning-tracker/internal/model"
	"github.com/xyd/web3-learning-tracker/internal/repository"
)

var (
	ErrDuplicateEmail    = errors.New("email already registered")
	ErrDuplicateUsername = errors.New("username already taken")
	ErrInvalidCreds      = errors.New("invalid email or password")
)

type AuthService struct {
	userRepo  *repository.UserRepo
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepo, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: jwtSecret}
}

func (s *AuthService) Register(username, email, password string) (*model.User, error) {
	if _, err := s.userRepo.FindByEmail(email); err == nil {
		return nil, ErrDuplicateEmail
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("register check email: %w", err)
	}
	if _, err := s.userRepo.FindByUsername(username); err == nil {
		return nil, ErrDuplicateUsername
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, fmt.Errorf("register check username: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt: %w", err)
	}

	u := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
	}
	id, err := s.userRepo.Create(u)
	if err != nil {
		return nil, err
	}
	u.ID = id
	return u, nil
}

func (s *AuthService) Login(email, password string) (string, *model.User, error) {
	u, err := s.userRepo.FindByEmail(email)
	if errors.Is(err, repository.ErrNotFound) {
		return "", nil, ErrInvalidCreds
	}
	if err != nil {
		return "", nil, fmt.Errorf("login: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCreds
	}
	token, err := s.generateJWT(u.ID)
	if err != nil {
		return "", nil, err
	}
	return token, u, nil
}

func (s *AuthService) generateJWT(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
