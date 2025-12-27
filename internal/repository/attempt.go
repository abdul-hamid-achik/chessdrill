package repository

import (
	"context"

	"github.com/abdul-hamid-achik/chessdrill/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AttemptRepository struct {
	collection *mongo.Collection
}

func NewAttemptRepository(db *mongo.Database) *AttemptRepository {
	return &AttemptRepository{
		collection: db.Collection("attempts"),
	}
}

func (r *AttemptRepository) Create(ctx context.Context, attempt *model.Attempt) error {
	result, err := r.collection.InsertOne(ctx, attempt)
	if err != nil {
		return err
	}
	attempt.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (r *AttemptRepository) FindBySessionID(ctx context.Context, sessionID bson.ObjectID) ([]model.Attempt, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var attempts []model.Attempt
	if err := cursor.All(ctx, &attempts); err != nil {
		return nil, err
	}
	return attempts, nil
}

func (r *AttemptRepository) GetSessionSummary(ctx context.Context, sessionID bson.ObjectID) (*model.DrillSessionSummary, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"session_id": sessionID}},
		{"$group": bson.M{
			"_id":             nil,
			"total_attempts":  bson.M{"$sum": 1},
			"correct":         bson.M{"$sum": bson.M{"$cond": []interface{}{"$correct", 1, 0}}},
			"avg_response_ms": bson.M{"$avg": "$response_ms"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		TotalAttempts int     `bson:"total_attempts"`
		Correct       int     `bson:"correct"`
		AvgResponseMs float64 `bson:"avg_response_ms"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &model.DrillSessionSummary{}, nil
	}

	// Calculate best streak
	bestStreak := r.calculateBestStreak(ctx, sessionID)

	return &model.DrillSessionSummary{
		TotalAttempts: results[0].TotalAttempts,
		Correct:       results[0].Correct,
		AvgResponseMs: int(results[0].AvgResponseMs),
		StreakBest:    bestStreak,
	}, nil
}

func (r *AttemptRepository) calculateBestStreak(ctx context.Context, sessionID bson.ObjectID) int {
	cursor, err := r.collection.Find(ctx, bson.M{"session_id": sessionID})
	if err != nil {
		return 0
	}
	defer cursor.Close(ctx)

	var attempts []model.Attempt
	if err := cursor.All(ctx, &attempts); err != nil {
		return 0
	}

	bestStreak := 0
	currentStreak := 0
	for _, attempt := range attempts {
		if attempt.Correct {
			currentStreak++
			if currentStreak > bestStreak {
				bestStreak = currentStreak
			}
		} else {
			currentStreak = 0
		}
	}
	return bestStreak
}

// GetSquareAccuracy returns accuracy stats for each square
func (r *AttemptRepository) GetSquareAccuracy(ctx context.Context, userID bson.ObjectID) ([]model.SquareAccuracy, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{
			"_id":     "$correct_answer",
			"total":   bson.M{"$sum": 1},
			"correct": bson.M{"$sum": bson.M{"$cond": []interface{}{"$correct", 1, 0}}},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		Square  string `bson:"_id"`
		Total   int    `bson:"total"`
		Correct int    `bson:"correct"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	var accuracies []model.SquareAccuracy
	for _, r := range results {
		accuracy := 0.0
		if r.Total > 0 {
			accuracy = float64(r.Correct) / float64(r.Total) * 100
		}
		accuracies = append(accuracies, model.SquareAccuracy{
			Square:   r.Square,
			Total:    r.Total,
			Correct:  r.Correct,
			Accuracy: accuracy,
		})
	}
	return accuracies, nil
}

// GetDrillStats returns stats for a specific drill type
func (r *AttemptRepository) GetDrillStats(ctx context.Context, userID bson.ObjectID, drillType model.DrillType) (*model.DrillStats, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID, "drill_type": drillType}},
		{"$group": bson.M{
			"_id":             nil,
			"total_attempts":  bson.M{"$sum": 1},
			"correct":         bson.M{"$sum": bson.M{"$cond": []interface{}{"$correct", 1, 0}}},
			"avg_response_ms": bson.M{"$avg": "$response_ms"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		TotalAttempts int     `bson:"total_attempts"`
		Correct       int     `bson:"correct"`
		AvgResponseMs float64 `bson:"avg_response_ms"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return &model.DrillStats{DrillType: string(drillType)}, nil
	}

	accuracy := 0.0
	if results[0].TotalAttempts > 0 {
		accuracy = float64(results[0].Correct) / float64(results[0].TotalAttempts) * 100
	}

	return &model.DrillStats{
		DrillType:       string(drillType),
		TotalAttempts:   results[0].TotalAttempts,
		CorrectAttempts: results[0].Correct,
		Accuracy:        accuracy,
		AvgResponseMs:   int(results[0].AvgResponseMs),
	}, nil
}

// GetOverallStats returns overall stats for a user
func (r *AttemptRepository) GetOverallStats(ctx context.Context, userID bson.ObjectID) (*model.OverallStats, error) {
	pipeline := []bson.M{
		{"$match": bson.M{"user_id": userID}},
		{"$group": bson.M{
			"_id":             nil,
			"total_attempts":  bson.M{"$sum": 1},
			"correct":         bson.M{"$sum": bson.M{"$cond": []interface{}{"$correct", 1, 0}}},
			"avg_response_ms": bson.M{"$avg": "$response_ms"},
		}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []struct {
		TotalAttempts int     `bson:"total_attempts"`
		Correct       int     `bson:"correct"`
		AvgResponseMs float64 `bson:"avg_response_ms"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	stats := &model.OverallStats{}
	if len(results) > 0 {
		stats.TotalAttempts = results[0].TotalAttempts
		if stats.TotalAttempts > 0 {
			stats.OverallAccuracy = float64(results[0].Correct) / float64(stats.TotalAttempts) * 100
		}
		stats.AvgResponseMs = int(results[0].AvgResponseMs)
	}

	return stats, nil
}
