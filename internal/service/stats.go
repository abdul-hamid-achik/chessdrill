package service

import (
	"context"

	"github.com/abdul-hamid-achik/chessdrill/internal/model"
	"github.com/abdul-hamid-achik/chessdrill/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type StatsService struct {
	attemptRepo      *repository.AttemptRepository
	drillSessionRepo *repository.DrillSessionRepository
}

func NewStatsService(attemptRepo *repository.AttemptRepository, drillSessionRepo *repository.DrillSessionRepository) *StatsService {
	return &StatsService{
		attemptRepo:      attemptRepo,
		drillSessionRepo: drillSessionRepo,
	}
}

// GetOverallStats returns overall stats for a user
func (s *StatsService) GetOverallStats(ctx context.Context, userID bson.ObjectID) (*model.OverallStats, error) {
	stats, err := s.attemptRepo.GetOverallStats(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Get session count
	sessionCount, err := s.drillSessionRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	stats.TotalSessions = int(sessionCount)

	// Get per-drill-type stats
	drillTypes := []model.DrillType{
		model.DrillTypeNameSquare,
		model.DrillTypeFindSquare,
		model.DrillTypePieceMovement,
		model.DrillTypeMoveNotation,
	}

	for _, dt := range drillTypes {
		drillStats, err := s.attemptRepo.GetDrillStats(ctx, userID, dt)
		if err != nil {
			continue
		}
		if drillStats.TotalAttempts > 0 {
			stats.DrillStats = append(stats.DrillStats, *drillStats)
		}
	}

	return stats, nil
}

// GetHeatmapData returns accuracy data for the heat map
func (s *StatsService) GetHeatmapData(ctx context.Context, userID bson.ObjectID) (*model.HeatmapData, error) {
	accuracies, err := s.attemptRepo.GetSquareAccuracy(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Fill in missing squares with zero data
	allSquares := make(map[string]bool)
	for _, acc := range accuracies {
		allSquares[acc.Square] = true
	}

	for _, file := range []string{"a", "b", "c", "d", "e", "f", "g", "h"} {
		for _, rank := range []string{"1", "2", "3", "4", "5", "6", "7", "8"} {
			square := file + rank
			if !allSquares[square] {
				accuracies = append(accuracies, model.SquareAccuracy{
					Square:   square,
					Total:    0,
					Correct:  0,
					Accuracy: 0,
				})
			}
		}
	}

	return &model.HeatmapData{
		Squares: accuracies,
	}, nil
}

// GetDrillStats returns stats for a specific drill type
func (s *StatsService) GetDrillStats(ctx context.Context, userID bson.ObjectID, drillType model.DrillType) (*model.DrillStats, error) {
	return s.attemptRepo.GetDrillStats(ctx, userID, drillType)
}
