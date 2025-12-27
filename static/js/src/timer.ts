export class Timer {
  private startTime: number = 0;
  private running: boolean = false;

  // Start the timer
  start(): void {
    this.startTime = performance.now();
    this.running = true;
  }

  // Stop the timer and return elapsed time in milliseconds
  stop(): number {
    if (!this.running) return 0;
    this.running = false;
    return this.elapsed();
  }

  // Get elapsed time without stopping
  elapsed(): number {
    if (this.startTime === 0) return 0;
    return Math.round(performance.now() - this.startTime);
  }

  // Reset the timer
  reset(): void {
    this.startTime = 0;
    this.running = false;
  }

  // Check if timer is running
  isRunning(): boolean {
    return this.running;
  }
}
