package model

import (
	"cncsmonster/gomoku/config"
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

// x up-down,y left-right

type Step struct {
	X,
	Y uint
}

type ChessBoard struct {
	gorm.Model
	Xsize         uint   `gorm:"column:x_size;getter:GetXSize;setter:SetXSize"`
	YSize         uint   `gorm:"column:y_size;getter:GetYSize;setter:SetYSize"`
	Player1ID     uint   `gorm:"column:player1_id;getter:GetPlayer1ID;setter:SetPlayer1ID"`
	Player2ID     uint   `gorm:"column:player2_id;getter:GetPlayer2ID;setter:SetPlayer2ID"`
	Player1Color  string `gorm:"column:player1_color;getter:GetPlayer1Color;setter:SetPlayer1Color"`
	Player2Color  string `gorm:"column:player2_color;getter:GetPlayer2Color;setter:SetPlayer2Color"`
	BoardColor    string `gorm:"column:board_color;getter:GetBoardColor;setter:SetBoardColor"`
	ChessesString string `gorm:"column:chesses"`
	StepsString   string `gorm:"column:steps"`
	chesses       [][]int
	steps         []Step
}

const (
	player1Color string = "black"
	player2Color string = "white"
	boardColor   string = "yellow"
)

var chessboardDB *gorm.DB

func init() {
	config.Connect()
	chessboardDB = config.GetDB()
	chessboardDB.AutoMigrate(&ChessBoard{})
}

// 初始化棋盘
func (board *ChessBoard) Init() {
	board.Player1Color = player1Color
	board.Player2Color = player2Color
	board.BoardColor = boardColor
	board.steps = make([]Step, 0)
	for i := 0; i < int(board.Xsize); i++ {
		for j := 0; j < int(board.YSize); j++ {
			board.chesses[i][j] = -1
		}
	}
}

// 获取当前下棋人，开始的时候player1下棋
func (board ChessBoard) Curplayer() uint {
	steps := len(board.steps)
	if steps%2 == 0 {
		return board.Player1ID
	} else {
		return board.Player2ID
	}
}

func NewChessBoard(xSize, ySize, player1ID, player2ID uint) ChessBoard {
	var board ChessBoard = ChessBoard{
		Xsize:        xSize,
		YSize:        ySize,
		Player1ID:    player1ID,
		Player2ID:    player2ID,
		Player1Color: player1Color,
		Player2Color: player2Color,
		BoardColor:   boardColor,
		steps:        make([]Step, 0),
	}
	board.chesses = make([][]int, xSize)
	for i := 0; i < int(xSize); i++ {
		board.chesses[i] = make([]int, ySize)
		for j := 0; j < int(ySize); j++ {
			board.chesses[i][j] = -1
		}
	}
	return board
}

func (board *ChessBoard) ChessAt(x, y uint) (int, error) {
	if x >= board.Xsize || y >= board.YSize {
		return -1, errors.New("Wrong")
	}
	return board.chesses[x][y], nil
}

func (board *ChessBoard) Play(playerId, x, y uint) (e error) {
	// 如果棋盘上某个地方已经有棋子了
	// 首先判断输入是否合法
	if x >= board.Xsize || y >= board.YSize {
		return errors.New("chess index out")
	}
	if board.chesses[x][y] >= 0 {
		return errors.New("unplayable position")
	}
	board.chesses[x][y] = int(board.Curplayer())
	board.steps = append(board.steps, Step{x, y})
	return nil
}

func (board *ChessBoard) Playable(playerId, x, y uint) bool {
	if x >= board.Xsize || y >= board.YSize || board.chesses[x][y] >= 0 {
		return false
	}
	if board.Curplayer() != playerId {
		return false
	}
	return true
}

func (board *ChessBoard) CheckWin() (winnerId uint, isWin bool) {
	// defer func() {
	// 	if p := recover(); p != nil {
	// 		var a int
	// 		fmt.Println(a)
	// 	}
	// }()
	// 使用三个三重循环，对每个行，列还有斜进行检测是否有连续五个
	//达成胜利条件需要的连续的棋子数量
	winCoe := 5
	isWin = false
	// 使用int类型的chan来返回判断的结果
	ch := make(chan int)
	stopChan := make(chan bool)
	// fmt.Println(winCoe)
	// TODO,使用协程并行计算三个循环

	// 处理连续棋子之间的状态跳转
	processState := func(lastChess, curChess, continuousNum int) (newLastChess, newContinuousNum int) {
		if curChess < 0 {
			newContinuousNum = 0
		} else if lastChess >= 0 && curChess != lastChess {
			newContinuousNum = 1
		} else if lastChess < 0 {
			newContinuousNum = 1
		} else {
			//否则就是上一个棋子与这个棋子角色相同
			newContinuousNum = continuousNum + 1
		}
		return curChess, newContinuousNum
	}

	go func() {
		// 先对n个对角线进行检查
		for i := uint(0); i < board.Xsize; i++ {
			// 使用管道判断是否关闭
			//往右下检查
			lastChess := board.chesses[i][0]
			continuousNum := 1
			if lastChess < 0 {
				continuousNum = 0
			}
			for j := uint(1); j < board.YSize && i+j < board.Xsize; j++ {
				// fmt.Println(i, j)
				select {
				case <-stopChan:
					return
				default:
				}
				curChess := board.chesses[i+j][j]
				lastChess, continuousNum = processState(lastChess, curChess, continuousNum)
				if continuousNum == winCoe {
					// 如果达成胜利条件
					ch <- curChess
				}
			}
			// 往右上检查
			lastChess = board.chesses[i][0]
			continuousNum = 1
			if lastChess < 0 {
				continuousNum = 0
			}
			for j := uint(1); j < board.YSize && int(i)-int(j) >= 0; j++ {
				select {
				case <-stopChan:
					return
				default:
				}
				curChess := board.chesses[i-j][j]
				lastChess, continuousNum = processState(lastChess, curChess, continuousNum)
				if continuousNum == winCoe {
					// 如果达成胜利条件
					ch <- curChess
					return
				}
			}
		}
		ch <- (-1)
	}()
	// tmp1 := <-ch
	go func() {
		// 对n行进行检查
		for i := uint(0); i < board.Xsize; i++ {
			lastChess := board.chesses[i][0]
			continuousNum := 1
			if lastChess < 0 {
				continuousNum = 0
			}
			for j := uint(1); j < board.YSize; j++ {
				select {
				case <-stopChan:
					return
				default:
				}
				curChess := board.chesses[i][j]
				lastChess, continuousNum = processState(lastChess, curChess, continuousNum)
				if continuousNum == winCoe {
					// 如果达成胜利条件
					ch <- curChess
					return
				}
			}
		}
		ch <- (-1)
	}()
	// tmp1 = <-ch
	go func() {
		// defer func() {
		// 	if p := recover(); p != nil {
		// 		var a int
		// 		fmt.Println(a)
		// 	}
		// }()
		// 对n列进行检查
		for i := uint(0); i < board.YSize; i++ {
			lastChess := board.chesses[0][i]
			continuousNum := 1
			if lastChess < 0 {
				continuousNum = 0
			}
			for j := uint(1); j < board.Xsize; j++ {
				select {
				case <-stopChan:
					return
				default:
				}
				curChess := board.chesses[j][i]
				lastChess, continuousNum = processState(lastChess, curChess, continuousNum)
				if continuousNum == winCoe {
					// 如果达成胜利条件
					ch <- curChess
					return
				}
			}
		}
		ch <- (-1)
	}()
	// tmp1 = <-ch
	// fmt.Println(tmp1)
	// 等待管道中所有信息处理完
	var times = 0
	for out := range ch {
		// fmt.Println(out, times)
		if out < 0 {
			times++
			if times == 3 {
				break
			}
			continue
		}
		close(stopChan)
		return uint(out), true
	}
	return 0, false
}

func (board *ChessBoard) Copy() {
	// TODO
	panic("not implement")
}

func (board *ChessBoard) Pack() {
	if body, err := json.Marshal(board.chesses); err == nil {
		board.ChessesString = string(body)
	}
	if body, err := json.Marshal(board.steps); err == nil {
		board.StepsString = string(body)
	}
}

func (board *ChessBoard) UnPack() {
	body := []byte(board.ChessesString)
	json.Unmarshal(body, &board.chesses)
	body = []byte(board.StepsString)
	json.Unmarshal(body, &board.steps)
}

// 用来保存到数据库中
func (board *ChessBoard) Save() {
	if board == nil {
		return
	}
	board.Pack()
	chessboardDB.Save(board)
}

/*
准备一些getter和setter用来辅助orm
*/

func (c *ChessBoard) SetXSize(size uint) {
	c.Xsize = size
}

func (c *ChessBoard) GetXSize() uint {
	return c.Xsize
}

func (c *ChessBoard) SetYSize(size uint) {
	c.YSize = size
}

func (c *ChessBoard) GetYSize() uint {
	return c.YSize
}

func (c *ChessBoard) SetPlayer1ID(id uint) {
	c.Player1ID = id
}

func (c *ChessBoard) GetPlayer1ID() uint {
	return c.Player1ID
}

func (c *ChessBoard) SetPlayer2ID(id uint) {
	c.Player2ID = id
}

func (c *ChessBoard) GetPlayer2ID() uint {
	return c.Player2ID
}

func (c *ChessBoard) SetPlayer1Color(color string) {
	c.Player1Color = color
}

func (c *ChessBoard) GetPlayer1Color() string {
	return c.Player1Color
}

func (c *ChessBoard) SetPlayer2Color(color string) {
	c.Player2Color = color
}

func (c *ChessBoard) GetPlayer2Color() string {
	return c.Player2Color
}

func (c *ChessBoard) SetBoardColor(color string) {
	c.BoardColor = color
}

func (c *ChessBoard) GetBoardColor() string {
	return c.BoardColor
}

// chesses
func (c *ChessBoard) SetChesses(chesses [][]int) {
	c.chesses = chesses
}
func (c *ChessBoard) GetChesses() [][]int {
	return c.chesses
}

// steps
func (c *ChessBoard) SetSteps(steps []Step) {
	c.steps = steps
}

func (c *ChessBoard) GetSteps() []Step {
	return c.steps
}
