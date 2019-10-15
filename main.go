package main

import (
	"fmt"
	"github.com/Diode222/GomokuGameReferee/conf"
	"github.com/Diode222/GomokuGameReferee/game"
	"github.com/Diode222/GomokuGameReferee/logging"
	"log"
	"os"
	"strconv"
)

func main() {
	var err error
	conf.PLAYER1_SERVICE_NAME = os.Getenv("PLAYER1_SERVICE_NAME")
	conf.PLAYER2_SERVICE_NAME = os.Getenv("PLAYER2_SERVICE_NAME")
	conf.PLAYER1_FIRST_HAND, err = strconv.Atoi(os.Getenv("PLAYER1_FIRST_HAND"))
	if err != nil {
		log.Println(fmt.Sprintf("PLAYER1_FIRST_HAND param wrong, PLAYER1_FIRST_HAND: %s", os.Getenv("PLAYER1_FIRST_HAND")))
		return
	}
	conf.MAX_THINKING_TIME, err = strconv.Atoi(os.Getenv("MAX_THINKING_TIME"))
	if err != nil {
		log.Println(fmt.Sprintf("MAX_THINKING_TIME param wrong, MAX_THINKING_TIME: %s", os.Getenv("MAX_THINKING_TIME")))
		return
	}

	game.StartGame()

	logging.TransLog()
}
