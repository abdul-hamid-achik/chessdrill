package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abdul-hamid-achik/chessdrill/internal/middleware"
	"github.com/abdul-hamid-achik/chessdrill/internal/model"
	"github.com/abdul-hamid-achik/chessdrill/internal/service"
	"github.com/abdul-hamid-achik/chessdrill/templates/partials"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type DrillHandler struct {
	drillService *service.DrillService
}

func NewDrillHandler(drillService *service.DrillService) *DrillHandler {
	return &DrillHandler{
		drillService: drillService,
	}
}

type StartDrillRequest struct {
	DrillType   string `json:"drill_type"`
	InputMethod string `json:"input_method"`
	Perspective string `json:"perspective"`
}

type StartDrillResponse struct {
	SessionID string          `json:"session_id"`
	Question  *model.Question `json:"question"`
}

func (h *DrillHandler) StartDrill(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req StartDrillRequest
	if err := r.ParseForm(); err == nil {
		req.DrillType = r.FormValue("drill_type")
		req.InputMethod = r.FormValue("input_method")
		req.Perspective = r.FormValue("perspective")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
	}

	if req.DrillType == "" {
		req.DrillType = "name_square"
	}
	if req.InputMethod == "" {
		req.InputMethod = "type"
	}
	if req.Perspective == "" {
		req.Perspective = "white"
	}

	session, question, err := h.drillService.StartSession(
		r.Context(),
		user.ID,
		model.DrillType(req.DrillType),
		model.InputMethod(req.InputMethod),
		req.Perspective,
	)
	if err != nil {
		http.Error(w, "Failed to start drill", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		partials.DrillQuestion(session.ID.Hex(), question).Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(StartDrillResponse{
		SessionID: session.ID.Hex(),
		Question:  question,
	})
}

type CheckAnswerRequest struct {
	SessionID  string `json:"session_id"`
	Target     string `json:"target"`
	Answer     string `json:"answer"`
	ResponseMs int    `json:"response_ms"`
	DrillType  string `json:"drill_type"`
}

func (h *DrillHandler) CheckAnswer(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CheckAnswerRequest
	if err := r.ParseForm(); err == nil {
		req.SessionID = r.FormValue("session_id")
		req.Target = r.FormValue("target")
		req.Answer = r.FormValue("answer")
		req.DrillType = r.FormValue("drill_type")
		if ms := r.FormValue("response_ms"); ms != "" {
			req.ResponseMs, _ = strconv.Atoi(ms)
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
	}

	sessionID, err := bson.ObjectIDFromHex(req.SessionID)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	if req.DrillType == "" {
		req.DrillType = "name_square"
	}

	correct, nextQuestion, err := h.drillService.CheckAnswer(
		r.Context(),
		sessionID,
		user.ID,
		model.DrillType(req.DrillType),
		req.Target,
		req.Target, // correctAnswer is the target
		req.Answer,
		req.ResponseMs,
	)
	if err != nil {
		http.Error(w, "Failed to check answer", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		message := "Incorrect!"
		if correct {
			message = "Correct!"
		}
		partials.Feedback(correct, message, req.SessionID, nextQuestion).Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"correct":       correct,
		"next_question": nextQuestion,
	})
}

func (h *DrillHandler) EndDrill(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionIDStr := r.FormValue("session_id")
	if sessionIDStr == "" {
		// Try JSON body
		var req struct {
			SessionID string `json:"session_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err == nil {
			sessionIDStr = req.SessionID
		}
	}

	sessionID, err := bson.ObjectIDFromHex(sessionIDStr)
	if err != nil {
		http.Error(w, "Invalid session ID", http.StatusBadRequest)
		return
	}

	summary, err := h.drillService.EndSession(r.Context(), sessionID)
	if err != nil {
		http.Error(w, "Failed to end drill", http.StatusInternalServerError)
		return
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		partials.SessionSummary(summary).Render(r.Context(), w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summary)
}

// GetLegalMoves returns legal moves for a piece at a given position
func (h *DrillHandler) GetLegalMoves(w http.ResponseWriter, r *http.Request) {
	// This is handled client-side by chessops
	// But we can provide a server-side fallback
	fen := r.URL.Query().Get("fen")
	square := r.URL.Query().Get("square")

	if fen == "" || square == "" {
		http.Error(w, "Missing fen or square parameter", http.StatusBadRequest)
		return
	}

	// For now, return empty - the client handles this with chessops
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"moves": []string{},
		"note":  "Legal moves are calculated client-side using chessops",
	})
}
