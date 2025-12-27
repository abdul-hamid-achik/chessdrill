interface SquareAccuracy {
  square: string;
  total: number;
  correct: number;
  accuracy: number;
}

interface HeatmapData {
  squares: SquareAccuracy[];
}

export class HeatmapRenderer {
  private canvas: HTMLCanvasElement;
  private ctx: CanvasRenderingContext2D;
  private squareSize: number;

  constructor(canvas: HTMLCanvasElement) {
    this.canvas = canvas;
    this.ctx = canvas.getContext('2d')!;
    this.squareSize = canvas.width / 8;
  }

  // Render the heatmap
  render(data: HeatmapData): void {
    this.clear();

    // Create a map for quick lookup
    const accuracyMap = new Map<string, SquareAccuracy>();
    for (const sq of data.squares) {
      accuracyMap.set(sq.square, sq);
    }

    // Draw each square
    for (let file = 0; file < 8; file++) {
      for (let rank = 0; rank < 8; rank++) {
        const square = String.fromCharCode('a'.charCodeAt(0) + file) + (rank + 1);
        const sqData = accuracyMap.get(square);
        
        // Calculate position (flip rank for display)
        const x = file * this.squareSize;
        const y = (7 - rank) * this.squareSize;

        // Draw base square color (light/dark)
        const isLight = (file + rank) % 2 === 1;
        this.ctx.fillStyle = isLight ? '#f0d9b5' : '#b58863';
        this.ctx.fillRect(x, y, this.squareSize, this.squareSize);

        // Overlay accuracy color
        if (sqData && sqData.total > 0) {
          const color = this.getAccuracyColor(sqData.accuracy);
          this.ctx.fillStyle = color;
          this.ctx.fillRect(x, y, this.squareSize, this.squareSize);
        }

        // Draw square label
        this.ctx.fillStyle = isLight ? '#b58863' : '#f0d9b5';
        this.ctx.font = '12px sans-serif';
        this.ctx.textAlign = 'center';
        this.ctx.textBaseline = 'middle';
        this.ctx.fillText(
          square,
          x + this.squareSize / 2,
          y + this.squareSize / 2
        );

        // Draw accuracy percentage if data exists
        if (sqData && sqData.total > 0) {
          this.ctx.fillStyle = '#000';
          this.ctx.font = 'bold 14px sans-serif';
          this.ctx.fillText(
            `${Math.round(sqData.accuracy)}%`,
            x + this.squareSize / 2,
            y + this.squareSize / 2 + 14
          );
        }
      }
    }

    // Draw border
    this.ctx.strokeStyle = '#333';
    this.ctx.lineWidth = 2;
    this.ctx.strokeRect(0, 0, this.canvas.width, this.canvas.height);
  }

  // Get color based on accuracy (0-100)
  private getAccuracyColor(accuracy: number): string {
    // Green to red gradient based on accuracy
    // 100% = green (rgba 76, 175, 80, 0.6)
    // 50% = yellow (rgba 255, 235, 59, 0.6)
    // 0% = red (rgba 244, 67, 54, 0.6)

    let r: number, g: number, b: number;

    if (accuracy >= 50) {
      // Green to yellow (50-100)
      const ratio = (accuracy - 50) / 50;
      r = Math.round(255 - ratio * (255 - 76));
      g = Math.round(235 + ratio * (175 - 235));
      b = Math.round(59 + ratio * (80 - 59));
    } else {
      // Yellow to red (0-50)
      const ratio = accuracy / 50;
      r = Math.round(244 + ratio * (255 - 244));
      g = Math.round(67 + ratio * (235 - 67));
      b = Math.round(54 + ratio * (59 - 54));
    }

    return `rgba(${r}, ${g}, ${b}, 0.6)`;
  }

  // Clear the canvas
  clear(): void {
    this.ctx.clearRect(0, 0, this.canvas.width, this.canvas.height);
  }

  // Update canvas size
  resize(size: number): void {
    this.canvas.width = size;
    this.canvas.height = size;
    this.squareSize = size / 8;
  }
}
