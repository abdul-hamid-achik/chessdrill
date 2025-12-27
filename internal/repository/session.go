package repository

import (
	"context"
	"errors"

	"github.com/abdul-hamid-achik/chessdrill/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrSessionNotFound = errors.New("session not found")

type SessionRepository struct {
	collection *mongo.Collection
}

func NewSessionRepository(db *mongo.Database) *SessionRepository {
	return &SessionRepository{
		collection: db.Collection("sessions"),
	}
}

func (r *SessionRepository) Create(ctx context.Context, session *model.AuthSession) error {
	_, err := r.collection.InsertOne(ctx, session)
	return err
}

func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*model.AuthSession, error) {
	var session model.AuthSession
	err := r.collection.FindOne(ctx, bson.M{"token": token}).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteByToken(ctx context.Context, token string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"token": token})
	return err
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID bson.ObjectID) error {
	_, err := r.collection.DeleteMany(ctx, bson.M{"user_id": userID})
	return err
}
