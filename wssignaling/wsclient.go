package wssignaling

import (
	"fmt"
	"log"
	"sync"
	"time"

	"vms/internal/models"

	"vms/reposity"

	"github.com/gorilla/websocket"
)

const (

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 64 * 1024 // 4096
)

// Client struct for websocket connection and message sending
type Client struct {
	ID   string
	Conn *websocket.Conn
	send chan SignalingMessage
}

var wg sync.WaitGroup

// NewClient creates a new client
func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{ID: id, Conn: conn, send: make(chan SignalingMessage, 4096)}
}

// Client goroutine to read messages from client
func (c *Client) Read() {

	defer func() {
		UnRegister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var recvMsg SignalingMessage
		err := c.Conn.ReadJSON(&recvMsg)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		// marshaledJsonMsg, err := json.MarshalIndent(recvMsg, "", "   ")
		// if err != nil {
		// 	fmt.Println("Error: ", err)
		// 	break
		// }

		//fmt.Println("WS Signaling: ===============================> Read From ID: ", c.ID, ", msg: ", string(marshaledJsonMsg))
		log.Printf("WS Signaling: ===============================> Read From ID: ", c.ID)
		// Send this message to destination client
		sendMsg := SignalingMessage{
			ID:          c.ID,
			Type:        recvMsg.Type,
			Sender:      c.ID,
			RecipientID: recvMsg.ID,
			SDP:         recvMsg.SDP,
			Channel:     recvMsg.Channel,
			StartTime:   recvMsg.StartTime,
			EndTime:     recvMsg.EndTime,
			Scale:       recvMsg.Scale,
			Content:     recvMsg.Content,
			Serial:      recvMsg.Serial,
			ViewType:    recvMsg.ViewType,
		}

		// if recvMsg.Type == "request" {
		// 	if len(sendMsg.RecipientID) >= 30 {
		// 		fmt.Println("len(sendMsg.RecipientID)", len(sendMsg.RecipientID))
		dto, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", sendMsg.RecipientID)
		if err != nil {
			fmt.Println("======> Err sending ID camera: ", err)
		}
		sendMsg.IPv4 = dto.IPAddress
		sendMsg.Host = dto.HttpPort
		sendMsg.User = dto.Username
		sendMsg.Password = dto.Password
		sendMsg.Streams = dto.Streams
		sendMsg.Serial = dto.SerialNumber

		// fmt.Println("WS Signaling: ===============================> start sending to client ID: ", sendMsg.RecipientID, ", type: ", sendMsg.Type, ", sdp: ", sendMsg.SDP)
		// fmt.Println("WS Signaling: ===============================> start sending to client: ", sendMsg)
		log.Printf("WS Signaling: ===============================> start sending to client ID: ", sendMsg.RecipientID)
		SendSignalingMessage(sendMsg)

		// 	} else {
		// 		defer wg.Done()

		// 		fmt.Println("len(sendMsg.RecipientID)", len(sendMsg.RecipientID))
		// 		dto, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("serial_number = ?", sendMsg.RecipientID)
		// 		if err != nil {
		// 			fmt.Println("======> Err sending Serial camera: ", err)
		// 		}
		// 		sendMsg.IPv4 = dto.IPAddress
		// 		sendMsg.Host = dto.HttpPort
		// 		sendMsg.User = dto.Username
		// 		sendMsg.Password = dto.Password
		// 		sendMsg.Streams = dto.Streams
		// 		sendMsg.Serial = dto.SerialNumber
		// 		sendMsg.ID = dto.ID.String()
		// 		fmt.Println("(String())", dto.ID.String())

		// 		fmt.Println("WS Signaling:M2===============================> start sending to client ID: ", sendMsg.RecipientID, ", type: ", sendMsg.Type, ", sdp: ", sendMsg.SDP)
		// 		fmt.Println("WS Signaling:M2===============================> start sending to client: ", sendMsg)
		// 		wg.Wait()
		// 		SendSignalingMessage(sendMsg)
		// 	}
		// }
		// if recvMsg.Type == "offer" || recvMsg.Type == "answer" {
		// 	dto, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", sendMsg.RecipientID)
		// 	if err != nil {
		// 		fmt.Println("======> Err sending ID camera: ", err)
		// 	}

		// 	sendMsg.IPv4 = dto.IPAddress
		// 	sendMsg.Host = dto.HttpPort
		// 	sendMsg.User = dto.Username
		// 	sendMsg.Password = dto.Password
		// 	sendMsg.Streams = dto.Streams
		// 	sendMsg.Serial = dto.SerialNumber

		// 	fmt.Println("WS Signaling: ===============================> start sending to client ID: ", sendMsg.RecipientID, ", type: ", sendMsg.Type, ", sdp: ", sendMsg.SDP)
		// 	fmt.Println("WS Signaling: ===============================> start sending to client: ", sendMsg)
		// 	SendSignalingMessage(sendMsg)
		// }

	}
}

// Client goroutine to write messages to client
func (c *Client) Write() {
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
				err := c.Conn.WriteJSON(message)
				if err != nil {
					fmt.Printf("wsclient id %s: ===============================> error: \n", err)
					break
				}
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Printf("wsclient id %s: ===============================> error: can not write ping message\n", c.ID)
				return
			}
		}

	}
}

// Client closing channel to unregister client
func (c *Client) Close() {
	close(c.send)
}
