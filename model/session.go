package model

import "github.com/google/uuid"

type Session struct {
	Id         uuid.UUID `json:"id"`
	Motd       string    `json:"motd"`
	CreateTime string    `json:"time"`
	Users      map[uuid.UUID]User    `json:"-"`
	Messages   []Message `json:"-"`
}
