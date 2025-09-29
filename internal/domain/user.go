package domain

import "time"

type User struct {
	ID       string    `json:"id"`
	Login    string    `json:"login"`
	Password string    `json:"-"`
	Created  time.Time `json:"created"`
}
