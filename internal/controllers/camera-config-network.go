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

// ReadCameraConfig godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         cameras
// @Produce      json
// @Param   	 id			path	string	true	"id Camera keyword"		minlength(1)  	maxlength(100)
// @Param        protocolType query    string  true  "Protocol Type" // assuming it's a query parameter, not path
// @Success      200   {object}  models.JsonDTORsp[models.DTO_NetworkConfig]
// @Router       /cameras/config/networkconfig/{id} [get]
// @Security	 BearerAuth
func GetNetworkConfigCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_NetworkConfig]()
	idCamera := c.Param("id")
	protocol := c.Query("protocolType")

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
		Cmd:          cmd_GetNetWorkConfig,
		EventTime:    time.Now().Format(time.RFC3339),
		ProtocolType: protocol,
		RequestUUID:  requestUUID,
		IPAddress:    dtoCamera.IPAddress,
		UserName:     dtoCamera.Username,
		Password:     dtoCamera.Password,
		OnvifPort:    dtoCamera.OnvifPort,
		HttpPort:     dtoCamera.HttpPort,
		CameraID:     idCamera,
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
			dtoCameraConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", dtoCamera.ConfigID)
			if err != nil {
				jsonRsp.Code = http.StatusInternalServerError
				jsonRsp.Message = err.Error()
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}

			dtoCameraNetworkConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NetworkConfig, models.NetworkConfig]("id = ?", dtoCameraConfig.NetworkConfigID)
			if err != nil {
				jsonRsp.Code = http.StatusInternalServerError
				jsonRsp.Message = err.Error()
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}
			jsonRsp.Data = dtoCameraNetworkConfig
			c.JSON(http.StatusOK, &jsonRsp)
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetOnSiteVideoConfig)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_GetOnSiteVideoConfig
			jsonRsp.Code = http.StatusRequestTimeout
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(requestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.NetworkConfig != nil {
					// Update entity from DTO
					_, err = reposity.UpdateItemByIDFromDTO[models.DTO_NetworkConfig, models.NetworkConfig](dtoCameraConfig.NetworkConfigID.String(), jsonRsp.Data)
					if err != nil {
						jsonRsp.Code = http.StatusInternalServerError
						jsonRsp.Message = err.Error()
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						return
					}
					jsonRsp.Data = *msg.PayLoad.NetworkConfig
					c.JSON(http.StatusOK, &jsonRsp)
					return
				} else {
					jsonRsp.Message = "No network config data received"
					jsonRsp.Code = http.StatusNoContent
					c.JSON(http.StatusNoContent, &jsonRsp)
					return
				}
			}
		}
	}
}

// ReadNVRConfig godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         cameras
// @Produce      json
// @Param   	 id			path	string	true	"id NetWorkConfig keyword"		minlength(1)  	maxlength(100)
// @Param   	 NetWorkConfigType		query	string	true	"items of NetWorkConfigType config"	Enums(tcpip,ddns,port,nat)
// @Param        NetWorkConfig  body      models.DTO_NetworkConfig  true  "NetWorkConfig JSON"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_NetworkConfig]
// @Router       /cameras/config/networkconfig/{id} [put]
// @Security	 BearerAuth
func UpdateNetworkConfigCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_NetworkConfig]()

	idCamera := c.Param("id")
	NetWorkConfigType := c.Query("NetWorkConfigType")
	fmt.Println("NetWorkConfigType: ", NetWorkConfigType)

	// Bind JSON request to DTO
	var dto models.DTO_NetworkConfig
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Invalid request body: " + err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Read CameraConfig by idNetworkConfig
	// Read Camera by cameraConfig
	dtoCamera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", idCamera)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to read Camera: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	camConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", dtoCamera.ConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to read CameraConfig: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Update NetworkConfig entity from DTO
	updatedDto, err := reposity.UpdateItemByIDFromDTO[models.DTO_NetworkConfig, models.NetworkConfig](camConfig.NetworkConfigID.String(), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to update NetworkConfig: " + err.Error()
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

	// Prepare the command based on NetWorkConfigType
	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:     dtoDevice.ModelID,
		Cmd:           cmd_UpdateNetworkConfig,
		EventTime:     time.Now().Format(time.RFC3339),
		IPAddress:     dtoCamera.IPAddress,
		UserName:      dtoCamera.Username,
		Password:      dtoCamera.Password,
		OnvifPort:     dtoCamera.OnvifPort,
		RequestUUID:   RequestUUID,
		HttpPort:      dtoCamera.HttpPort,
		NetworkConfig: updatedDto,
		ProtocolType:  dtoCamera.Protocol,
		EventType:     NetWorkConfigType,
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
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_UpdateNetworkConfig)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_UpdateNetworkConfig
			jsonRsp.Code = http.StatusRequestTimeout
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
				jsonRsp.Data = updatedDto
				jsonRsp.Code = http.StatusOK
				c.JSON(http.StatusOK, &jsonRsp)
				return
			}
		}
	}
}

