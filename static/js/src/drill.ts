import { ChessBoard } from './board';
import { ChessLogic, calculatePieceMoves, PieceType } from './chess';

interface Question {
  sessionId: string;
  target: string;
  prompt: string;
  fen: string;
  type: string;
}

interface DrillStats {
  total: number;
  correct: number;
  streak: number;
  bestStreak: number;
  totalResponseMs: number;
}

export class DrillController {
  private board: ChessBoard;
  private chess: ChessLogic;
  private drillType: string;
  private currentQuestion: Question | null = null;
  private sessionId: string | null = null;
  private stats: DrillStats = {
    total: 0,
    correct: 0,
    streak: 0,
    bestStreak: 0,
    totalResponseMs: 0,
  };

  constructor(board: ChessBoard, drillType: string) {
    this.board = board;
    this.chess = new ChessLogic();
    this.drillType = drillType;

    this.setupBoardEvents();
  }

  private setupBoardEvents(): void {
    // Listen for square selections
    document.addEventListener('chess:select', ((e: CustomEvent) => {
      const square = e.detail.square;
      this.handleSquareClick(square);
    }) as EventListener);
  }

  // Handle question ready event (from HTMX)
  handleQuestionReady(detail: Question): void {
    this.sessionId = detail.sessionId;
    this.setQuestion(detail);
  }

  // Handle next question event
  handleNextQuestion(detail: Question): void {
    this.setQuestion(detail);
    this.updateStatsDisplay();
  }

  // Set the current question
  setQuestion(question: Question): void {
    this.currentQuestion = question;

    // Update board based on drill type
    if (question.fen && question.fen !== '8/8/8/8/8/8/8/8') {
      this.board.setPosition(question.fen);
      this.chess.setPosition(question.fen);
    }

    // Clear previous highlights
    this.board.clearHighlights();

    // Set up based on drill type
    switch (this.drillType) {
      case 'name_square':
        // Highlight the target square
        this.board.highlightSquare(question.target);
        break;

      case 'find_square':
        // No highlighting, user must find the square
        // Enable click selection
        this.board.enableSelection((square) => {
          this.handleFindSquareAnswer(square);
        });
        break;

      case 'piece_movement':
        // Show piece and let user identify legal moves
        this.board.highlightSquare(question.target);
        break;

      case 'move_notation':
        // Parse the move and require clicking destination
        break;
    }

    // Focus the input if it exists
    const input = document.getElementById('answer-input') as HTMLInputElement;
    if (input) {
      input.value = '';
      input.focus();
    }
  }

  // Handle square click for find_square drill
  private handleFindSquareAnswer(square: string): void {
    if (!this.currentQuestion || !this.sessionId) return;
    if (this.drillType !== 'find_square') return;

    // Submit the answer
    this.submitAnswer(square);
  }

  // Handle general square click
  private handleSquareClick(square: string): void {
    if (!this.currentQuestion) return;

    if (this.drillType === 'find_square') {
      this.handleFindSquareAnswer(square);
    } else if (this.drillType === 'piece_movement') {
      // For piece movement, check if clicked square is a legal destination
      const prompt = this.currentQuestion.prompt || '';
      const words = prompt.toLowerCase().split(' ');
      const pieceType = words[3] as PieceType;
      if (pieceType) {
        const legalMoves = calculatePieceMoves(pieceType, this.currentQuestion.target);
        console.log('Legal moves for', pieceType, 'on', this.currentQuestion.target, ':', legalMoves);
      }
    }
  }

  // Submit an answer via HTMX
  private submitAnswer(answer: string): void {
    if (!this.currentQuestion || !this.sessionId) return;

    const form = document.getElementById('answer-form') as HTMLFormElement;
    if (form) {
      // Update answer input
      const answerInput = form.querySelector('[name="answer"]') as HTMLInputElement;
      if (answerInput) {
        answerInput.value = answer;
      }

      // Submit the form
      form.requestSubmit();
    } else {
      // Fallback: make direct fetch request
      fetch('/api/drill/check', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
          'HX-Request': 'true',
        },
        body: new URLSearchParams({
          session_id: this.sessionId,
          target: this.currentQuestion.target,
          answer: answer,
          drill_type: this.drillType,
          response_ms: '0',
        }),
      })
        .then(res => res.text())
        .then(html => {
          const feedbackArea = document.getElementById('feedback-area');
          if (feedbackArea) {
            feedbackArea.innerHTML = html;
          }
        });
    }
  }

  // Update stats after answer
  updateStats(correct: boolean, responseMs: number): void {
    this.stats.total++;
    this.stats.totalResponseMs += responseMs;

    if (correct) {
      this.stats.correct++;
      this.stats.streak++;
      if (this.stats.streak > this.stats.bestStreak) {
        this.stats.bestStreak = this.stats.streak;
      }
    } else {
      this.stats.streak = 0;
    }

    this.updateStatsDisplay();
  }

  // Update the stats display in the UI
  private updateStatsDisplay(): void {
    const scoreDisplay = document.getElementById('score-display');
    const streakDisplay = document.getElementById('streak-display');
    const timeDisplay = document.getElementById('time-display');

    if (scoreDisplay) {
      const accuracy = this.stats.total > 0 
        ? Math.round((this.stats.correct / this.stats.total) * 100) 
        : 0;
      scoreDisplay.textContent = `${this.stats.correct}/${this.stats.total} (${accuracy}%)`;
    }

    if (streakDisplay) {
      streakDisplay.textContent = String(this.stats.streak);
    }

    if (timeDisplay && this.stats.total > 0) {
      const avgMs = Math.round(this.stats.totalResponseMs / this.stats.total);
      timeDisplay.textContent = `${avgMs}ms`;
    }
  }

  // End the current session
  endSession(): void {
    if (!this.sessionId) return;

    fetch('/api/drill/end', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/x-www-form-urlencoded',
        'HX-Request': 'true',
      },
      body: new URLSearchParams({
        session_id: this.sessionId,
      }),
    })
      .then(res => res.text())
      .then(html => {
        const activeArea = document.getElementById('drill-active-area');
        if (activeArea) {
          activeArea.innerHTML = html;
        }
        // Hide end button
        const endBtn = document.getElementById('end-drill');
        if (endBtn) {
          endBtn.style.display = 'none';
        }
      });
  }

  // Get current stats
  getStats(): DrillStats {
    return { ...this.stats };
  }

  // Reset stats
  resetStats(): void {
    this.stats = {
      total: 0,
      correct: 0,
      streak: 0,
      bestStreak: 0,
      totalResponseMs: 0,
    };
    this.updateStatsDisplay();
  }
}
