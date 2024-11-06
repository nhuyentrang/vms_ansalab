package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"vms/internal/models"

	"vms/comongo/kafkaclient"

	"gorm.io/gorm"

	"vms/comongo/reposity"
)

var (
	M_deviceCommandFaceRegister string = ""
)

func RunDeviceDMPEventProcessing(
	kafkaConsumerName string,
	deviceCommandiviFaceRegister string) {

	// Save topic name
	M_deviceCommandFaceRegister = deviceCommandiviFaceRegister

	// Loop for process consumed kafka message
	go func() {
		for {
			// Waiting for msg in channel
			msg, topic, timestamp, err := kafkaclient.ConsumerReadMessage(kafkaConsumerName)
			if err != nil {
				continue
			}
			if msg != "" && topic != "" {
				//log.Println("Received kafka message from topic ", topic)
				start := time.Now()
				// Extract data from message
				var eventMessage models.KafkaJsonVMSMessage
				err := json.Unmarshal([]byte(msg), &eventMessage)
				if err != nil {
					log.Printf("==============> Failed to decode kafka message, error: %s, msg: %s, topic: %s, record timestamp: %s\n", err, msg, topic, timestamp.Format(time.RFC3339))
					continue
				}

				// Store message to map, for scanning device later
				if messageMap.Count() < MAX_MAP_MESSAGE {
					messageMap.Store(eventMessage.PayLoad.RequestUUID, &eventMessage)
				} else {
					log.Printf("==============> System is busy, cannot process new messages, reached max limit: %d\n", MAX_MAP_MESSAGE)
					continue
				}

				// Process data
				err = ProcessDeviceDMP(eventMessage)
				elapsed := time.Since(start)
				if err != nil {
					log.Printf("==============> Failed to process kafka message, error: %s, took %s\n", err, elapsed)
					continue
				}
			}
		}
	}()
}

// Process data from DMP event topic
func ProcessDeviceDMP(eventMessage models.KafkaJsonVMSMessage) error {
	if eventMessage.DeviceInfo != nil {
		dataRespDevice := models.DTO_Device{
			DeviceType: eventMessage.DeviceInfo.DeviceType,
			ModelID:    eventMessage.DeviceInfo.ID,
			DeviceCode: eventMessage.DeviceInfo.DeviceCode,
			Status:     eventMessage.DeviceInfo.Status,
			IPAddress:  eventMessage.DeviceInfo.IPAddress,
			MacAddress: eventMessage.DeviceInfo.MacAddress,
			NameDevice: "SmartBox_" + eventMessage.DeviceInfo.MacAddress,
		}
		dto, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Device, models.Device]("mac_address = ?", eventMessage.DeviceInfo.MacAddress)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// If the error is "record not found"
				dto, err = reposity.CreateItemFromDTO[models.DTO_Device, models.Device](dataRespDevice)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// Update the existing item
			dto, err = reposity.UpdateItemByIDFromDTO[models.DTO_Device, models.Device](dto.ID.String(), dataRespDevice)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
