package game

import (
	"context"
	"errors"
	"github.com/Diode222/GomokuGameReferee/client"
	"github.com/Diode222/GomokuGameReferee/conf"
	"github.com/Diode222/GomokuGameReferee/errorcode"
	pb "github.com/Diode222/GomokuGameReferee/probo"
	"github.com/Diode222/GomokuGameReferee/utils"
	"log"
	"time"
)

// winner
const (
	DRAW     = 0
	PLAYER1  = 1
	PLAYER2  = 2
	NOT_OVER = 3
)

const (
	BLANK = 0
	WHITE = 1
)

type Position struct {
	X int32
	Y int32
}

type Operation struct {
	Player   int
	Position *Position
	Type     int
}

var operations = []*Operation{}

type playerState struct {
	IsFirstHand bool
	PieceType int
}

var (
	PLAYER1_STATE *playerState
	PLAYER2_STATE *playerState
)

func StartGame() (int, int64, int64, []*Operation, error) {
	var err error

	// player can judge if stop move when timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(conf.MAX_THINKING_TIME))
	defer cancel()

	initPlayerStates()

	board, boardInfoMap := initBoard()

	player1ServiceClient := client.NewPlayerServiceClient(conf.PLAYER1_SERVICE_NAME)
	player2ServiceClient := client.NewPlayerServiceClient(conf.PLAYER2_SERVICE_NAME)

	err = waitPlayerToInit(ctx, player1ServiceClient, player2ServiceClient)
	if err != nil {
		t := time.Now().Unix()
		if err.Error() == errorcode.PLAYER1_TIMEOUT {
			return PLAYER2, t, t, operations, err
		} else {
			return PLAYER1, t, t, operations, err
		}
	}

	winner, startTime, endTime, err := playGame(ctx, player1ServiceClient, player2ServiceClient, board, boardInfoMap)
	return winner, startTime, endTime, operations, err
}

func initBoard() (*pb.Board, map[int]*pb.PiecePosition) {
	board := &pb.Board{
		ChessPositions: []*pb.PiecePosition{},
	}
	boardInfoMap := map[int]*pb.PiecePosition{}

	for x := 0; x < conf.BOARD_LENGTH; x += 1 {
		for y := 0; y < conf.BOARD_HEIGHT; y += 1 {
			tmpPiecePosition := &pb.PiecePosition{
				Type: pb.PieceType_NONE.Enum(),
				Position: &pb.Position{
					X: utils.GetPointerInt32(int32(x)),
					Y: utils.GetPointerInt32(int32(y)),
				},
			}
			board.ChessPositions = append(board.ChessPositions)
			boardInfoMap[y*conf.BOARD_LENGTH+x] = tmpPiecePosition
		}
	}

	return board, boardInfoMap
}

func initPlayerStates() {
	if conf.PLAYER1_FIRST_HAND == 1 {
		PLAYER1_STATE = &playerState{
			IsFirstHand: true,
			PieceType:   BLANK,
		}
		PLAYER2_STATE = &playerState{
			IsFirstHand: false,
			PieceType:   WHITE,
		}
	} else {
		PLAYER1_STATE = &playerState{
			IsFirstHand: false,
			PieceType:   WHITE,
		}
		PLAYER2_STATE = &playerState{
			IsFirstHand: true,
			PieceType:   BLANK,
		}
	}
}

func waitPlayerToInit(ctx context.Context, player1ServiceClient, player2ServiceClient pb.MakePieceServiceClient) error {
	var player1IsFirstHand bool
	var err error
	if conf.PLAYER1_FIRST_HAND == 1 {
		// player1 is first hand
		player1IsFirstHand = true
	} else {
		// player2 is first hand
		player1IsFirstHand = false
	}
	_, err = player1ServiceClient.Init(ctx, &pb.IsFirst{
		IsFirst: utils.GetPointerBool(player1IsFirstHand),
	})
	if err != nil {
		err = errors.New(errorcode.PLAYER1_TIMEOUT)
	}
	_, err = player2ServiceClient.Init(ctx, &pb.IsFirst{
		IsFirst: utils.GetPointerBool(!player1IsFirstHand),
	})
	if err != nil {
		err = errors.New(errorcode.PLAYER2_TIMEOUT)
	}
	return err
}

