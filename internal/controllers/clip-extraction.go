package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"vms/internal/models"

	"vms/comongo/minioclient"

	"vms/comongo/kafkaclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetClip		godoc
// @Summary      	Get all camera groups with query filter
// @Description  	Responds with the list of all camera as JSON.
// @Tags         	clips
// @Param   		idCamera			query	string	false	"camera name idCamera"		minlength(1)  	maxlength(100)
// @Param   		startTime			query	string	false	"camera name startTime"		minlength(1)  	maxlength(100)
// @Param   		endTime				query	string	false	"camera name endTime"		minlength(1)  	maxlength(100)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.ListVideoDownLoad]
// @Router       	/video/download [get]
// @Security		BearerAuth
func DownLoadVideo(c *gin.Context) {
	jsonRspDTOCamerasBasicInfos := models.NewJsonDTOListRsp[models.ListVideoDownLoad]()
	idCamera := c.Query("idCamera")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	// Parse and check start and end times
	startTimeParsed, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		jsonRspDTOCamerasBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCamerasBasicInfos.Message = "Invalid start time format: " + err.Error()
		c.JSON(http.StatusBadRequest, &jsonRspDTOCamerasBasicInfos)
		return
	}
	endTimeParsed, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		jsonRspDTOCamerasBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCamerasBasicInfos.Message = "Invalid end time format: " + err.Error()
		c.JSON(http.StatusBadRequest, &jsonRspDTOCamerasBasicInfos)
		return
	}

	currentTime := time.Now()
	if startTimeParsed.After(currentTime) || endTimeParsed.After(currentTime) {
		jsonRspDTOCamerasBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCamerasBasicInfos.Message = "Start time or end time cannot be in the future"
		c.JSON(http.StatusBadRequest, &jsonRspDTOCamerasBasicInfos)
		return
	}

	// Fetch Camera details
	dtoCamera, err := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](idCamera)
	if err != nil {
		jsonRspDTOCamerasBasicInfos.Code = http.StatusNotFound
		jsonRspDTOCamerasBasicInfos.Message = "Camera not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRspDTOCamerasBasicInfos)
		return
	}

	// Fetch NVR details
	dtoNVR, err := reposity.ReadItemByIDIntoDTO[models.DTO_NVR, models.NVR](dtoCamera.NVR.ID)
	if err != nil {
		jsonRspDTOCamerasBasicInfos.Code = http.StatusNotFound
		jsonRspDTOCamerasBasicInfos.Message = "NVR not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRspDTOCamerasBasicInfos)
		return
	}

	// Fetch Device details
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoCamera.Box.ID)
	if err != nil {
		jsonRspDTOCamerasBasicInfos.Code = http.StatusNotFound
		jsonRspDTOCamerasBasicInfos.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRspDTOCamerasBasicInfos)
		return
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		Cmd:          cmd_DownLoadClip,
		IPAddress:    dtoNVR.IPAddress,
		HttpPort:     dtoNVR.HttpPort,
		UserName:     dtoNVR.Username,
		Password:     dtoNVR.Password,
		Channel:      dtoCamera.NVR.Channel + "01",
		StartTime:    startTime,
		EndTime:      endTime,
		ProtocolType: dtoNVR.Protocol,
		RequestUUID:  RequestUUID,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	// Set timeout and ticker
	timeout := time.After(60 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Wait for response
	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_DownLoadClip)
			jsonRspDTOCamerasBasicInfos.Message = "No response from the device, timeout: " + cmd_DownLoadClip
			c.JSON(http.StatusRequestTimeout, &jsonRspDTOCamerasBasicInfos)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				video, _ := minioclient.GetPresignedURL(msg.PayLoad.VideoDownLoad.StorageBucket, msg.PayLoad.VideoDownLoad.Video)

				var listVideo models.ListVideoDownLoad
				newVideo := models.VideoDownLoad{
					Video: video,
				}
				fmt.Println("videoURL: ", video)
				listVideo.VideoDownLoad = append(listVideo.VideoDownLoad, newVideo)
				jsonRspDTOCamerasBasicInfos.Data = []models.ListVideoDownLoad{listVideo}
				jsonRspDTOCamerasBasicInfos.Message = "Success"
				jsonRspDTOCamerasBasicInfos.Code = http.StatusOK
				c.JSON(http.StatusOK, &jsonRspDTOCamerasBasicInfos)
				return
			}
		}
	}
}
