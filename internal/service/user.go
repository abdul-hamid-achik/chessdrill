package service

import (
	"context"

	"github.com/abdul-hamid-achik/chessdrill/internal/model"
	"github.com/abdul-hamid-achik/chessdrill/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) GetByID(ctx context.Context, userID bson.ObjectID) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}

func (s *UserService) UpdatePreferences(ctx context.Context, userID bson.ObjectID, prefs model.Preferences) error {
	return s.userRepo.UpdatePreferences(ctx, userID, prefs)
}
