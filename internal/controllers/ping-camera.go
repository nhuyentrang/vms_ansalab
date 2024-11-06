package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"vms/internal/models"
	"vms/wssignaling"

	"vms/comongo/kafkaclient"
	"vms/comongo/reposity"

	"github.com/google/uuid"
)

func processPingCamera(eventMessage models.KafkaJsonVMSMessage) {
	// Create a new slice to hold modified camera statuses
	var modifiedCameraStatus []models.DTO_CameraInfo

	// Process LisCameraStatus
	if len(eventMessage.PayLoad.ListCameraStatus) > 0 {
		//log.Println("Len camera: ", len(eventMessage.PayLoad.ListCameraStatus))
		for _, value := range eventMessage.PayLoad.ListCameraStatus {
			// Update device status in the database
			updateDeviceStatus("camera", value.ID, value.Status)
			dtoCameraInfo := value

			// Update the status value based on the status flag
			if value.Status {
				dtoCameraInfo.StatusValue.Name = "connected"
			} else {
				dtoCameraInfo.StatusValue.Name = "disconnected"
			}
			modifiedCameraStatus = append(modifiedCameraStatus, dtoCameraInfo)
		}
		// Send a notification message with the updated camera statuses
		wssignaling.SendNotifyMessage("cameraStatus", "cameraStatus", modifiedCameraStatus)
	}

	// Create a new slice to hold modified NVR statuses
	var modifiedNVRStatus []models.DTO_NVRInfo

	// Process ListNVRStatus
	if len(eventMessage.PayLoad.ListNVRStatus) > 0 {
		//log.Println("Len nvr: ", len(eventMessage.PayLoad.ListNVRStatus))
		for _, value := range eventMessage.PayLoad.ListNVRStatus {
			// Update device status in the database
			updateDeviceStatus("nvr", value.ID, value.Status)
			dtoNVRInfo := value

			// Update the status value based on the status flag
			if value.Status {
				dtoNVRInfo.StatusValue.Name = "connected"
			} else {
				dtoNVRInfo.StatusValue.Name = "disconnected"
			}
			modifiedNVRStatus = append(modifiedNVRStatus, dtoNVRInfo)
		}
		// Send a notification message with the updated NVR statuses
		wssignaling.SendNotifyMessage("nvrStatus", "nvrStatus", modifiedNVRStatus)
	}

	// Process SmartNVR status
	var modifiedSmartNVRStatus models.DTO_Device
	if eventMessage.PayLoad.CommandID != "" {
		// Update device status in the database
		updateDeviceStatus("smartnvr", eventMessage.PayLoad.CommandID, true)
		commandIDStr := eventMessage.PayLoad.CommandID
		commandID, err := uuid.Parse(commandIDStr)
		if err != nil {
			log.Printf("Failed to parse CommandID: %v", err)
			return
		}
		modifiedSmartNVRStatus.ID = commandID

		//TODO: Check Status
		// if value.Status {
		// 	modifiedSmartNVRStatus.StatusValue.Name = "connected"
		// } else {
		// 	modifiedSmartNVRStatus.StatusValue.Name = "disconnected"
		// }

		// You can apply any modification logic here
		//wssignaling.SendNotifyMessage("nvrStatus", "nvrStatus", modifiedNVRStatus)
	}
}

