package util

import (
	"os"
)

type Page struct {
	FileName string
	Body     []byte
}

// 从文件中读取对应html
func LoadPage(filename string) (*Page, error) {
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{FileName: filename, Body: body}, nil
}

func (p *Page) Save() error {
	return os.WriteFile(p.FileName, p.Body, 0600)
}
