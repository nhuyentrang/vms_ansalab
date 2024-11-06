package services

import (
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"
	"vms/internal/models"
	"vms/wsnode"

	"vms/comongo/reposity"

	"github.com/google/uuid"
)

/// Todo: when service first startup, it should read all camera ai property and send to aieagent
/// Todo: when service first startup, it should re-init all aieagent connection state to offline, before start wsnode to accept connection from aieagent

type agentInfo struct {
	Topic string
	ID    string
	State string
}

// Create map between clientID and agentID
var subscribedAgentList = make(map[string]agentInfo)

// Func send message to aieagent
func SendMsgToAIEAgent(topic string, action string, data string) error {
	// Sanity check
	if topic == "" {
		return errors.New("topic is empty")
	}

	msg := models.WSMsgAIEngine{
		ID:        uuid.New(),
		Timestamp: time.Now(),
		SenderID:  "VMS",
		//Type:      models.WSMsgTypeResponseFromVMS,
		Action: action,
		Data:   data,
	}
	packet, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	err = wsnode.Publish(topic, packet)
	return err
}

func WSHanleOnSubscribeAIEAgent(clientID string, topic string) {
	log.Printf("AIEngine %s subscribes on %s", clientID, topic)

	// Store this client to map
	subscribedAgentList[clientID] = agentInfo{
		Topic: topic,
		ID:    "",
		State: "connected",
	}
}

func WSHandleOnDisconnectAIEAgent(clientID string, code string) {
	log.Printf("AIEngine %s disconnected, disconnect: %s", clientID, code)
	// Update connection state of this agent in database
	agentInfo, ok := subscribedAgentList[clientID]
	if !ok {
		return
	}
	if agentInfo.ID != "" && agentInfo.State == "online" {
		UpdateAIEngineOnlineStatus(agentInfo.ID, false)
	}

	// Remove this agent from map
	delete(subscribedAgentList, clientID)
}

