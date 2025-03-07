package funcs

import (
	"encoding/json"
	"time"
)

type WebsocketData struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}
type STATUS struct {
	Status string `json:"status"`
	Since  string `json:"since"`
	Sender SUser  `json:"sender"`
}
type Private_Message struct {
	Type string 
	Id         int
	SenderID   int
	ReceiverID int
	Content    string
	Created_at time.Time
}

type Group_Mesaage struct {
	Type string 
	Vv string 
	Id         int
	SenderID   int
	Sender     SUser
	SenderName string
	Group      int
	Content    string
	Created_at time.Time
}

type SUser struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Nickname  string `json:"nickname"`
	Path      string `json:"path"`
}
