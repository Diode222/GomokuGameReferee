package conf

const (
	ETCD_ADDR = "127.0.0.1:2379"

	// BACKEND_SERVICE_NAME
	BACKEND_SERVICE_NAME = "backend_service"

	// LOG_PATH_NAME
	LOG_PATH_NAME = "game_log"
)

var (
	// player1 SERVICE_NAME
	PLAYER1_SERVICE_NAME string

	// player2 SERVICE_NAME
	PLAYER2_SERVICE_NAME string

	PLAYER1_FIRST_HAND int // 1: true, 0: false

	BOARD_LENGTH = 15
	BOARD_HEIGHT = 15

	// MAX_THINKING_TIME
	MAX_THINKING_TIME int

	GAME_ID int
)
