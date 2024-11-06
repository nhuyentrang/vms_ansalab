package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"vms/wssignaling"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Solve cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// ServeWS			godoc
// @Summary      	Process webrtc signaling between client and camnet gateway
// @Description  	Responds with the signaling message as JSON.
// @Tags         	webrtc-signalings
// @Param        	message  body   wssignaling.SignalingMessage  true  "Signaling message JSON"
// @Produce      	json
// @Success      	200	{object}	wssignaling.SignalingMessage
// @Router       	/ws/signaling/{id} [get]
// @Security		BearerAuth
func ServeWS(c *gin.Context) {
	// Save client
	wsclientID := c.Param("id")
	// Get query param
	channel := c.Param("channel")
	startTime, _ := strconv.Atoi(c.Query("startTime"))
	endTime, _ := strconv.Atoi(c.Query("endTime"))
	scale, _ := strconv.Atoi(c.Query("scale"))

	fmt.Println(
		"wsclientID: ", wsclientID,
		" - channel: ", channel,
		" - startTime: ", startTime,
		" - endTime: ", endTime,
		" - scale: ", scale)

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	wsclient := wssignaling.NewClient(wsclientID, ws)

	wssignaling.Register(wsclient)

	fmt.Println("===============================> New client, ID: ", wsclient.ID)

	//hub.register <- client
	go wsclient.Write()
	go wsclient.Read()
}

// ServeWS			godoc
// @Summary      	Process webrtc signaling between client and camnet gateway
// @Description  	Responds with the signaling message as JSON.
// @Tags         	webrtc-signalings
// @Param        	message  body   wssignaling.SignalingMessage  true  "Signaling message JSON"
// @Produce      	json
// @Success      	200	{object}	wssignaling.SignalingMessage
// @Router       	/ws/signaling/video/{id} [get]
// @Security		BearerAuth
func ServeWSS(c *gin.Context) {
	// Save client
	wsclientID := c.Param("id")
	fmt.Print(wsclientID)
	// Get query param
	channel := c.Param("channel")
	startTime, _ := strconv.Atoi(c.Query("startTime"))
	endTime, _ := strconv.Atoi(c.Query("endTime"))
	scale, _ := strconv.Atoi(c.Query("scale"))

	fmt.Println(
		"wsclientID: ", wsclientID,
		" - channel: ", channel,
		" - startTime: ", startTime,
		" - endTime: ", endTime,
		" - scale: ", scale)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// wsclient := wssignaling.NewClient(wsclientID, ws)

	// wssignaling.Register(wsclient)

	// fmt.Println("===============================> New client, ID: ", wsclient.ID)

	// //hub.register <- client
	// go wsclient.Write()
	// go wsclient.Read()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Lỗi khi đọc dữ liệu:", err)
			return
		}
		fmt.Printf("Nhận được thông điệp: %s\n", p)

		// Xử lý dữ liệu ở đây (ví dụ: gửi lại cùng thông điệp)
		err = conn.WriteMessage(messageType, p)
		if err != nil {
			fmt.Println("Lỗi khi gửi dữ liệu:", err)
			return
		}
	}

}

// ReadSensor	 godoc
// @Summary      Recieve single report
// @Description  Returns the new updated report for websocket client
// @Tags         websockets
// @Produce      json
// @Param        topic  path  string  true  "topic" Enums(reports, cabin-list, logs, notifications)
// @Param        client-uuid  path  string  true  "uuid of websocket client"
// @Success      200   {object}  models.DTO_WebsocketMessage
// @Router       /websocket/topic/{topic}/{client-uuid} [get]
// @Security	 BearerAuth
func WSServe(c *gin.Context) {
	Serve(c)
}

// Function to handle websocket connection and register client to hub and start goroutines
func Serve(c *gin.Context) {
	wsclientID := c.Param("id")
	wsclientTopic := c.Param("topic")

	if wsclientTopic == "" {
		cabinID := c.Param("cabin-uuid")
		wsclientTopic = "cabin-details_" + cabinID
	}

	fmt.Print(wsclientID)
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	wsclient := wssignaling.NewClientTopic(wsclientID, wsclientTopic, ws)

	wssignaling.RegisterTopic(wsclient)

	fmt.Printf("==============> New client, ID: %s, topic: %s\n", wsclient.ID, wsclient.Topic)

	go wsclient.WriteTopic()
	go wsclient.ReadTopic()
}
