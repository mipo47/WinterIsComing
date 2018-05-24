package core

const (
	TCP_PORT        = 8765
	TCP_TIMEOUT_SEC = 5    // wait for command up to 5 minutes
	TCP_SEND_ERRORS = true // command execution error
	TCP_SEND_RESULT = true // win/lose response to client

	SHOW_ZOMBIES_MS = 2000 // 2 seconds
	MOVE_ZOMBIES_MS = 1000 // 1 second

	BOARD_WIDTH  = 20
	BOARD_HEIGHT = 2
	ZOMBIE_COUNT = 2

	LOG_INFO        = true
	LOG_ERROR       = true
	LOG_TCP_SEND    = true
	LOG_TCP_RECEIVE = true
)
