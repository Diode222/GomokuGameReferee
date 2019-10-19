package model

type MatchDataModel struct {
	GameID           string       `json:"game_id"`
	BoardLength      int          `json:"borad_length"`
	BoardHeight      int          `json:"board_height"`
	Player1ID        string       `json:"player1_id"`
	Player2ID        string       `json:"player2_id"`
	Player1FirstHand bool         `json:"player1_first_hand"`
	MaxThingTime     int          `json:"max_thinking_time"`
	Winner           int          `json:"winner"`
	StartTime        int64        `json:"start_time"`
	EndTime          int64        `json:"end_time"`
	Operations       []*Operation `json:"operations"`
	FoulPlayer       int          `json:"foul_player"` // 0: no foul, 1: player1 foul, 2: player2 foul
}

const (
	// foul player
	NO_FOUL      = 0
	PLAYER1_FOUL = 1
	PLAYER2_FOUL = 2

	// operation type
	BLANK = 0
	WHITE = 1
	NONE  = 2
)

type Operation struct {
	Player   int       `json:"player"`
	Position *Position `json:"position"`
	Type     int       `json:"piece_type"`
}

type Position struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
}
