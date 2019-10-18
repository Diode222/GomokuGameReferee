package feedbackData

import (
	"github.com/Diode222/GomokuGameReferee/game"
	"log"
	"strconv"
)

func TransData(winner int, startTime int64, endTime int64, operations []*game.Operation, err error) {
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
		if operation.Type == game.BLANK {
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