// ProcessingEventWebsocketAIEAgent is a function to process event from aieagent
func WSHandleOnPublishAIEAgent(clientID string, topic string, data []byte) {
	var msg models.WSMsgAIEngine
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return
	}

	// We do not process message without senderID or senderID is VMS itself
	if msg.SenderID == "" || msg.SenderID == "VMS" {
		return
	}

	// Process ai-engine request
	//if msg.Type == models.WSMsgTypeRequestFromAIEngine {
	//log.Printf("Get request from AIEngine for action %s, data: %s", msg.Action, msg.Data)
	switch msg.Action {
	case models.AIEngineActionRegisterAgent:
		aieID, err := RegisterAIEngine(msg.Data, "1.0")
		if err != nil {
			SendMsgToAIEAgent(topic, msg.Action, "error: "+err.Error())
			log.Printf("Failed to register aieagent: %s", err)
		} else if aieID != "" {
			SendMsgToAIEAgent(topic, msg.Action, "success")
			// Store this client to map
			subscribedAgentList[clientID] = agentInfo{
				Topic: topic,
				ID:    aieID,
				State: "online",
			}

			log.Printf("Registered aieagent: %s, machine id: %s, preparing to sync ai config", aieID, msg.Data)
			// Todo: we should sync ai config after agent is registered

			// Currently, only sync if agent's machine id = "e6d6e90caab84d3475a78bc17b589fc907277f4bdf163e15129d24eb485ec039"
			// This is a test agent
			if msg.Data != "e6d6e90caab84d3475a78bc17b589fc907277f4bdf163e15129d24eb485ec039" {
				return
			}

			// Query all camera ai property and send to aieagent
			aiCamPropDTOs, count, err := reposity.ReadAllItemsIntoDTO[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty]("id")
			if err != nil {
				log.Printf("Failed to read camera ai property: %s", err)
				return
			}
			if count > 0 {
				wsAAICamConfigs := models.WSAICamProperties{}
				wsAAICamConfigs.WSAICamProperties = make([]models.WSAICamProperty, 0)
				for _, dto := range aiCamPropDTOs {
					camID := dto.CameraID.String()

					// Todo: check if this camera is belong to this aieagent, if not, skip this camera
					// Get info fromt this camera
					cam, err := reposity.ReadItemByIDIntoDTO[models.Camera, models.Camera](camID)
					if err != nil {
						continue
					}
					aiVideoStream, err := GetAIVideoStream(cam)
					if err != nil {
						continue
					}

					// Validate coordinate
					coordinate := "0.0;0.0;0.0"
					camCoordinates := strings.Split(cam.Coordinate, ";")
					if len(camCoordinates) == 3 {
						coordinate = cam.Coordinate
					}

					// Append to list
					wsAAICamConfigs.WSAICamProperties = append(wsAAICamConfigs.WSAICamProperties, models.WSAICamProperty{
						CameraName:       cam.Name,
						Description:      cam.Description,
						CameraMacAddress: cam.MACAddress,
						AICamProperty:    dto,
						VideoStream:      aiVideoStream,
						Location:         strings.Join([]string{cam.Lat, cam.Long, "0.0"}, ";"), // Locaction=lat;long;altitude
						Coordinate:       coordinate,
					})
				}

				if len(wsAAICamConfigs.WSAICamProperties) > 0 {
					msgData, err := json.Marshal(wsAAICamConfigs)
					if err != nil {
						log.Printf("Failed to marshal ai property for camera: %s", err)
						return
					}
					err = SendMsgToAIEAgent(topic, models.AIEngineActionSyncAIConfig, string(msgData))
					if err != nil {
						log.Printf("Failed to send sync ai property to aieagent: %s", err)
						return
					}
					log.Printf("Sync ai property to aieagent, data: %v", wsAAICamConfigs.WSAICamProperties)
				}
			}
		}
	case models.AIEngineActionReportUsage:
		// Check this agent is registered
		agentInfo, ok := subscribedAgentList[clientID]
		if !ok {
			err = SendMsgToAIEAgent(topic, msg.Action, "error: Agent is not registered")
			if err != nil {
				log.Printf("Failed to send ack to aieagent: %s", err)
			}
			return
		}
		if agentInfo.ID == "" || agentInfo.State != "online" {
			err = SendMsgToAIEAgent(topic, msg.Action, "error: Agent is not registered")
			if err != nil {
				log.Printf("Failed to send ack to aieagent: %s", err)
			}
			return
		}

		// Parse usage report
		var usageReport models.WSAIEngineUsageReport
		err := json.Unmarshal([]byte(msg.Data), &usageReport)
		if err != nil {
			log.Printf("Error when parsing usage report: %s", err)
			err = SendMsgToAIEAgent(topic, msg.Action, "error: "+err.Error())
			if err != nil {
				log.Printf("Failed to send ack to aieagent: %s", err)
			}
			return
		}
		err = UpdateAIEngineUsageReport(agentInfo.ID, usageReport)
		if err != nil {
			log.Printf("Failed to update usage report: %s", err)
			err = SendMsgToAIEAgent(topic, msg.Action, "error: "+err.Error())
			if err != nil {
				log.Printf("Failed to send ack to aieagent: %s", err)
			}
			return
		}
		// Send ack to aieagent
		err = SendMsgToAIEAgent(topic, msg.Action, "success")
		if err != nil {
			log.Printf("Failed to send ack to aieagent: %s", err)
		}
	case models.AIEngineActionSyncAIConfig:
		// AIEAgnent response after sync all ai cam properties
	case models.AIEngineActionUpdateAIConfig:
		// AIEAgnent response after update ai cam property
	case models.AIEngineActionUpdateAIModel:
		// AIEAgnent response after update ai model property
	default:
		log.Printf("Get request from AIEngine, unknown action: %s", msg.Action)
		return
	}
	//}
}

func RegisterAIEngine(machineID string, agentVersion string) (string, error) {

	// Check if this machine ID is existed in database, if not let create new one
	aie, err := reposity.ReadItemWithFilterIntoDTO[models.AIEngine, models.AIEngine]("machine_id = ?", machineID)
	if err != nil {
		// Create new aieagent
		aieNew := models.AIEngine{
			ID:            uuid.New(),
			MachineID:     machineID,
			LastMessageAt: time.Now(),
			ConnectedAt:   time.Now(),
			AgentVersion:  agentVersion,
			Online:        true,
			// Todo: update other info here (network, cpu, mem...)
		}
		aie, err := reposity.CreateItemFromDTO[models.AIEngine, models.AIEngine](aieNew)
		if err != nil {
			log.Printf("Failed to create aieagent: %s", err)
			return "", err
		}
		return aie.ID.String(), nil
	}

	// If existed, update the last active time
	aie.LastMessageAt = time.Now()
	aie.ConnectedAt = time.Now()
	aie.Online = true

	_, err = reposity.UpdateItemByIDFromDTO[models.AIEngine, models.AIEngine](aie.ID.String(), aie)
	if err != nil {
		//log.Printf("Failed to update register status for aieagent: %s, error: %s", aie.ID.String(), err.Error())
		return "", err
	}

	if aie.ID == uuid.Nil {
		return "", errors.New("aieagent ID is empty")
	}

	return aie.ID.String(), nil
}

