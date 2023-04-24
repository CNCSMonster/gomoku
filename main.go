package main

import (
	"cncsmonster/gomoku/model"
	"cncsmonster/gomoku/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

// 服务端代码
func main() {
	defer func() {
		if e := recover(); e != nil {
			model.SaveAll()
		}
	}()
	go func() {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT)
		// 等待接收中断信号
		<-interrupt
		model.SaveAll()
		panic("for SIGINT (CTRL-C)")
	}()
	r := mux.NewRouter()
	router.RegisterGomokuRoutes(r)
	http.Handle("/gomoku", r)
	fmt.Println("listen at 6363")
	log.Fatal(http.ListenAndServe(":6363", r))
}
