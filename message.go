package main

type MessageType uint8

const (
	InfoMessage MessageType = iota
	CriticalMessage
)

type Message struct {
	Type          MessageType
	ClientID      int
	MessageLength uint32
	Message       []byte
}