// Func update connection state for agent
func UpdateAIEngineOnlineStatus(agentID string, isOnline bool) error {
	aieNew := models.AIEngine{
		Online:        isOnline,
		LastMessageAt: time.Now(),
	}

	// Update to database
	_, err := reposity.UpdateItemByIDFromDTO[models.AIEngine, models.AIEngine](agentID, aieNew)
	if err != nil {
		log.Printf("Failed to update online status for aieagent: %s", err)
		return err
	}

	return nil
}

// Func update usage report for agent
func UpdateAIEngineUsageReport(agentID string, usageReport models.WSAIEngineUsageReport) error {
	// Sanity check
	if agentID == "" {
		return errors.New("agentID is empty")
	}

	aieNew := models.AIEngine{
		LastMessageAt:      time.Now(),
		CPUUtilization:     usageReport.CPUUsage,
		MEMUtilization:     usageReport.MEMUsage,
		StorageUtilization: usageReport.StorageUsage,
		// Todo: update other info here (network, cpu, mem...)
	}

	// Update to database
	_, err := reposity.UpdateItemByIDFromDTO[models.AIEngine, models.AIEngine](agentID, aieNew)
	if err != nil {
		log.Printf("Failed to update usage report for aieagent: %s", err)
		return err
	}

	return nil
}

// Func update camera ai property
func WSNotifyAIEngineUpdateCameraAIProperty(camID uuid.UUID, aiCamProperty models.DTO_Camera_AI_Property_BasicInfo) error {
	if camID == uuid.Nil {
		// When FE call PUT api, there is no Camera ID in request, we need to read it from dtoAICamProperty
		// The input dtoAICamProperties does not contain full infomation from database, we need to read full info from database
		dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty](aiCamProperty.ID.String())
		if err != nil {
			return errors.New("failed to read camera ai property: " + err.Error())
		}
		camID = dto.CameraID
	}
	// Note: always using aiCamProperty from FE request to send to ai-engine, do not use dtoAICamProperty from database
	// Only when data is sent to ai-engine successfully, then api will update aiCamProperty to database

	// Update camera ai property
	log.Printf("Update ai property for camera: %s", camID)

	// Find aieagent which manage this camera
	var isFoundAIEngine bool = false
	cam, err := reposity.BackRefManyToManyRetrieve[models.Camera](camID.String(), "AIEngines")
	if err != nil {
		return errors.New("failed to search associated aieagent for camera: " + err.Error())
	}

	if cam.MACAddress == "" {
		return errors.New("this camera do not have MAC address")
	}

	if cam.AIEngines != nil {
		if len(cam.AIEngines) > 0 {
			isFoundAIEngine = true
		}
	}

	if !isFoundAIEngine {
		// Associate this camera to first aieagent
		// Todo: need to update this logic, should have a way to select aieagent automatically
		aie, cnt, err := reposity.ReadAllItemsIntoDTO[models.AIEngine, models.AIEngine]("-created_at")
		if err != nil || cnt == 0 {
			return errors.New("no aieagent found for this camera")
		}
		err = AssociateCameraToAIEngine(camID.String(), aie[0])
		if err != nil {
			return errors.New("failed to associate camera to aieagent: " + err.Error())
		}
		cam.AIEngines = make([]*models.AIEngine, 1)
		cam.AIEngines = append(cam.AIEngines, &aie[0])
		log.Printf("Associated camera to aieagent: %s", aie[0].ID.String())
	}

	// Get video stream of this camera
	// First stream is main stream, second stream is sub stream, and custom streams are after that
	// When retrive stream, should check if channel contain "main" or "Main", if not take the first stream
	aiVideoStream, err := GetAIVideoStream(cam)
	if err != nil {
		return errors.New("failed to get main video stream for camera: " + cam.ID.String() + ", error: " + err.Error())
	}

	if cam.AIEngines == nil {
		return errors.New("no aieagent can be associated with this camera")
	}

	// Collect data to build message and send to aieagent
	for _, aiengine := range cam.AIEngines {
		if aiengine == nil {
			// Actutually this should not happen, but unfortunately it happens
			// Todo: need to investigate why this happen
			continue
		}

		// Validate coordinate
		coordinate := "0.0;0.0;0.0"
		camCoordinates := strings.Split(cam.Coordinate, ";")
		if len(camCoordinates) == 3 {
			coordinate = cam.Coordinate
		}

		// Build message to send to aieagent
		wsAIamConfig := models.WSAICamProperty{
			CameraName:       cam.Name,
			Description:      cam.Description,
			CameraMacAddress: cam.MACAddress,
			AICamProperty:    aiCamProperty,
			VideoStream:      aiVideoStream,
			Location:         strings.Join([]string{cam.Lat, cam.Long, "0.0"}, ";"), // Locaction=lat;long;altitude
			Coordinate:       coordinate,
		}
		msgData, err := json.Marshal(wsAIamConfig)
		if err != nil {
			return errors.New("failed to marshal ai property for camera: " + cam.ID.String() + ", error: " + err.Error())
		}

		// Search map to get topic if we know aieagent ID
		topic := ""
		for _, agent := range subscribedAgentList {
			if agent.ID == aiengine.ID.String() {
				topic = agent.Topic
				break
			}
		}
		if topic == "" {
			return errors.New("no associated ai-engine for this camera is online")
		}

		err = SendMsgToAIEAgent(topic, models.AIEngineActionUpdateAIConfig, string(msgData))
		if err != nil {
			return errors.New("failed to send update ai property to aiegine " + aiengine.ID.String() + ", error: " + err.Error())
		}
	}
	return nil
}

