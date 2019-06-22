package model

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	Id          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	SessionMotd string          `json:"sesionMotd"`
	IsWriting   bool            `json:"isWriting"`
	Connection  *websocket.Conn `json:"-"`
}

func NewUser(username string) User {
	return User{Id: uuid.New(), Name: username, IsWriting: false}
}
