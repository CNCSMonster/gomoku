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
	for i := 0; i < 5; i++ {
		chessboard.Play(1, 0, uint(i))
		chessboard.Play(2, 1, uint(i))
	}
	fmt.Println(chessboard.CheckWin())
	chessboard.Pack()
	db.Save(&chessboard)
}
