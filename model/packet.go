package model

type Packet struct {
	Name      string `json:"name"`
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
	Status    string `json:"status"`
}
