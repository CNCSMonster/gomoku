package main

import (
	"cncsmonster/gomoku/model"
	"fmt"
	"time"
)

func testSpeed() {
	//测试结果，使用goroutine优化前后，410->4 ms
	timeStart := time.Now()
	for i := 0; i < 1000000; i++ {
		board := model.NewChessBoard(10, 10, 1, 1)
		board.CheckWin()
	}
	timeEnd := time.Now()
	fmt.Println(timeEnd.Sub(timeStart))
}

func main() {
	// 测试，出现了死循环
	// board := model.NewChessBoard(10, 10, 1, 2)
	// winner, ifwinnere := board.CheckWin()
	// fmt.Println(winner, ifwinnere)
	testPlay()
	// testSpeed()
}