func playGame(ctx context.Context, player1ServiceClient, player2ServiceClient pb.MakePieceServiceClient, board *pb.Board, boardInfoMap map[int]*pb.PiecePosition) (int, int64, int64, error) {
	clients := []pb.MakePieceServiceClient{}
	if conf.PLAYER1_FIRST_HAND == 1 {
		clients[0] = player1ServiceClient
		clients[1] = player2ServiceClient
	} else {
		clients[0] = player2ServiceClient
		clients[1] = player1ServiceClient
	}

	startTime := time.Now().Unix()

	gameIsOver, winner := gameOver(boardInfoMap)
	var isPlayer1Current bool
	if conf.PLAYER1_FIRST_HAND == 1 {
		isPlayer1Current = true
	} else {
		isPlayer1Current = false
	}
	for !gameIsOver {
		err := makePiece(ctx, clients[0], board, boardInfoMap, isPlayer1Current)
		if err != nil {
			if isPlayer1Current {
				return PLAYER2, startTime, startTime, err
			} else {
				return PLAYER1, startTime, startTime, err
			}
		}
		gameIsOver, winner = gameOver(boardInfoMap)
		swapClient(clients)
		isPlayer1Current = !isPlayer1Current
	}

	endTime := time.Now().Unix()

	return winner, startTime, endTime, nil
}

func gameOver(boardInfoMap map[int]*pb.PiecePosition) (bool, int) {
	boardIsFull := true
	for i := 0; i < conf.BOARD_LENGTH; i += 1 {
		for j := 0; j < conf.BOARD_HEIGHT; j += 1 {
			if boardInfoMap[j*conf.BOARD_LENGTH+i].GetType() == pb.PieceType_NONE {
				boardIsFull = false
			}
			if isOver, winner := gameOverHelp(i, j, boardInfoMap); isOver {
				return isOver, winner
			}
		}
	}

	if boardIsFull {
		return true, DRAW
	}

	return false, NOT_OVER
}

func gameOverHelp(positionX, positionY int, boardInfoMap map[int]*pb.PiecePosition) (bool, int) {
	chessType := boardInfoMap[positionY*conf.BOARD_LENGTH+positionX].GetType()
	if chessType == pb.PieceType_NONE {
		return false, NOT_OVER
	}
	i := -1

	// right
	if positionX <= conf.BOARD_LENGTH-5 {
		for i = positionX; i < positionX+5; i++ {
			if boardInfoMap[positionY*conf.BOARD_LENGTH+i].GetType() != chessType {
				break
			}
		}
		if i == conf.BOARD_LENGTH {
			if chessType == pb.PieceType_BLANK {
				return getGameOverInfo(chessType)
			}
		}
	}

	// left
	if positionX >= 4 {
		for i = positionX; i > positionX-5; i-- {
			if boardInfoMap[positionY*conf.BOARD_LENGTH+i].GetType() != chessType {
				break
			}
		}
		if i == -1 {
			if chessType == pb.PieceType_WHITE {
				return getGameOverInfo(chessType)
			}
		}
	}

	// up
	if positionY <= conf.BOARD_HEIGHT-5 {
		for i = positionY; i < positionY+5; i++ {
			if boardInfoMap[i*conf.BOARD_LENGTH+positionX].GetType() != chessType {
				break
			}
		}
		if i == conf.BOARD_HEIGHT {
			return getGameOverInfo(chessType)
		}
	}

	// bottom
	if positionY >= 4 {
		for i = positionY; i > positionY-5; i-- {
			if boardInfoMap[i*conf.BOARD_LENGTH+positionX].GetType() != chessType {
				break
			}
		}
		if i == -1 {
			return getGameOverInfo(chessType)
		}
	}

	// right_up
	if positionX <= conf.BOARD_LENGTH-5 && positionY <= conf.BOARD_HEIGHT {
		for i = 0; i < 5; i++ {
			tmpX := positionX + i
			tmpY := positionY + i
			if boardInfoMap[tmpY*conf.BOARD_LENGTH+tmpX].GetType() != chessType {
				break
			}
		}
		if i == 5 {
			return getGameOverInfo(chessType)
		}
	}

	// right_bottom
	if positionX <= conf.BOARD_LENGTH-5 && positionY >= 4 {
		for i = 0; i < 5; i++ {
			tmpX := positionX + i
			tmpY := positionY - i
			if boardInfoMap[tmpY*conf.BOARD_LENGTH+tmpX].GetType() != chessType {
				break
			}
		}
		if i == 5 {
			return getGameOverInfo(chessType)
		}
	}

	// left_up
	if positionX >= 4 && positionY <= conf.BOARD_LENGTH-5 {
		for i = 0; i < 5; i++ {
			tmpX := positionX - i
			tmpY := positionY + i
			if boardInfoMap[tmpY*conf.BOARD_LENGTH+tmpX].GetType() != chessType {
				break
			}
		}
		if i == 5 {
			return getGameOverInfo(chessType)
		}
	}

	// left_bottom
	if positionX >= 4 && positionY >= 4 {
		for i = 0; i < 5; i++ {
			tmpX := positionX - i
			tmpY := positionY - i
			if boardInfoMap[tmpY*conf.BOARD_LENGTH+tmpX].GetType() != chessType {
				break
			}
		}
		if i == 5 {
			return getGameOverInfo(chessType)
		}
	}

	return getGameOverInfo(pb.PieceType_NONE)
}

