package router

import (
	"cncsmonster/gomoku/controller"
	"cncsmonster/gomoku/frontend"

	"github.com/gorilla/mux"
)

var RegisterGomokuRoutes = func(router *mux.Router) {
	// router
	router.HandleFunc("/gomoku/", frontend.ShowBoard)
	router.HandleFunc("/gomoku/board.html", frontend.ShowBoard).Methods("GET")
	router.HandleFunc("/gomoku/board.js", frontend.GiveJS).Methods("GET")
	router.HandleFunc("/gomoku/board.css", frontend.GiveCss).Methods("GET")
	router.HandleFunc("/gomoku/res/{graphName}", frontend.GiveImages).Methods("GET")
	router.HandleFunc("/gomoku/game/{gameID}/{playerID}", controller.HandleGame).Methods("GET")
	router.HandleFunc("/gomoku/game/{gameID}/{playerID}", controller.HandleDeleteGame).Methods("DELETE")
	router.HandleFunc("/gomoku/invite/{player}/{enemy}/{password}", controller.HandleEnterGame).Methods("GET")
	router.HandleFunc("/gomoku/game/{gameID}/{playerID}", controller.HandlePlayGame).Methods("POST")
	router.HandleFunc("/gomoku/game", controller.HandleInitGame).Methods("PUT")
}
