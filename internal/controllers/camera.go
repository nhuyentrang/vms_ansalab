package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	"vms/internal/models"

	"vms/comongo/kafkaclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateCamera		godoc
// @Summary      	Create a new camera
// @Description  	Takes a camera JSON and store in DB. Return saved JSON.
// @Tags         	cameras
// @Produce			json
// @Param        	camera  body   models.DTOCamera  true  "Camera JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTOCamera]
// @Router       	/cameras [post]
// @Security		BearerAuth
func CreateCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTOCamera]()

	// Bind the received JSON to the DTO
	var dto models.DTOCamera
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Invalid request body: " + err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	if dto.Protocol == "" {
		dto.Protocol = "ONVIF"
	}

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dto.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	// if dto.IsOfflineSetting == nil {
	// 	dto.IsOfflineSetting = false
	// }
	if dto.IsOfflineSetting == nil || !*dto.IsOfflineSetting {
		RequestUUID := uuid.New()
		cmd := models.DeviceCommand{
			CommandID:    dtoDevice.ModelID,
			Cmd:          cmd_GetDataConfig,
			RequestUUID:  RequestUUID,
			EventTime:    time.Now().Format(time.RFC3339),
			IPAddress:    dto.IPAddress,
			UserName:     dto.Username,
			Password:     dto.Password,
			HttpPort:     dto.HttpPort,
			ProtocolType: dto.Protocol,
			Channel:      "101",
		}
		cmsStr, _ := json.Marshal(cmd)
		kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

		// Set timeout and ticker
		timeout := time.After(30 * time.Second)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		// Wait for response
		for {
			select {
			case <-timeout:
				fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetDataConfig)
				jsonRsp.Message = "No response from the server device, timeout: " + cmd_GetDataConfig
				c.JSON(http.StatusRequestTimeout, &jsonRsp)
				return

			case <-ticker.C:
				if storedMsg, ok := messageMap.Load(RequestUUID); ok {
					msg := storedMsg.(*models.KafkaJsonVMSMessage)
					if msg.PayLoad.Status == "FAILURE" {
						jsonRsp.Message = "Error getting device's Configuration Data"
						jsonRsp.Data = dto
						c.JSON(http.StatusInternalServerError, &jsonRsp)
					} else {
						if (msg.PayLoad.NetworkConfig != nil) && (msg.PayLoad.VideoConfig != nil) {
							dtoNetWorkConfig, err := reposity.CreateItemFromDTO[models.DTO_NetworkConfig, models.NetworkConfig](*msg.PayLoad.NetworkConfig)
							if err != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = "Failed to create network config: " + err.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}

							// dtoVideoConfig, err := reposity.CreateItemFromDTO[models.DTO_VideoConfig, models.Video](*msg.PayLoad.VideoConfig)
							// if err != nil {
							// 	jsonRsp.Code = http.StatusInternalServerError
							// 	jsonRsp.Message = "Failed to create network config: " + err.Error()
							// 	c.JSON(http.StatusInternalServerError, &jsonRsp)
							// 	return
							// }
							//TODO make video config from videoconfig payload
							//dtoVideoConfig := mapStreamingChannelListToDTOVideoConfig(msg.PayLoad.VideoConfig)
							dtoVideoConfig := mapToDTOVideoConfig(msg.PayLoad.VideoConfig_DTO)
							dtoVideoConfig, err = reposity.CreateItemFromDTO[models.DTO_VideoConfig, models.VideoConfig](dtoVideoConfig)
							if err != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = "Failed to create network config: " + err.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}
							dtoImageConfig := mapToDTOImageConfig(msg.PayLoad.ImageConfig)
							dtoImageConfig, err = reposity.CreateItemFromDTO[models.DTO_ImageConfig, models.ImageConfig](dtoImageConfig)
							if err != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = "Failed to create network config: " + err.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}

							// Process data from video config and other configurations
							var dtoCameraConfig models.DTO_CameraConfig
							dtoCameraConfig.NetworkConfigID = dtoNetWorkConfig.ID // Done
							dtoCameraConfig.VideoConfigID = dtoVideoConfig.ID     // Done
							dtoCameraConfig.ImageConfigID = dtoImageConfig.ID
							dtoCameraConfig.StorageConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoCameraConfig.StreamingConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoCameraConfig.AIConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoCameraConfig.AudioConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoCameraConfig.RecordingScheduleID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoCameraConfig.PTZConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")

							// Check for existing camera by MAC address
							camEntry, errCameraCreate := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("mac_address = ?", dto.MACAddress)
							if errCameraCreate == nil && camEntry.MACAddress != "" {
								fmt.Println("\t\t> Failed to insert camera: Mac_address duplication ")
								jsonRsp.Message = "Failed to insert camera: Mac_address duplication"
								c.JSON(http.StatusBadRequest, &jsonRsp)
								return
							}

							// Create Camera to get the correct camera ID
							dto, errCamera := reposity.CreateItemFromDTO[models.DTOCamera, models.Camera](dto)
							if errCamera != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = "Failed to create camera: " + errCamera.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}

							// Create a map for camera channels
							cameraMap := make(map[string]models.ChannelCamera)
							streamingChannels := msg.PayLoad.VideoConfig.StreamingChannel
							videoStreamArray := models.VideoStreamArray{}
							for _, channel := range streamingChannels {
								cameraMap[channel.ChannelName] = models.ChannelCamera{
									OnDemand: channel.Video.ChannelCamera.OnDemand,
									Url:      channel.Video.ChannelCamera.Url,
									Codec:    strings.ToLower(channel.Video.ChannelCamera.Codec),
									Name:     channel.Video.ChannelCamera.Name,
								}
								channelType := "main"
								if !strings.Contains(strings.ToLower(channel.ChannelName), "main") {
									if !strings.Contains(strings.ToLower(channel.ChannelName), "sub") {
										channelType = strings.ToLower(channel.ChannelName)
									} else {
										channelType = "sub"
									}
								}
								stream := models.VideoStream{
									Name:      channel.ChannelName,
									Type:      "",
									URL:       channel.URI,
									IsProxied: false,
									IsDefault: true,
									Channel:   channelType,
									ID:        dto.ID.String(),
									Codec:     strings.ToLower(channel.Video.ChannelCamera.Codec),
								}
								videoStreamArray = append(videoStreamArray, stream)
							}

							dataFileConfig := map[string]models.ConfigCamera{}
							newCamera := models.ConfigCamera{
								NameCamera: dto.Name,
								IP:         dto.IPAddress,
								UserName:   dto.Username,
								PassWord:   dto.Password,
								HTTPPort:   strconv.Itoa(dtoNetWorkConfig.HTTP),
								RTSPPort:   strconv.Itoa(dtoNetWorkConfig.RTSP),
								OnvifPort:  strconv.Itoa(dtoNetWorkConfig.ONVIF),
								ChannelNVR: "",
								Channels:   cameraMap,
							}
							dataFileConfig[dto.ID.String()] = newCamera

							dtoCameraConfig, errCam := reposity.CreateItemFromDTO[models.DTO_CameraConfig, models.CameraConfig](dtoCameraConfig)
							if errCam != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = "Failed to create camera config: " + errCam.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}

							dto.ConfigID, _ = uuid.Parse(dtoCameraConfig.ID.String())

							dto.MACAddress = msg.PayLoad.NetworkConfig.MacAddress
							dto.Streams = videoStreamArray
							dto.Box = models.KeyValue{
								ID:   dtoDevice.ID.String(),
								Name: dtoDevice.NameDevice,
							}

							// Update the dto of camera with more information we have
							dto, errCamera = reposity.UpdateItemByIDFromDTO[models.DTOCamera, models.Camera](dto.ID.String(), dto)
							if errCamera != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = "Failed to update camera: " + errCamera.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}

							// Update Videoconfig cameraID
							dtoVideoConfig.CameraID = dto.ID
							dtoVideoConfig, errVidConfig := reposity.UpdateItemByIDFromDTO[models.DTO_VideoConfig, models.VideoConfig](dtoVideoConfig.ID.String(), dtoVideoConfig)
							if errVidConfig != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = "Failed to update camera: " + errCamera.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}
							requestUUID := uuid.New()
							// Add camera to file config
							cmdAddConfig := models.DeviceCommand{
								CommandID:    dtoDevice.ModelID,
								Cmd:          cmd_AddDataConfig,
								EventTime:    time.Now().Format(time.RFC3339),
								EventType:    "camera",
								ConfigCamera: dataFileConfig,
								ProtocolType: dto.Protocol,
								RequestUUID:  requestUUID,
							}
							cmsStr, _ := json.Marshal(cmdAddConfig)
							kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

							jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: "
							jsonRsp.Data = dto
							c.JSON(http.StatusCreated, &jsonRsp)
							return
						}
					}
				}
			}
		}
	} else {
		// Create the video streams from the existing dto.Streams
		// videoStreamArray := models.VideoStreamArray{}
		// for index, channel := range dto.Streams {
		// 	streamID := dto.Streams[index].ID
		// 	channelType := dto.Streams[index].Type

		// 	rtspURL := fmt.Sprintf("rtsp://%s:%s@%s:%s/Streaming/Channels/%s", dto.Username, dto.Password, dto.IPAddress, dto.HttpPort, streamID)

		// 	// Check if the RTSP URL format is valid
		// 	if !isRTSPFormatValid(rtspURL) {
		// 		jsonRsp.Code = http.StatusBadRequest
		// 		jsonRsp.Message = fmt.Sprintf("Invalid RTSP URL format for stream %s: %s", channel.Name, rtspURL)
		// 		c.JSON(http.StatusBadRequest, &jsonRsp)
		// 		return
		// 	}

		// 	stream := models.VideoStream{
		// 		Name:      channel.Name,
		// 		Type:      channelType,
		// 		URL:       rtspURL,
		// 		IsProxied: dto.Streams[index].IsProxied,
		// 		IsDefault: true,
		// 		Channel:   channelType,
		// 		ID:        dto.ID.String(),
		// 		Codec:     dto.Streams[index].Codec,
		// 	}
		// 	videoStreamArray = append(videoStreamArray, stream)
		// }
		// dto.Streams = videoStreamArray

		// Create Camera to get the correct camera ID
		dto, errCamera := reposity.CreateItemFromDTO[models.DTOCamera, models.Camera](dto)
		if errCamera != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to create camera: " + errCamera.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

		jsonRsp.Message = "Created Camera Successfully"
		jsonRsp.Data = dto
		c.JSON(http.StatusCreated, &jsonRsp)
		return
	}
}

