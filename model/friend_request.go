package model

type RequestType int

const (
	Outgoing RequestType = iota
	Incoming
)

type FriendRequest struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
	Type RequestType `json:"type" enums:"0,1"`
}
