package types

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	conn *websocket.Conn
	send chan []byte
}

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		conn: conn,
		send: make(chan []byte),
	}
}

func (c *Client) Read(manager *ClientManager) {
	defer func() {
		manager.Unregister <- c
		c.conn.Close()
	}()
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			manager.Unregister <- c
			c.conn.Close()
			break
		}
		jsonMessage, _ := json.Marshal(NewMessage(c.ID, "", string(message)))
		manager.Broadcast <- jsonMessage
	}
}

func (c *Client) Write() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			}
			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
