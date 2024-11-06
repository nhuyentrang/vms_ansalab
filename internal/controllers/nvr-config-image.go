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

// ReadNVRConfig godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         cameras
// @Produce      json
// @Param   	 idCamera			query	string	true	"id Camera keyword"		minlength(1)  	maxlength(100)
// @Success      200   {object}  models.JsonDTORsp[models.DTO_ImageConfig]
// @Router       /nvrs/config/imageconfig/{idNVR} [get]
// @Security	 BearerAuth
func GetImageConfigNVR(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_ImageConfig]()
	idCamera := c.Query("idCamera")
	idNVR := c.Param("idNVR")
	dtoNVR, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVR, models.NVR]("id = ?", idNVR)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	dtoCamera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", idCamera)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	dtoNVRConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVRConfig, models.NVRConfig]("id = ?", dtoNVR.ConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoNVR.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	requestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		Cmd:          cmd_GetOSDConfigNVR,
		EventTime:    time.Now().Format(time.RFC3339),
		ProtocolType: dtoNVR.Protocol,
		RequestUUID:  requestUUID,
		IPAddress:    dtoNVR.IPAddress,
		UserName:     dtoNVR.Username,
		Password:     dtoNVR.Password,
		OnvifPort:    dtoNVR.OnvifPort,
		HttpPort:     dtoNVR.HttpPort,
		Channel:      dtoCamera.NVR.Channel,
		CameraID:     idCamera,
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

			dtoCameraImageConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_ImageConfig, models.ImageConfig]("id = ?", dtoNVRConfig.ImageConfigID)
			if err != nil {
				jsonRsp.Code = http.StatusInternalServerError
				jsonRsp.Message = err.Error()
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}
			jsonRsp.Data = dtoCameraImageConfig
			c.JSON(http.StatusOK, &jsonRsp)
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetOSDConfigNVR)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_GetOSDConfigNVR
			jsonRsp.Code = http.StatusRequestTimeout
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(requestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.ImageConfig != nil {
					jsonRsp.Data = *msg.PayLoad.ImageConfig
					// Update entity from DTO
					_, err = reposity.UpdateItemByIDFromDTO[models.DTO_ImageConfig, models.ImageConfig](dtoNVRConfig.ImageConfigID.String(), jsonRsp.Data)
					if err != nil {
						jsonRsp.Code = http.StatusInternalServerError
						jsonRsp.Message = err.Error()
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						return
					}
					c.JSON(http.StatusOK, &jsonRsp)
					return
				} else {
					jsonRsp.Message = "No image config data received"
					jsonRsp.Code = http.StatusNoContent
					c.JSON(http.StatusNoContent, &jsonRsp)

				}
			}
		}
	}

}

// ReadNVRConfig godoc
// @Summary      Get single NVR by id
// @Description  Returns the NVR whose ID value matches the id.
// @Tags         nvrs
// @Produce      json
// @Param   	 idNetWorkConfig			query	string	true	"id NetWorkConfig keyword"		minlength(1)  	maxlength(100)
// @Param        NetWorkConfig  body      models.DTO_ImageConfig  true  "NetWorkConfig JSON"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_ImageConfig]
// @Router       /nvrs/config/imageconfig/{idimageconfig} [put]
// @Security	 BearerAuth
func UpdateImageConfig(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_ImageConfig]()
	idNetWorkConfig := c.Query("idNetWorkConfig")

	var dto models.DTO_ImageConfig
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_ImageConfig, models.ImageConfig](idNetWorkConfig, dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	dtoCamera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", dto.CameraID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	cmd := models.DeviceCommand{ //unused
		CommandID:    uuid.New().String(),
		Cmd:          "",
		EventTime:    time.Now().Format(time.RFC3339),
		ProtocolType: dtoCamera.Protocol,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	for {
		select {
		case <-time.After(30 * time.Second):
			fmt.Println("\t\t> error, waiting for command result timed out: ", "")
			jsonRsp.Message = "No response from the device, timeout: " + ""
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return
		case msg := <-UpdateImageConfigCameraChannelDataReceiving:
			jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
			jsonRsp.Data = dto
			c.JSON(http.StatusOK, &jsonRsp)
			return
		}
	}
}

// UpdateImageConfigNVR godoc
// @Summary      Update image configuration for a camera
// @Description  Updates the image configuration of a specific camera based on the provided ID and image configuration data.
// @Tags         cameras
// @Accept       json
// @Produce      json
// @Param        idCamera         query     string                   true  "Camera ID"
// @Param        imageConfig      body      models.DTO_ImageConfig   true  "Image configuration data"
// @Success      200              {object}  models.JsonDTORsp[models.DTO_ImageConfig]
// @Failure      400              {object}  models.JsonDTORsp[models.DTO_ImageConfig] "Bad Request"
// @Failure      404              {object}  models.JsonDTORsp[models.DTO_ImageConfig] "Not Found"
// @Failure      500              {object}  models.JsonDTORsp[models.DTO_ImageConfig] "Internal Server Error"
// @Router       /nvrs/config/imageconfigNVR/{idNVR} [put]
// @Security     BearerAuth
func UpdateImageConfigNVR(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_ImageConfig]()
	idCamera := c.Query("idCamera")
	idNVR := c.Param("idNVR")

	// Get new data from body
	var dto models.DTO_ImageConfig
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	dtoNVR, err := reposity.ReadItemByIDIntoDTO[models.DTO_NVR, models.NVR](idNVR)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
	}

	dtoCamera, errCam := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](idCamera)
	if errCam != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = errCam.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	dtoNVRConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVRConfig, models.NVRConfig]("id = ?", dtoNVR.ConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Update entity from DTO
	UpdatedDTO, err := reposity.UpdateItemByIDFromDTO[models.DTO_ImageConfig, models.ImageConfig](dtoNVRConfig.ImageConfigID.String(), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Map VideoConfigInfo to StreamingChannelList
	ImageConfigNVR := ConvertImageConfigToVideoOverlay(UpdatedDTO)

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoNVR.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:            dtoDevice.ModelID,
		Cmd:                  cmd_UpdateOSDConfigNVR,
		EventTime:            time.Now().Format(time.RFC3339),
		ProtocolType:         dtoNVR.Protocol,
		RequestUUID:          RequestUUID,
		IPAddress:            dtoNVR.IPAddress,
		UserName:             dtoNVR.Username,
		Password:             dtoNVR.Password,
		OnvifPort:            dtoNVR.OnvifPort,
		HttpPort:             dtoNVR.HttpPort,
		SetImageConfigCamera: ImageConfigNVR,
		Channel:              dtoCamera.NVR.Channel,
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
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Wait for response
	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_UpdateOSDConfigNVR)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_UpdateOSDConfigNVR
			jsonRsp.Code = http.StatusRequestTimeout
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
				jsonRsp.Data = UpdatedDTO
				jsonRsp.Code = http.StatusOK
				c.JSON(http.StatusOK, &jsonRsp)
				return
			}
		}
	}
}
