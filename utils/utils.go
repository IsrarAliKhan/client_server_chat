package utils

import (
	"chat-app/constants"
	"chat-app/model"
	"encoding/json"
	"log"
)

func Pack(packet model.Packet) []byte {
	b, err := json.Marshal(packet)
	if err != nil {
		packet.Status = constants.StatusNotAcceptable
		packet.Message = constants.MsgNotAcceptable

		b, _ = json.Marshal(packet)
		b = append(b, constants.Delimiter)

		log.Print(err)
		return b
	}

	b = append(b, constants.Delimiter)
	return b
}

func Unpack(message []byte) model.Packet {
	var packet model.Packet

	err := json.Unmarshal(message, &packet)
	if err != nil {
		packet.Status = constants.StatusNotAcceptable
		packet.Message = constants.MsgNotAcceptable
	}

	return packet
}
