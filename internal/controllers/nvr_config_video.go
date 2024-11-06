package controllers

import (
	"net/http"

	"vms/internal/models"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
)

// ReadNVRConfig godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         nvrs
// @Produce      json
// @Param   	 idNVR			query	string	true	"id nvr keyword"		minlength(1)  	maxlength(100)
// @Success      200   {object}  models.JsonDTORsp[models.DTO_VideoConfig]
// @Router       /nvrs/config/videoconfig/{idNVR} [get]
// @Security	 BearerAuth
func GetVideoConfig(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_VideoConfig]()
	nvrID := c.Query("idNVR")

	dtoNVR, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVR, models.NVR]("id = ?", nvrID)
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

	dtoNVRVideoConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_VideoConfig, models.VideoConfig]("id = ?", dtoNVRConfig.ImageConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Data = dtoNVRVideoConfig
	c.JSON(http.StatusOK, &jsonRsp)

}

// ReadCamera	 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         nvrs
// @Produce      json
// @Param   	 idNVR			query	string	true	"id NVR keyword"		minlength(1)  	maxlength(100)
// @Param        Event  body   models.DTO_VideoConfig  true  "DTO_VideoConfig JSON"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_VideoConfig]
// @Router       /nvrs/config/videoconfig [put]
// @Security	 BearerAuth
func UpdateVideoConfig(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_VideoConfig]()
	idNVR := c.Query("idNVR")

	// Get new data from body
	var dto models.DTO_VideoConfig
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_VideoConfig, models.NVR](idNVR, dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Message = "Ok"
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// // UpdateVideoConfigNVRs godoc
// // @Summary      Update video configurations for multiple cameras
// // @Description  Updates the video configurations of the cameras identified by the provided IDs.
// // @Tags         nvrs
// // @Produce      json
// // @Param        protocolType query    string  true  "Protocol Type"
// // @Param        videoConfig  body     []models.DTO_VideoConfig  true  "DTO_VideoConfig JSON Array"
// // @Success      200   {object}  models.JsonDTORsp[[]models.VideoConfigEditStatuses]  "Command sent successfully, returns status for each NVR."
// // @Failure      400   {object}  models.JsonDTORsp[[]models.DTO_VideoConfig]  "Bad Request"
// // @Failure      404   {object}  models.JsonDTORsp[[]models.DTO_VideoConfig]  "Device Not Found"
// // @Failure      408   {object}  models.JsonDTORsp[[]models.DTO_VideoConfig]  "Request Timeout - Data applied to database but camera did not respond"
// // @Failure      500   {object}  models.JsonDTORsp[[]models.DTO_VideoConfig]  "Internal Server Error"
// // @Router       /nvrs/config/videoconfigs [put]
// // @Security     BearerAuth
// func UpdateVideoConfigNVRs(c *gin.Context) {
// 	jsonRsp := models.NewJsonDTORsp[[]models.DTO_VideoConfig]()
// 	protocol := c.Query("protocolType")
// 	if protocol == "" {
// 		protocol = "HIKVISION"
// 	}
// 	// Get new data from body
// 	var dtos []models.DTO_VideoConfig
// 	if err := c.ShouldBindJSON(&dtos); err != nil {
// 		jsonRsp.Code = http.StatusBadRequest
// 		jsonRsp.Message = err.Error()
// 		c.JSON(http.StatusBadRequest, &jsonRsp)
// 		return
// 	}

// 	var updatedDTOs []models.DTO_VideoConfig
// 	var camerasVideoConfig []models.CamerasVideoConfig

// 	for _, dto := range dtos {
// 		dtoCamera, errCam := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](dto.ID.String())
// 		if errCam != nil {
// 			jsonRsp.Code = http.StatusInternalServerError
// 			jsonRsp.Message = errCam.Error()
// 			c.JSON(http.StatusInternalServerError, &jsonRsp)
// 			return
// 		}

// 		dtoCameraConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", dtoCamera.ConfigID)
// 		if err != nil {
// 			jsonRsp.Code = http.StatusInternalServerError
// 			jsonRsp.Message = err.Error()
// 			c.JSON(http.StatusInternalServerError, &jsonRsp)
// 			return
// 		}

// 		// Update entity from DTO
// 		UpdatedDTO, err := reposity.UpdateItemByIDFromDTO[models.DTO_VideoConfig, models.VideoConfig](dtoCameraConfig.VideoConfigID.String(), dto)
// 		if err != nil {
// 			jsonRsp.Code = http.StatusInternalServerError
// 			jsonRsp.Message = err.Error()
// 			c.JSON(http.StatusInternalServerError, &jsonRsp)
// 			return
// 		}

// 		updatedDTOs = append(updatedDTOs, UpdatedDTO)

// 		// Map VideoConfigInfo to StreamingChannelList
// 		streamingChannelList := mapVideoConfigToStreamingChannels(UpdatedDTO)

// 		cameraVideoConfig := models.CamerasVideoConfig{
// 			StreamingChannelListCameras: []models.StreamingChannelListCamera{streamingChannelList},
// 			IPAddress:                   dtoCamera.IPAddress,
// 			UserName:                    dtoCamera.Username,
// 			Password:                    dtoCamera.Password,
// 			HttpPort:                    dtoCamera.HttpPort,
// 			OnvifPort:                   dtoCamera.OnvifPort,
// 		}
// 		camerasVideoConfig = append(camerasVideoConfig, cameraVideoConfig)
// 	}

// 	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtos[0].CameraID.String())
// 	if err != nil {
// 		jsonRsp.Code = http.StatusNotFound
// 		jsonRsp.Message = "Device not found: " + err.Error()
// 		c.JSON(http.StatusNotFound, &jsonRsp)
// 		return
// 	}

// 	RequestUUID := uuid.New()
// 	cmd := models.DeviceCommand{
// 		CommandID:          dtoDevice.ModelID,
// 		Cmd:                cmd_SetVideoConfigOfMultiCamera,
// 		EventTime:          time.Now().Format(time.RFC3339),
// 		ProtocolType:       protocol,
// 		RequestUUID:        RequestUUID,
// 		CamerasVideoConfig: camerasVideoConfig,
// 	}
// 	cmsStr, _ := json.Marshal(cmd)

// 	cmsSStr, err := json.MarshalIndent(cmd, "", "  ")
// 	if err != nil {
// 		fmt.Println("Error marshalling JSON:", err)
// 	} else {
// 		fmt.Println(string(cmsSStr))
// 	}
// 	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

// 	// Set timeout and ticker
// 	timeout := time.After(30 * time.Second)
// 	ticker := time.NewTicker(2 * time.Second)
// 	defer ticker.Stop()

// 	// Wait for response
// 	for {
// 		select {
// 		case <-timeout:
// 			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_SetVideoConfigCamera)
// 			jsonRsp.Message = "Data has been applied to the database, but the camera did not respond within the expected time. " + cmd_SetVideoConfigCamera
// 			jsonRsp.Code = http.StatusRequestTimeout
// 			c.JSON(http.StatusRequestTimeout, &jsonRsp)
// 			return

// 		case <-ticker.C:
// 			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
// 				msg := storedMsg.(*models.KafkaJsonVMSMessage)
// 				jsonEditResp := models.NewJsonDTORsp[[]models.VideoConfigEditStatuses]()

// 				jsonEditResp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
// 				jsonEditResp.Data = msg.PayLoad.CameraStatuses
// 				jsonEditResp.Code = http.StatusOK
// 				c.JSON(http.StatusOK, &jsonEditResp)
// 				return
// 			}
// 		}
// 	}
// }
