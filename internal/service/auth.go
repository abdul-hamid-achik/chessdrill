package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/abdul-hamid-achik/chessdrill/internal/model"
	"github.com/abdul-hamid-achik/chessdrill/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	maxAge      int
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, maxAge int) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		maxAge:      maxAge,
	}
}

func (s *AuthService) Register(ctx context.Context, email, username, password string) (*model.User, string, error) {
	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user := model.NewUser(email, username, string(hash))
	if err := s.userRepo.Create(ctx, user); err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return nil, "", ErrUserExists
		}
		return nil, "", err
	}

	// Create session
	token, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*model.User, string, error) {
	// Find user
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", err
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Create session
	token, err := s.createSession(ctx, user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.DeleteByToken(ctx, token)
}

func (s *AuthService) ValidateSession(ctx context.Context, token string) (*model.User, error) {
	session, err := s.sessionRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if session.IsExpired() {
		_ = s.sessionRepo.DeleteByToken(ctx, token)
		return nil, errors.New("session expired")
	}

	user, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) createSession(ctx context.Context, userID bson.ObjectID) (string, error) {
	token, err := generateToken(32)
	if err != nil {
		return "", err
	}

	session := model.NewAuthSession(userID, token, s.maxAge)
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", err
	}

	return token, nil
}

func generateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