// UpdateNetworkConfigCameras godoc
// @Summary      Update network configurations for cameras
// @Description  Updates network configurations for multiple cameras. Accepts a list of network configurations and applies them to the corresponding cameras.
// @Tags         cameras
// @Produce      json
// @Param        protocolType     query   string   false  "Protocol type for the device. Default is 'ONVIF'"
// @Param        channel          query   string   false  "Channel to apply the changes"
// @Param        NetworkConfig    body    []models.DTO_NetworkConfig  true  "List of network configurations for the cameras"
// @Success      200   {object}   models.JsonDTORsp[[]models.ChangeNetworkStatuses] "Command sent successfully"
// @Failure      400   {object}   models.JsonDTORsp[string] "Bad Request"
// @Failure      404   {object}   models.JsonDTORsp[string] "Device not found"
// @Failure      408   {object}   models.JsonDTORsp[string] "Request Timeout"
// @Failure      500   {object}   models.JsonDTORsp[string] "Internal Server Error"
// @Router       /cameras/config/networkconfig/update [put]
func UpdateNetworkConfigCameras(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.DTO_NetworkConfig]()
	protocol := c.Query("protocolType")
	if protocol == "" {
		protocol = "ONVIF"
	}
	channel := c.Query("channel")

	// Get new data from body
	var dtos []models.DTO_NetworkConfig
	if err := c.ShouldBindJSON(&dtos); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	var deviceName string
	var updatedDTOs []models.DTO_NetworkConfig
	var CamerasNetworkConfig []models.CamerasNetworkConfig
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
		UpdatedDTO, err := reposity.UpdateItemByIDFromDTO[models.DTO_NetworkConfig, models.NetworkConfig](dtoCameraConfig.ImageConfigID.String(), dto)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		updatedDTOs = append(updatedDTOs, UpdatedDTO)

		// Map VideoConfigInfo to StreamingChannelList
		cameraNetworkConfig := models.CamerasNetworkConfig{
			NetworkConfigCameras: UpdatedDTO,
			IPAddress:            dtoCamera.IPAddress,
			UserName:             dtoCamera.Username,
			Password:             dtoCamera.Password,
			HttpPort:             dtoCamera.HttpPort,
			OnvifPort:            dtoCamera.OnvifPort,
		}
		CamerasNetworkConfig = append(CamerasNetworkConfig, cameraNetworkConfig)
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
		CommandID:            dtoDevice.ModelID,
		Cmd:                  cmd_UpdateNetWorkConfigSeries,
		EventTime:            time.Now().Format(time.RFC3339),
		ProtocolType:         protocol,
		RequestUUID:          RequestUUID,
		Channel:              channel,
		CamerasNetworkConfig: CamerasNetworkConfig,
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
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_UpdateNetWorkConfigSeries)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_UpdateNetWorkConfigSeries
			jsonRsp.Code = http.StatusRequestTimeout
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.Status != "FAILURE" {
					jsonEditResp := models.NewJsonDTORsp[[]models.ChangeNetworkStatuses]()
					jsonEditResp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
					jsonEditResp.Data = msg.PayLoad.NetworkConfigStatuses
					jsonEditResp.Code = http.StatusInternalServerError
					c.JSON(http.StatusOK, &jsonEditResp)
					return
				} else {
					jsonRsp.Message = "Failed to update Camera Image Config: " + msg.PayLoad.Cmd
					jsonRsp.Data = *msg.PayLoad.NetworkConfigs
					jsonRsp.Code = http.StatusInternalServerError
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return
				}
			}
		}
	}
}
