package domain

type DocumentMeta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}