func getGameOverInfo(chessType pb.PieceType) (bool, int) {
	if chessType == pb.PieceType_BLANK {
		if conf.PLAYER1_FIRST_HAND == 1 {
			return true, PLAYER1
		} else {
			return true, PLAYER2
		}
	} else if chessType == pb.PieceType_WHITE {
		if conf.PLAYER1_FIRST_HAND == 1 {
			return true, PLAYER2
		} else {
			return true, PLAYER1
		}
	} else {
		return false, NOT_OVER
	}
}

func makePiece(ctx context.Context, playerClient pb.MakePieceServiceClient, board *pb.Board, boardInfoMap map[int]*pb.PiecePosition, isPlayer1Current bool) error {
	piecePosition, err := playerClient.MakePiece(ctx, board)
	if err != nil {
		log.Println(err.Error())
		if isPlayer1Current {
			return errors.New(errorcode.PLAYER2_TIMEOUT)
		} else {
			return errors.New(errorcode.PLAYER1_TIMEOUT)
		}
	}

	// Add make piece operation to operations first, to know which piece operation is wrong
	var currentPlayer int
	var operationType int
	if isPlayer1Current {
		currentPlayer = PLAYER1
		operationType = PLAYER1_STATE.PieceType
	} else {
		currentPlayer = PLAYER2
		operationType = PLAYER2_STATE.PieceType
	}
	operations = append(operations, &Operation{
		Player:   currentPlayer,
		Position: &Position{X: piecePosition.GetPosition().GetX(), Y: piecePosition.GetPosition().GetY()},
		Type:     operationType,
	})

	if piecePosition.GetType() != pb.PieceType_BLANK || piecePosition.GetType() != pb.PieceType_WHITE {
		if currentPlayer == PLAYER1 {
			return errors.New(errorcode.PLAYER1_WRONG_OPERATION)
		} else {
			return errors.New(errorcode.PLAYER2_WRONG_OPERATION)
		}
	}

	if piecePosition.GetPosition() == nil {
		if currentPlayer == PLAYER1 {
			return errors.New(errorcode.PLAYER1_WRONG_OPERATION)
		} else {
			return errors.New(errorcode.PLAYER2_WRONG_OPERATION)
		}
	}

	if int(piecePosition.GetPosition().GetX()) > conf.BOARD_LENGTH || int(piecePosition.GetPosition().GetX()) < 0 ||
		int(piecePosition.GetPosition().GetY()) > conf.BOARD_HEIGHT || int(piecePosition.GetPosition().GetY()) < 0 {
		if currentPlayer == PLAYER1 {
			return errors.New(errorcode.PLAYER1_WRONG_OPERATION)
		} else {
			return errors.New(errorcode.PLAYER2_WRONG_OPERATION)
		}
	}

	var pieceType int
	if piecePosition.GetType() == pb.PieceType_BLANK {
		pieceType = BLANK
	} else {
		pieceType = WHITE
	}
	if currentPlayer == PLAYER1 && pieceType != PLAYER1_STATE.PieceType {
		return errors.New(errorcode.PLAYER1_WRONG_OPERATION)
	}
	if currentPlayer == PLAYER2 && pieceType != PLAYER2_STATE.PieceType {
		return errors.New(errorcode.PLAYER2_WRONG_OPERATION)
	}

	piecePositionOld := boardInfoMap[int(piecePosition.GetPosition().GetY())*conf.BOARD_LENGTH+int(piecePosition.GetPosition().GetX())]

	if piecePositionOld.GetType() != pb.PieceType_NONE {
		if currentPlayer == PLAYER1 {
			return errors.New(errorcode.PLAYER1_WRONG_OPERATION)
		} else {
			return errors.New(errorcode.PLAYER2_WRONG_OPERATION)
		}
	}

	piecePositionOld.Position.X = utils.GetPointerInt32(piecePosition.GetPosition().GetX())
	piecePositionOld.Position.Y = utils.GetPointerInt32(piecePosition.GetPosition().GetY())

	return nil
}

func swapClient(clients []pb.MakePieceServiceClient) {
	clients[0], clients[1] = clients[1], clients[0]
}
