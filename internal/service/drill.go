package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/abdul-hamid-achik/chessdrill/internal/model"
	"github.com/abdul-hamid-achik/chessdrill/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var (
	files = []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	ranks = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
)

type DrillService struct {
	drillSessionRepo *repository.DrillSessionRepository
	attemptRepo      *repository.AttemptRepository
}

func NewDrillService(drillSessionRepo *repository.DrillSessionRepository, attemptRepo *repository.AttemptRepository) *DrillService {
	return &DrillService{
		drillSessionRepo: drillSessionRepo,
		attemptRepo:      attemptRepo,
	}
}

func (s *DrillService) StartSession(ctx context.Context, userID bson.ObjectID, drillType model.DrillType, inputMethod model.InputMethod, perspective string) (*model.DrillSession, *model.Question, error) {
	session := model.NewDrillSession(userID, drillType, inputMethod, perspective)
	if err := s.drillSessionRepo.Create(ctx, session); err != nil {
		return nil, nil, err
	}

	question := s.GenerateQuestion(drillType, "")
	return session, question, nil
}

func (s *DrillService) CheckAnswer(ctx context.Context, sessionID, userID bson.ObjectID, drillType model.DrillType, question, correctAnswer, userAnswer string, responseMs int) (bool, *model.Question, error) {
	correctAnswer = strings.ToLower(strings.TrimSpace(correctAnswer))
	userAnswer = strings.ToLower(strings.TrimSpace(userAnswer))

	attempt := model.NewAttempt(sessionID, userID, drillType, question, correctAnswer, userAnswer, responseMs)
	if err := s.attemptRepo.Create(ctx, attempt); err != nil {
		return false, nil, err
	}

	nextQuestion := s.GenerateQuestion(drillType, "")
	return attempt.Correct, nextQuestion, nil
}

func (s *DrillService) EndSession(ctx context.Context, sessionID bson.ObjectID) (*model.DrillSessionSummary, error) {
	summary, err := s.attemptRepo.GetSessionSummary(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	if err := s.drillSessionRepo.EndSession(ctx, sessionID, *summary); err != nil {
		return nil, err
	}

	return summary, nil
}

func (s *DrillService) GenerateQuestion(drillType model.DrillType, pieceType string) *model.Question {
	switch drillType {
	case model.DrillTypeNameSquare:
		return s.generateNameSquareQuestion()
	case model.DrillTypeFindSquare:
		return s.generateFindSquareQuestion()
	case model.DrillTypePieceMovement:
		return s.generatePieceMovementQuestion(pieceType)
	case model.DrillTypeMoveNotation:
		return s.generateMoveNotationQuestion()
	default:
		return s.generateNameSquareQuestion()
	}
}

func (s *DrillService) RandomSquare() string {
	fileIdx, _ := rand.Int(rand.Reader, big.NewInt(8))
	rankIdx, _ := rand.Int(rand.Reader, big.NewInt(8))
	return files[fileIdx.Int64()] + ranks[rankIdx.Int64()]
}

func (s *DrillService) generateNameSquareQuestion() *model.Question {
	target := s.RandomSquare()
	return &model.Question{
		Type:   model.DrillTypeNameSquare,
		Target: target,
		Prompt: "",
		FEN:    "8/8/8/8/8/8/8/8 w - - 0 1",
	}
}

func (s *DrillService) generateFindSquareQuestion() *model.Question {
	target := s.RandomSquare()
	return &model.Question{
		Type:   model.DrillTypeFindSquare,
		Target: target,
		Prompt: target,
		FEN:    "8/8/8/8/8/8/8/8 w - - 0 1",
	}
}

func (s *DrillService) generatePieceMovementQuestion(pieceType string) *model.Question {
	if pieceType == "" {
		pieceTypes := []string{"knight", "bishop", "rook", "queen", "king"}
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(pieceTypes))))
		pieceType = pieceTypes[idx.Int64()]
	}

	square := s.RandomSquare()
	fen := s.generateSinglePieceFEN(pieceType, square)

	return &model.Question{
		Type:   model.DrillTypePieceMovement,
		Target: square,
		Prompt: fmt.Sprintf("Where can the %s move?", pieceType),
		FEN:    fen,
		Metadata: map[string]string{
			"piece_type": pieceType,
		},
	}
}

func (s *DrillService) generateMoveNotationQuestion() *model.Question {
	return s.generatePieceMovementQuestion("knight")
}

func (s *DrillService) generateSinglePieceFEN(pieceType, square string) string {
	pieceMap := map[string]rune{
		"knight": 'N',
		"bishop": 'B',
		"rook":   'R',
		"queen":  'Q',
		"king":   'K',
		"pawn":   'P',
	}

	piece := pieceMap[pieceType]
	if piece == 0 {
		piece = 'N'
	}

	file := int(square[0] - 'a')
	rank := int(square[1] - '1')

	var fenRows []string
	for r := 7; r >= 0; r-- {
		if r == rank {
			row := ""
			if file > 0 {
				row += fmt.Sprintf("%d", file)
			}
			row += string(piece)
			if file < 7 {
				row += fmt.Sprintf("%d", 7-file)
			}
			fenRows = append(fenRows, row)
		} else {
			fenRows = append(fenRows, "8")
		}
	}

	return strings.Join(fenRows, "/") + " w - - 0 1"
}

func (s *DrillService) GetSessionAttempts(ctx context.Context, sessionID bson.ObjectID) ([]model.Attempt, error) {
	return s.attemptRepo.FindBySessionID(ctx, sessionID)
}
