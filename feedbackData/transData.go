package feedbackData

import (
	"bufio"
	"encoding/json"
	"github.com/Diode222/GomokuGameReferee/client"
	"github.com/Diode222/GomokuGameReferee/conf"
	"github.com/Diode222/GomokuGameReferee/errorcode"
	"github.com/Diode222/GomokuGameReferee/model"
	"github.com/Diode222/GomokuGameReferee/utils"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"strconv"
)

func TransData1(winner int, startTime int64, endTime int64, operations []*model.Operation, err error) {
	log.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	log.Println(winner)
	log.Println(">>>>>>>>>>>>>>winner>>>>>>>>>>")
	log.Println(startTime)
	log.Println(">>>>>>>>>>>>>startTime>>>>>>>>")
	log.Println(endTime)
	log.Println(">>>>>>>>>>>>>endTime>>>>>>>>")
	log.Println()
	log.Println()
	for _, operation := range operations {
		var operationType string
		if operation.Type == model.BLANK {
			operationType = "blank"
		} else {
			operationType = "white"
		}
		log.Println("(" + strconv.Itoa(int(operation.Position.X)) + ", " + strconv.Itoa(int(operation.Position.Y)) + "):  " + operationType)
	}
	log.Println()
	log.Println()
	log.Println(">>>>>>>>>>>>>>operation>>>>>>>>>>>>>")
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("No error")
	}
	log.Println(">>>>>>>>>>>>>>err>>>>>>>>>>>>>>>>")
}

func TransData(winner int, startTime int64, endTime int64, operations []*model.Operation, err error) {
	var player1FirstHand bool
	if conf.PLAYER1_FIRST_HAND == "true" {
		player1FirstHand = true
	} else {
		player1FirstHand = false
	}
	var foulPlayer int
	if err == nil {
		foulPlayer = model.NO_FOUL
	} else {
		if utils.StringContains(err.Error(), []string{errorcode.PLAYER1_TIMEOUT, errorcode.PLAYER1_WRONG_OPERATION}) {
			foulPlayer = model.PLAYER1_FOUL
		} else {
			foulPlayer = model.PLAYER2_FOUL
		}
	}
	matchData := &model.MatchDataModel{
		GameID:           conf.GAME_ID,
		BoardLength:      conf.BOARD_LENGTH,
		BoardHeight:      conf.BOARD_HEIGHT,
		Player1ID:        conf.PLAYER1_ID,
		Player2ID:        conf.PLAYER2_ID,
		Player1FirstHand: player1FirstHand,
		MaxThingTime:     conf.MAX_THINKING_TIME,
		Winner:           winner,
		StartTime:        startTime,
		EndTime:          endTime,
		Operations:       operations,
		FoulPlayer:       foulPlayer,
	}

	binaryMatchData, err := json.Marshal(matchData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("marshal binaryMatchData failed.")
		TransServerError(err)
	}

	nsqProducer, err := client.NewNsqClient()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err.Error(),
		}).Fatal("get nsqProducer failed.")
		TransServerError(err)
	}

	transGameResultData(binaryMatchData, nsqProducer)
	transLogData(nsqProducer)
}

func transGameResultData(data []byte, nsqProducer *nsq.Producer) {
	err := nsqProducer.Publish(conf.NSQ_TOPIC_GAME_RESULT, data)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"topic":    conf.NSQ_TOPIC_GAME_RESULT,
			"nsq_addr": conf.NSQ_ADDR,
			"err":      err.Error(),
		}).Fatal("nsq publish failed.")
		TransServerError(err)
	}
}

func transLogData(nsqProducer *nsq.Producer) {
	transPlayer1LogData(nsqProducer)
	transPlayer2LogData(nsqProducer)
	transRefereeLogData(nsqProducer)
}

func transPlayer1LogData(nsqProducer *nsq.Producer) {
	logFile, err := os.Open(conf.LOG_VOLUME_ADDR_PLAYER1)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"player1_log_file_path": conf.LOG_VOLUME_ADDR_PLAYER1,
			"err":                   err.Error(),
		}).Fatal("read player1 log file failed.")
		TransServerError(err)
	}
	defer logFile.Close()
	buf := bufio.NewReader(logFile)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				for nsqProducer.Publish(conf.NSQ_TOPIC_LOG_PLAYER1, line) != nil {
				}
				break
			}
			continue
		}
	}
}

func transPlayer2LogData(nsqProducer *nsq.Producer) {
	logFile, err := os.Open(conf.LOG_VOLUME_ADDR_PLAYER2)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"player2_log_file_path": conf.LOG_VOLUME_ADDR_PLAYER2,
			"err":                   err.Error(),
		}).Fatal("read player2 log file failed.")
		TransServerError(err)
	}
	defer logFile.Close()
	buf := bufio.NewReader(logFile)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				for nsqProducer.Publish(conf.NSQ_TOPIC_LOG_PLAYER2, line) != nil {
				}
				break
			}
			continue
		}
	}
}

func transRefereeLogData(nsqProducer *nsq.Producer) {
	logFile, err := os.Open(conf.REFEREE_LOG_FILE)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"referee_log_file_path": conf.REFEREE_LOG_FILE,
			"err":                   err.Error(),
		}).Fatal("read referee log file failed.")
		TransServerError(err)
	}
	defer logFile.Close()
	buf := bufio.NewReader(logFile)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err == io.EOF {
				for nsqProducer.Publish(conf.NSQ_TOPIC_REFEREE_LOG, line) != nil {
				}
				break
			}
			continue
		}
	}
}
