package controller

import (
	"cncsmonster/gomoku/model"
	"cncsmonster/gomoku/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var originSource string = "http://localhost:6363/gomoku"

func HandleGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", originSource) //避免CORS策略导致浏览器不接受
	// 处理 /gomoku/game/{game id}的情况
	vars := mux.Vars(r)
	// 根据id获取游戏棋盘
	ID, err := strconv.ParseUint(vars["gameID"], 10, 32)
	uID := uint(ID)
	if err != nil {
		http.Error(w, "game not exists", http.StatusInternalServerError)
		return
	}
	chessboard, err := model.FindChessBoardByID(uID)
	if err != nil {
		http.Error(w, "game not exists", http.StatusInternalServerError)
		return
	}
	// 从chessboard获取boardcase传递给远程
	playerID := vars["playerID"]
	uPlayerID, err := strconv.ParseUint(playerID, 10, 32)
	if err != nil {
		http.Error(w, "playerID error!", http.StatusInternalServerError)
		return
	}
	boardCase := chessboard.GetBoardCase(uint(uPlayerID))
	body, _ := json.Marshal(boardCase)
	w.WriteHeader(http.StatusOK)                       //使用statusok表示请求完成
	w.Header().Set("Content-Type", "application/json") //使用application/json的contenttype表明传输的http回复体的内容
	fmt.Fprintf(w, "%s", string(body))
}

// 进入游戏
func HandleEnterGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", originSource) //避免CORS策略导致浏览器不接受
	// invite的时候的处理,传入了三个参数，账号密码以及敌人
	vars := mux.Vars(r)
	player1 := vars["player"]
	player2 := vars["enemy"]
	password := vars["password"]
	fmt.Println(player1, player2, password)
	// TODO,升级成一个websocket连接，然后互相传递密钥，建立加密连接
	// 查找账号，如果没有则自动创建
	var account1 *model.Account
	account1 = model.FindAccountByName(player1)
	if account1 == nil {
		var err error
		account1, err = model.CreateAccount(player1, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// 如果创建账号成功,则寻找到房间,返回房间号,或者等待房间建立，等待一段时间
	player1ID := account1.ID

	for i := 0; i < 50; i++ {
		// 循环等待
		account2 := model.FindAccountByName(player2)
		// 不会自动创建敌人账号
		if account2 == nil {
			time.Sleep(1 * time.Second)
			continue
		}
		player2ID := account2.ID

		chessboard := model.FindChessBoardByPlayers(player1ID, player2ID)
		if chessboard == nil {
			chessboard = model.FindChessBoardByPlayers(player2ID, player1ID)
			if chessboard != nil {
				// 棋盘成功,返回棋盘ID和player1ID
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, "%d-%d", chessboard.ID, player1ID)
				return
			}
		}
		if chessboard == nil {
			// 如果该棋盘不存在，创建新棋盘
			var err error
			chessboard, err = model.CreateChessBoard(player1ID, player2ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		// 棋盘成功,返回棋盘ID和player1ID
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%d-%d", chessboard.ID, player1ID)
		return
	}
	http.Error(w, "maybe your enemy not enter the game,try again later", http.StatusInternalServerError)
}

func HandlePlayGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", originSource) //避免CORS策略导致浏览器不接受
	// 收到发送来的请求
	x, y, err := util.ParseStep(r)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		return
	}
	vars := mux.Vars(r)
	ID, err := strconv.ParseUint(vars["gameID"], 10, 32)
	uID := uint(ID)
	if err != nil {
		http.Error(w, "game not exists", http.StatusInternalServerError)
		return
	}
	playerID := vars["playerID"]
	uPlayerID, err := strconv.ParseUint(playerID, 10, 32)
	if err != nil {
		http.Error(w, "playerID error!", http.StatusInternalServerError)
		return
	}
	chessboard, err := model.FindChessBoardByID(uID)
	if err != nil {
		http.Error(w, "game not exists", http.StatusInternalServerError)
		return
	}
	// 判断是player1还是player2以决定是否翻转
	if uPlayerID == uint64(chessboard.Player2ID) {
		y = chessboard.YSize - 1 - y
	}

	if chessboard.Playable(uint(uPlayerID), x, y) {
		chessboard.Play(uint(uPlayerID), x, y)
		// play成功设置回复状态为ok
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// 退出游戏
func HandleDeleteGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", originSource) //避免CORS策略导致浏览器不接受
	// 实现
	// 删除对应的房间，返回statusok如果删除成功
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	ugameID, err := strconv.ParseUint(gameID, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	playerID := vars["playerID"]
	uPlayerID, err := strconv.ParseUint(playerID, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	chessBoard := model.DeleteChessBoard(uint(ugameID))
	if chessBoard == nil {
		http.Error(w, "delete fail!", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	boardCase := chessBoard.GetBoardCase(uint(uPlayerID))
	body, _ := json.Marshal(boardCase)
	fmt.Fprintf(w, "%s", string(body))
}

// init记录
func HandleInitGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", originSource) //避免CORS策略导致浏览器不接受
	// 实现
	// 删除对应的房间，返回statusok如果删除成功
	gameID := r.FormValue("gameID")
	if gameID == "" {
		http.Error(w, "miss argument", http.StatusInternalServerError)
		return
	}
	ugameID, err := strconv.ParseInt(gameID, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if chessboard, err := model.FindChessBoardByID(uint(ugameID)); err == nil {
		chessboard.Init()
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
