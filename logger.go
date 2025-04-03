package main

import (
	"log"
)

type Color string

// terminal color enums
const (
	Reset   Color = "\033[0m"
	Red     Color = "\033[31m"
	Green   Color = "\033[32m"
	Yellow  Color = "\033[33m"
	Blue    Color = "\033[34m"
	Magenta Color = "\033[35m"
	Cyan    Color = "\033[36m"
	Gray    Color = "\033[37m"
	White   Color = "\033[97m"
)

type Logger struct {
	logger *log.Logger
}

func newLogger() *Logger {
	return &Logger{
		logger: log.Default(),
	}
}

func (l *Logger) info(msg string) {
	l.logger.Println("INFO: " + msg)
}

func (l *Logger) userMessage(msg string, client *Client) {
	l.logger.Println("[" + string(client.clientColor) + client.clientName + string(Reset) + "]: " + msg)
}

func (l *Logger) welcomeMessage(client *Client) {
	l.logger.Println("[Server]: Let us Welcome the User " + client.clientName + string(Reset) + " to the server!")
}

func (l *Logger) serverMessage(msg string) {
	l.logger.Println("[Server]: " + msg)
}