// Func associate camera to aieagent
func AssociateCameraToAIEngine(cameraID string, aie models.AIEngine) error {
	// Test add camera (id: df5304b5-aede-41f8-860a-09d4aa025111) to this agent

	// First check if this camera is existed in database
	cam, err := reposity.ReadItemByIDIntoDTO[models.Camera, models.Camera](cameraID)
	if err != nil {
		return err
	}

	// Next check if this aieagent is existed in database
	_, err = reposity.ReadItemByIDIntoDTO[models.AIEngine, models.AIEngine](aie.ID.String())
	if err != nil {
		return err
	}

	// Let associate camera with aieagent
	err = reposity.BackRefManyToManyAppend[models.AIEngine, models.Camera](cam, "AIEngines", aie)
	if err != nil {
		return err
	}

	return nil
}

// Func remove camera from aieagent
func RemoveCameraFromAIEngine(cameraID string, aie models.AIEngine) error {
	// First check if this camera is existed in database
	cam, err := reposity.ReadItemByIDIntoDTO[models.Camera, models.Camera](cameraID)
	if err != nil {
		return err
	}

	// Next check if this aieagent is existed in database
	_, err = reposity.ReadItemByIDIntoDTO[models.AIEngine, models.AIEngine](aie.ID.String())
	if err != nil {
		return err
	}

	// Let remove the association of camera with aieagent
	err = reposity.BackRefManyToManyRemove[models.AIEngine, models.Camera](cam, "AIEngines", aie)
	if err != nil {
		return err
	}

	return nil
}

// Func get main video stream from camera
func GetAIVideoStream(cam models.Camera) (models.VideoStream, error) {
	log.Printf("Get main video stream for camera: %s", cam.ID.String())
	// Print all camera's stream
	for _, stream := range cam.Streams {
		log.Printf("Stream name: %s, data: %v", stream.Name, stream)
	}

	// Get video stream of this camera
	// First stream is main stream, second stream is sub stream, and custom streams are after that
	// When retrive stream, should check if channel contain "main" or "Main", if not take the first stream
	aiVideoStream := models.VideoStream{}
	for _, stream := range cam.Streams {
		streamName := strings.ToLower(stream.Name)
		if strings.Contains(streamName, "aistream") {
			aiVideoStream = stream
		}
	}

	if aiVideoStream.URL == "" {
		return models.VideoStream{}, errors.New("this camera do not have video stream for AI")
	}

	return aiVideoStream, nil
}