func updateDeviceStatus(deviceType, id string, status bool) {
	var eventType, eventName, statusStr, types string

	// Determine event type, event name, status string, and type based on the status flag
	if status {
		eventType = "connected"
		eventName = "Đã kết nối"
		statusStr = "Đã xử lý"
		types = "Active"
	} else {
		eventType = "disconnected"
		eventName = "Mất kết nối"
		statusStr = "Chưa xử lý"
		types = "Deactive"
	}

	switch deviceType {
	case "camera":
		// Read camera information from the database
		camera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", id)
		if err != nil {
			log.Printf("Failed to read camera with ID %s: %v", id, err)
			return
		}
		// Update camera status
		camera.Status.Name = eventType
		if camera.LastPing.IsZero() || status {
			camera.LastPing = time.Now()
		}
		reposity.UpdateSingleColumn[models.Camera](id, "status", camera.Status)
		reposity.UpdateSingleColumn[models.Camera](id, "last_ping", camera.LastPing)

		// Track processed BoxIDs
		processedBoxIDs := make(map[string]bool)

		// Check if last ping is more than an hour ago
		if time.Since(camera.LastPing) > time.Hour {
			log.Printf("Camera with ID %s has not pinged for over an hour, synchronizing data for all cameras", id)

			// Fetch all cameras from the database
			allCameras, _, err := reposity.ReadAllItemsIntoDTO[models.DTOCamera, models.Camera]("-created_at")
			if err != nil {
				log.Printf("Failed to retrieve all cameras: %v", err)
				return
			}

			// Collect unique BoxIDs and update LastPing time for all cameras
			for _, cam := range allCameras {
				boxID := cam.Box.ID
				if _, processed := processedBoxIDs[boxID]; !processed {
					err := SynchronizeCameraDataConfigLogic(boxID)
					if err != nil {
						log.Printf("Failed to synchronize camera data config for BoxID %s: %v", boxID, err)
					}
					processedBoxIDs[boxID] = true
				}

				// Update LastPing time for the current camera
				cam.LastPing = time.Now()
				reposity.UpdateSingleColumn[models.Camera](cam.ID.String(), "last_ping", cam.LastPing)
			}
		}

		// Check if incident exists
		query := reposity.NewQuery[models.DTO_System_Incident_BasicInfo, models.SystemIncident]()
		query.AddConditionOfTextField("AND", "device_id", "=", id)
		incident, _, err := query.ExecWithPaging("-created_at", 1, 1)
		if err != nil {
			log.Printf("Failed to query incident for camera with ID %s: %v", id, err)
			return
		}

		if len(incident) > 0 {
			// If incident type is different from current type and event type is "disconnected"
			if incident[0].Type != types {
				if eventType == "disconnected" {
					dtoSystemIncident := models.DTO_System_Incident_BasicInfo{
						DeviceType: deviceType,
						DeviceID:   id,
						EventType:  eventType,
						EventName:  eventName,
						Status:     statusStr,
						Type:       types,
						Location:   camera.Location,
						Severity:   "Trung bình",
						Source:     camera.Name,
					}
					NotifiToTelegram(camera.Name, eventName, "camera")
					reposity.CreateItemFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](dtoSystemIncident)
				}
			}
			// If current type is "Active" and event type is "connected"
			if types == "Active" {
				if eventType == "connected" {
					dtoSystemIncident := models.DTO_System_Incident_BasicInfo{
						DeviceType: deviceType,
						DeviceID:   id,
						EventType:  eventType,
						EventName:  eventName,
						Status:     statusStr,
						Type:       types,
						Location:   camera.Location,
						Severity:   "Trung bình",
						Source:     camera.Name,
					}

					NotifiToTelegram(camera.Name, eventName, "camera")
					reposity.UpdateItemByIDFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](incident[0].ID.String(), dtoSystemIncident)
				}
			}
		} else {
			// If no incident exists and event type is "disconnected"
			if eventType == "disconnected" {
				dtoSystemIncident := models.DTO_System_Incident_BasicInfo{
					DeviceType: deviceType,
					DeviceID:   id,
					EventType:  eventType,
					EventName:  eventName,
					Status:     statusStr,
					Type:       types,
					Location:   camera.Location,
					Severity:   "Trung bình",
					Source:     camera.Name,
				}
				NotifiToTelegram(camera.Name, eventName, "camera")
				reposity.CreateItemFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](dtoSystemIncident)
			}
		}
	case "nvr":
		// Read NVR information from the database
		nvr, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVR, models.NVR]("id = ?", id)
		if err != nil {
			log.Printf("Failed to read NVR with ID %s: %v", id, err)
			return
		}
		// Update NVR status
		nvr.Status.Name = eventType
		if nvr.LastPing.IsZero() || status {
			nvr.LastPing = time.Now()
		}
		reposity.UpdateSingleColumn[models.NVR](id, "status", nvr.Status)
		reposity.UpdateSingleColumn[models.NVR](id, "last_ping", nvr.LastPing)

		// Track processed BoxIDs
		processedBoxIDs := make(map[string]bool)

		// Check if last ping is more than an hour ago
		if time.Since(nvr.LastPing) > time.Hour {
			log.Printf("NVR with ID %s has not pinged for over an hour, synchronizing data for all NVRs", id)

			// Fetch all NVRs from the database
			allNVRs, _, err := reposity.ReadAllItemsIntoDTO[models.DTO_NVR, models.NVR]("-created_at")
			if err != nil {
				log.Printf("Failed to retrieve all NVRs: %v", err)
				return
			}

			// Collect unique BoxIDs and update LastPing time for all NVRs
			for _, nvr := range allNVRs {
				boxID := nvr.Box.ID
				if _, processed := processedBoxIDs[boxID]; !processed {
					err := SynchronizeNVRDataConfigLogic(boxID)
					if err != nil {
						log.Printf("Failed to synchronize NVR data config for BoxID %s: %v", boxID, err)
					}
					processedBoxIDs[boxID] = true
				}

				// Update LastPing time for the current NVR
				nvr.LastPing = time.Now()
				reposity.UpdateSingleColumn[models.NVR](nvr.ID.String(), "last_ping", nvr.LastPing)
			}
		}

		// Check if incident exists
		query := reposity.NewQuery[models.DTO_System_Incident_BasicInfo, models.SystemIncident]()
		query.AddConditionOfTextField("AND", "device_id", "=", id)
		incident, _, err := query.ExecWithPaging("-created_at", 1, 1)
		if err != nil {
			log.Printf("Failed to query incident for NVR with ID %s: %v", id, err)
			return
		}

		if len(incident) > 0 {
			// If incident type is different from current type and event type is "disconnected"
			if incident[0].Type != types {
				if eventType == "disconnected" {
					dtoSystemIncident := models.DTO_System_Incident_BasicInfo{
						DeviceType: deviceType,
						DeviceID:   id,
						EventType:  eventType,
						EventName:  eventName,
						Status:     statusStr,
						Type:       types,
						Location:   nvr.Location,
						Severity:   "Trung bình",
						Source:     nvr.Name,
					}
					NotifiToTelegram(nvr.Name, eventName, "NVR")
					reposity.CreateItemFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](dtoSystemIncident)
				}
			}
			// If current type is "Active" and event type is "connected"
			if types == "Active" {
				if eventType == "connected" {
					dtoSystemIncident := models.DTO_System_Incident_BasicInfo{
						DeviceType: deviceType,
						DeviceID:   id,
						EventType:  eventType,
						EventName:  eventName,
						Status:     statusStr,
						Type:       types,
						Location:   nvr.Location,
						Severity:   "Trung bình",
						Source:     nvr.Name,
					}
					NotifiToTelegram(nvr.Name, eventName, "NVR")
					reposity.UpdateItemByIDFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](incident[0].ID.String(), dtoSystemIncident)
				}
			}
		} else {
			// If no incident exists and event type is "disconnected"
			if eventType == "disconnected" {
				dtoSystemIncident := models.DTO_System_Incident_BasicInfo{
					DeviceType: deviceType,
					DeviceID:   id,
					EventType:  eventType,
					EventName:  eventName,
					Status:     statusStr,
					Type:       types,
					Location:   nvr.Location,
					Severity:   "Trung bình",
					Source:     nvr.Name,
				}
				NotifiToTelegram(nvr.Name, eventName, "NVR")
				reposity.CreateItemFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](dtoSystemIncident)
			}
		}
	case "smartnvr":
		// Fetch all devices from the database
		devices, count, err := reposity.ReadAllItemsIntoDTO[models.DTO_Device, models.Device]("-created_at")
		if err != nil {
			log.Printf("Failed to read devices: %v", err)
			return
		}

		if count == 0 {
			log.Println("No devices found.")
			return
		}

		pingDevice, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Device, models.Device](" model_id = ?", id)
		if err != nil {
			log.Printf("Failed to read NVR with ID %s: %v", id, err)
			return
		}

		currentTime := time.Now()

		for _, device := range devices {
			var eventType, eventName, statusStr, types string

			// Check the time since the last ping for each device
			if device.LastPing.IsZero() || currentTime.Sub(device.LastPing) <= 5*time.Minute {
				// If last ping is within 5 minutes or it's zero, update the LastPing time
				device.LastPing = currentTime
				eventType = "connected"
				eventName = "Đã kết nối"
				statusStr = "Đã xử lý"
				types = "Active"
			} else {
				// If last ping is more than 5 minutes ago, set status to disconnected
				eventType = "disconnected"
				eventName = "Mất kết nối"
				statusStr = "Chưa xử lý"
				types = "Deactive"

				// Check if incident exists
				query := reposity.NewQuery[models.DTO_System_Incident_BasicInfo, models.SystemIncident]()
				query.AddConditionOfTextField("AND", "device_id", "=", device.ID.String())
				incident, _, err := query.ExecWithPaging("-created_at", 1, 1)
				if err != nil {
					log.Printf("Failed to query incident for device with ID %s: %v", device.ID.String(), err)
					return
				}

				if len(incident) > 0 {
					// Check if the latest incident is not already "disconnected"
					if incident[0].EventType != eventType {
						NotifiToTelegram(device.NameDevice, eventName, "SmartNVR")
						dtoSystemIncident := models.DTO_System_Incident_BasicInfo{
							DeviceType: "smartnvr",
							DeviceID:   device.ID.String(),
							EventType:  eventType,
							EventName:  eventName,
							Status:     statusStr,
							Type:       types,
							Location:   device.Location,
							Severity:   "Trung bình",
							Source:     device.NameDevice,
						}
						reposity.UpdateItemByIDFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](incident[0].ID.String(), dtoSystemIncident)
					}
				} else {
					dtoSystemIncident := models.DTO_System_Incident_BasicInfo{
						DeviceType: "smartnvr",
						DeviceID:   device.ID.String(),
						EventType:  eventType,
						EventName:  eventName,
						Status:     statusStr,
						Type:       types,
						Location:   device.Location,
						Severity:   "Trung bình",
						Source:     device.NameDevice,
					}
					NotifiToTelegram(device.NameDevice, eventName, "SmartNVR")
					reposity.CreateItemFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](dtoSystemIncident)
				}
			}

			// Update device status
			device.Status = eventType
			if device.ID == pingDevice.ID {
				pingDevice.Status = eventType
				pingDevice.LastPing = currentTime
				reposity.UpdateSingleColumn[models.Device](pingDevice.ID.String(), "status", pingDevice.Status)
				reposity.UpdateSingleColumn[models.Device](pingDevice.ID.String(), "last_ping", pingDevice.LastPing)
			} else {
				reposity.UpdateSingleColumn[models.Device](device.ID.String(), "status", device.Status)
				// if device.LastPing == currentTime {
				// 	reposity.UpdateSingleColumn[models.Device](device.ID.String(), "last_ping", device.LastPing)
				// }
			}
		}
	}
}