// ReadCamera		 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         cameras
// @Produce      json
// @Param        id  path  string  true  "Search camera by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTOCamera]
// @Router       /cameras/{id} [get]
// @Security		BearerAuth
func ReadCamera(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTOCamera]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// GetCameraSnapshot godoc
// @Summary      Retrieve a camera snapshot
// @Description  Retrieves a snapshot from the camera identified by the provided ID.
// @Tags         cameras
// @Produce      json
// @Param        id  path    string  true  "Camera ID"
// @Param        channel  query    string  true  "Camera ID"
// @Success      200  {object}  models.JsonDTORsp[models.CameraScreenshotResponse]
// @Failure      400  {object}  models.JsonDTORsp[models.CameraScreenshotResponse] "Bad Request"
// @Failure      404  {object}  models.JsonDTORsp[models.CameraScreenshotResponse] "Not Found"
// @Failure      408  {object}  models.JsonDTORsp[models.CameraScreenshotResponse] "Request Timeout"
// @Failure      500  {object}  models.JsonDTORsp[models.CameraScreenshotResponse] "Internal Server Error"
// @Router       /cameras/snapshot/{id} [get]
// @Security     BearerAuth
func GetCameraSnapshot(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.CameraScreenshotResponse]()
	idCamera := c.Param("id")
	channel := c.Query("channel")
	dtoCamera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", idCamera)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoCamera.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	requestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		Cmd:          cmd_GetCameraSnapshot,
		EventTime:    time.Now().Format(time.RFC3339),
		ProtocolType: dtoCamera.Protocol,
		RequestUUID:  requestUUID,
		IPAddress:    dtoCamera.IPAddress,
		UserName:     dtoCamera.Username,
		Password:     dtoCamera.Password,
		OnvifPort:    dtoCamera.OnvifPort,
		HttpPort:     dtoCamera.HttpPort,
		Channel:      channel,
		CameraID:     idCamera,
	}
	cmsStr, err := json.Marshal(cmd)
	if err != nil {
		fmt.Println(err)
	}
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	// Set timeout and ticker
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetCameraSnapshot)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_GetCameraSnapshot
			jsonRsp.Code = http.StatusRequestTimeout
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(requestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				jsonRsp.Data = msg.PayLoad.CameraScreenshot
				c.JSON(http.StatusOK, &jsonRsp)
				return
			}
		}
	}
}

