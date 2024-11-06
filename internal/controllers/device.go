package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"vms/internal/models"
	"vms/statuscode"

	"vms/comongo/kafkaclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetDevice		godoc
// @Summary      	Get all Device groups with query filter
// @Description  	Responds with the list of all Device as JSON.
// @Tags         	Device
// @Param   		keyword			query	string	false	"MacAddress"		minlength(1)  	maxlength(100)
// @Produce      	json
// @Param        idDevice      query   string  false    "idDevice keyword"       minlength(1)   maxlength(100)
// @Param        userID        query   string  false    "userID keyword"         minlength(1)   maxlength(100)
// @Param        nameDevice    query   string  false    "nameDevice keyword"     minlength(1)   maxlength(100)
// @Param        serial        query   string  false    "serial keyword"         minlength(1)   maxlength(100)
// @Param        deviceType    query   string  false    "deviceType keyword"     minlength(1)   maxlength(100)
// @Param        modelID       query   string  false    "modelID keyword"        minlength(1)   maxlength(100)
// @Param        appVersion    query   string  false    "appVersion keyword"     minlength(1)   maxlength(100)
// @Param        areaId        query   string  false    "areaId keyword"         minlength(1)   maxlength(100)
// @Param        areaName      query   string  false    "areaName keyword"       minlength(1)   maxlength(100)
// @Param        httpPort      query   string  false    "httpPort keyword"       minlength(1)   maxlength(100)
// @Param        onvifPort     query   string  false    "onvifPort keyword"      minlength(1)   maxlength(100)
// @Param        deviceCode    query   string  false    "deviceCode keyword"     minlength(1)   maxlength(100)
// @Param        status        query   string  false    "status keyword"         minlength(1)   maxlength(100)
// @Param        location      query   string  false    "location keyword"       minlength(1)   maxlength(100)
// @Param        ipAddress     query   string  false    "ipAddress keyword"      minlength(1)   maxlength(100)
// @Param        macAddress    query   string  false    "macAddress keyword"     minlength(1)   maxlength(100)
// @Param        mqttAccount   query   string  false    "mqttAccount keyword"    minlength(1)   maxlength(100)
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_Device]
// @Router       	/device [get]
// @Security		BearerAuth
func GetDevices(c *gin.Context) {
	jsonRsp := models.NewJsonDTOListRsp[models.DTO_Device]()

	// Get param
	idDevice := c.Query("idDevice")
	userID := c.Query("userID")
	nameDevice := c.Query("nameDevice")
	serial := c.Query("serial")
	deviceType := c.Query("deviceType")
	modelID := c.Query("modelID")
	appVersion := c.Query("appVersion")
	areaId := c.Query("areaId")
	areaName := c.Query("areaName")
	httpPort := c.Query("httpPort")
	onvifPort := c.Query("onvifPort")
	deviceCode := c.Query("deviceCode")
	status := c.Query("status")
	location := c.Query("location")
	ipAddress := c.Query("ipAddress")
	macAddress := c.Query("macAddress")
	mqttAccount := c.Query("mqttAccount")

	// Build query
	query := reposity.NewQuery[models.DTO_Device, models.Device]()

	// Search for keyword in name
	if idDevice != "" {
		query.AddConditionOfTextField("AND", "id", "=", idDevice)
	}
	if userID != "" {
		query.AddConditionOfTextField("AND", "user_id", "=", userID)
	}
	if nameDevice != "" {
		query.AddConditionOfTextField("AND", "name_device", "=", nameDevice)
	}
	if serial != "" {
		query.AddConditionOfTextField("AND", "serial", "=", serial)
	}
	if deviceType != "" {
		query.AddConditionOfTextField("AND", "device_type", "=", deviceType)
	}
	if modelID != "" {
		query.AddConditionOfTextField("AND", "model_id", "=", modelID)
	}
	if appVersion != "" {
		query.AddConditionOfTextField("AND", "app_version", "=", appVersion)
	}
	if areaId != "" {
		query.AddConditionOfTextField("AND", "area_id", "=", areaId)
	}
	if areaName != "" {
		query.AddConditionOfTextField("AND", "area_name", "=", areaName)
	}
	if httpPort != "" {
		query.AddConditionOfTextField("AND", "http_port", "=", httpPort)
	}
	if onvifPort != "" {
		query.AddConditionOfTextField("AND", "onvif_port", "=", onvifPort)
	}
	if deviceCode != "" {
		query.AddConditionOfTextField("AND", "device_code", "=", deviceCode)
	}
	if status != "" {
		query.AddConditionOfTextField("AND", "status", "=", status)
	}
	if location != "" {
		query.AddConditionOfTextField("AND", "location", "=", location)
	}
	if ipAddress != "" {
		query.AddConditionOfTextField("AND", "ip_address", "=", ipAddress)
	}
	if macAddress != "" {
		query.AddConditionOfTextField("AND", "mac_address", "=", macAddress)
	}
	if mqttAccount != "" {
		query.AddConditionOfTextField("AND", "mqtt_account", "=", mqttAccount)
	}

	// Exec query
	dtoDeviceBasics, count, err := query.ExecNoPaging("-created_at")
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	jsonRsp.Count = count
	jsonRsp.Data = dtoDeviceBasics
	c.JSON(http.StatusOK, &jsonRsp)
}