func SynchronizeCameraDataConfigLogic(modelID string) error {
	log.Printf("Received request to synchronize camera data config for model ID: %s", modelID)

	// Fetch device by model ID
	device, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Device, models.Device]("model_id = ?", modelID)
	if err != nil {
		log.Printf("Device not found for model ID: %s, error: %v", modelID, err)
		return fmt.Errorf("device not found: %v", err)
	}
	log.Printf("Device found: %+v", device)

	log.Printf("Updating status for Device with model ID: %s", modelID)
	updateDeviceStatus("smartnvr", device.ModelID, true)

	// Fetch all cameras with the device ID
	query := reposity.NewQuery[models.DTOCamera, models.Camera]()
	jsonCondition := fmt.Sprintf("{\"id\":\"%s\"}", device.ID)
	query.AddConditionOfTextField("AND", "box", "@>", jsonCondition)
	cameras, count, err := query.ExecNoPaging("-created_at")
	if err != nil {
		log.Printf("Failed to retrieve cameras for device ID: %s, error: %v", device.ID.String(), err)
		return fmt.Errorf("failed to retrieve cameras: %v", err)
	}

	if count == 0 {
		log.Printf("No cameras found for device ID: %s", device.ID.String())
		return fmt.Errorf("no cameras found for the device")
	}
	log.Printf("Found %d cameras for device ID: %s", count, device.ID.String())

	for _, camera := range cameras {
		log.Printf("Processing camera: %+v", camera)
		// Create a map for camera channels
		cameraMap := make(map[string]models.ChannelCamera)
		for _, stream := range camera.Streams {
			channel := models.ChannelCamera{
				OnDemand: true,
				Url:      stream.URL,
				Codec:    stream.Codec,
				Name:     stream.Name,
			}
			cameraMap[stream.Channel] = channel
		}
		dataFileConfig := map[string]models.ConfigCamera{}
		newCamera := models.ConfigCamera{
			NameCamera: camera.Name,
			IP:         camera.IPAddress,
			UserName:   camera.Username,
			PassWord:   camera.Password,
			HTTPPort:   camera.HttpPort,
			RTSPPort:   camera.HttpPort,
			OnvifPort:  camera.OnvifPort,
			ChannelNVR: camera.NVR.Channel,
			IDNVR:      camera.NVR.ID,
			Channels:   cameraMap,
		}
		dataFileConfig[camera.ID.String()] = newCamera

		// Send command to synchronize camera data config
		requestUUID := uuid.New()
		cmdAddConfig := models.DeviceCommand{
			CommandID:    device.ModelID,
			Cmd:          cmd_AddDataConfig,
			EventTime:    time.Now().Format(time.RFC3339),
			EventType:    "camera",
			ConfigCamera: dataFileConfig,
			ProtocolType: camera.Protocol,
			RequestUUID:  requestUUID,
		}
		cmsStr, err := json.Marshal(cmdAddConfig)
		if err != nil {
			log.Printf("Failed to marshal command for camera ID: %s, error: %v", camera.ID, err)
			return fmt.Errorf("failed to marshal command: %v", err)
		}
		log.Printf("Sending command to Kafka: %s", cmsStr)
		kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	}

	log.Println("Successfully synchronized camera data config")
	return nil
}

