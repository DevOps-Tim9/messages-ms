package dto

type MessageDto struct {
	Id           uint   `json:"id"`
	From         uint   `json:"from"`
	To           uint   `json:"to"`
	Text         string `json:"text"`
	Conversation uint   `json:"conversation"`
}
