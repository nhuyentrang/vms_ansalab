package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"vms/internal/models"

	"vms/comongo/kafkaclient"
	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetImageConfigCamera godoc
// @Summary      Retrieve image configuration for a specific camera
// @Description  Retrieves the image configuration settings for the camera identified by the given ID.
// @Tags         cameras
// @Produce      json
// @Param        idCamera     path      string  true  "Camera ID"
// @Param        channel     query     string                  true  "Channel of the camera"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_ImageConfig]
// @Failure      400   {object}  models.JsonDTORsp[models.DTO_ImageConfig] "Bad Request"
// @Failure      404   {object}  models.JsonDTORsp[models.DTO_ImageConfig] "Not Found"
// @Failure      500   {object}  models.JsonDTORsp[models.DTO_ImageConfig] "Internal Server Error"
// @Router       /cameras/config/imageconfig/{idCamera} [get]
// @Security     BearerAuth
func GetImageConfigCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_ImageConfig]()
	idCamera := c.Param("idCamera")
	protocol := c.Query("protocolType")
	channel := c.Query("channel")

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
		Cmd:          cmd_GetOSDConfig,
		EventTime:    time.Now().Format(time.RFC3339),
		ProtocolType: protocol,
		RequestUUID:  requestUUID,
		IPAddress:    dtoCamera.IPAddress,
		UserName:     dtoCamera.Username,
		Password:     dtoCamera.Password,
		OnvifPort:    dtoCamera.OnvifPort,
		HttpPort:     dtoCamera.HttpPort,
		CameraID:     idCamera,
		Channel:      channel,
	}
	cmsStr, err := json.Marshal(cmd)
	if err != nil {
		fmt.Println(err)
	}
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	// Set timeout and ticker
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
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

			dtoCameraImageConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_ImageConfig, models.ImageConfig]("id = ?", dtoCameraConfig.ImageConfigID)
			if err != nil {
				jsonRsp.Code = http.StatusInternalServerError
				jsonRsp.Message = err.Error()
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}
			jsonRsp.Data = dtoCameraImageConfig
			c.JSON(http.StatusOK, &jsonRsp)
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetOnSiteVideoConfig)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_GetOnSiteVideoConfig
			jsonRsp.Code = http.StatusRequestTimeout
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(requestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.ImageConfig != nil {
					// Update entity from DTO
					_, err = reposity.UpdateItemByIDFromDTO[models.DTO_ImageConfig, models.ImageConfig](dtoCameraConfig.ImageConfigID.String(), jsonRsp.Data)
					if err != nil {
						jsonRsp.Code = http.StatusInternalServerError
						jsonRsp.Message = err.Error()
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						return
					}
					jsonRsp.Data = *msg.PayLoad.ImageConfig
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

// UpdateImageConfigCamera godoc
// @Summary      Update image configuration for a camera
// @Description  Updates the image configuration of a specific camera based on the provided ID and image configuration data.
// @Tags         cameras
// @Accept       json
// @Produce      json
// @Param        id               path      string                  true  "Camera ID"
// @Param        protocolType     query     string                  false  "Protocol Type (e.g., Onvif, Hikvision)"
// @Param        channel     query     string                  true  "Channel of the camera"
// @Param        imageConfig      body      models.DTO_ImageConfig   true  "Image configuration data"
// @Success      200              {object}  models.JsonDTORsp[models.DTO_ImageConfig]
// @Failure      400              {object}  models.JsonDTORsp[models.DTO_ImageConfig] "Bad Request"
// @Failure      404              {object}  models.JsonDTORsp[models.DTO_ImageConfig] "Not Found"
// @Failure      500              {object}  models.JsonDTORsp[models.DTO_ImageConfig] "Internal Server Error"
// @Router       /cameras/config/imageconfig/{id} [put]
// @Security     BearerAuth
func UpdateImageConfigCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_ImageConfig]()
	idCamera := c.Param("id")
	protocol := c.Query("protocolType")
	channel := c.Query("channel")

	// Get new data from body
	var dto models.DTO_ImageConfig
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

	dtoCameraConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", dtoCamera.ConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Update entity from DTO
	UpdatedDTO, err := reposity.UpdateItemByIDFromDTO[models.DTO_ImageConfig, models.ImageConfig](dtoCameraConfig.ImageConfigID.String(), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Map VideoConfigInfo to StreamingChannelList
	ImageConfigCamera := ConvertImageConfigToVideoOverlay(UpdatedDTO)

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoCamera.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:            dtoDevice.ModelID,
		Cmd:                  cmd_UpdateOSDConfig,
		EventTime:            time.Now().Format(time.RFC3339),
		ProtocolType:         protocol,
		RequestUUID:          RequestUUID,
		IPAddress:            dtoCamera.IPAddress,
		UserName:             dtoCamera.Username,
		Password:             dtoCamera.Password,
		OnvifPort:            dtoCamera.OnvifPort,
		HttpPort:             dtoCamera.HttpPort,
		SetImageConfigCamera: ImageConfigCamera,
		Channel:              channel,
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
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_UpdateOSDConfig)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_UpdateOSDConfig
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
					jsonRsp.Message = "Failed to update Camera Image Config: " + msg.PayLoad.Cmd
					jsonRsp.Data = UpdatedDTO
					jsonRsp.Code = http.StatusInternalServerError
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return
				}
			}
		}
	}
}

