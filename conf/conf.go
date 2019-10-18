package conf

const (
	// BACKEND_SERVICE_NAME
	BACKEND_SERVICE_NAME = "backend_service"

	// LOG_PATH_NAME
	LOG_PATH_NAME = "game_log"
)

type PlayerState struct {
	IsFirstHand bool
	PieceType int
}

var (
	PLAYER1_STATE *PlayerState
	PLAYER2_STATE *PlayerState

	FIRST_HAND_PORT = 47777
	BACK_HAND_PORT = 47778
)

var (
	ETCD_ADDR string

	IP = "127.0.0.1"

	PLAYER1_PORT string

	PLAYER2_PORT string

	PLAYER1_FIRST_HAND string // "true": true, "false": false

	BOARD_LENGTH = 15
	BOARD_HEIGHT = 15

	// max thinking time of each player
	MAX_THINKING_TIME int

	GAME_ID string
)
