package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// DrillType represents the type of drill
type DrillType string

const (
	DrillTypeNameSquare    DrillType = "name_square"
	DrillTypeFindSquare    DrillType = "find_square"
	DrillTypePieceMovement DrillType = "piece_movement"
	DrillTypeMoveNotation  DrillType = "move_notation"
)

// InputMethod represents how user provides answers
type InputMethod string

const (
	InputMethodType       InputMethod = "type"
	InputMethodClick      InputMethod = "click"
	InputMethodGrid       InputMethod = "grid"
	InputMethodBoardClick InputMethod = "board_click"
)

// DrillSessionSummary contains aggregated stats for a session
type DrillSessionSummary struct {
	TotalAttempts int `bson:"total_attempts" json:"total_attempts"`
	Correct       int `bson:"correct" json:"correct"`
	AvgResponseMs int `bson:"avg_response_ms" json:"avg_response_ms"`
	StreakBest    int `bson:"streak_best" json:"streak_best"`
}

// DrillSession represents a practice session
type DrillSession struct {
	ID          bson.ObjectID       `bson:"_id,omitempty" json:"id"`
	UserID      bson.ObjectID       `bson:"user_id" json:"user_id"`
	DrillType   DrillType           `bson:"drill_type" json:"drill_type"`
	InputMethod InputMethod         `bson:"input_method" json:"input_method"`
	Perspective string              `bson:"perspective" json:"perspective"`
	StartedAt   time.Time           `bson:"started_at" json:"started_at"`
	EndedAt     *time.Time          `bson:"ended_at,omitempty" json:"ended_at"`
	Summary     DrillSessionSummary `bson:"summary" json:"summary"`
}

func NewDrillSession(userID bson.ObjectID, drillType DrillType, inputMethod InputMethod, perspective string) *DrillSession {
	return &DrillSession{
		UserID:      userID,
		DrillType:   drillType,
		InputMethod: inputMethod,
		Perspective: perspective,
		StartedAt:   time.Now(),
		Summary:     DrillSessionSummary{},
	}
}

// AttemptMetadata contains additional info for certain drill types
type AttemptMetadata struct {
	PieceType  string `bson:"piece_type,omitempty" json:"piece_type,omitempty"`
	FromSquare string `bson:"from_square,omitempty" json:"from_square,omitempty"`
	FEN        string `bson:"fen,omitempty" json:"fen,omitempty"`
}

// Attempt represents a single question-answer attempt
type Attempt struct {
	ID            bson.ObjectID   `bson:"_id,omitempty" json:"id"`
	SessionID     bson.ObjectID   `bson:"session_id" json:"session_id"`
	UserID        bson.ObjectID   `bson:"user_id" json:"user_id"`
	DrillType     DrillType       `bson:"drill_type" json:"drill_type"`
	Question      string          `bson:"question" json:"question"`
	CorrectAnswer string          `bson:"correct_answer" json:"correct_answer"`
	UserAnswer    string          `bson:"user_answer" json:"user_answer"`
	Correct       bool            `bson:"correct" json:"correct"`
	ResponseMs    int             `bson:"response_ms" json:"response_ms"`
	AnsweredAt    time.Time       `bson:"answered_at" json:"answered_at"`
	Metadata      AttemptMetadata `bson:"metadata,omitempty" json:"metadata,omitempty"`
}

func NewAttempt(sessionID, userID bson.ObjectID, drillType DrillType, question, correctAnswer, userAnswer string, responseMs int) *Attempt {
	return &Attempt{
		SessionID:     sessionID,
		UserID:        userID,
		DrillType:     drillType,
		Question:      question,
		CorrectAnswer: correctAnswer,
		UserAnswer:    userAnswer,
		Correct:       correctAnswer == userAnswer,
		ResponseMs:    responseMs,
		AnsweredAt:    time.Now(),
	}
}

// Question represents a drill question sent to the client
type Question struct {
	Type     DrillType         `json:"type"`
	Target   string            `json:"target"`
	Prompt   string            `json:"prompt"`
	FEN      string            `json:"fen,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
}
