package model

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
)

// 缓存用到的列表
var games map[uint]*ChessBoard
var gameslock = sync.RWMutex{}

const max_games_cache_num = 10000

func init() {
	games = make(map[uint]*ChessBoard)
	go func() {
		level := 0
		// 执行一个协程，没过一段时间检查下数据量
		for {
			time.Sleep(30 * time.Minute)
			var num int64
			chessboardDB.Unscoped().Count(&num)
			var num2 int64
			accountDB.Unscoped().Count(&num2)
			num += num2
			if level == 0 {
				if num > 10000 {
					level += 1
					fmt.Println("level 1:", num)
				}
			} else if level == 1 {
				if num > 100000 {
					level += 1
					fmt.Println("level 2:", num)
				} else if num < 10000 {
					level -= 1
				}
			} else if level == 2 {
				if num > 1000000 {
					level += 1
					fmt.Println("level 3:", num)
				} else if num < 100000 {
					level -= 1
				}
			} else {
				panic("to many record! num:" + strconv.FormatInt(num, 10))
			}
		}
	}()
}

// 把所有内容写入数据库
func SaveAll() {
	// 所有当前游戏写入数据库
	gameslock.Lock()
	for _, v := range games {
		v.Save()
	}
	gameslock.Unlock()
}

// 去掉c个最新更新时间最晚的缓存
func shrinkGames(c, n uint) {
	//
	arr := make([]*ChessBoard, 0, n)
	for _, v := range games {
		arr = append(arr, v)
	}
	// 对arr进行排序
	sort.Slice(arr, func(i, j int) bool {
		return arr[i].UpdatedAt.Before(arr[j].UpdatedAt)
	})
	// 取出arr中更新时间最晚的m个元素,写入数据库，并且删除
	gameslock.Lock()
	for i := 0; i < int(c) && i < len(arr); i++ {
		delete(games, arr[i].ID)
		arr[i].Save()
	}
	gameslock.Unlock()
}

// 一些使用model进行增删查改的操作放在这里

// 通过player1和player2来进行查找chessboard	,查找成功返回其指针，失败返回nil
func FindChessBoardByPlayers(player1ID, player2ID uint) *ChessBoard {
	var out ChessBoard
	result := chessboardDB.Where("player1_id=? AND player2_id=?", player1ID, player2ID).Find(&out)
	if result.RowsAffected == 0 {
		return nil
	}
	gameslock.Lock()
	defer gameslock.Unlock()
	if a, ok := games[out.ID]; ok {
		return a
	}
	out.UnPack()
	games[out.ID] = &out
	return &out
}

func FindChessBoardByID(gameID uint) (*ChessBoard, error) {
	// 如果找到ChessBaord就返回引用，否则就返回err
	// 首先在games中寻找,
	gameslock.Lock()
	defer gameslock.Unlock()
	game, ok := games[gameID]
	if ok {
		return game, nil
	}
	// 然后games中找不到再到数据库中寻找
	var out ChessBoard
	result := chessboardDB.Where("ID=?", gameID).Find(&out)
	if result.RowsAffected == 0 {
		return nil, errors.New("not found")
	}
	out.UnPack()
	games[out.ID] = &out
	return &out, nil
}

func CreateChessBoard(player1ID, player2ID uint) (*ChessBoard, error) {
	var out ChessBoard = NewChessBoard(10, 10, player1ID, player2ID)
	out.Pack()
	result := chessboardDB.Create(&out)
	if result.RowsAffected == 0 {
		return nil, errors.New("not found")
	}
	gameslock.Lock()
	games[out.ID] = &out
	gameslock.Unlock()

	if len(games) >= max_games_cache_num {
		go shrinkGames(max_games_cache_num/2+1, uint(len(games)))
	}
	return &out, nil
}

// 删除棋盘，返回删掉的棋盘,删除失败返回nil
func DeleteChessBoard(ID uint) *ChessBoard {
	var mode ChessBoard
	chessboard, err := FindChessBoardByID(ID)
	if err != nil {
		return nil
	}
	result := chessboardDB.Where("ID=?", ID).Delete(&mode)
	if result.Error != nil {
		return nil
	}
	return chessboard
}