// ReadDevice godoc
// @Summary      Get single device by id
// @Description  Returns the device whose ID value matches the idDevice.
// @Tags         Device
// @Produce      json
// @Param        id      path   string  true    "Device ID"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_Device]
// @Router       /device/search/{id} [get]
// @Security     BearerAuth
func ReadDevice(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_Device]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateDevice		 	godoc
// @Summary      	Update single Device by id
// @Description  	Updates and returns a single Device whose ID value matches the id. New data must be passed in the body.
// @Tags         	Device
// @Produce      	json
// @Param        	id  path  string  true  "Update Device by id"
// @Param        	Device  body      models.DTO_Device  true  "Device JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_Device]
// @Router       	/device/{id} [put]
// @Security		BearerAuth
func UpdateDevice(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_Device]()

	// Get new data from body
	var dto models.DTO_Device
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_Device, models.Device](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Override ID
	dto.ID, _ = uuid.Parse(c.Param("id"))

	// Return
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// CreateDevice		godoc
// @Summary      	Create a new Device
// @Description  	Takes a Device JSON and store in DB. Return saved JSON.
// @Tags         	Device
// @Produce			json
// @Param        	Device  body   models.DTO_Device  true  "Device JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_Device]
// @Router       	/device [post]
// @Security		BearerAuth
func CreateDevice(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_Device]()

	// Call BindJSON to bind the received JSON to
	var dtoResp models.DTO_Device
	if err := c.BindJSON(&dtoResp); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	dto, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Device, models.Device]("mac_address = ?", dtoResp.MacAddress)
	if err != nil {
		dto, err = reposity.CreateItemFromDTO[models.DTO_Device, models.Device](dtoResp)
		if err != nil {
			return
		}
	} else {
		dto, err = reposity.UpdateItemByIDFromDTO[models.DTO_Device, models.Device](dto.ID.String(), dtoResp)
		if err != nil {
			return
		}
	}

	jsonRsp.Data = dto
	c.JSON(http.StatusCreated, &jsonRsp)
}

