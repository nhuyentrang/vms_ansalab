package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"vms/internal/models"

	"vms/comongo/kafkaclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadCamera		 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         device
// @Produce      json
// @Param        device  body      models.ScanDevice  true  "Device JSON"
// @Param   	 keyword			query	string	false	"device name keyword"		minlength(1)  	maxlength(100)
// @Param 		 protocol query string false "Device Protocol"
// @Param   	 device_type	query	string	false	"Type of device"	Enums(nvr, camera)
// @Success      200   {object}  models.JsonDTORsp[models.ScanDevice]
// @Router       /device/onvif/scandevice [post]
// @Security	 BearerAuth
func DeviceScanOnvif(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.ScanDevice]()

	// Get new data from body
	var dto models.ScanDevice
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dto.IDBox)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		Cmd:          cmd_ScanDevice,
		EventTime:    time.Now().Format(time.RFC3339),
		UserName:     dto.UserName,
		Password:     dto.PassWord,
		ProtocolType: dto.Protocol,
		RequestUUID:  RequestUUID,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_ScanDevice)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_ScanDevice
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return
		//case msg := <-DeviceScanOnvifChannelDataReceiving:
		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.DeviceScan == nil {
					jsonRsp.Message = "Device does not exist"
					c.JSON(http.StatusNoContent, &jsonRsp)
					return
				} else {
					jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
					jsonRsp.Data = *msg.PayLoad.DeviceScan
					c.JSON(http.StatusOK, &jsonRsp)
					return
				}
			}
		}
	}
}

// ReadCamera		 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         device
// @Produce      json
// @Param        device  body      models.ScanDevice  true  "Device JSON"
// @Param   	 keyword			query	string	false	"device name keyword"		minlength(1)  	maxlength(100)
// @Param 		 protocol query string false "Device Protocol"
// @Param   	 device_type	query	string	false	"Type of device"	Enums(nvr, camera)
// @Success      200   {object}  models.JsonDTORsp[models.ScanDevice]
// @Router       /device/scandevice [post]
// @Security	 BearerAuth
func DeviceScan(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.ScanDevice]()

	// Get new data from body
	var dto models.ScanDevice
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dto.IDBox)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	if dto.Protocol != "ONVIF" && dto.Protocol != "HIKVISION" {
		fmt.Println("\t\t> error, Protocol is not yet supported")
		jsonRsp.Message = "Error, Protocol is not yet supported"
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		RequestUUID:  RequestUUID,
		Cmd:          cmd_ScanDevice,
		EventTime:    time.Now().Format(time.RFC3339),
		UserName:     dto.UserName,
		Password:     dto.PassWord,
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
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_ScanDevice)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_ScanDevice
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.DeviceScan == nil {
					jsonRsp.Message = "Device does not exist"
					c.JSON(http.StatusNoContent, &jsonRsp)
					return
				} else {
					jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
					jsonRsp.Data = *msg.PayLoad.DeviceScan
					c.JSON(http.StatusOK, &jsonRsp)
					return
				}
			}
		}
	}
}

// ReadCamera		 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         device
// @Produce      json
// @Param        device  body      models.ScanDevice  true  "Device JSON"
// @Param   	 keyword	query	string	false	"device name keyword"		minlength(1)  	maxlength(100)
// @Param 		 protocol query string false "Device Protocol"
// @Success      200   {object}  models.JsonDTORsp[models.ScanDevice]
// @Router       /device/onvif/scandevicestaticip [post]
// @Security	 BearerAuth
func DeviceScanStaticIP(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.ScanDevice]()

	// Get new data from body
	var dto models.ScanDevice
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dto.IDBox)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		RequestUUID:  RequestUUID,
		Cmd:          cmd_ScanDeviceIP,
		EventTime:    time.Now().Format(time.RFC3339),
		IPAddress:    dto.URL,
		HttpPort:     dto.Host,
		UserName:     dto.UserName,
		Password:     dto.PassWord,
		ProtocolType: "ONVIF",
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_ScanDeviceIP)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_ScanDeviceIP
			jsonRsp.Data = []models.ScanDevice{
				{
					URL:      dto.URL,
					Host:     dto.Host,
					Name:     "",
					Type:     "",
					Location: "",
					StartURL: "",
				},
			}
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.Status == "FAILURE" {
					jsonRsp.Message = "Device doesn't support ONVIF Protocol"
					c.JSON(http.StatusNoContent, &jsonRsp)
					return
				} else {
					if msg.PayLoad.DeviceScan == nil {
						jsonRsp.Message = "Device does not exist"
						c.JSON(http.StatusNoContent, &jsonRsp)
						return
					} else {
						jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
						jsonRsp.Data = *msg.PayLoad.DeviceScan
						c.JSON(http.StatusOK, &jsonRsp)
						return
					}
				}

			}
		}
	}
}

