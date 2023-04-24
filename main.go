package main

import (
	"cncsmonster/gomoku/model"
	"cncsmonster/gomoku/router"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// 服务端代码
func main() {
	defer func() {
		if e := recover(); e != nil {
			model.SaveAll()
			fmt.Println("ggg???")
		}
	}()
	r := mux.NewRouter()
	router.RegisterGomokuRoutes(r)
	http.Handle("/gomoku", r)
	log.Fatal(http.ListenAndServe(":6363", r))
}
