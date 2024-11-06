package wssignaling

import (
	"fmt"
	"log"

	"vms/internal/models"

	"vms/comongo/reposity"
)

// Hub is a struct that holds all the clients and the messages that are sent to them
type Hub struct {
	// Registered clients.
	clients map[string]map[*Client]bool
	//Unregistered clients.
	unregister chan *Client
	// Register requests from the clients.
	register chan *Client
	// Inbound messages from the clients.
	signalingChannel chan SignalingMessage
}

type SignalingMessage struct {
	ID          string `json:"id,omitempty"`
	Type        string `json:"type,omitempty"`
	SDP         string `json:"sdp,omitempty"`
	Channel     string `json:"channel,omitempty"`
	StartTime   string `json:"startTime,omitempty"`
	EndTime     string `json:"endTime,omitempty"`
	Scale       string `json:"scale,omitempty"`
	Sender      string `json:"sender,omitempty"`
	RecipientID string `json:"recipientID,omitempty"`
	Content     string `json:"content,omitempty"`
	Serial      string `json:"serial,omitempty"`
	ViewType    string `json:"viewType,omitempty"`

	IPv4       string                  `json:"ipv4,omitempty"`
	Host       string                  `json:"host,omitempty"`
	User       string                  `json:"user,omitempty"`
	Password   string                  `json:"password,omitempty"`
	MacAddress string                  `json:"macAddress,omitempty"`
	IPv4Device string                  `json:"iPv4Device,omitempty"`
	Streams    models.VideoStreamArray `json:"streams,omitempty"`
}

type Data struct {
	Playback    Playback `json:"playback,omitempty"`
	Live        Live     `json:"live,omitempty"`
	Capture     Capture  `json:"capture,omitempty"`
	ID          string   `json:"id,omitempty"`
	RecipientID string   `json:"recipientID,omitempty"`
	Content     string   `json:"content,omitempty"`
	SDP         string   `json:"sdp,omitempty"`
	Sender      string   `json:"sender,omitempty"`
	Type        string   `json:"type,omitempty"`
}

type Playback struct {
	Track     string     `json:"track,omitempty"`
	StartTime string     `json:"startTime,omitempty"`
	EndTime   string     `json:"endTime,omitempty"`
	Scale     string     `json:"scale,omitempty"`
	MetaData  []MetaData `json:"metaData,omitempty"`
}

type MetaData struct {
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
	Type      string `json:"type,omitempty"`
}

type Live struct {
	Channel string `json:"channel,omitempty"`
}

type Capture struct {
	AtTime      string `json:"atTime,omitempty"`
	ImageBase64 string `json:"imageBase64,omitempty"`
}

var webRTCWebSocketSignalingHUB *Hub = nil

/*
// Message struct to hold message data
type WSHubMessage struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
	Content   string `json:"content"`
}
*/

func NewHub() *Hub {
	return &Hub{
		clients:          make(map[string]map[*Client]bool),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		signalingChannel: make(chan SignalingMessage),
	}
}

// Core function to run the hub
func (h *Hub) Run() {
	for {
		select {
		// Register a client.
		case client := <-h.register:
			h.RegisterNewClient(client)
			// Unregister a client.
		case client := <-h.unregister:
			h.RemoveClient(client)
			// Send signaling message to recipient client.
		case message := <-h.signalingChannel:
			h.HandleWebRTCSignalingMessage(message)
		}
	}
}

// function check if room exists and if not create it and add client to it
func (h *Hub) RegisterNewClient(client *Client) {
	connections := h.clients[client.ID]
	if connections == nil {
		connections = make(map[*Client]bool)
		h.clients[client.ID] = connections
	}
	h.clients[client.ID][client] = true

	fmt.Println("Size of clients: ", len(h.clients[client.ID]))
	fmt.Println("Total number of clients: ", len(h.clients))
}

// function to remvoe client from room
func (h *Hub) RemoveClient(client *Client) {
	if _, ok := h.clients[client.ID]; ok {
		delete(h.clients[client.ID], client)
		close(client.send)
		fmt.Println("Removed client")
	}
}

// function to handle message based on type of message
func (h *Hub) HandleWebRTCSignalingMessage(message SignalingMessage) {
	fmt.Printf("HandleWebRTCSignalingMessage: ===============================> from id: %s, to id: %s \n", message.Sender, message.RecipientID)

	// only get get recipientID, then clear not necessary feilds
	senderID := message.Sender
	recipientID := message.RecipientID
	typeStatus := "status"
	content := "not available"
	recipientClients := h.clients[recipientID]

	if len(recipientID) <= 30 {
		dto, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", recipientID)
		if err != nil {
			fmt.Println("======> Err sending ID camera: ", err)
		}
		fmt.Println("WS Signaling: ===============================> ID camera: ", dto.ID)
		fmt.Println("WS Signaling: ===============================> ID camera: ", dto.ID)
		if dto.ID.String() != "00000000-0000-0000-0000-000000000000" {
			recipientClients = h.clients[dto.ID.String()]
		}

	}

	log.Printf("======> recipientID", recipientID)
	log.Printf("======> senderID", senderID)
	log.Printf("======> recipientClients", recipientClients)

	if len(recipientClients) == 0 {
		// Send back the message to sender that camera is not found
		senderClients := h.clients[senderID]
		if len(senderClients) > 0 {
			fmt.Println("===")
			replyMsg := SignalingMessage{
				ID:          recipientID,
				Sender:      senderID,
				RecipientID: recipientID,
				Type:        typeStatus,
				Content:     content,
			}
			for client := range senderClients {
				select {
				case client.send <- replyMsg:
				default:
					// If this recipientID is not found, close send and detele that client
					close(client.send)
					delete(h.clients[senderID], client)
				}
			}
		}
	} else {
		log.Println("======")
		// Forward signaling message to recipient
		for client := range recipientClients {
			select {
			case client.send <- message:
			default:
				// If this recipientID is not found, close send and detele that client
				close(client.send)
				delete(h.clients[recipientID], client)
			}
		}
	}

	/*
	   //Check if the message is a type of "message"

	   	if message.Type == "message" {
	   		clients := h.clients[message.ID]
	   		for client := range clients {
	   			select {
	   			case client.send <- message:
	   			default:
	   				close(client.send)
	   				delete(h.clients[message.ID], client)
	   			}
	   		}
	   	}

	   //Check if the message is a type of "notification"

	   	if message.Type == "notification" {
	   		fmt.Println("Notification: ", message.Content)
	   		clients := h.clients[message.Recipient]
	   		for client := range clients {
	   			select {
	   			case client.send <- message:
	   			default:
	   				close(client.send)
	   				delete(h.clients[message.Recipient], client)
	   			}
	   		}
	   	}
	*/
}

func Start() {
	//create new Hub and run it
	webRTCWebSocketSignalingHUB = NewHub()
	go webRTCWebSocketSignalingHUB.Run()
}

func Register(c *Client) {
	fmt.Println("======>Register webRTCWebSocketSignalingHUB mess: ", c)
	if webRTCWebSocketSignalingHUB != nil && c != nil {
		webRTCWebSocketSignalingHUB.register <- c
		fmt.Println("======>Register webRTCWebSocketSignalingHUB ok")
	}
}

func UnRegister(c *Client) {
	webRTCWebSocketSignalingHUB.unregister <- c
}

func SendSignalingMessage(msg SignalingMessage) {
	webRTCWebSocketSignalingHUB.signalingChannel <- msg
}
