package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"vms/internal/models"
	"vms/wssignaling"

	"vms/comongo/minioclient"

	"github.com/go-co-op/gocron"

	"vms/comongo/kafkaclient"
	"vms/comongo/reposity"
)

func RunAICoreEventProcessing(consumerName string) {
	// Loop for comsuming message queue for aicore event
	go func() {
		for {
			// Waiting for msg in channel
			msg, topic, timestamp, err := kafkaclient.ConsumerReadMessage(consumerName)
			if err != nil {
				continue
			}
			if msg != "" && topic != "" {
				start := time.Now()
				// Extract data from message
				var eventMessage models.KafkaJsonAIEventMessage
				err := json.Unmarshal([]byte(msg), &eventMessage)
				if err != nil {
					log.Printf("==============> Failed to decode kafka message, error: %s, msg: %s, topic: %s, record timestamp: %s\n", err, msg, topic, timestamp.Format(time.RFC3339))
					continue
				}
				err = ProcessAIVMSEventData(eventMessage)
				elapsed := time.Since(start)
				if err != nil {
					log.Printf("==============> Failed to  process kafka message, error: %s, took %s\n", err, elapsed)
					continue
				}
			}
		}
	}()
}

func formatMACAddress(mac string) (string, error) {
	if len(mac) != 12 {
		return "", fmt.Errorf("invalid MAC address: %s", mac)
	}

	var sb strings.Builder
	for i := 0; i < len(mac); i += 2 {
		if i > 0 {
			sb.WriteString(":")
		}
		sb.WriteString(mac[i : i+2])
	}
	return sb.String(), nil
}

// Process data from ai vms event topic
func ProcessAIVMSEventData(eventMessage models.KafkaJsonAIEventMessage) error {
	var dtoAIWaring models.DTO_AIWaring

	//get MAC Address from message
	// Split the string by the underscore character
	parts := strings.Split(eventMessage.CameraId, "_")
	if len(parts) < 1 {
		return fmt.Errorf("invalid camera ID: %s", eventMessage.CameraId)
	}

	macPartIndex := len(parts) - 1
	cameraMacAddress, err := formatMACAddress(parts[macPartIndex])
	if err != nil {
		return err
	}

	//get real camera ID from mac address
	var cameraID string = ""
	var cameraName string = ""

	dtoCamera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("mac_address = ?", cameraMacAddress)
	if err == nil {
		cameraID = dtoCamera.ID.String()
		cameraName = dtoCamera.Name
	}

	// Get lat, long from event message location
	lat, long := 0.0, 0.0
	locations := strings.Split(eventMessage.Location, ";")
	if len(locations) >= 2 {
		lat, err = strconv.ParseFloat(locations[0], 64)
		if err != nil {
			log.Println("error converting Lat of camera  to float64: %w, set default lat for that", err)
			lat = 21.028511
		}
		long, err = strconv.ParseFloat(locations[1], 64)
		if err != nil {
			log.Println("error converting Lat of camera  to float64: %w, set default long for that", err)
			long = 105.804817
		}
	} else {
		lat = 21.028511
		long = 105.804817
	}

	// Get blacklist's name from memberID if it's blacklist event
	if eventMessage.EventType == "AI_EVENT_BLACKLIST_FACE_RECOGNITION" {
		// Get blacklist's name from memberID
		blacklistName := ""
		if eventMessage.MemberID != "" {
			blacklistData, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_BlackList, models.BlackList]("message_id = ?", eventMessage.MemberID)
			if err == nil {
				blacklistName = blacklistData.Name
			}
		}
		eventMessage.Result = blacklistName
	}

	// Convert timestamp to datetime
	timestampMs := int64(eventMessage.Timestamp)
	seconds := timestampMs / 1000
	nanoseconds := (timestampMs % 1000) * 1000000
	datetime := time.Unix(seconds, nanoseconds)

	// Create AI warning DTO
	dtoAIWaring.MessageID = eventMessage.ID

	dtoAIWaring.MsVersion = eventMessage.MsVersion
	dtoAIWaring.SensorID = eventMessage.SensorID
	dtoAIWaring.Description = eventMessage.Description
	dtoAIWaring.Timestamp = eventMessage.Timestamp
	dtoAIWaring.TimeStart = eventMessage.TimeStart
	dtoAIWaring.TimeEnd = eventMessage.TimeEnd
	dtoAIWaring.Image = eventMessage.Image
	dtoAIWaring.ImageResult = eventMessage.ImageResult
	dtoAIWaring.Video = eventMessage.Video
	dtoAIWaring.StorageBucket = eventMessage.StorageBucket
	dtoAIWaring.EventType = eventMessage.EventType
	dtoAIWaring.Location = eventMessage.CameraId
	dtoAIWaring.CamIP = eventMessage.CamIP
	dtoAIWaring.CamName = cameraName
	dtoAIWaring.CameraId = cameraID
	dtoAIWaring.MemberID = eventMessage.MemberID
	dtoAIWaring.Result = eventMessage.Result
	dtoAIWaring.Latitude = lat
	dtoAIWaring.Longtitude = long
	dtoAIWaring.ImageObject = eventMessage.ImageObject
	dtoAIWaring.EventTypeString = ConvertAIEventString(eventMessage.EventType)
	dtoAIWaring.Status = "NEW"
	dtoAIWaring.ConverTimestamp = datetime

	// Use `CreateItemWithPartitionFromDTO` to handle partition creation and data insertion
	_, errCreate := reposity.CreateItemWithPartitionFromDTO[models.DTO_AIWaring, models.AIWaring](dtoAIWaring, "ai_warning")
	if errCreate != nil {
		log.Println("Error inserting AI warning:", errCreate)
		return errCreate
	}
	//log.Println("Image: ", dtoAIWaring.Image, dtoAIWaring.ImageResult, dtoAIWaring.ImageObject)
	image, _ := minioclient.GetPresignedURL(dtoAIWaring.StorageBucket, dtoAIWaring.Image)
	imageObject, _ := minioclient.GetPresignedURL(dtoAIWaring.StorageBucket, dtoAIWaring.ImageObject)
	imageResult, _ := minioclient.GetPresignedURL(dtoAIWaring.StorageBucket, dtoAIWaring.ImageResult)

	dtoAIWaring.Image = image
	dtoAIWaring.ImageObject = imageObject
	dtoAIWaring.ImageResult = imageResult

	wssignaling.SendNotifyMessage("list_ai_event", "list_ai_event", dtoAIWaring)
	return nil
}

func GetPartitionedTableName(timestamp time.Time) string {
	year, month := timestamp.Year(), timestamp.Month()
	return fmt.Sprintf("ai_warning_%d%02d", year, int(month))
}

func StartCronJob() {
	// Initialize the scheduler
	scheduler := gocron.NewScheduler(time.UTC)

	// Schedule RemoveDuplicateAIWarnings to run every 5 minutes
	scheduler.Every(5).Minutes().Do(func() {
		partitionedTableName := GetPartitionedTableName(time.Now())
		err := reposity.RemoveDuplicateAIWarnings(partitionedTableName)
		if err != nil {
			log.Println("Error in RemoveDuplicateAIWarnings:", err)
		}
	})

	// Start the scheduler
	scheduler.StartAsync()
}
