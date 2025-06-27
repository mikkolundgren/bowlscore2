# Bowling Score Calculator

A web-based bowling score calculator with multiplayer support and persistent storage using Go backend and SQLite database.

Generated pretty much with Warp AI https://www.warp.dev/warp-ai

## Features

- **Interactive Bowling Scorecard**: Real-time score calculation with proper bowling rules
- **Dynamic Multiple Series**: Add unlimited bowling series with "Add New Series" button
- **Custom Player IDs**: User-defined player names/IDs for each series
- **Persistent Storage**: Save scores to SQLite database with timestamps
- **Score Management**: View and delete saved scores
- **Real-time Updates**: Automatic score calculation as you enter throws
- **Proper Bowling Rules**: Handles strikes, spares, and 10th frame special scoring
- **HTTPS Security**: Secure connections with self-signed certificates
- **HTTP Redirect**: Automatic redirect from HTTP to HTTPS

## Project Structure

```
bowlscore/
├── main.go              # Go backend server
├── go.mod               # Go module dependencies
├── certs/
│   ├── server.crt       # SSL certificate
│   └── server.key       # SSL private key
├── static/
│   ├── index.html       # Frontend HTML/CSS/JavaScript
│   ├── favicon.svg      # Bowling ball favicon
│   └── favicon.ico      # Traditional favicon
├── start.sh             # Startup script
├── README.md            # This file
└── bowling_scores.db    # SQLite database (created automatically)
```

## API Endpoints

- `POST /api/scores` - Save a new bowling score
- `GET /api/scores` - List all saved scores (newest first)
- `DELETE /api/scores/{id}` - Delete a specific score by ID

## Requirements

- Go 1.21 or later
- SQLite3 (included with macOS)

## Installation & Usage

1. **Start the server:**
   ```bash
   ./start.sh
   ```
   
   Or manually:
   ```bash
   go run main.go
   ```

2. **Open your browser:**
   Navigate to `https://localhost:8443`
   
   **Note:** You'll see a security warning due to the self-signed certificate. This is normal for development. Click "Advanced" and "Proceed to localhost" (or similar) to continue.

3. **Use the calculator:**
   - Click "Add New Series" to create a new bowling scorecard
   - Enter a custom Player ID for each series
   - Enter bowling scores using standard notation:
     - Numbers 0-10 for pins knocked down
     - "X" for strikes
     - "/" for spares
     - "-" for gutter balls
   - Click "Save Score" on each series to store it individually in the database
   - Each series is saved separately with its own timestamp

## Database Schema

The SQLite database contains a single table:

```sql
CREATE TABLE bowling_scores (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id TEXT NOT NULL,
    frames TEXT NOT NULL,           -- JSON array of throws
    total_score INTEGER NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

## Bowling Rules Implemented

- **Regular Frames (1-9)**: Sum of two throws, max 10 pins
- **Strike**: All 10 pins in first throw, bonus = next 2 throws
- **Spare**: All 10 pins in two throws, bonus = next 1 throw  
- **10th Frame**: Up to 3 throws if strike/spare, proper bonus calculation
- **Input Validation**: Prevents invalid entries (e.g., total > 10 pins per frame)

## Technologies Used

- **Backend**: Go with Gorilla Mux router
- **Database**: SQLite3
- **Frontend**: Vanilla HTML/CSS/JavaScript
- **Styling**: Custom CSS with responsive design

## Development

To modify the frontend, edit `static/index.html`. To modify the backend API, edit `main.go`.

The server automatically serves static files from the `static/` directory and provides CORS headers for development.
