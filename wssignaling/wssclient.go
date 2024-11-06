package wssignaling

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// Client struct for websocket connection and message sending
type ClientTopic struct {
	ID    string
	Conn  *websocket.Conn
	send  chan string
	Topic string
}

// NewClient creates a new client
func NewClientTopic(id string, topic string, conn *websocket.Conn) *ClientTopic {
	return &ClientTopic{ID: id, Topic: topic, Conn: conn, send: make(chan string, 4096)}
}

// Client goroutine to read messages from client
func (c *ClientTopic) ReadTopic() {

	defer func() {
		UnRegisterTopic(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		fmt.Printf("==============> WS read, messageType: %d, topic: %s, msg: %s\n", messageType, c.Topic, string(p))
	}
}

// Client goroutine to write messages to client
func (c *ClientTopic) WriteTopic() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			} else {
				//err := c.Conn.WriteJSON(message)
				// Remove char \ in string
				//data := strings.ReplaceAll(message, "\\", "")
				err := c.Conn.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					fmt.Println("Error: ", err)
					break
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// Client closing channel to unregister client
func (c *ClientTopic) CloseTopic() {
	close(c.send)
}
