package model

type Message struct {
	Text      string `json:"text"`
	Time      string `json:"time"`
	UserId    int    `json:"userId"`
	SystemMsg bool   `json:"systemMsg"`
}
