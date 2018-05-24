package core

const (
	TCP_PORT        = 8765
	TCP_TIMEOUT_SEC = 60   // wait for command up to one minute
	TCP_SEND_ERRORS = true // command execution error
	TCP_SEND_RESULT = true // win/lose response to client

	SHOW_ZOMBIES_MS = 2000 // 2 seconds
	MOVE_ZOMBIES_MS = 1000 // 1 second

	BOARD_WIDTH  = 10
	BOARD_HEIGHT = 4
	ZOMBIE_COUNT = 2

	LOG_INFO        = true
	LOG_ERROR       = true
	LOG_TCP_SEND    = true
	LOG_TCP_RECEIVE = true
)
