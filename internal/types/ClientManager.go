package types

import "encoding/json"

type ClientManager struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
}

func NewClientManager() ClientManager {
	return ClientManager{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan []byte),
	}
}

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.Register:
			manager.Clients[conn] = true
			jsonMessage, _ := json.Marshal(NewMessage("", "", "A new socket has connected."))
			manager.send(jsonMessage, conn)
		case conn := <-manager.Unregister:
			if _, ok := manager.Clients[conn]; ok {
				close(conn.send)
				delete(manager.Clients, conn)
				jsonMessage, _ := json.Marshal(NewMessage("", "", "A socket has disconnected."))
				manager.send(jsonMessage, conn)
			}
		case message := <-manager.Broadcast:
			for conn := range manager.Clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		if conn != ignore {
			conn.send <- message
		}
	}
}
