package main

type ClientMessage struct {
	client  *Client
	payload []byte
}

type Server struct {
	// Logger for the server.
	// This is a simple logger that writes to stdout.
	logger *Logger

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Custom channel with both client and payload
	informOthers chan *ClientMessage

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newServer(logger *Logger, clients map[*Client]bool) *Server {
	return &Server{
		logger:       logger,
		clients:      clients,
		broadcast:    make(chan []byte),
		informOthers: make(chan *ClientMessage),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
	}
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.clients[client] = true
			s.logger.serverMessage("Client registered: " + client.clientName)
			client.send <- []byte("Welcome " + client.clientName + "!")
			for others := range s.clients {
				if others != client {
					others.send <- []byte(client.clientName + " has joined the chat!")
				}
			}
			s.logger.welcomeMessage(client)
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				close(client.send)
				s.logger.serverMessage("Client closed connection: " + client.clientName)
			}
		case message := <-s.broadcast:
			for client := range s.clients {
				select {
				case client.send <- message:
					s.logger.userMessage(string(message), client)
				default:
					close(client.send)
					delete(s.clients, client)
					s.logger.serverMessage("Client closed connection: " + client.clientName)
				}
			}
		case clientInfo := <-s.informOthers:
			// Inform all clients except the sender about the message.
			for client := range s.clients {
				if client != clientInfo.client {
					client.send <- clientInfo.payload
					s.logger.serverMessage(string(clientInfo.payload))
				}
			}
		}
	}
}
