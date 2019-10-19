package main

import (
	"errors"
	"fmt"
	"github.com/Diode222/GomokuGameReferee/conf"
	"github.com/Diode222/GomokuGameReferee/feedbackData"
	"github.com/Diode222/GomokuGameReferee/game"
	"github.com/Diode222/logS"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

func main() {
	initLogHook(conf.REFEREE_LOG_FILE)

	var err error
	conf.PLAYER1_FIRST_HAND = os.Getenv("PLAYER1_FIRST_HAND")
	conf.MAX_THINKING_TIME, err = strconv.Atoi(os.Getenv("MAX_THINKING_TIME"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"MAX_THINKING_TIME": os.Getenv("MAX_THINKING_TIME"),
		}).Fatal("MAX_THINKING_TIME param wrong.")
		feedbackData.TransError(errors.New(fmt.Sprintf("MAX_THINKING_TIME param wrong, MAX_THINKING_TIME: %s", os.Getenv("MAX_THINKING_TIME"))))
		return
	}
	conf.GAME_ID = os.Getenv("GAME_ID")
	conf.NSQ_ADDR = os.Getenv("NSQ_ADDR")
	conf.NSQ_TOPIC_GAME_RESULT = conf.GAME_ID + "_game_result"
	conf.NSQ_TOPIC_LOG = conf.GAME_ID + "_log"
	conf.LOG_VOLUME_ADDR_PLAYER1 = os.Getenv("LOG_VOLUME_ADDR_PLAYER1")
	conf.LOG_VOLUME_ADDR_PLAYER2 = os.Getenv("LOG_VOLUME_ADDR_PLAYER2")

	//conf.PLAYER1_FIRST_HAND = "true"
	//conf.MAX_THINKING_TIME = 5

	winner, startTime, endTime, operations, err := game.StartGame()

	feedbackData.TransData(winner, startTime, endTime, operations, err)
}

func initLogHook(logFileAddr string) {
	hook := logS.NewHook(logFileAddr)
	if hook == nil {
		logrus.WithFields(logrus.Fields{
			"logFileAddr": logFileAddr,
		}).Fatal("File log hook created failed.")
		panic(errors.New("File log hook created failed."))
	}
	logrus.AddHook(hook)
}
