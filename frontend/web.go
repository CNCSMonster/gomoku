package frontend

import (
	"cncsmonster/gomoku/util"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var (
	boardPage,
	backGroundImg,
	backGroundDeepImg,
	blackStoneImg,
	whiteStoneImg *util.Page
)

func initPage(filename string) *util.Page {
	out, err := util.LoadPage(filename)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return out
}

// 初始化常用资源
func init() {
	boardPage = initPage("frontend/board.html")
	backGroundImg = initPage("frontend/res/background.gif")
	backGroundDeepImg = initPage("frontend/res/background-deep.png")
	blackStoneImg = initPage("frontend/res/blackStone.gif")
	whiteStoneImg = initPage("frontend/res/whiteStone.gif")
}

func ShowBoard(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%s", string(boardPage.Body))
}

func GiveJS(w http.ResponseWriter, r *http.Request) {
	body, _ := os.ReadFile("frontend/board.js")
	w.Header().Set("Content-Type", "text/js")
	fmt.Fprintf(w, "%s", string(body))
}
func GiveCss(w http.ResponseWriter, r *http.Request) {
	body, _ := os.ReadFile("frontend/board.css")
	w.Header().Set("Content-Type", "text/css")
	fmt.Fprintf(w, "%s", string(body))
}

func GiveImages(w http.ResponseWriter, r *http.Request) {
	// str := r.URL.Path
	// fmt.Println(str)
	vars := mux.Vars(r)
	str := vars["graphName"]
	fmt.Println(str)
	if strings.HasPrefix(str, "background.gif") {
		w.Header().Set("Content-Type", "image/gif")
		fmt.Fprintf(w, "%s", string(backGroundImg.Body))
	} else if strings.HasPrefix(str, "background-deep.png") {
		w.Header().Set("Content-Type", "image/png")
		fmt.Fprintf(w, "%s", string(backGroundDeepImg.Body))
	} else if strings.HasPrefix(str, "whiteStone.gif") {
		w.Header().Set("Content-Type", "image/gif")
		fmt.Fprintf(w, "%s", string(whiteStoneImg.Body))
	} else if strings.HasPrefix(str, "blackStone.gif") {
		w.Header().Set("Content-Type", "image/gif")
		fmt.Fprintf(w, "%s", string(blackStoneImg.Body))
	}
}
