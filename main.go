package main

import (
	"errors"
	"fmt"
	"github.com/Diode222/GomokuGameReferee/conf"
	"github.com/Diode222/GomokuGameReferee/feedbackData"
	"github.com/Diode222/GomokuGameReferee/game"
	"log"
	"os"
	"strconv"
)

func main() {
	var err error
	conf.ETCD_ADDR = os.Getenv("ETCD_ADDR")
	conf.PLAYER1_PORT = os.Getenv("PLAYER1_PORT")
	conf.PLAYER2_PORT = os.Getenv("PLAYER2_PORT")
	conf.PLAYER1_FIRST_HAND = os.Getenv("PLAYER1_FIRST_HAND")
	conf.MAX_THINKING_TIME, err = strconv.Atoi(os.Getenv("MAX_THINKING_TIME"))
	if err != nil {
		log.Println(fmt.Sprintf("MAX_THINKING_TIME param wrong, MAX_THINKING_TIME: %s", os.Getenv("MAX_THINKING_TIME")))
		feedbackData.TransError(errors.New(fmt.Sprintf("MAX_THINKING_TIME param wrong, MAX_THINKING_TIME: %s", os.Getenv("MAX_THINKING_TIME"))))
		return
	}
	conf.GAME_ID = os.Getenv("GAME_ID")

	//conf.PLAYER1_PORT = "10001"
	//conf.PLAYER2_PORT = "10002"
	//conf.PLAYER1_FIRST_HAND = "true"
	//conf.MAX_THINKING_TIME = 5

	winner, startTime, endTime, operations, err := game.StartGame()

	feedbackData.TransData(winner, startTime, endTime, operations, err)
}
