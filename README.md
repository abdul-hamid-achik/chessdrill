# ChessDrill

A web application that teaches users to recognize chess board squares using algebraic notation and learn legal piece movements.

## Features

- **Name the Square** - A square is highlighted, you type its algebraic notation
- **Find the Square** - Given a notation, click the correct square on the board
- **Piece Movement** - Learn where each piece can legally move
- **Move Notation** - Read algebraic notation and identify moves
- **Progress Tracking** - Accuracy stats, response times, heat maps
- **User Accounts** - Save your progress and track improvement over time

## Tech Stack

- **Backend**: Go 1.25+, Chi router, templ templates
- **Database**: MongoDB
- **Frontend**: TypeScript, Tailwind CSS v4, HTMX
- **Chess Board**: Lichess chessground + chessops
- **Build Tools**: esbuild, Tailwind CLI

## Prerequisites

- Go 1.25+
- Node.js 18+
- MongoDB 7+ (or use Docker)
- Task (optional, for task runner)

## Quick Start

### 1. Clone and Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install templ CLI
go install github.com/a-h/templ/cmd/templ@latest

# Install npm dependencies
npm install
```

### 2. Start MongoDB

Using Docker:
```bash
docker compose up -d mongo
```

Or connect to an existing MongoDB instance by updating `.env`.

### 3. Build and Run

```bash
# Generate templ files
templ generate

# Build frontend
npm run build

# Build and run server
go build -o bin/chessdrill ./cmd/server
./bin/chessdrill
```

Or use the Taskfile:
```bash
task setup    # Install all dependencies
task build    # Build everything
task dev      # Start with hot reload
```

### 4. Open the App

Visit http://localhost:8080

## Development

### Hot Reload

Install air for Go hot reload:
```bash
go install github.com/air-verse/air@latest
```

Run with hot reload:
```bash
# In one terminal: watch frontend
npm run watch

# In another terminal: watch Go
air
```

Or use:
```bash
task dev
```

### Project Structure

```
chessdrill/
├── cmd/server/          # Entry point
├── internal/
│   ├── config/          # Configuration
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # Auth & logging
│   ├── model/           # Data models
│   ├── mongo/           # Database client
│   ├── repository/      # Data access
│   ├── server/          # Router setup
│   └── service/         # Business logic
├── templates/
│   ├── components/      # Reusable components
│   ├── pages/           # Full page templates
│   └── partials/        # HTMX partials
├── static/
│   ├── css/             # Tailwind CSS
│   └── js/src/          # TypeScript source
└── bin/                 # Build output
```

### Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
PORT=8080
ENV=development
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=chessdrill
SESSION_SECRET=your-secret-key-min-32-chars
SESSION_MAX_AGE=604800
LOG_LEVEL=debug
```

## API Routes

### Pages (SSR)
- `GET /` - Landing page
- `GET /login` - Login form
- `GET /register` - Registration form
- `GET /dashboard` - User stats (auth required)
- `GET /drill` - Drill selection (auth required)
- `GET /drill/:type` - Active drill (auth required)
- `GET /stats` - Detailed analytics (auth required)
- `GET /settings` - User preferences (auth required)

### Auth
- `POST /auth/register` - Create account
- `POST /auth/login` - Login
- `POST /auth/logout` - Logout

### Drill API
- `POST /api/drill/start` - Start session
- `POST /api/drill/check` - Check answer
- `POST /api/drill/end` - End session

### Stats API
- `GET /api/stats/heatmap` - Square accuracy data

## License

MIT
