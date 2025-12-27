import { Chessground } from '@lichess-org/chessground';
import type { Api } from '@lichess-org/chessground/api';
import type { Config } from '@lichess-org/chessground/config';
import type { Key, Color } from '@lichess-org/chessground/types';

export class ChessBoard {
  private ground: Api;
  private element: HTMLElement;
  private orientation: Color = 'white';
  private showCoordinates: boolean = true;

  constructor(element: HTMLElement, config?: Partial<Config>) {
    this.element = element;
    
    const defaultConfig: Config = {
      fen: '8/8/8/8/8/8/8/8',
      orientation: this.orientation,
      coordinates: this.showCoordinates,
      movable: {
        free: false,
        color: undefined,
      },
      draggable: {
        enabled: false,
      },
      selectable: {
        enabled: true,
      },
      events: {
        select: (key) => this.onSquareSelect(key),
      },
      highlight: {
        lastMove: false,
        check: false,
      },
    };

    this.ground = Chessground(element, { ...defaultConfig, ...config });
  }

  // Get the chessground API
  getApi(): Api {
    return this.ground;
  }

  // Set position from FEN
  setPosition(fen: string): void {
    this.ground.set({ fen });
  }

  // Highlight a specific square
  highlightSquare(square: string): void {
    this.ground.set({
      drawable: {
        autoShapes: [
          { 
            orig: square as Key, 
            brush: 'yellow' 
          }
        ],
      },
    });
  }

  // Clear all highlights
  clearHighlights(): void {
    this.ground.set({
      drawable: {
        autoShapes: [],
      },
    });
  }

  // Show multiple squares highlighted (for piece movement drill)
  highlightSquares(squares: string[], brush: string = 'green'): void {
    this.ground.set({
      drawable: {
        autoShapes: squares.map(sq => ({
          orig: sq as Key,
          brush: brush,
        })),
      },
    });
  }

  // Toggle board orientation
  toggleOrientation(): void {
    this.orientation = this.orientation === 'white' ? 'black' : 'white';
    this.ground.set({ orientation: this.orientation });
  }

  // Set board orientation
  setOrientation(color: Color): void {
    this.orientation = color;
    this.ground.set({ orientation: color });
  }

  // Toggle coordinates display
  toggleCoordinates(): void {
    this.showCoordinates = !this.showCoordinates;
    this.ground.set({ coordinates: this.showCoordinates });
  }

  // Enable click selection mode
  enableSelection(callback: (square: string) => void): void {
    this.onSquareSelectCallback = callback;
    this.ground.set({
      selectable: {
        enabled: true,
      },
    });
  }

  // Disable selection
  disableSelection(): void {
    this.onSquareSelectCallback = null;
    this.ground.set({
      selectable: {
        enabled: false,
      },
    });
  }

  // Set a piece on the board
  setPiece(square: string, piece: { role: string; color: Color }): void {
    const pieces = new Map();
    pieces.set(square as Key, piece);
    this.ground.set({ fen: '8/8/8/8/8/8/8/8' }); // Clear first
    this.ground.setPieces(pieces);
  }

  // Callback for square selection
  private onSquareSelectCallback: ((square: string) => void) | null = null;

  private onSquareSelect(key: Key): void {
    // Emit custom event for drill handlers
    document.dispatchEvent(new CustomEvent('chess:select', { 
      detail: { square: key } 
    }));

    if (this.onSquareSelectCallback) {
      this.onSquareSelectCallback(key);
    }
  }

  // Get current orientation
  getOrientation(): Color {
    return this.orientation;
  }

  // Destroy the board
  destroy(): void {
    this.ground.destroy();
  }
}
