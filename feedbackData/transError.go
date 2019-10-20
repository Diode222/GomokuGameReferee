package feedbackData

import (
	"encoding/json"
	"github.com/Diode222/GomokuGameReferee/client"
	"github.com/Diode222/GomokuGameReferee/conf"
	"github.com/Diode222/GomokuGameReferee/model"
	"time"
)

// TODO
// Should stop running when transmission finished, I use select to stop.
func TransServerError(err error) {
	var player1FirstHand bool
	if conf.PLAYER1_FIRST_HAND == "true" {
		player1FirstHand = true
	} else {
		player1FirstHand = false
	}

	matchErrData := &model.MatchDataModel{
		GameID:           conf.GAME_ID,
		BoardLength:      conf.BOARD_LENGTH,
		BoardHeight:      conf.BOARD_HEIGHT,
		Player1ID:        conf.PLAYER1_ID,
		Player2ID:        conf.PLAYER2_ID,
		Player1FirstHand: player1FirstHand,
		MaxThingTime:     conf.MAX_THINKING_TIME,
		Winner:           -1,
		StartTime:        time.Now().Unix(),
		EndTime:          time.Now().Unix(),
		Operations:       []*model.Operation{},
		FoulPlayer:       0,
		ServerError:      true,
	}

	binaryData, err1 := json.Marshal(matchErrData)
	if err1 != nil {
		select {}
	}

	producer, err1 := client.NewNsqClient()
	if err1 != nil {
		select {}
	}

	i := 0
	for ;producer.Publish(conf.NSQ_TOPIC_GAME_RESULT, binaryData) != nil && i < 10; i++ {}

	select {}
}
