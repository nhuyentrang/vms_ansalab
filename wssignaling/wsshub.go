package wssignaling

import (
	"encoding/json"
	"fmt"

	"vms/internal/models"
)

// Hub is a struct that holds all the clients and the messages that are sent to them
type HubTopic struct {
	// Registered clients.
	clients map[string]map[*ClientTopic]bool
	//Unregistered clients.
	unregister chan *ClientTopic
	// Register requests from the clients.
	register chan *ClientTopic
	// Inbound messages from the clients.
	notifyChanel chan NotifyMessage
}

type NotifyMessage struct {
	Topic   string `json:"topic"`
	Message string `json:"message,omitempty"`
}

var defaultWSHUB *HubTopic = nil

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

func NewHubTopic() *HubTopic {
	return &HubTopic{
		clients:      make(map[string]map[*ClientTopic]bool),
		register:     make(chan *ClientTopic),
		unregister:   make(chan *ClientTopic),
		notifyChanel: make(chan NotifyMessage),
	}
}

// Core function to run the hub
func (h *HubTopic) RunTopic() {
	for {
		select {
		// Register a client.
		case client := <-h.register:
			h.RegisterNewClientTopic(client)
			// Unregister a client.
		case client := <-h.unregister:
			h.RemoveClientTopic(client)
		case message := <-h.notifyChanel:
			//Check if the message is a type of "message"
			h.HandleNotifyMessage(message)
		}
	}
}

// function check if room exists and if not create it and add client to it
func (h *HubTopic) RegisterNewClientTopic(client *ClientTopic) {
	connections := h.clients[client.Topic]
	if connections == nil {
		connections = make(map[*ClientTopic]bool)
		h.clients[client.Topic] = connections
	}
	h.clients[client.Topic][client] = true

	fmt.Println("Size of clients: ", len(h.clients[client.Topic]))
}

// function to remvoe client from room
func (h *HubTopic) RemoveClientTopic(client *ClientTopic) {
	if _, ok := h.clients[client.Topic]; ok {
		delete(h.clients[client.Topic], client)
		close(client.send)
		fmt.Println("Removed client")
	}
}

// function to handle message based on type of message
func (h *HubTopic) HandleNotifyMessage(notify NotifyMessage) {
	//fmt.Printf("\t\tWSHub HandleNotifyMessage, topic: %s\n", notify.Topic)
	clients := h.clients[notify.Topic]
	for client := range clients {
		select {
		case client.send <- notify.Message:
		default:
			close(client.send)
			delete(h.clients[notify.Topic], client)
		}
	}
}

func StartTopic() {
	//create new Hub and run it
	defaultWSHUB = NewHubTopic()
	go defaultWSHUB.RunTopic()
}

func RegisterTopic(c *ClientTopic) {
	if defaultWSHUB != nil && c != nil {
		defaultWSHUB.register <- c
	}
}

func UnRegisterTopic(c *ClientTopic) {
	defaultWSHUB.unregister <- c
}

func SendNotifyMessage(topic, dataType string, data interface{}) {
	dataStr, err := json.Marshal(data)
	if err != nil {
		return
	}
	wsMessage := models.DTO_WebsocketMessage{
		DataType: dataType,
		Data:     string(dataStr),
	}
	wsMessageStr, err := json.Marshal(wsMessage)
	if err != nil {
		return
	}
	notify := NotifyMessage{
		Topic:   topic,
		Message: string(wsMessageStr),
	}
	defaultWSHUB.notifyChanel <- notify
}