// ChangeCameraDevice godoc
// @Summary       Change the device (box) of a camera
// @Description   Change the device (box) of a camera. If the device doesn't exist, it will attach it. If it exists, it will replace it.
// @Tags          cameras
// @Accept        json
// @Produce       json
// @Param         id         path  string  true  "Camera ID"
// @Param         deviceID   body  string  true  "New Device (Box) ID"
// @Success       200  {object}  models.JsonDTORsp[models.DTOCamera]  "Successfully changed the device"
// @Failure       400  "Bad Request - Invalid input"
// @Failure       404  "Not Found - Camera or Device not found"
// @Failure       500  "Internal Server Error - Unable to change the device"
// @Router        /cameras/{id}/device [put]
// @Security      BearerAuth
func ChangeCameraDevice(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTOCamera]()

	cameraID := c.Param("id")
	var request struct {
		DeviceID string `json:"deviceID"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Bad Request - Invalid input: " + err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Fetch the camera to ensure it exists
	camera, err := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](cameraID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Camera not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	device, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](request.DeviceID)
	if err != nil {
		newDevice := models.DTO_Device{
			ID: uuid.MustParse(request.DeviceID),
		}
		if _, err := reposity.CreateItemFromDTO[models.DTO_Device, models.Device](newDevice); err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to create new device: " + err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		device = newDevice
	}

	// Update the camera's device (box)
	camera.Box.ID = device.ID.String()
	if _, err := reposity.UpdateItemByIDFromDTO[models.DTOCamera, models.Camera](cameraID, camera); err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to update camera device: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Data = camera
	jsonRsp.Code = http.StatusOK
	jsonRsp.Message = "Successfully changed the device"
	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateCamera godoc
// @Summary       Update single camera by id
// @Description   Updates and returns a single camera whose ID value matches the id. New data must be passed in the body.
// @Tags          cameras
// @Produce       json
// @Param         id      path    string           true "Update camera by id"
// @Param         camera  body    models.DTOCamera true "Camera JSON"
// @Success       200     {object} models.JsonDTORsp[models.DTOCamera]
// @Router        /cameras/{id} [put]
// @Security      BearerAuth
func UpdateCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTOCamera]()

	// Get new data from body
	var dto models.DTOCamera
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	cameraID := c.Param("id")
	getDataCamera, err := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](cameraID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Camera not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	query := reposity.NewQuery[models.DTO_CameraGroup, models.CameraGroup]()
	jsonCondition := fmt.Sprintf("[{\"id\":\"%s\"}]", cameraID)
	query.AddConditionOfTextField("AND", "cameras", "@>", jsonCondition)

	// Execute query to get associated cameras
	cameraGroups, _, err := query.ExecNoPaging("-created_at")
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to fetch cameras: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	if dto.Box.ID != "" {
		device, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dto.Box.ID)
		if err != nil {
			jsonRsp.Code = http.StatusNotFound
			jsonRsp.Message = "Device not found: " + err.Error()
			c.JSON(http.StatusNotFound, &jsonRsp)
			return
		}

		// Update the camera's device (box)
		dto.Box.ID = device.ID.String()
	}

	// Check if IP address or HTTP port has changed
	if getDataCamera.IPAddress != dto.IPAddress || getDataCamera.HttpPort != dto.HttpPort {
		responseUUID := uuid.New()
		cmd := models.DeviceCommand{
			CommandID:    dto.Box.ID,
			Cmd:          cmd_UpdateIPandPortHTTP,
			EventType:    "camera",
			EventTime:    time.Now().Format(time.RFC3339),
			IPAddress:    dto.IPAddress,
			CameraID:     getDataCamera.ID.String(),
			HttpPort:     dto.HttpPort,
			ProtocolType: dto.Protocol,
			RequestUUID:  responseUUID,
		}
		cmsStr, _ := json.Marshal(cmd)
		kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	} // Update camera data
	if dto.Lat == "" {
		dto.Lat = "0"
	}
	if dto.Long == "" {
		dto.Long = "0"
	}
	dto, err = reposity.UpdateItemByIDFromDTO[models.DTOCamera, models.Camera](cameraID, dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Update camera name in group
	for _, group := range cameraGroups {
		for i, camera := range group.Cameras {
			if camera.ID == cameraID {
				group.Cameras[i].DeviceName = dto.Name
			}
		}
		_, err := reposity.UpdateItemByIDFromDTO[models.DTO_CameraGroup, models.CameraGroup](group.ID.String(), group)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to update camera group: " + err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
	}

	// Override ID
	dto.ID, _ = uuid.Parse(cameraID)

	// Return
	jsonRsp.Data = dto
	jsonRsp.Code = http.StatusOK
	jsonRsp.Message = "Successfully updated the camera"
	c.JSON(http.StatusOK, &jsonRsp)
}

// DeleteCamera godoc
// @Summary      Remove single camera by id
// @Description  Delete a single entry from the repository based on id.
// @Tags         cameras
// @Produce      json
// @Param        id  path  string  true  "Delete camera by id"
// @Success      200     {object} models.JsonDTORsp[models.DTOCamera]
// @Failure      404     {object} models.JsonDTORsp[models.DTOCamera] "Camera not found"
// @Failure      408     {object} models.JsonDTORsp[models.DTOCamera] "Request timeout"
// @Failure      500     {object} models.JsonDTORsp[models.DTOCamera] "Internal server error"
// @Router       /cameras/{id} [delete]
// @Security     BearerAuth
func DeleteCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTOCamera]()
	id := c.Param("id")

	// Fetch Camera details
	dtoCamera, err := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](id)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Camera not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	// Attempt to fetch Device details, but continue if it doesn't exist
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoCamera.Box.ID)
	if err != nil {
		log.Printf("Device not found for camera %s: %v. Continuing with camera deletion.", id, err)
	}

	// Send command to delete the camera stream if device exists
	if dtoDevice.ID != uuid.Nil {
		RequestUUID := uuid.New()
		cmd := models.DeviceCommand{
			CommandID:    dtoDevice.ModelID,
			Cmd:          cmd_DeleteFileStream,
			RequestUUID:  RequestUUID,
			EventTime:    time.Now().Format(time.RFC3339),
			EventType:    "camera",
			CameraID:     id,
			ProtocolType: dtoCamera.Protocol,
		}
		cmsStr, _ := json.Marshal(cmd)
		kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

		// Set timeout and ticker
		timeout := time.After(30 * time.Second)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		// Wait for response
		for {
			select {
			case <-timeout:
				fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_DeleteFileStream)
				jsonRsp.Message = "Camera deleted, NVR asynchronous device timed out: " + cmd_DeleteFileStream
				handleCameraDeletion(c, dtoCamera, jsonRsp)
				c.JSON(http.StatusOK, &jsonRsp)
				return

			case <-ticker.C:
				if storedMsg, ok := messageMap.Load(RequestUUID); ok {
					msg := storedMsg.(*models.KafkaJsonVMSMessage)
					if msg != nil {
						if strings.ToLower(msg.PayLoad.Status) == "success" {
							jsonRsp.Code = http.StatusOK
							jsonRsp.Message = "Camera Deleted Successfully"
						} else {
							jsonRsp.Code = http.StatusOK
							jsonRsp.Message = "Camera Deleted Successfully, NVR asynchronous"
						}
						handleCameraDeletion(c, dtoCamera, jsonRsp)
						return
					}
				}
			}
		}
	} else {
		handleCameraDeletion(c, dtoCamera, jsonRsp)
	}
}

func handleCameraDeletion(c *gin.Context, dtoCamera models.DTOCamera, jsonRsp *models.JsonDTORsp[models.DTOCamera]) {
	// Update currently active system incident
	queryCamStats := reposity.NewQuery[models.DTO_System_Incident_BasicInfo, models.SystemIncident]()
	queryCamStats.AddConditionOfTextField("AND", "status", "=", "Chưa xử lý")
	queryCamStats.AddConditionOfTextField("AND", "source", "=", dtoCamera.Name)

	// Execute query to get associated cameras
	cameraGroups, _, errCamStats := queryCamStats.ExecNoPaging("-created_at")
	if errCamStats != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to fetch cameras: " + errCamStats.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Loop through the cameraGroups and update the Status and Type fields
	for _, cameraGroup := range cameraGroups {
		cameraGroup.Status = "Đã xóa camera"
		if cameraGroup.Type == "Deactive" || cameraGroup.Type == "deactive" {
			cameraGroup.Type = "Active"
		}

		// Update the system incident in the repository
		if _, err := reposity.UpdateItemByIDFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](cameraGroup.ID.String(), cameraGroup); err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to update system incident: " + err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
	}

	// Delete CameraConfig
	dtoCameraConfig, err := reposity.ReadItemByIDIntoDTO[models.DTO_CameraConfig, models.CameraConfig](dtoCamera.ConfigID.String())
	if err == nil {
		// Delete Network Config
		err = reposity.DeleteItemByID[models.NetworkConfig](dtoCameraConfig.NetworkConfigID.String())
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to delete camera network config: " + err.Error()
			c.JSON(http.StatusInternalServerError, jsonRsp)
			return
		}

		// Delete Video Config
		err = reposity.DeleteItemByID[models.VideoConfig](dtoCameraConfig.VideoConfigID.String())
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to delete camera video config: " + err.Error()
			c.JSON(http.StatusInternalServerError, jsonRsp)
			return
		}
	}

	aiEngineList, count, err := reposity.ReadAllItemsIntoDTO[models.AIEngine, models.AIEngine](dtoCameraConfig.ID.String())
	if err != nil || count == 0 {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "no aiEngine associated with this camera: " + err.Error()
		c.JSON(http.StatusInternalServerError, jsonRsp)
		return
	}

	cameraModel, err := reposity.ReadItemByIDIntoDTO[models.Camera, models.Camera](dtoCamera.ID.String())
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to get camera model: " + err.Error()
		c.JSON(http.StatusInternalServerError, jsonRsp)
		return
	}

	for _, aiEngine := range aiEngineList {
		//with each AI engine => delete relationship
		err = reposity.BackRefManyToManyRemove[models.AIEngine, models.Camera](cameraModel, "AIEngines", aiEngine)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to remove relationship from AIEngine to camera: " + err.Error()
			c.JSON(http.StatusInternalServerError, jsonRsp)
			return
		}
	}

	// TODO: Delete more config (Sound, Storage...)

	// Delete CameraConfig
	errCamConf := reposity.DeleteItemByID[models.CameraConfig](dtoCameraConfig.ID.String())
	if errCamConf != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Camera Config Not Found Or Doesn't Exist Continue Deleteing Camera " + errCamConf.Error()
		c.JSON(http.StatusInternalServerError, jsonRsp)
	}

	// Delete the camera
	errCam := reposity.DeleteItemByID[models.Camera](dtoCamera.ID.String())
	if errCam != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to delete camera: " + errCam.Error()
		c.JSON(http.StatusInternalServerError, jsonRsp)
		return
	}

	// Update groups containing the camera
	query := reposity.NewQuery[models.DTO_CameraGroup, models.CameraGroup]()
	querynvr := reposity.NewQuery[models.DTO_NVR, models.NVR]()
	sort := "-created_at"
	jsonCondition := fmt.Sprintf("[{\"id\":\"%s\"}]", dtoCamera.ID.String())
	query.AddConditionOfTextField("AND", "cameras", "@>", jsonCondition)

	// Execute Query
	groupCameraGroupDTO, _, err := query.ExecNoPaging(sort)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to retrieve camera groups: " + err.Error()
		c.JSON(http.StatusInternalServerError, jsonRsp)
		return
	}

	for _, groupDTO := range groupCameraGroupDTO {
		// Update camera list in Camera group
		updatedCameras := []models.CameraKeyValue{}
		for _, camera := range groupDTO.Cameras {
			if camera.ID != dtoCamera.ID.String() {
				updatedCameras = append(updatedCameras, camera)
			}
		}
		groupDTO.Cameras = updatedCameras

		// Update Camera group
		if _, err := reposity.UpdateItemByIDFromDTO[models.DTO_CameraGroup, models.CameraGroup](groupDTO.ID.String(), models.DTO_CameraGroup(groupDTO)); err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to update camera group: " + err.Error()
			c.JSON(http.StatusInternalServerError, jsonRsp)
			return
		}
	}

	// Execute Query
	groupNVRDTO, _, err := querynvr.ExecNoPaging(sort)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to retrieve camera groups: " + err.Error()
		c.JSON(http.StatusInternalServerError, jsonRsp)
		return
	}

	for _, groupDTO := range groupNVRDTO {
		// Update camera list in NVR group
		updatedCameras := models.KeyValueArray{}
		for _, camera := range *groupDTO.Cameras {
			if camera.ID != dtoCamera.ID.String() {
				updatedCameras = append(updatedCameras, camera)
			}
		}
		groupDTO.Cameras = &updatedCameras

		// Update NVR Cameras group
		if _, err := reposity.UpdateItemByIDFromDTO[models.DTO_NVR, models.NVR](groupDTO.ID.String(), models.DTO_NVR(groupDTO)); err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to update camera group: " + err.Error()
			c.JSON(http.StatusInternalServerError, jsonRsp)
			return
		}
	}
	jsonRsp.Code = http.StatusOK
	c.JSON(http.StatusOK, jsonRsp)
}

// GetCameras		godoc
// @Summary      	Get all camera groups with query filter
// @Description  	Responds with the list of all camera as JSON.
// @Tags         	cameras
// @Param   		keyword			query	string	false	"camera name keyword"		minlength(1)  	maxlength(100)
// @Param   		sort  			query	string	false	"sort"						default(+created_at)
// @Param   		limit			query	int     false  	"limit"          			minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          			minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_Camera_Read_BasicInfo]
// @Router       	/cameras [get]
// @Security		BearerAuth
func GetCameras(c *gin.Context) {
	jsonRspDTOCamerasBasicInfos := models.NewJsonDTOListRsp[models.DTO_Camera_Read_BasicInfo]()

	// Get param
	keyword := c.Query("keyword")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	fmt.Println(
		"keyword: ", keyword,
		" - sort: ", sort,
		" - limit: ", limit,
		" - page: ", page)

	// Build query
	query := reposity.NewQuery[models.DTO_Camera_Read_BasicInfo, models.Camera]()

	// Search for keyword in name
	if keyword != "" {
		query.AddConditionOfTextField("AND", "name", "LIKE", keyword)
	}

	// Exec query
	dtoCameraBasics, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTOCamerasBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCamerasBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCamerasBasicInfos)
		return
	}

	jsonRspDTOCamerasBasicInfos.Count = count
	jsonRspDTOCamerasBasicInfos.Data = dtoCameraBasics
	jsonRspDTOCamerasBasicInfos.Page = int64(page)
	jsonRspDTOCamerasBasicInfos.Size = int64(len(dtoCameraBasics))
	c.JSON(http.StatusOK, &jsonRspDTOCamerasBasicInfos)
}

// GetCameraConfigOptions		godoc
// @Summary      	Get items of camera protocols for select box
// @Description  	Responds with the list of item for camera protocol.
// @Tags         	cameras
// @Produce      	json
// @Success      	200  		{object}  		models.JsonDTORsp[[]models.KeyValue]
// @Router       	/cameras/options/protocol-types [get]
// @Security		BearerAuth
func GetCameraProtocolTypes(c *gin.Context) {

	jsonConfigTypeRsp := models.NewJsonDTORsp[[]models.KeyValue]()
	jsonConfigTypeRsp.Data = make([]models.KeyValue, 0)

	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-onvif",
		Name: "ONVIF",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-ivi",
		Name: "IVI",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-hikvision",
		Name: "HIKVISION",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-dahua",
		Name: "DAHUA",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-axis",
		Name: "AXIS",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-bosch",
		Name: "BOSCH",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-hanwha",
		Name: "HANWHA",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-panasonic",
		Name: "PANASONIC",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-unv",
		Name: "UNIVIEW",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-kbvision",
		Name: "KBVISION",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-agivilon",
		Name: "AGIVILON",
	})

	jsonConfigTypeRsp.Code = 0
	c.JSON(http.StatusOK, &jsonConfigTypeRsp)
}

// GetCameraConfigOptions		godoc
// @Summary      	Get items of camera stream for select box
// @Description  	Responds with the list of item for camera stream type.
// @Tags         	cameras
// @Produce      	json
// @Success      	200  		{object}  		models.JsonDTORsp[[]models.KeyValue]
// @Router       	/cameras/options/stream-types [get]
// @Security		BearerAuth
func GetCameraStreamTypes(c *gin.Context) {

	jsonConfigTypeRsp := models.NewJsonDTORsp[[]models.KeyValue]()
	jsonConfigTypeRsp.Data = make([]models.KeyValue, 0)

	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-stream-ondemand",
		Name: "ON Demand",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-stream-persistant",
		Name: "Persistant",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-stream-p2p",
		Name: "P2P",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-stream-proxied",
		Name: "Proxied",
	})

	jsonConfigTypeRsp.Code = 0
	c.JSON(http.StatusOK, &jsonConfigTypeRsp)
}

// GetCameraConfigOptions		godoc
// @Summary      	Get items of camera types for select box
// @Description  	Responds with the list of item for camera type.
// @Tags         	cameras
// @Produce      	json
// @Success      	200  		{object}  		models.JsonDTORsp[[]models.KeyValue]
// @Router       	/cameras/options/types [get]
// @Security		BearerAuth
func GetCameraTypes(c *gin.Context) {

	jsonConfigTypeRsp := models.NewJsonDTORsp[[]models.KeyValue]()
	jsonConfigTypeRsp.Data = make([]models.KeyValue, 0)

	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "ip-camera",
		Name: "IP Camera",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "ptz-camera",
		Name: "PTZ",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "ai-camera",
		Name: "AI",
	})

	jsonConfigTypeRsp.Code = 0
	c.JSON(http.StatusOK, &jsonConfigTypeRsp)
}

// ReadCamera		 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         cameras
// @Produce      json
// @Param        serial  path  string  true  "Search id by serialnumber"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_Camera_Serial]
// @Router       /cameras/serial/{serial} [get]
// @Security		BearerAuth
func ReadSerialCamera(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_Camera_Serial]()

	dto, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Camera_Serial, models.Camera]("serial_number = ?", c.Param("serial"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// ImportCamerasFromCSV godoc
// @Summary       Import cameras from a CSV file
// @Description   Reads a CSV file containing camera data and inserts them into the database.
// @Tags          cameras
// @Produce       json
// @Param         file  formData  file  true  "CSV file with cameras"
// @Param         idbox query     string true  "Id of the box"  minlength(1)   maxlength(100)
// @Success       200   {object}  models.JsonDTORsp[[]models.DTOCamera]
// @Router        /cameras/import [post]
// @Security      BearerAuth
func ImportCamerasFromCSV(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.DTOCameraImport]()
	idbox := c.Query("idbox")
	file, err := c.FormFile("file")
	if err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Failed to get CSV file: " + err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	f, err := file.Open()
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to open CSV file: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	var cameras []models.DTOCameraImport
	var errors []string

	headers, err := r.Read()
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to read CSV headers: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Fetch Device details
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](idbox)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "failed to retrieve device information"
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, "Failed to read CSV record: "+err.Error())
			continue
		}

		dto, err := parseCameraRecord(headers, record)
		if err != nil {
			errors = append(errors, "Failed to parse CSV record: "+err.Error())
			continue
		} else if strings.ToLower(dto.Protocol) != "onvif" && strings.ToLower(dto.Protocol) != "hikvision" {
			errors = append(errors, "Entry protocol is invalid")
			continue
		}

		dto, err = insertCameraIntoDB(dto, dtoDevice)
		if err != nil {
			log.Printf("Failed to insert camera: %v\n", err)
			errors = append(errors, fmt.Sprintf("Failed to insert camera %s: %v", dto.IPAddress, err))
		}

		cameras = append(cameras, dto)
	}

	jsonRsp.Data = cameras
	if len(errors) > 0 {
		jsonRsp.Message = "Some errors occurred: " + strings.Join(errors, "; ")
		jsonRsp.Code = http.StatusPartialContent
	} else {
		jsonRsp.Message = "All records processed successfully"
		jsonRsp.Code = http.StatusOK
	}
	c.JSON(http.StatusOK, &jsonRsp)
}

func createCameraLogic(dto models.DTOCameraImport, device models.DTO_Device) (models.DTOCameraImport, error) {

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    device.ModelID,
		Cmd:          cmd_GetDataConfig,
		RequestUUID:  RequestUUID,
		EventTime:    time.Now().Format(time.RFC3339),
		IPAddress:    dto.IPAddress,
		UserName:     dto.Username,
		Password:     dto.Password,
		HttpPort:     dto.HttpPort,
		ProtocolType: dto.Protocol,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return dto, fmt.Errorf("timeout waiting for command result: %s", cmd_GetDataConfig)
		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.NetworkConfig != nil && msg.PayLoad.VideoConfig != nil {
					dtoNetWorkConfig, err := reposity.CreateItemFromDTO[models.DTO_NetworkConfig, models.NetworkConfig](*msg.PayLoad.NetworkConfig)
					if err != nil {
						return dto, fmt.Errorf("failed to create network config: %w", err)
					}

					dtoVideoConfig := mapToDTOVideoConfig(msg.PayLoad.VideoConfig_DTO)
					dtoVideoConfig, err = reposity.CreateItemFromDTO[models.DTO_VideoConfig, models.VideoConfig](dtoVideoConfig)
					if err != nil {
						return dto, fmt.Errorf("failed to create video config: %w", err)

					}

					var dtoCameraConfig models.DTO_CameraConfig
					dtoCameraConfig.NetworkConfigID = dtoNetWorkConfig.ID
					dtoCameraConfig.VideoConfigID = dtoVideoConfig.ID
					dtoCameraConfig.ImageConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoCameraConfig.StorageConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoCameraConfig.StreamingConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoCameraConfig.AIConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoCameraConfig.AudioConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoCameraConfig.RecordingScheduleID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoCameraConfig.PTZConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")

					// Check if the entry already exists in the main camera table based on MACAddress
					_, errCameraCreate := reposity.ReadItemWithFilterIntoDTO[models.DTOCameraImport, models.Camera]("mac_address = ?", dto.MACAddress)
					if errCameraCreate == nil {
						return dto, fmt.Errorf("entry already exists in main camera table")
					}

					// Assign additional fields
					dto.ConfigID = dtoCameraConfig.ID
					dto.MACAddress = msg.PayLoad.NetworkConfig.MacAddress

					// Create the camera entry
					dto, errCamera := reposity.CreateItemFromDTO[models.DTOCameraImport, models.Camera](dto)
					if errCamera != nil {
						return dto, fmt.Errorf("failed to create camera: %w", errCamera)
					}

					// Process streaming channels and create video streams
					cameraMap := make(map[string]models.ChannelCamera)
					streamingChannels := msg.PayLoad.VideoConfig.StreamingChannel
					videoStreamArray := models.VideoStreamArray{}
					for _, channel := range streamingChannels {
						cameraMap[channel.ChannelName] = models.ChannelCamera{
							OnDemand: channel.Video.ChannelCamera.OnDemand,
							Url:      channel.Video.ChannelCamera.Url,
							Codec:    strings.ToLower(channel.Video.ChannelCamera.Codec),
							Name:     channel.Video.ChannelCamera.Name,
						}
						channelType := "main"
						if !strings.Contains(strings.ToLower(channel.ChannelName), "main") {
							if !strings.Contains(strings.ToLower(channel.ChannelName), "sub") {
								channelType = strings.ToLower(channel.ChannelName)
							} else {
								channelType = "sub"
							}
						}
						stream := models.VideoStream{
							Name:      channel.ChannelName,
							Type:      "",
							URL:       channel.URI,
							IsProxied: false,
							IsDefault: true,
							Channel:   channelType,
							ID:        dto.ID.String(),
							Codec:     strings.ToLower(channel.Video.ChannelCamera.Codec),
						}
						videoStreamArray = append(videoStreamArray, stream)
					}

					dataFileConfig := map[string]models.ConfigCamera{}
					newCamera := models.ConfigCamera{
						NameCamera: dto.Name,
						IP:         dto.IPAddress,
						UserName:   dto.Username,
						PassWord:   dto.Password,
						HTTPPort:   strconv.Itoa(dtoNetWorkConfig.HTTP),
						RTSPPort:   strconv.Itoa(dtoNetWorkConfig.RTSP),
						OnvifPort:  strconv.Itoa(dtoNetWorkConfig.ONVIF),
						ChannelNVR: "",
						Channels:   cameraMap,
					}
					dataFileConfig[dto.ID.String()] = newCamera

					// Create camera config
					dtoCameraConfig, errCam := reposity.CreateItemFromDTO[models.DTO_CameraConfig, models.CameraConfig](dtoCameraConfig)
					if errCam != nil {
						return dto, fmt.Errorf("failed to create camera config: %w", errCam)
					}

					// Update the camera entry with streams
					dto.ConfigID, _ = uuid.Parse(dtoCameraConfig.ID.String())

					dto.MACAddress = msg.PayLoad.NetworkConfig.MacAddress
					dto.Streams = videoStreamArray
					dto.Box = models.KeyValue{
						ID:   device.ID.String(),
						Name: device.NameDevice,
					}
					dto, errCamera = reposity.UpdateItemByIDFromDTO[models.DTOCameraImport, models.Camera](dto.ID.String(), dto)
					if errCamera != nil {
						return dto, fmt.Errorf("failed to update camera: %w", errCamera)
					}

					requestUUID := uuid.New()
					cmdAddConfig := models.DeviceCommand{
						CommandID:    device.ModelID,
						Cmd:          cmd_AddDataConfig,
						EventTime:    time.Now().Format(time.RFC3339),
						EventType:    "camera",
						ConfigCamera: dataFileConfig,
						ProtocolType: dto.Protocol,
						RequestUUID:  requestUUID,
					}
					cmsStr, _ := json.Marshal(cmdAddConfig)
					kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
					//TODO: Check if this command get error
					return dto, nil
				}
			}
		}
	}
}

func insertCameraIntoDB(dto models.DTOCameraImport, device models.DTO_Device) (models.DTOCameraImport, error) {

	// Send command to scan device
	RequestUUID := uuid.New()
	cmdScan := models.DeviceCommand{
		CommandID:    device.ModelID,
		Cmd:          cmd_ScanDeviceIP,
		EventTime:    time.Now().Format(time.RFC3339),
		IPAddress:    dto.IPAddress,
		HttpPort:     dto.HttpPort,
		UserName:     dto.Username,
		Password:     dto.Password,
		ProtocolType: dto.Protocol,
		RequestUUID:  RequestUUID,
	}
	cmsStr, _ := json.Marshal(cmdScan)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	// Set timeout and ticker
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	scanResultNotDone := 1
	// Wait for scan result
	for scanResultNotDone > 0 {
		select {
		case <-timeout:
			return dto, fmt.Errorf("timeout waiting for scan command result: %s", cmd_ScanDeviceIP)
		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.DeviceScan != nil && len(*msg.PayLoad.DeviceScan) > 0 {
					scanResult := (*msg.PayLoad.DeviceScan)[0]
					dto.ID = uuid.New()
					dto.Name = scanResult.Name
					dto.Type.ID = scanResult.Type
					dto.Location = scanResult.Location
					dto.MACAddress = scanResult.MacAddress
					dto.FirmwareVersion = scanResult.FirmwareVersion
					dto.ExportStatus = false
					//dto.Protocol = scanResult.Protocol
					scanResultNotDone = 0
				} else {
					return dto, fmt.Errorf("failed to retrieve camera information")
				}
			}
		}
	}

	dto.InsertStatus = false // Default status to false
	importedDto, err := reposity.CreateItemFromDTO[models.DTOCameraImport, models.CameraImport](dto)
	if err != nil {
		return dto, fmt.Errorf("failed to insert into import table: %w", err)
	}

	// Attempt to process and insert the camera into the main Camera table
	processedDto, err := createCameraLogic(dto, device)
	if err != nil {
		// Update the import record to reflect the failure
		importedDto.InsertStatus = false
		_, updateErr := reposity.UpdateItemByIDFromDTO[models.DTOCameraImport, models.CameraImport](importedDto.ID.String(), importedDto)
		if updateErr != nil {
			return dto, fmt.Errorf("failed to update import record status after camera insert failure: %w", updateErr)
		}
		return dto, fmt.Errorf("failed to create camera: %w", err)
	}

	// Update the import record to reflect the success
	dto.InsertStatus = true
	_, updateErr := reposity.UpdateItemByIDFromDTO[models.DTOCameraImport, models.CameraImport](dto.ID.String(), dto)
	if updateErr != nil {
		return dto, fmt.Errorf("failed to update import record status after successful camera insert: %w", updateErr)
	}

	return processedDto, nil
}

func parseCameraRecord(headers, record []string) (models.DTOCameraImport, error) {
	var dto models.DTOCameraImport

	for i, header := range headers {
		switch header {
		case "name":
			dto.Name = record[i]
		case "protocol":
			dto.Protocol = record[i]
		case "ipAddress":
			dto.IPAddress = record[i]
		case "macAddress":
			dto.MACAddress = record[i]
		case "httpPort":
			dto.HttpPort = record[i]
		case "onvifPort":
			dto.OnvifPort = record[i]
		case "managementPort":
			dto.ManagementPort = record[i]
		case "username":
			dto.Username = record[i]
		case "password":
			dto.Password = record[i]
		default:
			return dto, fmt.Errorf("unexpected header: %s", header)
		}
	}

	return dto, nil
}

// GenerateSampleCameraData godoc
// @Summary       Generate and download sample camera data
// @Description   Provides a sample CSV file for users to know what data to input.
// @Tags          cameras
// @Produce       text/csv
// @Success       200   "Sample data file"
// @Router        /cameras/sample-data [get]
// @Security      BearerAuth
func GenerateSampleCameraData(c *gin.Context) {
	sampleData := [][]string{
		{"ipAddress", "httpPort", "username", "password"},
		{"192.168.1.2", "80", "admin", "password123"},
		{"192.168.1.3", "8080", "user", "userpass"},
		{"192.168.1.4", "8081", "admin", "adminpass"},
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=sample_camera_data.csv")
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	for _, record := range sampleData {
		if err := writer.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV data"})
			return
		}
	}
}

// DownloadImportedCameras godoc
// @Summary       Download the imported camera list as an Excel file
// @Description   Provides the imported camera list in Excel format.
// @Tags          cameras
// @Produce       octet-stream
// @Success       200   "Excel file"
// @Router        /cameras/imported/download [get]
// @Security      BearerAuth
func DownloadImportedCameras(c *gin.Context) {
	// Fetch imported cameras from the database
	importedCameras, err := fetchImportedCameras()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch imported cameras"})
		return
	}

	// Create a new CSV file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")

	// Get filename from query or use default
	date := time.Now().Format("2006-01-02")
	defaultFilename := fmt.Sprintf("Cameras_%s.csv", date)
	filename := c.DefaultQuery("filename", defaultFilename)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write headers
	headers := []string{"id", "name", "indexNVR", "type", "protocol", "model", "serial", "firmwareVersion", "ipAddress", "macAddress", "httpPort", "onvifPort", "managementPort", "username", "password", "useTLS", "isOfflineSetting", "iframeURL", "lat", "long", "insertstatus", "location", "coordinate", "position", "faceRecognition", "licensePlateRecognition", "configID"}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV headers"})
		return
	}

	// Write data
	for _, camera := range importedCameras {
		record := []string{
			camera.ID.String(),
			camera.Name,
			camera.IndexNVR,
			camera.Type.ID,
			camera.Protocol,
			camera.Model,
			camera.SerialNumber,
			camera.FirmwareVersion,
			camera.IPAddress,
			camera.MACAddress,
			camera.HttpPort,
			camera.OnvifPort,
			camera.ManagementPort,
			camera.Username,
			camera.Password,
			strconv.FormatBool(camera.UseTLS),
			strconv.FormatBool(camera.IsOfflineSetting),
			camera.IFrameURL,
			camera.Lat,
			camera.Long,
			strconv.FormatBool(camera.InsertStatus),
			camera.Location,
			camera.Coordinate,
			camera.Position,
			strconv.FormatBool(camera.FaceRecognition),
			strconv.FormatBool(camera.LicensePlateRecognition),
			camera.ConfigID.String(),
		}
		if err := writer.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV data"})
			return
		}
	}

	// Update export status to true for all exported cameras
	for _, camera := range importedCameras {
		camera.InsertStatus = true
		_, updateErr := reposity.UpdateItemByIDFromDTO[models.DTOCameraImport, models.CameraImport](camera.ID.String(), camera)
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update export status for camera %s: %v", camera.ID.String(), updateErr)})
			return
		}
	}
}

func fetchImportedCameras() ([]models.DTOCameraImport, error) {
	query := reposity.NewQuery[models.DTOCameraImport, models.CameraImport]()
	query.AddConditionOfTextField("AND", "export_status", "=", false)
	cameras, _, err := query.ExecWithPaging("+created_at", 9999, 1)
	if err != nil {
		return nil, err
	}

	// Update the exportstatus to true for all fetched cameras
	for _, camera := range cameras {
		camera.ExportStatus = true
		_, updateErr := reposity.UpdateItemByIDFromDTO[models.DTOCameraImport, models.CameraImport](camera.ID.String(), camera)
		if updateErr != nil {
			return nil, fmt.Errorf("failed to update export status for camera %s: %w", camera.ID, updateErr)
		}
	}

	return cameras, nil
}

func isRTSPFormatValid(rtspURL string) bool {
	rtspRegex := regexp.MustCompile(`^rtsp:\/\/.+:.+@.+:.+\/Streaming\/Channels\/.+$`)
	return rtspRegex.MatchString(rtspURL)
}

func mapStreamingChannelListToDTOVideoConfig(videoConfigList *models.StreamingChannelList) models.DTO_VideoConfig {
	var videoConfigInfos models.VideoConfigInfoArr

	for _, channel := range videoConfigList.StreamingChannel {
		videoConfigInfo := models.VideoConfigInfo{
			DataStreamType: channel.ChannelName, // Assuming ChannelName corresponds to DataStreamType
			Resolution:     fmt.Sprintf("%dx%d", channel.Video.VideoResolutionWidth, channel.Video.VideoResolutionHeight),
			BitrateType:    channel.Video.VideoQualityControlType,
			VideoQuality:   fmt.Sprintf("%d", channel.Video.FixedQuality),
			FrameRate:      fmt.Sprintf("%d", channel.Video.MaxFrameRate),
			MaxBitrate:     fmt.Sprintf("%d", channel.Video.VbrUpperCap),
			VideoEncoding:  channel.Video.VideoCodecType,
			H265:           channel.Video.H265Profile,
		}

		videoConfigInfos = append(videoConfigInfos, videoConfigInfo)
	}

	return models.DTO_VideoConfig{
		ID:              uuid.New(), // Generate a new UUID or use an existing one
		CameraID:        uuid.New(), // This should be set based on the actual camera ID
		VideoConfigInfo: videoConfigInfos,
	}
}

func mapToDTOVideoConfig(videoConfigList *models.DTO_VideoConfig) models.DTO_VideoConfig {
	var videoConfigInfos models.VideoConfigInfoArr

	// Check if videoConfigList is nil and initialize it if needed
	if videoConfigList == nil {
		videoConfigList = &models.DTO_VideoConfig{
			VideoConfigInfo: models.VideoConfigInfoArr{},
		}
	}

	for _, channel := range videoConfigList.VideoConfigInfo {
		videoConfigInfo := models.VideoConfigInfo{
			DataStreamType: channel.DataStreamType,
			Resolution:     channel.Resolution,
			BitrateType:    channel.BitrateType,
			VideoQuality:   channel.VideoQuality,
			FrameRate:      channel.FrameRate,
			MaxBitrate:     channel.MaxBitrate,
			VideoEncoding:  channel.VideoEncoding,
			H265:           channel.H265,
		}

		videoConfigInfos = append(videoConfigInfos, videoConfigInfo)
	}

	return models.DTO_VideoConfig{
		ID:              uuid.New(), // Generate a new UUID or use an existing one
		VideoConfigInfo: videoConfigInfos,
	}
}

func mapToDTOImageConfig(imageConfig *models.DTO_ImageConfig) models.DTO_ImageConfig {
	// Check if imageConfig is nil and initialize it if needed
	if imageConfig == nil {
		imageConfig = &models.DTO_ImageConfig{}
	}

	return models.DTO_ImageConfig{
		CameraID:    imageConfig.CameraID,
		DisableName: imageConfig.DisableName,
		DisableDate: imageConfig.DisableDate,
		DisableWeek: imageConfig.DisableWeek,
		DateFormat:  imageConfig.DateFormat,
		TimeFormat:  imageConfig.TimeFormat,
		NameX:       imageConfig.NameX,
		NameY:       imageConfig.NameY,
		DateX:       imageConfig.DateX,
		DateY:       imageConfig.DateY,
		WeekX:       imageConfig.WeekX,
		WeekY:       imageConfig.WeekY,
	}
}
