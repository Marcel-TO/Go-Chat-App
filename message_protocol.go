package main

import (
	"bytes"
)

type MessageType int

// MessageType represents the type of message being sent in the byte stream.
const (
	InfoMessage       MessageType = iota // 0: Info message from Server to Client
	TextMessage                          // 1: Text message from Client to everyone
	CommandMessage                       // 2: Command message from Client to change client info
	AllClientsMessage                    // 3: Message from Server about all clients
)

var (
	appSignature = []byte{23, 10, 0}
)

// Message Code Protocol
type MessageCodeProtocol struct {
	appSignature  []byte
	messageType   MessageType
	messageLength uint16
	message       []byte
}

func newMessageCodeProtocol(messageType MessageType, message []byte) *MessageCodeProtocol {
	return &MessageCodeProtocol{
		appSignature:  appSignature,
		messageType:   messageType,
		messageLength: uint16(len(message)),
		message:       message,
	}
}

func (m *MessageCodeProtocol) encode() []byte {
	// Encode the message into a byte slice.
	encodedMessage := make([]byte, 0, len(m.appSignature)+2+len(m.message))
	encodedMessage = append(encodedMessage, m.appSignature...)
	encodedMessage = append(encodedMessage, byte(m.messageType))
	encodedMessage = append(encodedMessage, byte(m.messageLength>>8), byte(m.messageLength&0xFF))
	encodedMessage = append(encodedMessage, m.message...)
	return encodedMessage
}

func decode(data []byte) *MessageCodeProtocol {
	if len(data) < 5 {
		return nil // Not enough data to decode
	}
	if bytes.Equal(data[:3], appSignature) {
		messageType := MessageType(data[3])
		messageLength := uint16(data[4])<<8 | uint16(data[5])
		if len(data) < int(messageLength)+6 {
			return nil // Not enough data to decode
		}
		message := data[6 : 6+messageLength]
		return &MessageCodeProtocol{
			appSignature:  appSignature,
			messageType:   messageType,
			messageLength: messageLength,
			message:       message,
		}
	}
	return nil // Invalid signature
}

func encodeMessage(server *Server, client *Client, messageType MessageType, message []byte) []byte {
	// Create a new MessageCodeProtocol instance and encode it.
	protocol := newMessageCodeProtocol(messageType, message)
	return protocol.encode()
}
