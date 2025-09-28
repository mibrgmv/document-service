package handlers

import "github.com/mibrgmv/document-service/internal/domain"

type AuthRequest struct {
	Login string `json:"login"`
	Pswd  string `json:"pswd"`
}

type RegisterRequest struct {
	Token string `json:"token"`
	Login string `json:"login"`
	Pswd  string `json:"pswd"`
}

type DocumentMeta struct {
	Name   string   `json:"name"`
	File   bool     `json:"file"`
	Public bool     `json:"public"`
	Mime   string   `json:"mime"`
	Grant  []string `json:"grant"`
}

func (m *DocumentMeta) ToDomain() *domain.DocumentMeta {
	return &domain.DocumentMeta{
		Name:   m.Name,
		File:   m.File,
		Public: m.Public,
		Mime:   m.Mime,
		Grant:  m.Grant,
	}
}
