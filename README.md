# Bowling Score Calculator

A web-based bowling score calculator with multiplayer support and persistent storage using Go backend with JSON file storage.

## Features

- **Interactive Bowling Scorecard**: Real-time score calculation with proper bowling rules
- **Dynamic Multiple Series**: Add unlimited bowling series with "Add New Series" button  
- **Custom Player IDs**: User-defined player names/IDs for each series
- **Persistent Storage**: Save scores to JSON file with timestamps
- **Score Management**: View and delete saved scores
- **Real-time Updates**: Automatic score calculation as you enter throws
- **Proper Bowling Rules**: Handles strikes, spares, and 10th frame special scoring
- **HTTPS Security**: Secure connections with self-signed certificates
- **HTTP Redirect**: Automatic redirect from HTTP to HTTPS

## API Endpoints

- `POST /api/scores` - Save a new bowling score
- `GET /api/scores` - List all saved scores (newest first)  
- `DELETE /api/scores/{id}` - Delete a specific score by ID
- `GET /api/player-progress` - Get player progress statistics

## Requirements

- Go 1.21 or later

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
   
   **Note:** You'll see a security warning due to the self-signed certificate. This is normal for development. Click "Advanced" and "Proceed to localhost" to continue.

3. **Use the calculator:**
   - Click "Add New Series" to create a new bowling scorecard
   - Enter a custom Player ID for each series
   - Enter bowling scores using standard notation:
     - Numbers 0-10 for pins knocked down
     - "X" for strikes
     - "/" for spares
     - "-" for gutter balls
   - Click "Save Score" on each series to store it
   - Each series is saved separately with its own timestamp

## Data Storage Format

Scores are stored in `scores.json` with the following format:

```json
[
  {
    "id": 1,
    "player_id": "player1",
    "frames": "[10,7,3,9,0,10,0,8,8,2,0,6,10,10,10,8,1]",
    "total_score": 167,
    "timestamp": "2023-11-15T14:30:00Z"
  }
]
```

## Bowling Rules Implemented

- **Regular Frames (1-9)**: Sum of two throws, max 10 pins
- **Strike**: All 10 pins in first throw, bonus = next 2 throws  
- **Spare**: All 10 pins in two throws, bonus = next 1 throw
- **10th Frame**: Up to 3 throws if strike/spare, proper bonus calculation
- **Input Validation**: Prevents invalid entries (e.g., total > 10 pins per frame)

## Technologies Used

- **Backend**: Go with Gorilla Mux router
- **Storage**: JSON file with thread-safe operations
- **Frontend**: Vanilla HTML/CSS/JavaScript  
- **Styling**: Custom CSS with responsive design

## Development

- Modify frontend: Edit `static/index.html`
- Modify backend: Edit `main.go` and files in `database/` directory
- The server automatically serves static files from `static/` directory