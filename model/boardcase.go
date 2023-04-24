package model

import "github.com/mohae/deepcopy"

// 准备用来发送给player1和player2d的信息

type BoardCase struct {
	Curplayer uint    `json:"curPlayer"`
	Winner    int     `json:"winner"`
	Chesses   [][]int `json:"chesses"`
}

func (board *ChessBoard) GetBoardCase(playerID uint) *BoardCase {
	if playerID != board.Player1ID && playerID != board.Player2ID {
		return nil
	}
	winner, isWinner := board.CheckWin()

	out := &BoardCase{
		Curplayer: board.Curplayer(),
		Chesses:   deepcopy.Copy(board.chesses).([][]int),
	}
	if !isWinner {
		out.Winner = -1
	} else if int(winner) == int(playerID) {
		out.Winner = 1
	} else {
		out.Winner = 2
	}
	if out.Curplayer == playerID {
		out.Curplayer = 1
	} else {
		out.Curplayer = 2
	}
	if playerID == board.Player2ID {
		for i := 0; i < len(out.Chesses); i++ {
			for j := 0; j < len(out.Chesses[i])/2; j++ {
				tmp := out.Chesses[i][j]
				out.Chesses[i][j] = out.Chesses[i][len(out.Chesses[i])-1-j]
				out.Chesses[i][len(out.Chesses[i])-1-j] = tmp
			}
		}
	}
	for i := 0; i < len(out.Chesses); i++ {
		for j := 0; j < len(out.Chesses[i]); j++ {
			if out.Chesses[i][j] < 0 {
				out.Chesses[i][j] = 0
			} else if out.Chesses[i][j] == int(playerID) {
				out.Chesses[i][j] = 1
			} else {
				out.Chesses[i][j] = 2
			}
		}
	}

	return out
}
