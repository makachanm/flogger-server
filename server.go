package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

type ClientConnetionStatus int

const (
	Disconnected ClientConnetionStatus = iota
	Connected
)

type FloggerServer struct {
	connection net.Listener
	buffer     chan Message
	terminate  chan bool
	status     ClientConnetionStatus
}

func NewFloggerServer(buf chan Message, term chan bool) *FloggerServer {
	sock, err := net.Listen("unix", "/tmp/flogger.sock")
	if err != nil {
		panic(err)
	}

	return &FloggerServer{
		status:     Disconnected,
		buffer:     buf,
		terminate:  term,
		connection: sock,
	}
}

func (fs *FloggerServer) StartServer() {
	go func() {
		<-fs.terminate
		fs.connection.Close()
		close(fs.terminate)
		close(fs.buffer)

	}()

	go func() {
		for {
			connection, err := fs.connection.Accept()
			fmt.Println("Client attached!")
			if err != nil {
				// connection closed
				break
			}

			fs.status = Connected
			go fs.handleConnection(connection)
		}

	}()
}

// handleConnection handles an incoming client connection.
func (fs *FloggerServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		var messageType MessageType
		err := binary.Read(conn, binary.BigEndian, &messageType)
		if err != nil {
			if err != io.EOF {
				log.Println("type read error:", err)
			}
			break
		}

		var clientID int32
		err = binary.Read(conn, binary.BigEndian, &clientID)
		if err != nil {
			if err != io.EOF {
				log.Println("id read error:", err)
			}
			break
		}

		var msgLen uint32
		err = binary.Read(conn, binary.BigEndian, &msgLen)
		if err != nil {
			if err != io.EOF {
				log.Println("len read error:", err)
			}
			break
		}

		msg := make([]byte, msgLen)
		_, err = io.ReadFull(conn, msg)
		if err != nil {
			if err != io.EOF {
				log.Println("msg read error:", err)
			}
			break
		}

		fs.buffer <- Message{
			Type:          MessageType(messageType),
			ClientID:      int(clientID),
			MessageLength: msgLen,
			Message:       msg,
		}
	}
}
