import { Chess } from 'chessops/chess';
import { parseFen, makeFen } from 'chessops/fen';
import { parseSquare, makeSquare } from 'chessops/util';

export type PieceType = 'knight' | 'bishop' | 'rook' | 'queen' | 'king' | 'pawn';

export class ChessLogic {
  private position: Chess | null = null;

  constructor() {
    this.position = null;
  }

  // Set position from FEN
  setPosition(fen: string): boolean {
    // Skip empty or invalid FENs silently
    if (!fen || fen === '8/8/8/8/8/8/8/8' || fen === '8/8/8/8/8/8/8/8 w - - 0 1') {
      this.position = null;
      return true;
    }

    const setup = parseFen(fen);
    if (setup.isErr) {
      // Silently ignore invalid FENs for simple drills
      this.position = null;
      return false;
    }

    const pos = Chess.fromSetup(setup.value);
    if (pos.isErr) {
      // Position may be invalid for chess rules but FEN parsed ok
      // This is fine for drill display purposes
      this.position = null;
      return false;
    }

    this.position = pos.value;
    return true;
  }

  // Get legal moves for a piece on a given square
  getLegalMoves(square: string): string[] {
    if (!this.position) return [];

    const sq = parseSquare(square);
    if (sq === undefined) return [];

    const dests = this.position.dests(sq);
    const moves: string[] = [];

    for (const dest of dests) {
      moves.push(makeSquare(dest));
    }

    return moves;
  }

  // Get all legal moves as a map
  getAllLegalMoves(): Map<string, string[]> {
    const result = new Map<string, string[]>();
    if (!this.position) return result;

    for (const [from, squares] of this.position.allDests()) {
      const fromKey = makeSquare(from);
      const toKeys: string[] = [];

      for (const to of squares) {
        toKeys.push(makeSquare(to));
      }

      if (toKeys.length > 0) {
        result.set(fromKey, toKeys);
      }
    }

    return result;
  }

  // Get current FEN
  getFen(): string {
    if (!this.position) return '8/8/8/8/8/8/8/8 w - - 0 1';
    return makeFen(this.position.toSetup());
  }

  // Check if a move is legal
  isMoveLegal(from: string, to: string): boolean {
    const legalMoves = this.getLegalMoves(from);
    return legalMoves.includes(to);
  }

  // Get piece at square
  getPieceAt(square: string): { type: string; color: 'white' | 'black' } | null {
    if (!this.position) return null;

    const sq = parseSquare(square);
    if (sq === undefined) return null;

    const piece = this.position.board.get(sq);
    if (!piece) return null;

    const roleMap: { [key: string]: string } = {
      'king': 'k',
      'queen': 'q',
      'rook': 'r',
      'bishop': 'b',
      'knight': 'n',
      'pawn': 'p',
    };

    return {
      type: roleMap[piece.role] || piece.role,
      color: piece.color,
    };
  }
}

// Calculate legal moves for a piece type on a given square (simplified, ignoring other pieces)
export function calculatePieceMoves(pieceType: PieceType, square: string): string[] {
  const file = square.charCodeAt(0) - 'a'.charCodeAt(0);
  const rank = parseInt(square[1]) - 1;
  const moves: string[] = [];

  const isValidSquare = (f: number, r: number): boolean => {
    return f >= 0 && f < 8 && r >= 0 && r < 8;
  };

  const toSquare = (f: number, r: number): string => {
    return String.fromCharCode('a'.charCodeAt(0) + f) + (r + 1);
  };

  switch (pieceType) {
    case 'knight':
      const knightMoves = [
        [-2, -1], [-2, 1], [-1, -2], [-1, 2],
        [1, -2], [1, 2], [2, -1], [2, 1]
      ];
      for (const [df, dr] of knightMoves) {
        const newFile = file + df;
        const newRank = rank + dr;
        if (isValidSquare(newFile, newRank)) {
          moves.push(toSquare(newFile, newRank));
        }
      }
      break;

    case 'bishop':
      for (const [df, dr] of [[-1, -1], [-1, 1], [1, -1], [1, 1]]) {
        for (let i = 1; i < 8; i++) {
          const newFile = file + df * i;
          const newRank = rank + dr * i;
          if (!isValidSquare(newFile, newRank)) break;
          moves.push(toSquare(newFile, newRank));
        }
      }
      break;

    case 'rook':
      for (const [df, dr] of [[-1, 0], [1, 0], [0, -1], [0, 1]]) {
        for (let i = 1; i < 8; i++) {
          const newFile = file + df * i;
          const newRank = rank + dr * i;
          if (!isValidSquare(newFile, newRank)) break;
          moves.push(toSquare(newFile, newRank));
        }
      }
      break;

    case 'queen':
      // Queen = Rook + Bishop
      for (const [df, dr] of [
        [-1, -1], [-1, 1], [1, -1], [1, 1],
        [-1, 0], [1, 0], [0, -1], [0, 1]
      ]) {
        for (let i = 1; i < 8; i++) {
          const newFile = file + df * i;
          const newRank = rank + dr * i;
          if (!isValidSquare(newFile, newRank)) break;
          moves.push(toSquare(newFile, newRank));
        }
      }
      break;

    case 'king':
      for (const [df, dr] of [
        [-1, -1], [-1, 0], [-1, 1],
        [0, -1], [0, 1],
        [1, -1], [1, 0], [1, 1]
      ]) {
        const newFile = file + df;
        const newRank = rank + dr;
        if (isValidSquare(newFile, newRank)) {
          moves.push(toSquare(newFile, newRank));
        }
      }
      break;

    case 'pawn':
      // White pawn moves (simplified - forward only)
      if (rank < 7) {
        moves.push(toSquare(file, rank + 1));
        if (rank === 1) {
          moves.push(toSquare(file, rank + 2));
        }
      }
      break;
  }

  return moves;
}