func SynchronizeNVRDataConfigLogic(modelID string) error {
	log.Printf("Received request to synchronize NVR data config for model ID: %s", modelID)

	// Fetch device by model ID
	device, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Device, models.Device]("model_id = ?", modelID)
	if err != nil {
		log.Printf("Device not found for model ID: %s, error: %v", modelID, err)
		return fmt.Errorf("device not found: %w", err)
	}
	log.Printf("Device found: %+v", device)
	log.Printf("Updating status for Device with model ID: %s", modelID)
	updateDeviceStatus("smartnvr", device.ModelID, true)

	// Fetch all NVRs with the device ID
	query := reposity.NewQuery[models.DTO_NVR, models.NVR]()
	jsonCondition := fmt.Sprintf("{\"id\":\"%s\"}", device.ID.String())
	query.AddConditionOfTextField("AND", "box", "@>", jsonCondition)
	nvrs, count, err := query.ExecNoPaging("-created_at")
	if err != nil {
		log.Printf("Failed to retrieve NVRs for device ID: %s, error: %v", device.ID.String(), err)
		return fmt.Errorf("failed to retrieve NVRs: %w", err)
	}

	if count == 0 {
		log.Printf("No NVRs found for device ID: %s", device.ID.String())
		return fmt.Errorf("no NVRs found for the device")
	}
	log.Printf("Found %d NVRs for device ID: %s", count, device.ID.String())

	for _, nvr := range nvrs {
		log.Printf("Processing NVR: %+v", nvr)
		// Create a map for nvr's cameras

		cameras := models.KeyValueArray{}
		for _, camera := range *nvr.Cameras {
			cameras = append(cameras, models.KeyValue{
				ID:      camera.ID,
				Name:    camera.Name,
				Channel: camera.Channel,
			})
		}

		// Create a map of camera for NVR
		dataFileConfig := map[string]models.ConfigNVR{}
		newNVR := models.ConfigNVR{
			NameCamera: nvr.Name,
			IP:         nvr.IPAddress,
			UserName:   nvr.Username,
			PassWord:   nvr.Password,
			HTTPPort:   nvr.HttpPort,
			RTSPPort:   nvr.HttpPort, // Assuming RTSPPort is same as HttpPort
			OnvifPort:  nvr.OnvifPort,
			Cameras:    cameras,
		}
		dataFileConfig[nvr.ID.String()] = newNVR

		// Send command to synchronize NVR data config
		requestUUID := uuid.New()
		cmdAddConfig := models.DeviceCommand{
			CommandID:    device.ModelID,
			Cmd:          cmd_AddDataConfig,
			EventTime:    time.Now().Format(time.RFC3339),
			EventType:    "nvr",
			ConfigNVR:    dataFileConfig,
			ProtocolType: nvr.Protocol,
			RequestUUID:  requestUUID,
		}
		cmsStr, err := json.Marshal(cmdAddConfig)
		if err != nil {
			log.Printf("Failed to marshal command for NVR ID: %s, error: %v", nvr.ID.String(), err)
			return fmt.Errorf("failed to marshal command: %w", err)
		}
		log.Printf("Sending command to Kafka: %s", cmsStr)
		kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	}

	log.Println("Successfully synchronized NVR data config")
	return nil
}
