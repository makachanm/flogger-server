package main

type LoggerUI struct {
	messages chan Message
}

func NewLoggerUI(msgChan chan Message) *LoggerUI {
	return &LoggerUI{
		messages: msgChan,
	}
}

func (ui *LoggerUI) Start() {
	for msg := range ui.messages {
		// Here you would update your UI with the new message
		// For simplicity, we'll just print it to the console
		switch msg.Type {
		case InfoMessage:
			println("(", msg.ClientID, ") INFO:", string(msg.Message))
		case CriticalMessage:
			println("(", msg.ClientID, ") CRITICAL:", string(msg.Message))
		default:
			println("UNKNOWN TYPE:", string(msg.Message))
		}
	}
}
