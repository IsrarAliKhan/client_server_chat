package main

import (
	"bufio"
	"chat-app/constants"
	"chat-app/model"
	"chat-app/utils"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var connections []net.Conn

func main() {
	// create socket
	listener, err := net.Listen(constants.Network, constants.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup(listener)
	log.Println("Server started on " + constants.Port)

	// handle interrupt call
	serverInterrupt := make(chan os.Signal)
	signal.Notify(serverInterrupt, os.Interrupt, syscall.SIGTERM)
	go interrupt(serverInterrupt, listener)

	for {
		// accept client request
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		fmt.Println("Client " + conn.RemoteAddr().String() + " connected.")

		// save client connection
		connections = append(connections, conn)

		// thread to handle each client
		go clientHandler(conn)
	}
}

func clientHandler(conn net.Conn) {
	for {
		// read message from channel
		message, err := bufio.NewReader(conn).ReadBytes(constants.Delimiter)
		if err != nil {
			fmt.Printf("client [%s] left\n", conn.RemoteAddr())
			index := find(conn.RemoteAddr().String())
			if index != constants.StatusNotFound {
				remove(index)
			}
			conn.Close()
			return
		}

		// unmarshal the message to packet
		packet := utils.Unpack(message)

		// packet not acceptable
		if packet.Status != constants.StatusAcceptable {
			packet.Status = constants.StatusNotAcceptable
			conn.Write(utils.Pack(packet))
			continue
		}

		// log at server
		fmt.Print(fmt.Sprintf("[%s]: %s", packet.Name, packet.Message))

		// check if recipient exits
		index := find(packet.Recipient)
		if index == constants.StatusNotFound {
			packet.Name = constants.ServerName
			packet.Status = constants.StatusRecipientNotFound
			packet.Message = constants.MsgRecipientNotFound
			conn.Write(utils.Pack(packet))
			continue
		}

		// send message to recipient
		for _, connection := range connections {
			if connection.RemoteAddr().String() == packet.Recipient {
				connection.Write(utils.Pack(packet))
			}
		}
	}
}

func cleanup(listener net.Listener) {
	// draft packet
	packet := model.Packet{
		Name:    constants.ServerName,
		Status:  constants.StatusInterrupted,
		Message: constants.MsgInterrupted,
	}

	// notify each client
	for _, connection := range connections {
		connection.Write(utils.Pack(packet))
	}

	// close listener
	listener.Close()
}

func interrupt(serverInterrupt chan os.Signal, listener net.Listener) {
	<-serverInterrupt
	cleanup(listener)
	os.Exit(1)
}

func find(address string) int {
	cIndex := constants.StatusNotFound
	for index, connection := range connections {
		if connection.RemoteAddr().String() == address {
			cIndex = index
		}
	}

	return cIndex
}

func remove(index int) {
	connections[index] = connections[len(connections)-1]
	connections = connections[:len(connections)-1]
}
