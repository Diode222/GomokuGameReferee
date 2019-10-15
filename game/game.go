package game

import (
	"context"
	"github.com/Diode222/GomokuGameReferee/client"
	"github.com/Diode222/GomokuGameReferee/conf"
	pb "github.com/Diode222/GomokuGameReferee/probo"
	"github.com/Diode222/GomokuGameReferee/utils"
	"time"
)

func StartGame() {
	// player can judge if stop move when timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * time.Duration(conf.MAX_THINKING_TIME))
	defer cancel()

	board, boardInfoMap := initBoard()

	player1FirstHand := conf.PLAYER1_FIRST_HAND

	play1ServiceClient := client.NewPlayerServiceClient(conf.PLAYER1_SERVICE_NAME)
	play2ServiceClient := client.NewPlayerServiceClient(conf.PLAYER2_SERVICE_NAME)

	waitPlayerToInit(ctx, play1ServiceClient, play2ServiceClient, player1FirstHand)

	playGame(ctx, play1ServiceClient, play2ServiceClient, player1FirstHand)
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
			board.ChessPositions = append(board.ChessPositions, )
			boardInfoMap[y * conf.BOARD_LENGTH + x] = tmpPiecePosition
		}
	}

	return board, boardInfoMap
}

func waitPlayerToInit(ctx context.Context, play1ServiceClient, play2ServiceClient pb.MakePieceServiceClient, player1FirstHand int) {
	if player1FirstHand == 1 {
		// player1 is first hand
		play1ServiceClient.Init(ctx, &pb.IsFirst{
			IsFirst:              utils.GetPointerBool(true),
		})
		play2ServiceClient.Init(ctx, &pb.IsFirst{
			IsFirst:              utils.GetPointerBool(false),
		})
	} else {
		// player2 is first hand
		play1ServiceClient.Init(ctx, &pb.IsFirst{
			IsFirst:              utils.GetPointerBool(false),
		})
		play2ServiceClient.Init(ctx, &pb.IsFirst{
			IsFirst:              utils.GetPointerBool(true),
		})
	}
}

func playGame(ctx context.Context, play1ServiceClient, play2ServiceClient pb.MakePieceServiceClient, player1FirstHand int, board pb.Board, boardInfoMap map[int]*pb.PiecePosition) {
	clients := []pb.MakePieceServiceClient{}
	if player1FirstHand == 1 {
		clients[0] = play1ServiceClient
		clients[1] = play2ServiceClient
	} else {
		clients[0] = play2ServiceClient
		clients[1] = play1ServiceClient
	}

	for ;!gameOver(board); {

	}
}

// TODO
func gameOver(board pb.Board) bool {
	return true
}