// ReadCamera		 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         device
// @Produce      json
// @Param        device  body      models.ScanDevice  true  "Device JSON"
// @Param 		 protocol query string false "Device Protocol"
// @Param        idbox  path  string  true  "IDBox"
// @Success      200   {object}  models.JsonDTORsp[models.ScanDevice]
// @Router       /device/onvif/devicescanhikivision [post]
// @Security	 BearerAuth
func DeviceScanHikivision(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.ScanDevice]()

	// Get new data from body
	var dto models.ScanDevice
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dto.IDBox)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		RequestUUID:  RequestUUID,
		Cmd:          cmd_ScanDevice,
		EventTime:    time.Now().Format(time.RFC3339),
		UserName:     dto.UserName,
		Password:     dto.PassWord,
		ProtocolType: dto.Protocol,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(120 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_ScanDevice)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_ScanDevice
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.DeviceScan == nil {
					jsonRsp.Message = "Device does not exist"
					c.JSON(http.StatusOK, &jsonRsp)
					return
				} else {
					jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
					jsonRsp.Data = *msg.PayLoad.DeviceScan
					c.JSON(http.StatusOK, &jsonRsp)
					return
				}
			}
		}
	}
}

// ReadCamera		 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         device
// @Produce      json
// @Param        device  body      models.ScanDevice  true  "Device JSON"
// @Param        idbox  path  string  true  "IDBox"
// @Param 		 protocol query string false "Device Protocol"
// @Param   	 device_type	query	string	false	"Type of device"	Enums(nvr, camera)
// @Success      200   {object}  models.JsonDTORsp[models.ScanDevice]
// @Router       /device/onvif/devicescanlistip [post]
// @Security	 BearerAuth
func DeviceScanIPList(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.ScanDevice]()

	// Get new data from body
	var dto models.ScanDevice
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dto.IDBox)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		RequestUUID:  RequestUUID,
		Cmd:          cmd_ScanDeviceListIP,
		EventTime:    time.Now().Format(time.RFC3339),
		UserName:     dto.UserName,
		Password:     dto.PassWord,
		ProtocolType: "ONVIF",
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_ScanDeviceListIP)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_ScanDeviceListIP
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.Status == "FAILURE" {
					jsonRsp.Message = "Device doesn't support ONVIF Protocol"
					c.JSON(http.StatusNoContent, &jsonRsp)
					return
				} else {
					if msg.PayLoad.DeviceScan == nil {
						jsonRsp.Message = "Device does not exist"
						c.JSON(http.StatusNoContent, &jsonRsp)
						return
					} else {
						filteredMsgs := models.KafkaJsonVMSMessage{
							PayLoad: models.PayLoad{
								DeviceScan: &[]models.ScanDevice{},
							},
						}
						StartURL := net.ParseIP(dto.StartURL)
						EndURL := net.ParseIP(dto.EndURL)
						if dto.StartURL != "" && dto.EndURL != "" && dto.StartHost != 0 && dto.EndHost != 0 {
							for _, message := range *msg.PayLoad.DeviceScan {
								messageHostInt, err := strconv.Atoi(message.Host)
								trial := net.ParseIP(message.URL)
								if err != nil {
									continue
								}
								if bytes.Compare(trial, StartURL) >= 0 && bytes.Compare(trial, EndURL) <= 0 &&
									messageHostInt >= dto.StartHost && messageHostInt <= dto.EndHost {
									*filteredMsgs.PayLoad.DeviceScan = append(*filteredMsgs.PayLoad.DeviceScan, message)
								}
							}
						}
						jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
						jsonRsp.Data = *filteredMsgs.PayLoad.DeviceScan
						c.JSON(http.StatusOK, &jsonRsp)
						return
					}
				}
			}
		}
	}
}

// ReadCamera		 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         device
// @Produce      json
// @Param        device  body      models.ScanDevice  true  "Device JSON"
// @Success      200   {object}  models.JsonDTORsp[models.ScanDevice]
// @Router       /device/onvif/devicescanhikivision [post]
// @Security	 BearerAuth
func UpdateOTABox(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.ScanDevice]()

	// Get new data from body
	var dto models.ScanDevice
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	cmd := models.DeviceCommand{ //unused
		CommandID: dto.IDBox,
		Cmd:       "",
		EventTime: time.Now().Format(time.RFC3339),
		UserName:  dto.UserName,
		Password:  dto.PassWord,
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
		case msg := <-ChannelDataReceiving:
			if msg.PayLoad.DeviceScan == nil {
				jsonRsp.Message = "Device does not exist "
				c.JSON(http.StatusOK, &jsonRsp)
				return
			} else {
				jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
				jsonRsp.Data = dto
				c.JSON(http.StatusOK, &jsonRsp)
				return
			}
		}
	}
}
