package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"vms/internal/models"

	"vms/internal/models/hikivision"

	"vms/comongo/kafkaclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetStorage		godoc
// @Summary      	Get all storage playback with query filter
// @Description  	Responds with the list of all storage as JSON.
// @Tags         	playback
// @Param   		days			query	string	false	"camera name keyword"		minlength(1)  	maxlength(100)
// @Param   		startTime			query	string	false	"camera startTime keyword"		minlength(1)  	maxlength(100)
// @Param   		endTime			query	string	false	"camera endTime keyword"		minlength(1)  	maxlength(100)
// @Param        	id  path  string  true  "Get camera by id"
// @Produce      	json
// @Success      	200
// @Router       	/playback/camera/{id} [get]
// @Security		BearerAuth
func GetStoragePlayback(c *gin.Context) {
	jsonRspStoragePlayback := models.NewJsonDTORsp[*hikivision.CMSearchResult]()
	id := c.Param("id")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	// Fetch Camera details
	dtoCamera, err := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](id)
	if err != nil {
		jsonRspStoragePlayback.Code = http.StatusNotFound
		jsonRspStoragePlayback.Message = "Camera not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRspStoragePlayback)
		return
	}

	if dtoCamera.NVR.Channel == "" {
		jsonRspStoragePlayback.Code = http.StatusNotFound
		jsonRspStoragePlayback.Message = "The camera has not been configured for storage"
		c.JSON(http.StatusNotFound, &jsonRspStoragePlayback)
		return
	}

	// Fetch NVR details
	dtoNVR, err := reposity.ReadItemByIDIntoDTO[models.DTO_NVR, models.NVR](dtoCamera.NVR.ID)
	if err != nil {
		jsonRspStoragePlayback.Code = http.StatusNotFound
		jsonRspStoragePlayback.Message = "NVR not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRspStoragePlayback)
		return
	}

	// Fetch Device details
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoNVR.Box.ID)
	if err != nil {
		jsonRspStoragePlayback.Code = http.StatusNotFound
		jsonRspStoragePlayback.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRspStoragePlayback)
		return
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		RequestUUID:  RequestUUID,
		Cmd:          cmd_GetDataStoragePlayback,
		EventTime:    time.Now().Format(time.RFC3339),
		IPAddress:    dtoNVR.IPAddress,
		UserName:     dtoNVR.Username,
		Password:     dtoNVR.Password,
		HttpPort:     dtoNVR.HttpPort,
		StartTime:    startTime,
		EndTime:      endTime,
		Channel:      dtoCamera.NVR.Channel + "01",
		NVRID:        dtoNVR.ID.String(),
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
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetDataStoragePlayback)
			jsonRspStoragePlayback.Message = "No response from the device, timeout: " + cmd_GetDataStoragePlayback
			c.JSON(http.StatusRequestTimeout, &jsonRspStoragePlayback)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.CMSearchResult != nil {
					jsonRspStoragePlayback.Data = msg.PayLoad.CMSearchResult
					c.JSON(http.StatusOK, &jsonRspStoragePlayback)
					return
				}

				if msg.PayLoad.CMSearchResult == nil {
					jsonRspStoragePlayback.Message = "Format error: " + msg.PayLoad.ResponseStatus.StatusString
					c.JSON(http.StatusNoContent, &jsonRspStoragePlayback)
					return
				}
			}
		}
	}
}

// PostCalender		godoc
// @Summary      	Get all calender playback with query filter
// @Description  	Responds with the list of all storage as JSON.
// @Tags         	playback
// @Param        	id  path  string  true  "Get calender by id"
// @Param        	calender  body   hikivision.TrackDailyParam  true  "calender JSON"
// @Produce      	json
// @Success      	200
// @Router       	/playback/camera/{id} [post]
// @Security		BearerAuth
func GetCalenderPlayback(c *gin.Context) {
	jsonRspStoragePlayback := models.NewJsonDTORsp[hikivision.TrackDailyDistribution]()
	id := c.Param("id")

	// Bind JSON request to DTO
	var dtoCalender hikivision.TrackDailyParam
	if err := c.BindJSON(&dtoCalender); err != nil {
		jsonRspStoragePlayback.Code = http.StatusBadRequest
		jsonRspStoragePlayback.Message = "Invalid request body: " + err.Error()
		c.JSON(http.StatusBadRequest, &jsonRspStoragePlayback)
		return
	}

	// Fetch Camera details
	dtoCamera, err := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](id)
	if err != nil {
		jsonRspStoragePlayback.Code = http.StatusNotFound
		jsonRspStoragePlayback.Message = "Camera not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRspStoragePlayback)
		return
	}

	if dtoCamera.NVR.Channel == "" {
		jsonRspStoragePlayback.Code = http.StatusNotFound
		jsonRspStoragePlayback.Message = "The camera has not been configured for storage"
		c.JSON(http.StatusNotFound, &jsonRspStoragePlayback)
		return
	}

	// Fetch NVR details
	dtoNVR, err := reposity.ReadItemByIDIntoDTO[models.DTO_NVR, models.NVR](dtoCamera.NVR.ID)
	if err != nil {
		jsonRspStoragePlayback.Code = http.StatusNotFound
		jsonRspStoragePlayback.Message = "NVR not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRspStoragePlayback)
		return
	}

	// Fetch Device details
	dtoDevice, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Device, models.Device]("id = ?", dtoNVR.Box.ID)
	if err != nil {
		jsonRspStoragePlayback.Code = http.StatusNotFound
		jsonRspStoragePlayback.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRspStoragePlayback)
		return
	}

	// Prepare the command
	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:       dtoDevice.ModelID,
		RequestUUID:     RequestUUID,
		Cmd:             cmd_GetCalendelPlayback,
		EventTime:       time.Now().Format(time.RFC3339),
		IPAddress:       dtoNVR.IPAddress,
		UserName:        dtoNVR.Username,
		Password:        dtoNVR.Password,
		HttpPort:        dtoNVR.HttpPort,
		Channel:         dtoCamera.NVR.Channel + "01",
		TrackDailyParam: dtoCalender,
		ProtocolType:    dtoNVR.Protocol,
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
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetCalendelPlayback)
			jsonRspStoragePlayback.Message = "No response from the device, timeout: " + cmd_GetCalendelPlayback
			c.JSON(http.StatusRequestTimeout, &jsonRspStoragePlayback)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.TrackDailyDistribution != nil {
					jsonRspStoragePlayback.Data = *msg.PayLoad.TrackDailyDistribution
					c.JSON(http.StatusOK, &jsonRspStoragePlayback)
					return
				}

				if msg.PayLoad.TrackDailyDistribution == nil {
					jsonRspStoragePlayback.Message = "No data"
					c.JSON(http.StatusNoContent, &jsonRspStoragePlayback)
					return
				}
			}
		}
	}
}
