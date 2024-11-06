package controllers

import (
	"encoding/json"
	"log"
	"time"

	"vms/internal/models"

	"vms/comongo/kafkaclient"
)

const MAX_MAP_MESSAGE int = 8192

var messageMap = NewSafeOrderedMap()

var (
	M_deviceCommand string = ""

	ChannelDataReceiving                          chan *models.KafkaJsonVMSMessage
	DataStorageReceiving                          chan *models.KafkaJsonVMSMessage
	ChannelDeleteFileStreamReceiving              chan *models.KafkaJsonVMSMessage
	DeviceScanOnvifChannelDataReceiving           chan *models.KafkaJsonVMSMessage
	DeviceScanStaticIPChannelDataReceiving        chan *models.KafkaJsonVMSMessage
	GetNetworkConfigChannelDataReceiving          chan *models.KafkaJsonVMSMessage
	DeviceScanIPListChannelDataReceiving          chan *models.KafkaJsonVMSMessage
	UpdateImageConfigCameraChannelDataReceiving   chan *models.KafkaJsonVMSMessage //Unused
	ChannelDataConfig                             chan *models.KafkaJsonVMSMessage
	ChangePasswordChannelDataReceiving            chan *models.KafkaJsonVMSMessage
	GetCalenderPlaybackChannelDataReceiving       chan *models.KafkaJsonVMSMessage
	ChangePasswordSeriesChannelDataReceiving      chan *models.KafkaJsonVMSMessage
	AddCameraToNVRChannelDataReceiving            chan *models.KafkaJsonVMSMessage
	UpdateNetworkConfigCameraChannelDataReceiving chan *models.KafkaJsonVMSMessage
	GetImageConfigChannelDataReceiving            chan *models.KafkaJsonVMSMessage
	EditIPandPortHTTPChannelDataReceiving         chan *models.KafkaJsonVMSMessage
	//DownLoadVideoChannelDataReceiving             chan *models.KafkaJsonVMSMessage
	GetVideoConfigChannelDataReceiving chan *models.KafkaJsonVMSMessage
	DownloadClipChannelDataReceiving   chan *models.KafkaJsonVMSMessage
	ExtractChipChannelDataReceiving    chan *models.KafkaJsonVMSMessage
)

func RunDeviceEventProcessing(
	kafkaConsumerName string,
	deviceCommandivi string) {

	// Save topic name
	M_deviceCommand = deviceCommandivi

	// Init channels
	ChannelDeleteFileStreamReceiving = make(chan *models.KafkaJsonVMSMessage, 10)              // buffers 10 values without blocking.
	ChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)                          // buffers 10 values without blocking.
	ChannelDataConfig = make(chan *models.KafkaJsonVMSMessage, 10)                             // buffers 10 values without blocking.
	DeviceScanStaticIPChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)        // buffers 10 values without blocking.
	DeviceScanOnvifChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)           // buffers 10 values without blocking.
	UpdateImageConfigCameraChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)   // buffers 10 values without blocking.
	UpdateNetworkConfigCameraChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10) // buffers 10 values without blocking.
	DataStorageReceiving = make(chan *models.KafkaJsonVMSMessage, 10)                          // buffers 10 values without blocking.
	GetNetworkConfigChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)          // buffers 10 values without blocking.
	GetImageConfigChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)            // buffers 10 values without blocking.
	ChangePasswordSeriesChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)      // buffers 10 values without blocking.
	GetCalenderPlaybackChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)       // buffers 10 values without blocking.
	DeviceScanIPListChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)          // buffers 10 values without blocking.
	GetVideoConfigChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)            // buffers 10 values without blocking.
	ChangePasswordChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)            // buffers 10 values without blocking.
	AddCameraToNVRChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)            // buffers 10 values without blocking.
	EditIPandPortHTTPChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)         // buffers 10 values without blocking.
	DownloadClipChannelDataReceiving = make(chan *models.KafkaJsonVMSMessage, 10)              // buffers 10 values without blocking.

	// Loop for processing consumed Kafka messages
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
					// Optionally, send a response or notification about the busy state
					continue
				}

				// Process message
				err = ProcessDeviceEventData(eventMessage)
				elapsed := time.Since(start)
				if err != nil {
					log.Printf("==============> Failed to process kafka message, error: %s, took %s\n", err, elapsed)
					continue
				}
			}
		}
	}()
}

// Process data from device event topic
func ProcessDeviceEventData(eventMessage models.KafkaJsonVMSMessage) error {
	//log.Println("CMD: ", eventMessage.PayLoad.Cmd)
	switch eventMessage.PayLoad.Cmd {

	case cmd_PingCamera:
		processPingCamera(eventMessage)

	}
	return nil
}
