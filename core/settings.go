package core

const (
	TCP_PORT        = 8765
	TCP_SEND_ERRORS = true // command execution error
	TCP_SEND_RESULT = true // win/lose response to client

	SHOW_ZOMBIES_MS = 2000 // 2 seconds
	MOVE_ZOMBIES_MS = 1000 // 1 second
	SHOOT_SPEED_MS  = 1500 // 1.5 seconds, archers reload time

	BOARD_WIDTH  = 30
	BOARD_HEIGHT = 10
	ZOMBIE_COUNT = 10

	LOG_INFO        = true
	LOG_ERROR       = true
	LOG_TCP_SEND    = true
	LOG_TCP_RECEIVE = true
)
