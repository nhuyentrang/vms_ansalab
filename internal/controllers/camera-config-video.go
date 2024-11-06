package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"vms/internal/models"

	"vms/comongo/kafkaclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadNVRConfig godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         cameras
// @Produce      json
// @Param   	 id			path	string	true	"id nvr keyword"		minlength(1)  	maxlength(100)
// @Success      200   {object}  models.JsonDTORsp[models.DTO_VideoConfig]
// @Router       /cameras/config/videoconfig/{id} [get]
// @Security	 BearerAuth
func GetVideoConfigCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_VideoConfig]()
	idCamera := c.Param("id")

	dtoCamera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", idCamera)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	dtoCameraConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", dtoCamera.ConfigID)
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
		Cmd:          cmd_GetOnSiteVideoConfig,
		EventTime:    time.Now().Format(time.RFC3339),
		ProtocolType: dtoCamera.Protocol,
		RequestUUID:  requestUUID,
		IPAddress:    dtoCamera.IPAddress,
		UserName:     dtoCamera.Username,
		Password:     dtoCamera.Password,
		OnvifPort:    dtoCamera.OnvifPort,
		HttpPort:     dtoCamera.HttpPort,
		CameraID:     idCamera,
	}
	cmsStr, err := json.Marshal(cmd)
	if err != nil {
		fmt.Println(err)
	}
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	// Set timeout and ticker
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)

	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			dtoCameraConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", dtoCamera.ConfigID)
			if err != nil {
				jsonRsp.Code = http.StatusInternalServerError
				jsonRsp.Message = err.Error()
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}

			dtoCameraVideoConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_VideoConfig, models.VideoConfig]("id = ?", dtoCameraConfig.VideoConfigID)
			if err != nil {
				jsonRsp.Code = http.StatusInternalServerError
				jsonRsp.Message = err.Error()
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}
			jsonRsp.Data = dtoCameraVideoConfig
			c.JSON(http.StatusOK, &jsonRsp)
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetOnSiteVideoConfig)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_GetOnSiteVideoConfig
			jsonRsp.Code = http.StatusRequestTimeout
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(requestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.VideoConfigCamera != nil {
					// Update entity from DTO
					dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_VideoConfig, models.VideoConfig](dtoCameraConfig.VideoConfigID.String(), *msg.PayLoad.VideoConfigCamera)
					if err != nil {
						jsonRsp.Code = http.StatusInternalServerError
						jsonRsp.Message = err.Error()
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						return
					}
					jsonRsp.Data = dto
					c.JSON(http.StatusOK, &jsonRsp)
					return
				} else {
					jsonRsp.Message = "No video config data received"
					jsonRsp.Code = http.StatusNoContent
					c.JSON(http.StatusNoContent, &jsonRsp)

				}
			}
		}
	}

}

