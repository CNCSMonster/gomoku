package main

import (
	"cncsmonster/gomoku/model"
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testPlay() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		// 处理连接错误
		panic("wrong connect")
	}
	// 自动创建 user 表
	db.AutoMigrate(&model.ChessBoard{})
	chessboard := model.NewChessBoard(10, 10, 1, 2)
	for i := 0; i < 4; i++ {
		chessboard.Play(1, 0, uint(i))
		chessboard.Play(2, 1, uint(i))
	}
	chessboard.Play(1, 0, 4)
	fmt.Println(chessboard.GetBoardCase(1))
	fmt.Println(chessboard.GetBoardCase(2))
	fmt.Println(chessboard.CheckWin())
	chessboard.Pack()
	db.Save(&chessboard)
}