// DeleteDevice godoc
// @Summary      Remove single device by id
// @Description  Delete a single entry from the repository based on id, and also delete associated cameras.
// @Tags         Device
// @Produce      json
// @Param        id  path  string  true  "Delete device by id"
// @Success      204
// @Router       /device/{id} [delete]
// @Security     BearerAuth
func DeleteDevice(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_Device]()

	deviceID := c.Param("id")

	// Fetch the device to ensure it exists
	_, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](deviceID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	// Build query to fetch all cameras associated with the device's box ID
	query := reposity.NewQuery[models.DTOCamera, models.Camera]()
	jsonCondition := fmt.Sprintf("[{\"id\":\"%s\"}]", deviceID)
	query.AddConditionOfTextField("AND", "box", "@>", jsonCondition)

	// Execute query to get associated cameras
	cameras, count, err := query.ExecNoPaging("-created_at")
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to fetch cameras: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	if count > 0 {
		// Delete each camera associated with the device
		for _, camera := range cameras {
			err := reposity.DeleteItemByID[models.Camera](camera.ID.String())
			if err != nil {
				jsonRsp.Code = http.StatusInternalServerError
				jsonRsp.Message = "Failed to delete camera: " + err.Error()
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}
		}
	}

	// Delete the device
	err = reposity.DeleteItemByID[models.Device](deviceID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	c.JSON(http.StatusNoContent, &jsonRsp)
}

// ReadDevice		 godoc
// @Summary      Get single Device by id
// @Description  Returns the Device whose ID value matches the id.
// @Tags         Device
// @Produce      json
// @Param   	 idNVR			query	string	true	"id nvr keyword"		minlength(1)  	maxlength(100)
// @Success      200
// @Router       /device/hik [get]
// @Security	 BearerAuth
func ReadNetworkInterfaces(c *gin.Context) {
	jsonRsp := models.NewJsonDTOListRsp[models.DTO_NetworkInterfaces]()
	nvrConfigID := c.Query("idNVR")

	// Build query
	query := reposity.NewQuery[models.DTO_NetworkInterfaces, models.NetworkInterfaces]()
	if nvrConfigID != "" {
		query.AddConditionOfTextField("AND", "nvr_config_id", "=", nvrConfigID)
	}
	// Exec query
	dtoDeviceBasics, count, err := query.ExecNoPaging("-created_at")
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Count = count
	jsonRsp.Data = dtoDeviceBasics
	c.JSON(http.StatusOK, &jsonRsp)
}

// ReadDevice		 godoc
// @Summary      Get single Device by id
// @Description  Returns the Device whose ID value matches the id.
// @Tags         Device
// @Produce      json
// @Param        id  path  string  true  "Update network config by id"
// @Param        command  path  string  true  "Update network config by cmd"
// @Param        device  body      models.DTO_Update_NetworkInterfaces  true  "Network JSON"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_Update_NetworkInterfaces]
// @Router       /device/hik/{id}/{command} [put]
// @Security	 BearerAuth
func UpdateNetworkInterfaces(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_Update_NetworkInterfaces]()

	// Get new data from body
	var dto models.DTO_Update_NetworkInterfaces
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_Update_NetworkInterfaces, models.NetworkInterfaces](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	//TODO need to add cameraID or NVR ID to this network config ID
	//dtoCamera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", dto.CameraID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	cmd := models.DeviceCommand{
		CommandID: uuid.New().String(),
		Cmd:       c.Param("command"),
		EventTime: time.Now().Format(time.RFC3339),
		//ProtocolType: dtoDevice.Protocol,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	for {
		select {
		case <-time.After(10 * time.Second):
			fmt.Println("\t\t> error, waiting for command update firmware result timed out")
			jsonRsp.Message = "No response from the device, timeout"
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return
		case msg := <-UpdateNetworkConfigCameraChannelDataReceiving:
			jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
			jsonRsp.Data = dto
			c.JSON(http.StatusOK, &jsonRsp)
			return
		}
	}
}