// UpdateVideoConfigCamera godoc
// @Summary      Update video configuration for a specific camera
// @Description  Updates the video configuration of the camera identified by the provided ID.
// @Tags         cameras
// @Produce      json
// @Param        id           path     string  true  "Camera ID"
// @Param        protocolType query    string  true  "Protocol Type" // assuming it's a query parameter, not path
// @Param        videoConfig  body     models.DTO_VideoConfig  true  "DTO_VideoConfig JSON"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_VideoConfig]
// @Failure      400   {object}  models.JsonDTORsp[models.DTO_VideoConfig]  "Bad Request"
// @Failure      404   {object}  models.JsonDTORsp[models.DTO_VideoConfig]  "Device Not Found"
// @Failure      408   {object}  models.JsonDTORsp[models.DTO_VideoConfig]  "Request Timeout"
// @Failure      500   {object}  models.JsonDTORsp[models.DTO_VideoConfig]  "Internal Server Error"
// @Router       /cameras/config/videoconfig/{id} [put]
// @Security     BearerAuth
func UpdateVideoConfigCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_VideoConfig]()
	idCamera := c.Param("id")
	protocol := c.Query("protocolType")

	var dto models.DTO_VideoConfig
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dtoCamera, errCam := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](idCamera)
	if errCam != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = errCam.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	if protocol == "" {
		protocol = dtoCamera.Protocol
	}

	dtoCameraConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", dtoCamera.ConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	UpdatedDTO, err := reposity.UpdateItemByIDFromDTO[models.DTO_VideoConfig, models.VideoConfig](dtoCameraConfig.VideoConfigID.String(), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	streamingChannelList := mapVideoConfigToStreamingChannels(UpdatedDTO)

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoCamera.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:                  dtoDevice.ModelID,
		Cmd:                        cmd_SetVideoConfigCamera,
		EventTime:                  time.Now().Format(time.RFC3339),
		ProtocolType:               protocol,
		RequestUUID:                RequestUUID,
		IPAddress:                  dtoCamera.IPAddress,
		UserName:                   dtoCamera.Username,
		Password:                   dtoCamera.Password,
		OnvifPort:                  dtoCamera.OnvifPort,
		HttpPort:                   dtoCamera.HttpPort,
		StreamingChannelListCamera: streamingChannelList,
	}
	cmsStr, _ := json.Marshal(cmd)

	cmsSStr, err := json.MarshalIndent(cmd, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	} else {
		fmt.Println(string(cmsSStr))
	}
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_SetVideoConfigCamera)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_SetVideoConfigCamera
			jsonRsp.Code = http.StatusRequestTimeout
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.Status != "FAILURE" {
					jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
					jsonRsp.Data = UpdatedDTO
					jsonRsp.Code = http.StatusOK
					c.JSON(http.StatusOK, &jsonRsp)
					return
				} else {
					jsonRsp.Message = "Failed to update Camera Video Config: " + msg.PayLoad.Cmd
					jsonRsp.Data = UpdatedDTO
					jsonRsp.Code = http.StatusInternalServerError
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return
				}

			}
		}
	}
}

func mapVideoConfigToStreamingChannels(dtoConfig models.DTO_VideoConfig) models.StreamingChannelListCamera {
	var mainStream []models.StreamingChannelUpdateCamera
	var subStreams []models.StreamingChannelUpdateCamera

	for _, vci := range dtoConfig.VideoConfigInfo {
		width, height := parseResolution(vci.Resolution)

		var id int
		switch vci.DataStreamType {
		case "main":
			id = 101
		case "sub":
			id = 102
		default:
			id = 103
		}

		streamingChannel := models.StreamingChannelUpdateCamera{
			ID: id,
			Transport: models.Transport{
				MaxPacketSize: 1000,
				ControlProtocolList: models.ControlProtocolList{
					ControlProtocols: []models.ControlProtocol{
						{StreamingTransport: "RTSP"},
						{StreamingTransport: "HTTP"},
						{StreamingTransport: "SHTTP"},
					},
				},
				Unicast: models.Unicast{
					Enabled:          true,
					RTPTransportType: "RTP/TCP",
				},
				Multicast: models.Multicast{
					Enabled:         true,
					DestIPAddress:   "0.0.0.0",
					VideoDestPortNo: 8860,
					AudioDestPortNo: 8862,
				},
				Security: models.Security{
					Enabled:         true,
					CertificateType: "digest/basic",
					SecurityAlgorithm: models.SecurityAlgorithm{
						AlgorithmType: "MD5/SHA256",
					},
				},
			},
			Enabled: true,
			Video: models.SetVideoConfig{
				VideoCodecType:          vci.VideoEncoding,
				VideoResolutionWidth:    width,
				VideoResolutionHeight:   height,
				VideoQualityControlType: vci.BitrateType,
				FixedQuality:            stringToInt(vci.VideoQuality),
				VbrUpperCap:             stringToInt(vci.MaxBitrate),
				MaxFrameRate:            stringToInt(vci.FrameRate),
				H265Profile:             vci.H265,
				Smoothing:               50,
			},
		}

		if id == 101 {
			mainStream = append(mainStream, streamingChannel)
		} else {
			subStreams = append(subStreams, streamingChannel)
		}
	}

	streamingChannels := append(mainStream, subStreams...)

	return models.StreamingChannelListCamera{
		StreamingChannel: streamingChannels,
	}
}

func parseResolution(resolution string) (int, int) {
	parts := strings.Split(resolution, "x")
	if len(parts) != 2 {
		return 0, 0
	}

	width, _ := strconv.Atoi(parts[0])
	height, _ := strconv.Atoi(parts[1])

	return width, height
}