func ConvertImageConfigToVideoOverlay(imageConfig models.DTO_ImageConfig) models.VideoOverlays {
	return models.VideoOverlays{
		// Mapping DateTimeOverlay from DTO_ImageConfig
		DateTimeOverlay: models.DateTimeOverlay{
			Enabled:     !imageConfig.DisableDate,
			PositionX:   stringToInt(imageConfig.DateX),
			PositionY:   stringToInt(imageConfig.DateY),
			DateStyle:   imageConfig.DateFormat,
			TimeStyle:   imageConfig.TimeFormat,
			DisplayWeek: !imageConfig.DisableWeek,
		},

		// Mapping ChannelNameOverlay from DTO_ImageConfig
		ChannelNameOverlay: models.ChannelNameOverlay{
			Enabled:   !imageConfig.DisableName,
			PositionX: stringToInt(imageConfig.NameX),
			PositionY: stringToInt(imageConfig.NameY),
		},
	}
}

// UpdateImageConfigCameras godoc
// @Summary      Update image configuration for multiple cameras
// @Description  Updates the image configuration for a list of cameras, sending a command to the device for each camera in the list.
// @Tags         cameras
// @Accept       json
// @Produce      json
// @Param        protocolType  query    string                       false  "Protocol Type (e.g., Onvif, Hikvision)"
// @Param        channel       query    string                       false  "Channel of the camera"
// @Param        body          body     []models.DTO_ImageConfig      true   "Array of Image configuration data for each camera"
// @Success      200           {object} models.JsonDTORsp[[]models.DTO_ImageConfig]
// @Failure      400           {object} models.JsonDTORsp[[]models.DTO_ImageConfig] "Bad Request"
// @Failure      404           {object} models.JsonDTORsp[[]models.DTO_ImageConfig] "Not Found"
// @Failure      500           {object} models.JsonDTORsp[[]models.DTO_ImageConfig] "Internal Server Error"
// @Router       /cameras/config/imageconfigs [put]
// @Security     BearerAuth
func UpdateImageConfigCameras(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.DTO_ImageConfig]()
	protocol := c.Query("protocolType")
	if protocol == "" {
		protocol = "ONVIF"
	}
	channel := c.Query("channel")

	// Get new data from body
	var dtos []models.DTO_ImageConfig
	if err := c.ShouldBindJSON(&dtos); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	var deviceName string
	var updatedDTOs []models.DTO_ImageConfig
	var CamerasImageConfig []models.CamerasImageConfig
	for _, dto := range dtos {
		dtoCamera, errCam := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](dto.ID.String())
		if errCam != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = errCam.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		deviceName = dtoCamera.Box.Name
		dtoCameraConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", dtoCamera.ConfigID)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

		// Update entity from DTO
		UpdatedDTO, err := reposity.UpdateItemByIDFromDTO[models.DTO_ImageConfig, models.ImageConfig](dtoCameraConfig.ImageConfigID.String(), dto)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		updatedDTOs = append(updatedDTOs, UpdatedDTO)

		// Map VideoConfigInfo to StreamingChannelList
		streamingChannelList := ConvertImageConfigToVideoOverlay(UpdatedDTO)

		cameraImageConfig := models.CamerasImageConfig{
			VideoOverlayCamera: []models.VideoOverlays{streamingChannelList},
			IPAddress:          dtoCamera.IPAddress,
			UserName:           dtoCamera.Username,
			Password:           dtoCamera.Password,
			HttpPort:           dtoCamera.HttpPort,
			OnvifPort:          dtoCamera.OnvifPort,
		}
		CamerasImageConfig = append(CamerasImageConfig, cameraImageConfig)
	}
	dtoDevice, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Device, models.Device]("name_device = ?", deviceName)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:          dtoDevice.ModelID,
		Cmd:                cmd_UpdateOSDConfigs,
		EventTime:          time.Now().Format(time.RFC3339),
		ProtocolType:       protocol,
		RequestUUID:        RequestUUID,
		Channel:            channel,
		CamerasImageConfig: CamerasImageConfig,
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
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_UpdateOSDConfig)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_UpdateOSDConfig
			jsonRsp.Code = http.StatusRequestTimeout
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.Status != "FAILURE" {
					jsonEditResp := models.NewJsonDTORsp[[]models.ImageConfigEditStatuses]()
					jsonEditResp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
					jsonEditResp.Data = msg.PayLoad.CameraImageStatuses
					jsonEditResp.Code = http.StatusInternalServerError
					c.JSON(http.StatusOK, &jsonEditResp)
					return
				} else {
					jsonRsp.Message = "Failed to update Camera Image Config: " + msg.PayLoad.Cmd
					jsonRsp.Data = *msg.PayLoad.ImageConfigs
					jsonRsp.Code = http.StatusInternalServerError
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return
				}
			}
		}
	}
}
