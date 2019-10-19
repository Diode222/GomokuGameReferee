package conf

const (
	// BACKEND_SERVICE_NAME
	BACKEND_SERVICE_NAME = "backend_service"

	// LOG_PATH_NAME
	LOG_PATH_NAME = "game_log"
)

type PlayerState struct {
	IsFirstHand bool
	PieceType   int
}

var (
	PLAYER1_STATE *PlayerState
	PLAYER2_STATE *PlayerState

	REFEREE_LOG_FILE = "referee.log"
)

var (
	IP = "127.0.0.1"

	PLAYER1_PORT = "10001"

	PLAYER2_PORT = "10002"

	PLAYER1_FIRST_HAND string // "true": true, "false": false

	BOARD_LENGTH = 15
	BOARD_HEIGHT = 15

	// max thinking time of each player
	MAX_THINKING_TIME int

	GAME_ID string

	PLAYER1_ID string

	PLAYER2_ID string

	NSQ_ADDR string

	NSQ_TOPIC_LOG_PLAYER1 string

	NSQ_TOPIC_LOG_PLAYER2 string

	NSQ_TOPIC_REFEREE_LOG string

	NSQ_TOPIC_GAME_RESULT string

	LOG_VOLUME_ADDR_PLAYER1 string

	LOG_VOLUME_ADDR_PLAYER2 string
)