func stringToInt(str string) int {
	value, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0
	}

	return int(value)
}

// UpdateVideoConfigCameras godoc
// @Summary      Update video configurations for multiple cameras
// @Description  Updates the video configurations of the cameras identified by the provided IDs.
// @Tags         cameras
// @Produce      json
// @Param        protocolType query    string  true  "Protocol Type"
// @Param        videoConfig  body     []models.DTO_VideoConfig  true  "DTO_VideoConfig JSON Array"
// @Success      200   {object}  models.JsonDTORsp[[]models.VideoConfigEditStatuses]  "Command sent successfully, returns status for each camera."
// @Failure      400   {object}  models.JsonDTORsp[[]models.DTO_VideoConfig]  "Bad Request"
// @Failure      404   {object}  models.JsonDTORsp[[]models.DTO_VideoConfig]  "Device Not Found"
// @Failure      408   {object}  models.JsonDTORsp[[]models.DTO_VideoConfig]  "Request Timeout - Data applied to database but camera did not respond"
// @Failure      500   {object}  models.JsonDTORsp[[]models.DTO_VideoConfig]  "Internal Server Error"
// @Router       /cameras/config/videoconfigs [put]
// @Security     BearerAuth
func UpdateVideoConfigCameras(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.DTO_VideoConfig]()
	protocol := c.Query("protocolType")
	if protocol == "" {
		protocol = "HIKVISION"
	}
	// Get new data from body
	var dtos []models.DTO_VideoConfig
	if err := c.ShouldBindJSON(&dtos); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dtoCamera, errCam := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](dtos[0].CameraID.String())
	if errCam != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = errCam.Error()
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

	var updatedDTOs []models.DTO_VideoConfig
	var camerasVideoConfig []models.CamerasVideoConfig

	for _, dto := range dtos {
		dtoCamera, errCam := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](dto.ID.String())
		if errCam != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = errCam.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

		dtoCameraConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", dtoCamera.ConfigID)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

		// Update entity from DTO
		UpdatedDTO, err := reposity.UpdateItemByIDFromDTO[models.DTO_VideoConfig, models.VideoConfig](dtoCameraConfig.VideoConfigID.String(), dto)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

		updatedDTOs = append(updatedDTOs, UpdatedDTO)

		// Map VideoConfigInfo to StreamingChannelList
		streamingChannelList := mapVideoConfigToStreamingChannels(UpdatedDTO)

		cameraVideoConfig := models.CamerasVideoConfig{
			StreamingChannelListCameras: []models.StreamingChannelListCamera{streamingChannelList},
			IPAddress:                   dtoCamera.IPAddress,
			UserName:                    dtoCamera.Username,
			Password:                    dtoCamera.Password,
			HttpPort:                    dtoCamera.HttpPort,
			OnvifPort:                   dtoCamera.OnvifPort,
		}
		camerasVideoConfig = append(camerasVideoConfig, cameraVideoConfig)
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:          dtoDevice.ModelID,
		Cmd:                cmd_SetVideoConfigOfMultiCamera,
		EventTime:          time.Now().Format(time.RFC3339),
		ProtocolType:       protocol,
		RequestUUID:        RequestUUID,
		CamerasVideoConfig: camerasVideoConfig,
	}
	cmsStr, _ := json.Marshal(cmd)

	cmsSStr, err := json.MarshalIndent(cmd, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	} else {
		fmt.Println(string(cmsSStr))
	}
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	// Set timeout and ticker
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Wait for response
	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_SetVideoConfigCamera)
			jsonRsp.Message = "Data has been applied to the database, but the camera did not respond within the expected time. " + cmd_SetVideoConfigCamera
			jsonRsp.Code = http.StatusRequestTimeout
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				jsonEditResp := models.NewJsonDTORsp[[]models.VideoConfigEditStatuses]()

				jsonEditResp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
				jsonEditResp.Data = msg.PayLoad.CameraStatuses
				jsonEditResp.Code = http.StatusOK
				c.JSON(http.StatusOK, &jsonEditResp)
				return
			}
		}
	}
}
