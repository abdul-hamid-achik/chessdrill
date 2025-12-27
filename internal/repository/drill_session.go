package repository

import (
	"context"
	"errors"
	"time"

	"github.com/abdul-hamid-achik/chessdrill/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var ErrDrillSessionNotFound = errors.New("drill session not found")

type DrillSessionRepository struct {
	collection *mongo.Collection
}

func NewDrillSessionRepository(db *mongo.Database) *DrillSessionRepository {
	return &DrillSessionRepository{
		collection: db.Collection("drill_sessions"),
	}
}

func (r *DrillSessionRepository) Create(ctx context.Context, session *model.DrillSession) error {
	result, err := r.collection.InsertOne(ctx, session)
	if err != nil {
		return err
	}
	session.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (r *DrillSessionRepository) FindByID(ctx context.Context, id bson.ObjectID) (*model.DrillSession, error) {
	var session model.DrillSession
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&session)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrDrillSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *DrillSessionRepository) EndSession(ctx context.Context, id bson.ObjectID, summary model.DrillSessionSummary) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"ended_at": now,
			"summary":  summary,
		},
	}
	result, err := r.collection.UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return ErrDrillSessionNotFound
	}
	return nil
}

func (r *DrillSessionRepository) FindByUserID(ctx context.Context, userID bson.ObjectID, limit int) ([]model.DrillSession, error) {
	opts := &struct {
		Sort  bson.M
		Limit int64
	}{
		Sort:  bson.M{"started_at": -1},
		Limit: int64(limit),
	}

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	_ = opts // We'll use aggregation for sorting/limiting in production

	var sessions []model.DrillSession
	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

func (r *DrillSessionRepository) CountByUserID(ctx context.Context, userID bson.ObjectID) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"user_id": userID})
}
