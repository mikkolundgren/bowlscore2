package models


type BowlingScore struct {
	ID         int       `json:"id"`
	PlayerID   string    `json:"player_id"`
	Frames     string    `json:"frames"` // JSON string of frame data
	TotalScore int       `json:"total_score"`
	Timestamp  string	 `json:"timestamp"`
}

type ScoreRequest struct {
	PlayerID   string `json:"player_id"`
	Frames     string `json:"frames"`
	TotalScore int    `json:"total_score"`
}

type DailyProgress struct {
	Date        string  `json:"date"`
	Average     float64 `json:"average"`
	GamesPlayed int     `json:"games_played"`
	StrikesRatio	float64 `json:"strikes_ratio"`
	SparesRatio		float64 `json:"spares_ratio"`
	MissesRatio		float64 `json:"misses_ratio"`
}
