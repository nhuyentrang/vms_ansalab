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
// @Tags         nvrs
// @Produce      json
// @Param   	 idNVR			query	string	true	"id nvr keyword"		minlength(1)  	maxlength(100)
// @Success      200   {object}  models.JsonDTORsp[models.DTO_NetworkConfig]
// @Router       /nvrs/config/networkconfig/{idNVR} [get]
// @Security	 BearerAuth
func GetNetworkConfig(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_NetworkConfig]()
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

	dtoNVRNetworkConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NetworkConfig, models.NetworkConfig]("id = ?", dtoNVRConfig.NetworkConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Data = dtoNVRNetworkConfig
	c.JSON(http.StatusOK, &jsonRsp)

}

// ReadNVRConfig godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         nvrs
// @Produce      json
// @Param   	 idNetWorkConfig			query	string	true	"id NetWorkConfig keyword"		minlength(1)  	maxlength(100)
// @Param   	 NetWorkConfigType		query	string	true	"items of NetWorkConfigType config"	Enums(tcpip,ddns,port,nat)
// @Param        NetWorkConfig  body      models.DTO_NetworkConfig  true  "NetWorkConfig JSON"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_NetworkConfig]
// @Router       /nvrs/config/networkconfig/{idNetWorkConfig} [put]
// @Security	 BearerAuth
func UpdateNetworkConfig(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_NetworkConfig]()
	idNetWorkConfig := c.Query("idNetWorkConfig")
	NetWorkConfigType := c.Query("NetWorkConfigType")

	var dto models.DTO_NetworkConfig
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Fetch NVRConfig from idNetworkConfig
	nvrConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVRConfig, models.NVRConfig]("network_config_id = ?", idNetWorkConfig)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Fetch NVR from NVRConfig
	dtoNVR, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVR, models.NVR]("config_id = ?", nvrConfig.ID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoNVR.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, errNetUpd := reposity.UpdateItemByIDFromDTO[models.DTO_NetworkConfig, models.NetworkConfig](idNetWorkConfig, dto)
	if errNetUpd != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = errNetUpd.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	requestUUID := uuid.New()
	// Prepare the command based on the NetWorkConfigType
	cmd := models.DeviceCommand{
		CommandID:     dtoDevice.ModelID,
		Cmd:           cmd_UpdateNetworkConfig,
		EventTime:     time.Now().Format(time.RFC3339),
		EventType:     NetWorkConfigType,
		IPAddress:     dtoNVR.IPAddress,
		UserName:      dtoNVR.Username,
		Password:      dtoNVR.Password,
		OnvifPort:     dtoNVR.OnvifPort,
		HttpPort:      dtoNVR.HttpPort,
		NetworkConfig: dto,
		ProtocolType:  dtoNVR.Protocol,
		RequestUUID:   requestUUID,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Wait for response
	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_UpdateNetworkConfig)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_UpdateNetworkConfig
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(requestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
				jsonRsp.Data = dto
				c.JSON(http.StatusOK, &jsonRsp)
				return
			}
		}
	}
}