// GetCamerasByDeviceID godoc
// @Summary      List all cameras attached to a device
// @Description  Returns the list of all cameras associated with the specified device ID.
// @Tags         Device
// @Produce      json
// @Param        deviceID  path   string  true   "Device ID"
// @Param        sort         query   string  false "Sort"                      default(-created_at)
// @Success      200       {object}  models.JsonDTOListRsp[models.DTOCamera]
// @Router       /device/{deviceID}/cameras [get]
// @Security     BearerAuth
func GetCamerasByDeviceID(c *gin.Context) {
	deviceID := c.Param("deviceID")
	jsonRsp := models.NewJsonDTOListRsp[models.DTOCamera]()
	sort := c.Query("sort")

	// Build query
	query := reposity.NewQuery[models.DTOCamera, models.Camera]()

	if deviceID != "" {
		// Create JSON condition for querying JSONB field
		jsonCondition := fmt.Sprintf(`{"id":"%s"}`, deviceID)
		query.AddConditionOfTextField("WHERE", "box", "@>", jsonCondition)
	}

	// Exec query
	dtos, count, err := query.ExecNoPaging(sort)
	if err != nil {
		jsonRsp.Code = statuscode.StatusSearchItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Count = count
	jsonRsp.Data = dtos
	jsonRsp.Size = int64(len(dtos))
	c.JSON(http.StatusOK, &jsonRsp)
}

// GetActiveDevices godoc
// @Summary      Get all active devices
// @Description  Responds with the list of all active devices (status = connected) as JSON.
// @Tags         Device
// @Produce      json
// @Success      200  {object}  models.JsonDTOListRsp[models.DTO_Device]
// @Param        status         query   string  false "device status"                      default(connected)
// @Router       /device/active [get]
// @Security     BearerAuth
func GetActiveDevices(c *gin.Context) {
	jsonRsp := models.NewJsonDTOListRsp[models.DTO_Device]()
	deviceStatus := c.Query("status")

	// Build query
	query := reposity.NewQuery[models.DTO_Device, models.Device]()
	if deviceStatus != "" {
		query.AddConditionOfTextField("AND", "status", "=", deviceStatus)
	}
	// Exec query
	dtoDeviceBasics, count, err := query.ExecNoPaging("-created_at")
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Count = count
	jsonRsp.Data = dtoDeviceBasics
	c.JSON(http.StatusOK, &jsonRsp)
}

// func ReadNetworkInterfaces(c *gin.Context) {
// 	jsonRsp := models.NewJsonDTOListRsp[models.DTO_NetworkInterfaces]()
// 	nvrConfigID := c.Query("idNVR")

// 	// Build query
// 	query := reposity.NewQuery[models.DTO_NetworkInterfaces, models.NetworkInterfaces]()
// 	if nvrConfigID != "" {
// 		query.AddConditionOfTextField("AND", "nvr_config_id", "=", nvrConfigID)
// 	}
// 	// Exec query
// 	dtoDeviceBasics, count, err := query.ExecNoPaging("-created_at")
// 	if err != nil {
// 		jsonRsp.Code = http.StatusInternalServerError
// 		jsonRsp.Message = err.Error()
// 		c.JSON(http.StatusInternalServerError, &jsonRsp)
// 		return
// 	}

// 	jsonRsp.Count = count
// 	jsonRsp.Data = dtoDeviceBasics
// 	c.JSON(http.StatusOK, &jsonRsp)
// }

// SynchronizeCameraDataConfig godoc
// @Summary      Synchronize camera data config
// @Description  Synchronize camera data config based on model ID
// @Tags         Device
// @Produce      json
// @Param        model_id  query  string  true  "Model ID of the device"
// @Success      200  {object}  models.JsonDTORsp[models.DTOCamera]  "Successfully synchronized camera data config"
// @Failure      400  "Bad Request - Invalid input"
// @Failure      404  "Not Found - Device or Camera not found"
// @Failure      500  "Internal Server Error - Unable to synchronize camera data config"
// @Router       /device/cameras/synchronize [get]
// @Security     BearerAuth
func SynchronizeCameraDataConfig(c *gin.Context) {
	modelID := c.Query("model_id")
	jsonRsp := models.NewJsonDTOListRsp[models.DTOCamera]()

	log.Printf("Received request to synchronize camera data config for model ID: %s", modelID)

	// Fetch device by model ID
	device, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Device, models.Device]("model_id = ?", modelID)
	if err != nil {
		log.Printf("Device not found for model ID: %s, error: %v", modelID, err)
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	log.Printf("Device found: %+v", device)

	log.Printf("Updating status for Device with model ID: %s", modelID)
	updateDeviceStatus("smartnvr", device.ModelID, true)

	// Fetch all cameras with the device ID
	query := reposity.NewQuery[models.DTOCamera, models.Camera]()
	jsonCondition := fmt.Sprintf("{\"id\":\"%s\"}", device.ID)
	query.AddConditionOfTextField("AND", "box", "@>", jsonCondition)
	cameras, count, err := query.ExecNoPaging("-created_at")
	if err != nil {
		log.Printf("Failed to retrieve cameras for device ID: %s, error: %v", device.ID.String(), err)
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to retrieve cameras: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	if count == 0 {
		log.Printf("No cameras found for device ID: %s", device.ID.String())
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "No cameras found for the device"
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	log.Printf("Found %d cameras for device ID: %s", count, device.ID.String())

	for _, camera := range cameras {
		log.Printf("Processing camera: %+v", camera)
		// Create a map for camera channels
		cameraMap := make(map[string]models.ChannelCamera)
		for _, stream := range camera.Streams {
			channel := models.ChannelCamera{
				OnDemand: true,
				Url:      stream.URL,
				Codec:    stream.Codec,
				Name:     stream.Name,
			}
			cameraMap[stream.Channel] = channel
		}
		dataFileConfig := map[string]models.ConfigCamera{}
		newCamera := models.ConfigCamera{
			NameCamera: camera.Name,
			IP:         camera.IPAddress,
			UserName:   camera.Username,
			PassWord:   camera.Password,
			HTTPPort:   camera.HttpPort,
			RTSPPort:   camera.HttpPort, // Assuming RTSPPort is same as HttpPort
			OnvifPort:  camera.OnvifPort,
			ChannelNVR: camera.NVR.Channel,
			Channels:   cameraMap,
			IDNVR:      camera.NVR.ID,
		}
		dataFileConfig[camera.ID.String()] = newCamera

		// Send command to synchronize camera data config
		requestUUID := uuid.New()
		cmdAddConfig := models.DeviceCommand{
			CommandID:    device.ModelID,
			Cmd:          cmd_AddDataConfig,
			EventTime:    time.Now().Format(time.RFC3339),
			EventType:    "camera",
			ConfigCamera: dataFileConfig,
			ProtocolType: camera.Protocol,
			RequestUUID:  requestUUID,
		}
		cmsStr, err := json.Marshal(cmdAddConfig)
		if err != nil {
			log.Printf("Failed to marshal command for camera ID: %s, error: %v", camera.ID, err)
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to marshal command: " + err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		log.Printf("Sending command to Kafka: %s", cmsStr)
		kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	}

	log.Println("Successfully synchronized camera data config")
	jsonRsp.Message = "Successfully synchronized camera data config"
	c.JSON(http.StatusOK, &jsonRsp)
}

// SynchronizeCameraDataConfigForCamera godoc
// @Summary      Synchronize camera data config for a specific camera
// @Description  Synchronize camera data config based on camera ID
// @Tags         Device
// @Produce      json
// @Param        camera_id  query  string  true  "Camera ID"
// @Success      200  {object}  models.JsonDTORsp[models.DTOCamera]  "Successfully synchronized camera data config"
// @Failure      400  "Bad Request - Invalid input"
// @Failure      404  "Not Found - Device or Camera not found"
// @Failure      500  "Internal Server Error - Unable to synchronize camera data config"
// @Router       /device/camera/synchronize [get]
// @Security     BearerAuth
func SynchronizeCameraDataConfigForCamera(c *gin.Context) {
	cameraID := c.Query("camera_id")
	jsonRsp := models.NewJsonDTOListRsp[models.DTOCamera]()

	log.Printf("Received request to synchronize camera data config for camera ID: %s", cameraID)

	// Fetch camera by camera ID
	camera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", cameraID)
	if err != nil {
		log.Printf("Camera not found for camera ID: %s, error: %v", cameraID, err)
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Camera not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	log.Printf("Camera found: %+v", camera)

	// Fetch device by box ID from camera
	device, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](camera.Box.ID)
	if err != nil {
		log.Printf("Device not found for box ID: %s, error: %v", camera.Box.ID, err)
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	log.Printf("Device found: %+v", device)

	log.Printf("Updating status for Device with model ID: %s", device.ModelID)
	updateDeviceStatus("smartnvr", device.ModelID, true)

	log.Printf("Processing camera: %+v", camera)
	// Create a map for camera channels
	cameraMap := make(map[string]models.ChannelCamera)
	for _, stream := range camera.Streams {
		channel := models.ChannelCamera{
			OnDemand: true,
			Url:      stream.URL,
			Codec:    stream.Codec,
			Name:     stream.Name,
		}
		cameraMap[stream.Channel] = channel
	}
	dataFileConfig := map[string]models.ConfigCamera{}
	newCamera := models.ConfigCamera{
		NameCamera: camera.Name,
		IP:         camera.IPAddress,
		UserName:   camera.Username,
		PassWord:   camera.Password,
		HTTPPort:   camera.HttpPort,
		RTSPPort:   camera.HttpPort, // Assuming RTSPPort is same as HttpPort
		OnvifPort:  camera.OnvifPort,
		ChannelNVR: camera.NVR.Channel,
		IDNVR:      camera.NVR.ID,
		Channels:   cameraMap,
	}
	dataFileConfig[camera.ID.String()] = newCamera

	// Send command to synchronize camera data config
	requestUUID := uuid.New()
	cmdAddConfig := models.DeviceCommand{
		CommandID:    device.ModelID,
		Cmd:          cmd_AddDataConfig,
		EventTime:    time.Now().Format(time.RFC3339),
		EventType:    "camera",
		ConfigCamera: dataFileConfig,
		ProtocolType: camera.Protocol,
		RequestUUID:  requestUUID,
	}
	cmsStr, err := json.Marshal(cmdAddConfig)
	if err != nil {
		log.Printf("Failed to marshal command for camera ID: %s, error: %v", camera.ID, err)
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to marshal command: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	log.Printf("Sending command to Kafka: %s", cmsStr)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	log.Println("Successfully synchronized camera data config")
	jsonRsp.Message = "Successfully synchronized camera data config"
	c.JSON(http.StatusOK, &jsonRsp)
}

// SynchronizeNVRDataConfigForNVR godoc
// @Summary      Synchronize NVR data config for a specific NVR
// @Description  Synchronize NVR data config based on NVR ID
// @Tags         Device
// @Produce      json
// @Param        nvr_id  query  string  true  "NVR ID"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_NVR]  "Successfully synchronized NVR data config"
// @Failure      400  "Bad Request - Invalid input"
// @Failure      404  "Not Found - Device or NVR not found"
// @Failure      500  "Internal Server Error - Unable to synchronize NVR data config"
// @Router       /device/nvr/synchronize [get]
// @Security     BearerAuth
func SynchronizeNVRDataConfigForNVR(c *gin.Context) {
	nvrID := c.Query("nvr_id")
	jsonRsp := models.NewJsonDTOListRsp[models.DTO_NVR]()
	cameras := models.KeyValueArray{}

	log.Printf("Received request to synchronize NVR data config for NVR ID: %s", nvrID)

	// Fetch NVR by NVR ID
	nvr, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVR, models.NVR]("id = ?", nvrID)
	if err != nil {
		log.Printf("NVR not found for NVR ID: %s, error: %v", nvrID, err)
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "NVR not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	log.Printf("NVR found: %+v", nvr)

	// Fetch device by box ID from NVR
	device, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Device, models.Device]("id = ?", nvr.Box.ID)
	if err != nil {
		log.Printf("Device not found for model ID: %s, error: %v", nvr.Box.ID, err)
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	log.Printf("Device found: %+v", device)

	log.Printf("Updating status for Device with model ID: %s", device.ModelID)
	updateDeviceStatus("smartnvr", device.ModelID, true)

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    device.ModelID,
		RequestUUID:  RequestUUID,
		Cmd:          cmd_GetNetWorkConfig,
		EventTime:    time.Now().Format(time.RFC3339),
		IPAddress:    nvr.IPAddress,
		UserName:     nvr.Username,
		Password:     nvr.Password,
		HttpPort:     nvr.HttpPort,
		OnvifPort:    "80",
		Channel:      "101",
		Track:        "103",
		ProtocolType: nvr.Protocol,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetNetWorkConfig)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_GetNetWorkConfig
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.NetworkConfig != nil {
					dtoNetWorkConfig, err := reposity.CreateItemFromDTO[models.DTO_NetworkConfig, models.NetworkConfig](*msg.PayLoad.NetworkConfig)
					if err != nil {
						jsonRsp.Code = http.StatusInternalServerError
						jsonRsp.Message = err.Error()
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						return
					}

					log.Printf("Processing NVR: %+v", nvr)

					// Create a map for NVR's cameras
					if nvr.Cameras != nil {
						for _, camera := range *nvr.Cameras {
							cam := models.KeyValue{
								ID:      camera.ID,
								Name:    camera.Name,
								Channel: camera.Channel,
							}
							cameras = append(cameras, cam)
						}
					}

					// Create a map of camera for NVR
					dataFileConfig := map[string]models.ConfigNVR{}
					newNVR := models.ConfigNVR{
						NameCamera: nvr.Name,
						IP:         nvr.IPAddress,
						UserName:   nvr.Username,
						PassWord:   nvr.Password,
						HTTPPort:   strconv.Itoa(dtoNetWorkConfig.HTTP),
						RTSPPort:   strconv.Itoa(dtoNetWorkConfig.RTSP),
						OnvifPort:  nvr.OnvifPort,
						Cameras:    cameras,
					}
					dataFileConfig[nvr.ID.String()] = newNVR

					// Send command to synchronize NVR data config
					requestUUID := uuid.New()
					cmdAddConfig := models.DeviceCommand{
						CommandID:    device.ModelID,
						Cmd:          cmd_AddDataConfig,
						EventTime:    time.Now().Format(time.RFC3339),
						EventType:    "nvr",
						ConfigNVR:    dataFileConfig,
						ProtocolType: nvr.Protocol,
						RequestUUID:  requestUUID,
					}
					cmsStr, err := json.Marshal(cmdAddConfig)
					if err != nil {
						log.Printf("Failed to marshal command for NVR ID: %s, error: %v", nvr.ID.String(), err)
						jsonRsp.Code = http.StatusInternalServerError
						jsonRsp.Message = "Failed to marshal command: " + err.Error()
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						return
					}
					log.Printf("Sending command to Kafka: %s", cmsStr)
					kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

					log.Println("Successfully synchronized NVR data config")
					jsonRsp.Message = "Successfully synchronized NVR data config"
					c.JSON(http.StatusOK, &jsonRsp)
					return // Add this return to stop further execution after success
				} else {
					jsonRsp.Message = "Failed to synchronize NVR data config"
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return // Ensure the function exits after sending an error response
				}
			}
		}
	}
}

// SynchronizeNVRDataConfig godoc
// @Summary      Synchronize NVR data config
// @Description  Synchronize NVR data config based on model ID
// @Tags         Device
// @Produce      json
// @Param        model_id  query  string  true  "Model ID of the device"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_NVR]  "Successfully synchronized NVR data config"
// @Failure      400  "Bad Request - Invalid input"
// @Failure      404  "Not Found - Device or NVR not found"
// @Failure      500  "Internal Server Error - Unable to synchronize NVR data config"
// @Router       /device/nvrs/synchronize [get]
// @Security     BearerAuth
func SynchronizeNVRDataConfig(c *gin.Context) {
	modelID := c.Query("model_id")
	jsonRsp := models.NewJsonDTOListRsp[models.DTO_NVR]()

	log.Printf("Received request to synchronize NVR data config for model ID: %s", modelID)

	// // Fetch device by model ID
	device, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Device, models.Device]("model_id = ?", modelID)
	if err != nil {
		log.Printf("Device not found for model ID: %s, error: %v", modelID, err)
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	log.Printf("Device found: %+v", device)
	log.Printf("Updating status for Device with model ID: %s", modelID)
	updateDeviceStatus("smartnvr", modelID, true)

	// Fetch all NVRs with the device ID
	query := reposity.NewQuery[models.DTO_NVR, models.NVR]()
	jsonCondition := fmt.Sprintf("{\"id\":\"%s\"}", device.ID.String())
	query.AddConditionOfTextField("AND", "box", "@>", jsonCondition)
	nvrs, count, err := query.ExecNoPaging("-created_at")
	if err != nil {
		log.Printf("Failed to retrieve NVRs for device ID: %s, error: %v", device.ID.String(), err)
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to retrieve NVRs: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	if count == 0 {
		log.Printf("No NVRs found for device ID: %s", device.ID.String())
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "No NVRs found for the device"
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	log.Printf("Found %d NVRs for device ID: %s", count, device.ID.String())

	for _, nvr := range nvrs {
		log.Printf("Processing NVR: %+v", nvr)
		// Create a map for nvr's cameras

		cameras := models.KeyValueArray{}
		for _, camera := range *nvr.Cameras {

			camera := models.KeyValue{
				ID:      camera.ID,
				Name:    camera.Name,
				Channel: camera.Channel,
			}
			cameras = append(cameras, camera)

		}
		// Create a map of camera for NVR
		dataFileConfig := map[string]models.ConfigNVR{}
		newCamera := models.ConfigNVR{
			NameCamera: nvr.Name,
			IP:         nvr.IPAddress,
			UserName:   nvr.Username,
			PassWord:   nvr.Password,
			HTTPPort:   nvr.HttpPort,
			RTSPPort:   nvr.RtspPort,
			OnvifPort:  nvr.OnvifPort,
			Cameras:    cameras,
		}
		dataFileConfig[nvr.ID.String()] = newCamera

		// Send command to synchronize NVR data config
		requestUUID := uuid.New()
		cmdAddConfig := models.DeviceCommand{
			CommandID:    modelID,
			Cmd:          cmd_AddDataConfig,
			EventTime:    time.Now().Format(time.RFC3339),
			EventType:    "nvr",
			ConfigNVR:    dataFileConfig,
			ProtocolType: nvr.Protocol,
			RequestUUID:  requestUUID,
		}
		cmsStr, err := json.Marshal(cmdAddConfig)
		if err != nil {
			log.Printf("Failed to marshal command for NVR ID: %s, error: %v", nvr.ID.String(), err)
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to marshal command: " + err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		log.Printf("Sending command to Kafka: %s", cmsStr)
		kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	}

	log.Println("Successfully synchronized NVR data config")
	jsonRsp.Message = "Successfully synchronized NVR data config"
	c.JSON(http.StatusOK, &jsonRsp)
}
