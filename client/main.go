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
	"time"
)

var serverInterrupt chan int

func main() {
	// create channel
	serverInterrupt = make(chan int)

	for {
		// establish connection
		conn, err := net.Dial(constants.Network, constants.Port)
		if err != nil {
			log.Println("Failed to connect to server. Trying again in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		log.Println("Connection established at " + constants.Port)

		// create packet
		packet := model.Packet{
			Status: constants.StatusAcceptable,
		}

		// get client name
		fmt.Print("Please enter your name:\t")
		if _, err = fmt.Scanln(&packet.Name); err != nil || packet.Name == "" {
			log.Fatal(err)
			return
		}

		// get recipient address
		fmt.Print("Please enter recipient address:\t")
		if _, err = fmt.Scanln(&packet.Recipient); err != nil || packet.Recipient == "" {
			log.Fatal(err)
			return
		}

		// thread for sending messages
		go sender(conn, packet)

		// thread for receiving messages
		go receiver(conn)

		// wait for server interrupt
		<-serverInterrupt
	}
}

func sender(conn net.Conn, packet model.Packet) {
	for {
		// read message
		message, err := bufio.NewReader(os.Stdin).ReadString(constants.Delimiter)
		if err != nil {
			log.Fatal(err)
		}

		// write message to packet
		packet.Message = message

		// write in channel
		conn.Write(utils.Pack(packet))
	}
}

func receiver(conn net.Conn) {
	for {
		// read message from channel
		message, err := bufio.NewReader(conn).ReadBytes(constants.Delimiter)
		if err != nil {
			continue
		}

		// unmarshal the message to packet
		packet := utils.Unpack(message)

		// server interrupts
		if packet.Status == constants.StatusInterrupted {
			conn.Close()
			serverInterrupt <- constants.ServerInterruptCall
		}

		// message not acceptable
		if packet.Status == constants.StatusNotAcceptable {
			continue
		}

		// response to client
		fmt.Printf("[%s] %s", packet.Name, constants.NotificationSound+packet.Message)
	}
}
