package model

// SquareAccuracy represents accuracy data for a specific square
type SquareAccuracy struct {
	Square   string  `json:"square"`
	Total    int     `json:"total"`
	Correct  int     `json:"correct"`
	Accuracy float64 `json:"accuracy"`
}

// DrillStats represents aggregated stats for a drill type
type DrillStats struct {
	DrillType       string  `json:"drill_type"`
	TotalAttempts   int     `json:"total_attempts"`
	CorrectAttempts int     `json:"correct_attempts"`
	Accuracy        float64 `json:"accuracy"`
	AvgResponseMs   int     `json:"avg_response_ms"`
	BestStreak      int     `json:"best_streak"`
	CurrentStreak   int     `json:"current_streak"`
}

// OverallStats represents user's overall performance
type OverallStats struct {
	TotalSessions   int          `json:"total_sessions"`
	TotalAttempts   int          `json:"total_attempts"`
	OverallAccuracy float64      `json:"overall_accuracy"`
	AvgResponseMs   int          `json:"avg_response_ms"`
	BestStreak      int          `json:"best_streak"`
	DrillStats      []DrillStats `json:"drill_stats"`
}

// HeatmapData represents accuracy data for the heat map visualization
type HeatmapData struct {
	Squares []SquareAccuracy `json:"squares"`
}
