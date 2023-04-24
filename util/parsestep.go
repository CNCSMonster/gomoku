package util

import (
	"cncsmonster/gomoku/model"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// 把发来的下棋请求,也就是位置,进行解析

func ParseStep(r *http.Request) (uint, uint, error) {
	//r的body实现了可读接口
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		var step model.Step
		if json.Unmarshal(body, &step) == nil {
			return step.X, step.Y, nil
		}
	}
	return 0, 0, errors.New("parse error")
}
