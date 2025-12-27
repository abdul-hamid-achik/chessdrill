import { ChessBoard } from './board';
import { DrillController } from './drill';
import { Timer } from './timer';
import { HeatmapRenderer } from './heatmap';

// Global app state
interface AppState {
  board: ChessBoard | null;
  drill: DrillController | null;
  timer: Timer | null;
}

const app: AppState = {
  board: null,
  drill: null,
  timer: null,
};

// Initialize when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  initializeBoard();
  initializeDrill();
  initializeHeatmap();
  setupEventListeners();
});

function initializeBoard(): void {
  const boardElement = document.getElementById('board');
  if (boardElement) {
    app.board = new ChessBoard(boardElement);
    
    // Expose for debugging
    (window as any).chessBoard = app.board;
  }
}

function initializeDrill(): void {
  const drillContainer = document.getElementById('drill-container');
  if (drillContainer && app.board) {
    const drillType = drillContainer.dataset.drillType || 'name_square';
    app.drill = new DrillController(app.board, drillType);
    app.timer = new Timer();
    
    // Expose for debugging
    (window as any).chessDrill = app.drill;
  }
}

function initializeHeatmap(): void {
  const canvas = document.getElementById('heatmap-canvas') as HTMLCanvasElement;
  const dataScript = document.getElementById('heatmap-data');
  
  if (canvas && dataScript) {
    try {
      const data = JSON.parse(dataScript.textContent || '{}');
      const renderer = new HeatmapRenderer(canvas);
      renderer.render(data);
    } catch (e) {
      console.error('Failed to render heatmap:', e);
    }
  }
}

function setupEventListeners(): void {
  // Flip board button
  const flipButton = document.getElementById('flip-board');
  if (flipButton && app.board) {
    flipButton.addEventListener('click', () => {
      app.board?.toggleOrientation();
    });
  }
  
  // Toggle coordinates button
  const coordsButton = document.getElementById('toggle-coords');
  if (coordsButton && app.board) {
    coordsButton.addEventListener('click', () => {
      app.board?.toggleCoordinates();
    });
  }
  
  // End drill button
  const endButton = document.getElementById('end-drill');
  if (endButton) {
    endButton.addEventListener('click', () => {
      app.drill?.endSession();
    });
  }
  
  // Listen for question ready events (from HTMX)
  window.addEventListener('chessdrill:questionReady', ((e: CustomEvent) => {
    const detail = e.detail;
    if (app.drill && app.board) {
      app.drill.handleQuestionReady(detail);
      app.timer?.start();
    }
  }) as EventListener);
  
  // Listen for next question events
  window.addEventListener('chessdrill:nextQuestion', ((e: CustomEvent) => {
    const detail = e.detail;
    if (app.drill && app.board) {
      app.drill.handleNextQuestion(detail);
      app.timer?.start();
    }
  }) as EventListener);
  
  // File/rank button clicks for name_square drill
  document.addEventListener('click', (e) => {
    const target = e.target as HTMLElement;
    
    if (target.classList.contains('file-btn')) {
      const file = target.dataset.file;
      if (file) {
        handleFileClick(file);
      }
    }
    
    if (target.classList.contains('rank-btn')) {
      const rank = target.dataset.rank;
      if (rank) {
        handleRankClick(rank);
      }
    }
  });
}

// State for button input method
let selectedFile = '';

function handleFileClick(file: string): void {
  selectedFile = file;
  
  // Highlight selected file button
  document.querySelectorAll('.file-btn').forEach(btn => {
    btn.classList.remove('selected');
  });
  document.querySelector(`.file-btn[data-file="${file}"]`)?.classList.add('selected');
}

function handleRankClick(rank: string): void {
  if (selectedFile) {
    const square = selectedFile + rank;
    const input = document.getElementById('answer-input') as HTMLInputElement;
    if (input) {
      input.value = square;
      // Auto-submit
      const form = document.getElementById('answer-form') as HTMLFormElement;
      if (form) {
        // Set response time
        const responseMs = document.getElementById('response-ms') as HTMLInputElement;
        if (responseMs && app.timer) {
          responseMs.value = String(app.timer.elapsed());
        }
        form.requestSubmit();
      }
    }
    selectedFile = '';
    // Clear button selections
    document.querySelectorAll('.file-btn, .rank-btn').forEach(btn => {
      btn.classList.remove('selected');
    });
  }
}

// Handle form submission to add response time
document.addEventListener('submit', (e) => {
  const form = e.target as HTMLFormElement;
  if (form.id === 'answer-form') {
    const responseMs = form.querySelector('#response-ms') as HTMLInputElement;
    if (responseMs && app.timer) {
      responseMs.value = String(app.timer.elapsed());
    }
  }
});

// Export for use in templates
(window as any).chessDrill = {
  setQuestion: (question: any) => {
    if (app.drill) {
      app.drill.setQuestion(question);
    }
  },
  getBoard: () => app.board,
};
