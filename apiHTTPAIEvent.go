package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

// Event struct để lưu thông tin sự kiện
type Event struct {
	ID            string `json:"id"`
	Description   string `json:"description"`
	Timestamp     int64  `json:"timestamp"`
	Image         string `json:"image"`
	StorageBucket string `json:"storageBucket"`
	EventType     string `json:"eventType"`
	CameraID      string `json:"cameraId"`
	Result        string `json:"result"`
	Location      string `json:"location"`
}

type DTO_Event struct {
	ID            string `json:"id"`
	Description   string `json:"description"`
	Timestamp     int64  `json:"timestamp"`
	Image         string `json:"image"`
	StorageBucket string `json:"storageBucket"`
	EventType     string `json:"eventType"`
	CameraID      string `json:"cameraId"`
	Result        string `json:"result"`
	Location      string `json:"location"`
}

// Global slice để lưu trữ các sự kiện từ RabbitMQ
var eventList []Event

var amqpURI = "amqp://guest:guest@localhost:5672/" // Địa chỉ RabbitMQ, thay đổi theo cấu hình của bạn

func ConnectRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(amqpURI)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to connect to RabbitMQ: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to open a channel: %s", err)
	}

	return conn, ch, nil
}

// Hàm xử lý tin nhắn từ RabbitMQ và lưu vào PostgreSQL
func ListenAndProcessRabbitMQ() {
	conn, ch, err := ConnectRabbitMQ()
	if err != nil {
		log.Fatalf("Error connecting to RabbitMQ: %s", err)
	}
	defer conn.Close()
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"ai_events_queue", // Tên queue
		true,              // Đảm bảo tin nhắn được lưu trữ
		false,             // Không xóa queue khi không có consumer
		false,             // Không có queue exclusive
		false,             // Không chờ tin nhắn
		nil,               // Các option thêm vào
	)
	if err != nil {
		log.Fatalf("Error declaring the queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name, // Queue name
		"",     // Consumer name
		true,   // Auto-acknowledge
		false,  // Exclusive
		false,  // No local
		false,  // No wait
		nil,    // Arguments
	)
	if err != nil {
		log.Fatalf("Error consuming messages: %s", err)
	}

	// Lắng nghe tin nhắn và xử lý
	for msg := range msgs {
		var event DTO_Event
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Printf("Error unmarshalling message: %s", err)
			continue
		}

		// Chuyển đổi dữ liệu và lưu vào PostgreSQL
		_, err := CreateItemFromDTO[DTO_Event, Event](event)
		if err != nil {
			log.Printf("Error creating item in DB: %s", err)
		}
	}
}

func SearchAIEvent(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := NewJsonDTOListRsp[DTO_Event]()

	// Get param
	keywords := strings.Split(c.Query("keyword"), ",")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	if page < 1 {
		page = 1
	}
	query := NewQuery[DTO_Event, Event]()
	if len(keywords) > 0 {
		for _, keyword := range keywords {
			query.AddConditionOfTextField("OR", "event_type_string", "LIKE", keyword)
			query.AddConditionOfTextField("OR", "description", "LIKE", keyword)
			query.AddConditionOfTextField("OR", "cam_name", "LIKE", keyword)
			query.AddConditionOfTextField("OR", "location", "LIKE", keyword)
			query.AddConditionOfTextField("OR", "result", "LIKE", keyword)
		}
	}
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)

	// Exec query
	dtoCabinBasics, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}

	if count == 0 {
		jsonRspDTOCabinsBasicInfos.Data = []DTO_Event{}
	} else {
		jsonRspDTOCabinsBasicInfos.Data = dtoCabinBasics
	}
	jsonRspDTOCabinsBasicInfos.Size = int64(limit)
	jsonRspDTOCabinsBasicInfos.Page = int64(page)
	jsonRspDTOCabinsBasicInfos.Count = int64(count)

	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)
}

func ReadCabinEventAI(c *gin.Context) {

	jsonRsp := NewJsonDTORsp[DTO_Event]()

	dto, err := ReadItemByIDIntoDTO[DTO_Event, Event](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}
