package domain

import "time"

type Document struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Mime    string    `json:"mime"`
	File    bool      `json:"file"`
	Public  bool      `json:"public"`
	Created time.Time `json:"created"`
	Grant   []string  `json:"grant"`
	Owner   string    `json:"-"`
	Data    []byte    `json:"-"`
	JSON    string    `json:"-"`
}
