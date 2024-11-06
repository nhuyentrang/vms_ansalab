package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"vms/internal/models"
	"vms/statuscode"

	"vms/internal/models/hikivision"

	"vms/comongo/kafkaclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateNVR		godoc
// @Summary      	Create a new camera
// @Description  	Takes a camera JSON and store in DB. Return saved JSON.
// @Tags         	nvrs
// @Produce			json
// @Param        	camera  body   models.DTO_NVR  true  "NVR JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_NVR]
// @Router       	/nvrs [post]
// @Security		BearerAuth
func CreateNVR(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_NVR]()

	var dto models.DTO_NVR
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dto.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	if dto.IsOfflineSetting == nil || !*dto.IsOfflineSetting {
		// Scan data network config
		RequestUUID := uuid.New()
		cmd := models.DeviceCommand{
			CommandID:    dtoDevice.ModelID,
			RequestUUID:  RequestUUID,
			Cmd:          cmd_GetDataConfig,
			EventTime:    time.Now().Format(time.RFC3339),
			IPAddress:    dto.IPAddress,
			UserName:     dto.Username,
			Password:     dto.Password,
			HttpPort:     dto.HttpPort,
			OnvifPort:    "80",
			Channel:      "101",
			Track:        "103",
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
				fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetDataConfig)
				jsonRsp.Message = "No response from the device, timeout: " + cmd_GetDataConfig
				c.JSON(http.StatusRequestTimeout, &jsonRsp)
				return
			case <-ticker.C:
				if storedMsg, ok := messageMap.Load(RequestUUID); ok {
					msg := storedMsg.(*models.KafkaJsonVMSMessage)
					if msg.PayLoad.Status == "FAILURE" {
						jsonRsp.Message = "Error getting device's Configuration Data"
						jsonRsp.Data = dto
						c.JSON(http.StatusInternalServerError, &jsonRsp)
					} else {
						if msg.PayLoad.NetworkConfig != nil { //TODO FIX ME
							dtoNetWorkConfig, err := reposity.CreateItemFromDTO[models.DTO_NetworkConfig, models.NetworkConfig](*msg.PayLoad.NetworkConfig)
							if err != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = err.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}

							dtoVideoConfig := mapToDTOVideoConfig(msg.PayLoad.VideoConfig_DTO)
							dtoVideoConfig, err = reposity.CreateItemFromDTO[models.DTO_VideoConfig, models.VideoConfig](dtoVideoConfig)
							if err != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = "Failed to create network config: " + err.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}
							var dtoNVRConfig models.DTO_NVRConfig
							dtoNVRConfig.NetworkConfigID = dtoNetWorkConfig.ID
							dtoNVRConfig.VideoConfigID = dtoVideoConfig.ID
							dtoNVRConfig.ImageConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoNVRConfig.StorageConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoNVRConfig.StreamingConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoNVRConfig.AIConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoNVRConfig.AudioConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoNVRConfig.RecordingScheduleID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
							dtoNVRConfig.PTZConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")

							dtoNVRConfig, err = reposity.CreateItemFromDTO[models.DTO_NVRConfig, models.NVRConfig](dtoNVRConfig)
							if err != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = err.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}

							if dto.ConfigID == uuid.Nil {
								dto.ConfigID, _ = uuid.Parse(dtoNVRConfig.ID.String())
							}

							nvrEntry, errCameraCreate := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVR, models.NVR]("mac_address = ?", dto.MACAddress)
							if errCameraCreate == nil && nvrEntry.MACAddress != "" {
								fmt.Println("\t\t> Failed to insert NVR: Mac_address duplication ")
								jsonRsp.Message = "Failed to insert NVR: Mac_address duplication"
								c.JSON(http.StatusBadRequest, &jsonRsp)
								return
							}
							dto.OnvifPort = strconv.Itoa(dtoNetWorkConfig.ONVIF)
							dto.ManagementPort = strconv.Itoa(dtoNetWorkConfig.Server)
							dto.RtspPort = strconv.Itoa(dtoNetWorkConfig.RTSP)
							dto.Box.ID = dtoDevice.ID.String()
							dto, err = reposity.CreateItemFromDTO[models.DTO_NVR, models.NVR](dto)
							if err != nil {
								jsonRsp.Code = http.StatusInternalServerError
								jsonRsp.Message = err.Error()
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}

							dataFileConfig := map[string]models.ConfigNVR{}
							newNVR := models.ConfigNVR{
								NameCamera: dto.Name,
								IP:         dto.IPAddress,
								UserName:   dto.Username,
								PassWord:   dto.Password,
								HTTPPort:   strconv.Itoa(dtoNetWorkConfig.HTTP),
								RTSPPort:   strconv.Itoa(dtoNetWorkConfig.RTSP),
								OnvifPort:  dto.OnvifPort,
							}
							dataFileConfig[dto.ID.String()] = newNVR
							requestUUID := uuid.New()
							cmdAddConfig := models.DeviceCommand{
								CommandID:    dtoDevice.ModelID,
								Cmd:          cmd_AddDataConfig,
								EventTime:    time.Now().Format(time.RFC3339),
								EventType:    "nvr",
								ConfigNVR:    dataFileConfig,
								ProtocolType: dto.Protocol,
								RequestUUID:  requestUUID,
							}
							cmsStr, _ = json.Marshal(cmdAddConfig)
							kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

							jsonRsp.Data = dto
							jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
							c.JSON(http.StatusCreated, &jsonRsp)
							return
						} else {
							jsonRsp.Data = dto
							jsonRsp.Message = "Insert NVR Failed"
							c.JSON(http.StatusInternalServerError, &jsonRsp)
							return
						}
					}
				}
			}
		}
	} else {
		// Create Camera to get the correct camera ID
		dto, errNVR := reposity.CreateItemFromDTO[models.DTO_NVR, models.NVR](dto)
		if errNVR != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Failed to create NVR: " + errNVR.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

		jsonRsp.Message = "Created NVR Successfully"
		jsonRsp.Data = dto
		c.JSON(http.StatusCreated, &jsonRsp)
		return
	}
}

// ReadNVR		 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         nvrs
// @Produce      json
// @Param        id  path  string  true  "Search camera by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_NVR]
// @Router       /nvrs/{id} [get]
// @Security		BearerAuth
func ReadNVR(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_NVR]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_NVR, models.NVR](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateNVR		 	godoc
// @Summary      	Update single camera by id
// @Description  	Updates and returns a single camera whose ID value matches the id. New data must be passed in the body.
// @Tags         	nvrs
// @Produce      	json
// @Param        	id  path  string  true  "Update camera by id"
// @Param        	camera  body      models.DTO_NVR  true  "NVR JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_NVR]
// @Router       	/nvrs/{id} [put]
// @Security		BearerAuth
func UpdateNVR(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_NVR]()
	// Get new data from body
	var dto models.DTO_NVR
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Retrieve the original NVR
	dtoNVR, errNVR := reposity.ReadItemByIDIntoDTO[models.DTO_NVR, models.NVR](c.Param("id"))
	if errNVR != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "NVR not found: " + errNVR.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	if dto.Cameras != nil {
		pendingIDs := extractNVRCameraIDs(dto)

		// Detect duplication in pendingIDs
		if detectDuplication(pendingIDs) {
			jsonRsp.Code = statuscode.StatusCreateItemFailed
			jsonRsp.Message = "Camera already exists in the NVR"
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		// Identify removed cameras
		var removedCameras []models.KeyValue
		originalCamerasMap := make(map[string]models.KeyValue)
		for _, camera := range *dtoNVR.Cameras {
			originalCamerasMap[camera.ID] = camera
		}
		updatedCamerasMap := make(map[string]models.KeyValue)
		for _, camera := range *dto.Cameras {
			updatedCamerasMap[camera.ID] = camera
		}
		for id, camera := range originalCamerasMap {
			if _, exists := updatedCamerasMap[id]; !exists {
				removedCameras = append(removedCameras, camera)
			}
		}

		if removedCameras != nil {
			// Process removal of each identified camera
			for _, removedCamera := range removedCameras {
				camera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", removedCamera.ID)
				if err != nil {
					jsonRsp.Code = http.StatusInternalServerError
					jsonRsp.Message = err.Error()
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return
				}
				cameraconfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", camera.ConfigID)
				if err != nil {
					jsonRsp.Code = http.StatusInternalServerError
					jsonRsp.Message = err.Error()
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return
				}
				networkconfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NetworkConfig, models.NetworkConfig]("id = ?", cameraconfig.NetworkConfigID)
				if err != nil {
					jsonRsp.Code = http.StatusInternalServerError
					jsonRsp.Message = err.Error()
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return
				}
				dataInputProxy := hikivision.InputProxyChannel{
					ID: "0",
					SourceInputPort: hikivision.SourceInputPortDescriptor{
						AdminProtocol:        camera.Protocol,
						AddressingFormatType: "hostname",
						HostName:             camera.IPAddress,
						IPAddress:            camera.IPAddress,
						ManagePortNo:         networkconfig.Server,
						SrcInputPort:         "1",
						UserName:             camera.Username,
						Password:             camera.Password,
						StreamType:           "auto",
					},
				}
				dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoNVR.Box.ID)
				if err != nil {
					jsonRsp.Code = http.StatusNotFound
					jsonRsp.Message = err.Error()
					c.JSON(http.StatusNotFound, &jsonRsp)
					return
				}
				// Logic to handle the removal of each camera from the NVR
				RequestUUID := uuid.New()
				cmd := models.DeviceCommand{
					CommandID:     dtoDevice.ModelID,
					RequestUUID:   RequestUUID,
					Cmd:           cmd_RemoveCameraFromNVR,
					EventTime:     time.Now().Format(time.RFC3339),
					IPAddress:     dtoNVR.IPAddress,
					HttpPort:      dtoNVR.HttpPort,
					ProtocolType:  dtoNVR.Protocol,
					OnvifPort:     dtoNVR.OnvifPort,
					UserName:      dtoNVR.Username,
					Password:      dtoNVR.Password,
					CameraID:      camera.ID.String(),
					NVRID:         dtoNVR.ID.String(),
					SetInputProxy: dataInputProxy,
				}
				cmsStr, _ := json.Marshal(cmd)
				kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

				timeout := time.After(30 * time.Second)
				ticker := time.NewTicker(2 * time.Second)
				defer ticker.Stop()

				for {
					select {
					case <-timeout:
						fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_RemoveCameraFromNVR)
						jsonRsp.Message = "device timed out: " + cmd_RemoveCameraFromNVR
						c.JSON(http.StatusRequestTimeout, &jsonRsp)
						return

					case <-ticker.C:
						if storedMsg, ok := messageMap.Load(RequestUUID); ok {
							msg := storedMsg.(*models.KafkaJsonVMSMessage)
							if msg.PayLoad.ResponseStatus.StatusCode == 1 {
								// Log the state of removed camera
								log.Printf("Camera %s removed successfully from NVR", removedCamera.ID)

								var updatedCameras models.KeyValueArray
								for _, camera := range *dtoNVR.Cameras {
									if camera.ID != removedCamera.ID {
										updatedCameras = append(updatedCameras, camera)
									}
								}
								dtoNVR.Cameras = &updatedCameras

								// Update the NVR in the database
								if _, err := reposity.UpdateItemByIDFromDTO[models.DTO_NVR, models.NVR](dtoNVR.ID.String(), dtoNVR); err != nil {
									jsonRsp.Code = http.StatusInternalServerError
									jsonRsp.Message = err.Error()
									c.JSON(http.StatusInternalServerError, &jsonRsp)
									return
								}
								// Return
								jsonRsp.Data = dto
								c.JSON(http.StatusOK, &jsonRsp)
								return
							} else {
								jsonRsp.Code = int64(msg.PayLoad.ResponseStatus.StatusCode)
								jsonRsp.Message = msg.PayLoad.ResponseStatus.StatusString + ", " + msg.PayLoad.ResponseStatus.SubStatusCode
								c.JSON(http.StatusInternalServerError, &jsonRsp)
								return
							}
						}
					}
				}
			}
		} else {
			// Update entity from DTO
			dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_NVR, models.NVR](c.Param("id"), dto)
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
			jsonRsp.Message = "NVR Successfully Updated"
			c.JSON(http.StatusOK, &jsonRsp)
			return
		}
	} else {
		// Update entity from DTO
		dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_NVR, models.NVR](c.Param("id"), dto)
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
		jsonRsp.Message = "NVR Successfully Updated"

		c.JSON(http.StatusOK, &jsonRsp)
	}
}

// Utility function to extract camera IDs from a DTO_CameraGroup_Create instance
func extractNVRCameraIDs(dto models.DTO_NVR) []string {
	ids := make([]string, len(*dto.Cameras))

	// Iterate over the Cameras slice within the DTO
	for i, camera := range *dto.Cameras {
		ids[i] = camera.ID
	}

	return ids
}

func detectDuplication(ids []string) bool {
	seen := make(map[string]bool)
	for _, id := range ids {
		if _, found := seen[id]; found {
			// Duplicate found
			return true
		}
		seen[id] = true
	}
	return false
}

// Utility function to find unique IDs in the first slice that aren't in the second
func findNVRUniqueIDs(slice1, slice2 []string) []string {
	unique := make([]string, 0)
	set := make(map[string]struct{})
	for _, id := range slice2 {
		set[id] = struct{}{}
	}
	for _, id := range slice1 {
		if _, found := set[id]; !found {
			unique = append(unique, id)
		}
	}
	return unique
}

// DeleteNVR godoc
// @Summary      Remove NVR by id
// @Description  Delete a single NVR entry from the repository based on id.
// @Tags         nvrs
// @Produce      json
// @Param        id  path  string  true  "Delete NVR by id"
// @Success      200     {object} models.JsonDTORsp[models.DTO_NVR] "NVR deleted successfully"
// @Failure      404     {object} models.JsonDTORsp[models.DTO_NVR] "NVR not found"
// @Failure      408     {object} models.JsonDTORsp[models.DTO_NVR] "Request timeout"
// @Failure      500     {object} models.JsonDTORsp[models.DTO_NVR] "Internal server error"
// @Router       /nvrs/{id} [delete]
// @Security     BearerAuth
func DeleteNVR(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_NVR]()

	dtoNVR, err := reposity.ReadItemByIDIntoDTO[models.DTO_NVR, models.NVR](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "NVR not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoNVR.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	errDel := reposity.DeleteItemByID[models.NVR](c.Param("id"))
	if errDel != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to delete NVR: " + errDel.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Check if there are any cameras to remove
	if dtoNVR.Cameras == nil || len(*dtoNVR.Cameras) == 0 {
		// No cameras associated with this NVR, proceed to NVR deletion steps
		performNVRDeletion(c, jsonRsp, &dtoNVR, &dtoDevice)
		return
	}

	// Collect all cameras to be removed
	var inputProxyChannels []hikivision.InputProxyChannel
	for _, removedCamera := range *dtoNVR.Cameras {
		camera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", removedCamera.ID)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		cameraconfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", camera.ConfigID)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		networkconfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NetworkConfig, models.NetworkConfig]("id = ?", cameraconfig.NetworkConfigID)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		dataInputProxy := hikivision.InputProxyChannel{
			ID: "0",
			SourceInputPort: hikivision.SourceInputPortDescriptor{
				AdminProtocol:        camera.Protocol,
				AddressingFormatType: "hostname",
				HostName:             camera.IPAddress,
				IPAddress:            camera.IPAddress,
				ManagePortNo:         networkconfig.Server,
				SrcInputPort:         "1",
				UserName:             camera.Username,
				Password:             camera.Password,
				StreamType:           "auto",
			},
		}
		inputProxyChannels = append(inputProxyChannels, dataInputProxy)
	}

	if len(inputProxyChannels) > 0 {
		RequestUUID := uuid.New()
		cmd := models.DeviceCommand{
			CommandID:       dtoDevice.ModelID,
			RequestUUID:     RequestUUID,
			Cmd:             cmd_RemoveCamerasFromNVR,
			EventTime:       time.Now().Format(time.RFC3339),
			IPAddress:       dtoNVR.IPAddress,
			HttpPort:        dtoNVR.HttpPort,
			ProtocolType:    dtoNVR.Protocol,
			OnvifPort:       dtoNVR.OnvifPort,
			UserName:        dtoNVR.Username,
			Password:        dtoNVR.Password,
			NVRID:           dtoNVR.ID.String(),
			SetInputProxies: inputProxyChannels,
		}
		cmsStr, _ := json.Marshal(cmd)
		kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

		timeout := time.After(30 * time.Second)
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-timeout:
				fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_RemoveCameraFromNVR)
				jsonRsp.Message = "Unable to remove NVR camera, device timed out: " + cmd_RemoveCameraFromNVR
				c.JSON(http.StatusRequestTimeout, &jsonRsp)
				return

			case <-ticker.C:
				if storedMsg, ok := messageMap.Load(RequestUUID); ok {
					msg := storedMsg.(*models.KafkaJsonVMSMessage)
					if msg.PayLoad.ResponseStatus.StatusCode == 1 {
						for _, removedCamera := range *dtoNVR.Cameras {
							log.Printf("Camera %s removed successfully from NVR", removedCamera.ID)
						}
						break
					} else {
						jsonRsp.Code = int64(msg.PayLoad.ResponseStatus.StatusCode)
						jsonRsp.Message = msg.PayLoad.ResponseStatus.StatusString + ", " + msg.PayLoad.ResponseStatus.SubStatusCode
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						break
					}
				}
			}
		}
	}

	performNVRDeletion(c, jsonRsp, &dtoNVR, &dtoDevice)
}

func performNVRDeletion(c *gin.Context, jsonRsp *models.JsonDTORsp[models.DTO_NVR], dtoNVR *models.DTO_NVR, dtoDevice *models.DTO_Device) {
	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		RequestUUID:  RequestUUID,
		Cmd:          cmd_DeleteFileStream,
		EventTime:    time.Now().Format(time.RFC3339),
		EventType:    "nvr",
		NVRID:        c.Param("id"),
		ProtocolType: dtoNVR.Protocol,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_DeleteFileStream)
			jsonRsp.Message = "NVR deleted, SmartNVR asynchronous, device timed out: " + cmd_DeleteFileStream
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				jsonRsp.Data = *dtoNVR // Dereference the pointer to assign the value
				jsonRsp.Message = "NVR Deleted and SmartNVR synched successfully " + msg.PayLoad.Cmd
				c.JSON(http.StatusOK, &jsonRsp)
				return
			}
		}
	}
}

// GetNVRs		godoc
// @Summary      	Get all camera groups with query filter
// @Description  	Responds with the list of all camera as JSON.
// @Tags         	nvrs
// @Param   		keyword			query	string	false	"camera name keyword"		minlength(1)  	maxlength(100)
// @Param   		sort  			query	string	false	"sort"						default(+created_at)
// @Param   		limit			query	int     false  	"limit"          			minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          			minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_NVR_Read_BasicInfo]
// @Router       	/nvrs [get]
// @Security		BearerAuth
func GetNVRs(c *gin.Context) {
	jsonRspDTONVRsBasicInfos := models.NewJsonDTOListRsp[models.DTO_NVR_Read_BasicInfo]()

	// Get param
	keyword := c.Query("keyword")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	fmt.Println(
		"keyword: ", keyword,
		" - sort: ", sort,
		" - limit: ", limit,
		" - page: ", page)

	// Build query
	query := reposity.NewQuery[models.DTO_NVR_Read_BasicInfo, models.NVR]()

	// Search for keyword in name
	if keyword != "" {
		query.AddConditionOfTextField("AND", "name", "LIKE", keyword)
	}

	// Exec query
	dtoNVRBasics, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTONVRsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTONVRsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTONVRsBasicInfos)
		return
	}

	jsonRspDTONVRsBasicInfos.Count = count
	jsonRspDTONVRsBasicInfos.Data = dtoNVRBasics
	jsonRspDTONVRsBasicInfos.Page = int64(page)
	jsonRspDTONVRsBasicInfos.Size = int64(len(dtoNVRBasics))
	c.JSON(http.StatusOK, &jsonRspDTONVRsBasicInfos)
}

// GetNVRConfigOptions		godoc
// @Summary      	Get items of camera protocols for select box
// @Description  	Responds with the list of item for camera protocol.
// @Tags         	nvrs
// @Produce      	json
// @Success      	200  		{object}  		models.JsonDTORsp[[]models.KeyValue]
// @Router       	/nvrs/options/protocol-types [get]
// @Security		BearerAuth
func GetNVRProtocolTypes(c *gin.Context) {

	jsonConfigTypeRsp := models.NewJsonDTORsp[[]models.KeyValue]()
	jsonConfigTypeRsp.Data = make([]models.KeyValue, 0)

	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-onvif",
		Name: "ONVIF",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-ivi",
		Name: "IVI",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-hikvision",
		Name: "HIKVISION",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-dahua",
		Name: "DAHUA",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-axis",
		Name: "AXIS",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-bosch",
		Name: "BOSCH",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-hanwha",
		Name: "HANWHA",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-panasonic",
		Name: "PANASONIC",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-unv",
		Name: "UNIVIEW",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-kbvision",
		Name: "KBVISION",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-protocol-agivilon",
		Name: "AGIVILON",
	})

	jsonConfigTypeRsp.Code = 0
	c.JSON(http.StatusOK, &jsonConfigTypeRsp)
}

// GetNVRConfigOptions		godoc
// @Summary      	Get items of camera stream for select box
// @Description  	Responds with the list of item for camera stream type.
// @Tags         	nvrs
// @Produce      	json
// @Success      	200  		{object}  		models.JsonDTORsp[[]models.KeyValue]
// @Router       	/nvrs/options/stream-types [get]
// @Security		BearerAuth
func GetNVRStreamTypes(c *gin.Context) {

	jsonConfigTypeRsp := models.NewJsonDTORsp[[]models.KeyValue]()
	jsonConfigTypeRsp.Data = make([]models.KeyValue, 0)

	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-stream-ondemand",
		Name: "ON Demand",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-stream-persistant",
		Name: "Persistant",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-stream-p2p",
		Name: "P2P",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "camera-stream-proxied",
		Name: "Proxied",
	})

	jsonConfigTypeRsp.Code = 0
	c.JSON(http.StatusOK, &jsonConfigTypeRsp)
}

// GetNVRConfigOptions		godoc
// @Summary      	Get items of camera types for select box
// @Description  	Responds with the list of item for camera type.
// @Tags         	nvrs
// @Produce      	json
// @Success      	200  		{object}  		models.JsonDTORsp[[]models.KeyValue]
// @Router       	/nvrs/options/types [get]
// @Security		BearerAuth
func GetNVRTypes(c *gin.Context) {

	jsonConfigTypeRsp := models.NewJsonDTORsp[[]models.KeyValue]()
	jsonConfigTypeRsp.Data = make([]models.KeyValue, 0)

	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "ip-camera",
		Name: "IP NVR",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "ptz-camera",
		Name: "PTZ",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "ai-camera",
		Name: "AI",
	})

	jsonConfigTypeRsp.Code = 0
	c.JSON(http.StatusOK, &jsonConfigTypeRsp)
}

// UpdateNVR		 	godoc
// @Summary      	Update single camera by id
// @Description  	Updates and returns a single camera whose ID value matches the id. New data must be passed in the body.
// @Tags         	nvrs
// @Produce      	json
// @Param        	id  path  string  true  "IDNVR, add camera to NVR"
// @Param        	camera  body      models.DTO_NVR  true  "NVR JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_NVR]
// @Router       	/nvrs/camera/{id} [put]
// @Security		BearerAuth
func AddCameraToNVR(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_NVR]()

	var dto models.DTO_NVR
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Ensure there are cameras in the DTO
	if len(*dto.Cameras) == 0 {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "No cameras provided"
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	lastElement := (*dto.Cameras)[len(*dto.Cameras)-1]

	// Check if the camera already exists in dto.Cameras except for the last element
	for _, camera := range (*dto.Cameras)[:len(*dto.Cameras)-1] {
		if camera.ID == lastElement.ID {
			jsonRsp.Code = http.StatusConflict
			jsonRsp.Message = "Camera already exists in the NVR"
			c.JSON(http.StatusConflict, &jsonRsp)
			return
		}
		// else if camera. != dto.Box{
		// 	jsonRsp.Code = http.StatusConflict
		// 	jsonRsp.Message = "Camera and NVR are not in "
		// 	c.JSON(http.StatusConflict, &jsonRsp)
		// 	return
		// }

	}

	dtoCamera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", lastElement.ID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	dtoNVR, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVR, models.NVR]("id = ?", c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	camera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", lastElement.ID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	if camera.Box.Name != dtoNVR.Box.Name {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Camera must be attached to the same SmartNVR as NVR before being added to the NVR"
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	cameraconfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", camera.ConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	networkconfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NetworkConfig, models.NetworkConfig]("id = ?", cameraconfig.NetworkConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	dataInputProxy := hikivision.InputProxyChannel{
		ID: "0",
		SourceInputPort: hikivision.SourceInputPortDescriptor{
			AdminProtocol:        camera.Protocol,
			AddressingFormatType: "hostname",
			HostName:             dtoCamera.IPAddress,
			IPAddress:            dtoCamera.IPAddress,
			ManagePortNo:         networkconfig.Server,
			SrcInputPort:         "1",
			UserName:             dtoCamera.Username,
			Password:             dtoCamera.Password,
			StreamType:           "auto",
		},
	}
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoNVR.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:     dtoDevice.ModelID,
		RequestUUID:   RequestUUID,
		Cmd:           cmd_AddCameratoNVR,
		EventTime:     time.Now().Format(time.RFC3339),
		IPAddress:     dtoNVR.IPAddress,
		HttpPort:      dtoNVR.HttpPort,
		ProtocolType:  dtoNVR.Protocol,
		OnvifPort:     dtoNVR.OnvifPort,
		UserName:      dtoNVR.Username,
		Password:      dtoNVR.Password,
		CameraID:      dtoCamera.ID.String(),
		NVRID:         dtoNVR.ID.String(),
		SetInputProxy: dataInputProxy,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_AddCameratoNVR)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_AddCameratoNVR
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.ResponseStatus.StatusCode == 1 {
					nvrInfo := models.KeyValue{
						ID:      dtoNVR.ID.String(),
						Name:    dtoNVR.Name,
						Channel: strconv.Itoa(msg.PayLoad.ResponseStatus.ID),
					}
					(*dto.Cameras)[len(*dto.Cameras)-1].Channel = strconv.Itoa(msg.PayLoad.ResponseStatus.ID)

					// Log the state of dtoNVR.Cameras before the update
					log.Printf("Updating NVR Cameras: %+v", dtoNVR.Cameras)

					// Update the NVR with the new camera list
					dtoNVR.Cameras = dto.Cameras

					log.Printf("NVR Cameras before DB update: %+v", dtoNVR.Cameras)

					_, err := reposity.UpdateItemByIDFromDTO[models.DTO_NVR, models.NVR](c.Param("id"), dtoNVR)
					if err != nil {
						jsonRsp.Code = http.StatusInternalServerError
						jsonRsp.Message = err.Error()
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						return
					}

					// Log the state after the update
					log.Println("NVR Cameras updated successfully in the DB")

					// Update the camera's NVR information
					err = reposity.UpdateSingleColumn[models.Camera]((*dto.Cameras)[len(*dto.Cameras)-1].ID, "nvr", nvrInfo)
					if err != nil {
						fmt.Println("======> Updated NVR configuration failed !!!")
						return
					}

					jsonRsp.Message = msg.PayLoad.ResponseStatus.StatusString + ", " + msg.PayLoad.ResponseStatus.SubStatusCode
					c.JSON(http.StatusOK, &jsonRsp)
					return
				} else {
					jsonRsp.Code = int64(msg.PayLoad.ResponseStatus.StatusCode)
					jsonRsp.Message = msg.PayLoad.ResponseStatus.StatusString + ", " + msg.PayLoad.ResponseStatus.SubStatusCode
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return
				}
			}
		}
	}
}

// ImportNVRsFromCSV godoc
// @Summary       Import NVRs from a CSV file
// @Description   Reads a CSV file containing NVR data and inserts them into the database.
// @Tags          nvrs
// @Produce       json
// @Param         file  formData  file  true  "CSV file with NVRs"
// @Param         idbox query     string true  "Id of the box"  minlength(1)   maxlength(100)
// @Success       200   {object}  models.JsonDTORsp[[]models.DTONVRImport]
// @Router        /nvrs/import [post]
// @Security      BearerAuth
func ImportNVRsFromCSV(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.DTONVRImport]()
	idbox := c.Query("idbox")
	file, err := c.FormFile("file")
	if err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Failed to get CSV file: " + err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	f, err := file.Open()
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to open CSV file: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	var nvrs []models.DTONVRImport
	var errors []string

	headers, err := r.Read()
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to read CSV headers: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](idbox)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "failed to retrieve device information"
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, "Failed to read CSV record: "+err.Error())
			continue
		}

		dto, err := parseNVRRecord(headers, record)
		if err != nil {
			errors = append(errors, "Failed to parse CSV record: "+err.Error())
			continue
		}

		dto, err = insertNVRIntoDB(dto, dtoDevice)
		if err != nil {
			log.Printf("Failed to insert NVR: %v\n", err)
			errors = append(errors, fmt.Sprintf("Failed to insert NVR %s: %v", dto.IPAddress, err))
		}

		nvrs = append(nvrs, dto)
	}

	jsonRsp.Data = nvrs
	if len(errors) > 0 {
		jsonRsp.Message = "Some errors occurred: " + strings.Join(errors, "; ")
		jsonRsp.Code = http.StatusPartialContent // Indicates that some records were not processed successfully
	} else {
		jsonRsp.Message = "All records processed successfully"
		jsonRsp.Code = http.StatusOK
	}
	c.JSON(http.StatusOK, &jsonRsp)
}

func createNVRLogic(dto models.DTONVRImport, device models.DTO_Device) (models.DTONVRImport, error) {

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    device.ModelID,
		Cmd:          cmd_GetDataConfig,
		RequestUUID:  RequestUUID,
		EventTime:    time.Now().Format(time.RFC3339),
		IPAddress:    dto.IPAddress,
		UserName:     dto.Username,
		Password:     dto.Password,
		HttpPort:     dto.HttpPort,
		ProtocolType: dto.Protocol,
		Channel:      "101",
		Track:        "103",
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return dto, fmt.Errorf("timeout waiting for command result: %s", cmd_GetDataConfig)
		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.NetworkConfig != nil && msg.PayLoad.VideoConfig != nil {
					dtoNetWorkConfig, err := reposity.CreateItemFromDTO[models.DTO_NetworkConfig, models.NetworkConfig](*msg.PayLoad.NetworkConfig)
					if err != nil {
						return dto, fmt.Errorf("failed to create network config: %w", err)
					}

					var dtoNVRConfig models.DTO_NVRConfig
					dtoNVRConfig.NetworkConfigID = dtoNetWorkConfig.ID
					dtoNVRConfig.VideoConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoNVRConfig.ImageConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoNVRConfig.StorageConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoNVRConfig.StreamingConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoNVRConfig.AIConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoNVRConfig.AudioConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoNVRConfig.RecordingScheduleID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")
					dtoNVRConfig.PTZConfigID, _ = uuid.Parse("11111111-1111-1111-1111-111111111111")

					// Check if the entry already exists in the main NVR table based on MACAddress
					_, errNVRCreate := reposity.ReadItemWithFilterIntoDTO[models.DTONVRImport, models.NVR]("mac_address = ?", dto.MACAddress)
					if errNVRCreate == nil {
						return dto, fmt.Errorf("entry already exists in main NVR table")
					}

					// Create NVR config
					dtoNVRConfig, errNVR := reposity.CreateItemFromDTO[models.DTO_NVRConfig, models.NVRConfig](dtoNVRConfig)
					if errNVR != nil {
						return dto, fmt.Errorf("failed to create NVR config: %w", errNVR)
					}

					// Assign additional fields
					dto.OnvifPort = strconv.Itoa(dtoNetWorkConfig.ONVIF)
					dto.ManagementPort = strconv.Itoa(dtoNetWorkConfig.Server)
					dto.RtspPort = strconv.Itoa(dtoNetWorkConfig.RTSP)
					dto.Box.ID = device.ID.String()

					// Create the NVR entry
					dto, errNVR = reposity.CreateItemFromDTO[models.DTONVRImport, models.NVR](dto)
					if errNVR != nil {
						return dto, fmt.Errorf("failed to create NVR: %w", errNVR)
					}

					dataFileConfig := map[string]models.ConfigCamera{}
					newCamera := models.ConfigCamera{
						NameCamera: dto.Name,
						IP:         dto.IPAddress,
						UserName:   dto.Username,
						PassWord:   dto.Password,
						HTTPPort:   strconv.Itoa(dtoNetWorkConfig.HTTP),
						RTSPPort:   strconv.Itoa(dtoNetWorkConfig.RTSP),
						OnvifPort:  dto.OnvifPort,
						Channels: map[string]models.ChannelCamera{
							"Main": {
								OnDemand: true,
								Url:      fmt.Sprintf("rtsp://%s:%s@%s:%s/Streaming/Channels/101", dto.Username, dto.Password, dto.IPAddress, strconv.Itoa(dtoNetWorkConfig.RTSP)),
								Codec:    "h264",
								Name:     "Main",
							},
							"Sub": {
								OnDemand: true,
								Url:      fmt.Sprintf("rtsp://%s:%s@%s:%s/Streaming/Channels/102", dto.Username, dto.Password, dto.IPAddress, strconv.Itoa(dtoNetWorkConfig.RTSP)),
								Codec:    "h264",
								Name:     "Sub",
							},
						},
					}
					dataFileConfig[dto.ID.String()] = newCamera

					// Additional logic to move from NVRImport to NVR
					requestUUID := uuid.New()
					cmdAddConfig := models.DeviceCommand{
						CommandID:    device.ModelID,
						Cmd:          cmd_AddDataConfig,
						EventTime:    time.Now().Format(time.RFC3339),
						EventType:    "nvr",
						ConfigCamera: dataFileConfig,
						ProtocolType: dto.Protocol,
						RequestUUID:  requestUUID,
					}
					cmsStr, _ := json.Marshal(cmdAddConfig)
					kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

					return dto, nil
				}
				return dto, fmt.Errorf("network config is nil")
			}
		}
	}
}

func insertNVRIntoDB(dto models.DTONVRImport, device models.DTO_Device) (models.DTONVRImport, error) {
	// Send command to scan device
	RequestUUID := uuid.New()
	cmdScan := models.DeviceCommand{
		CommandID:    device.ModelID,
		Cmd:          cmd_ScanDeviceIP,
		EventTime:    time.Now().Format(time.RFC3339),
		IPAddress:    dto.IPAddress,
		HttpPort:     dto.HttpPort,
		UserName:     dto.Username,
		Password:     dto.Password,
		ProtocolType: dto.Protocol,
		RequestUUID:  RequestUUID,
	}
	cmsStr, _ := json.Marshal(cmdScan)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	// Set timeout and ticker
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	scanResultNotDone := 1
	// Wait for scan result
	for scanResultNotDone > 0 {
		select {
		case <-timeout:
			return dto, fmt.Errorf("timeout waiting for scan command result: %s", cmd_ScanDeviceIP)
		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.DeviceScan != nil && len(*msg.PayLoad.DeviceScan) > 0 {
					scanResult := (*msg.PayLoad.DeviceScan)[0]
					dto.ID = uuid.New()
					dto.Name = scanResult.Name
					dto.Type.ID = scanResult.Type
					dto.Location = scanResult.Location
					dto.MACAddress = scanResult.MacAddress
					dto.FirmwareVersion = scanResult.FirmwareVersion
					dto.ExportStatus = false
					scanResultNotDone = 0
				} else {
					return dto, fmt.Errorf("failed to retrieve NVR information")
				}
			}
		}
	}

	dto.InsertStatus = false // Default status to false
	importedDto, err := reposity.CreateItemFromDTO[models.DTONVRImport, models.NVRImport](dto)
	if err != nil {
		return dto, fmt.Errorf("failed to insert into import table: %w", err)
	}

	// Attempt to process and insert the NVR into the main NVR table
	processedDto, err := createNVRLogic(dto, device)
	if err != nil {
		// Update the import record to reflect the failure
		importedDto.InsertStatus = false
		_, updateErr := reposity.UpdateItemByIDFromDTO[models.DTONVRImport, models.NVRImport](importedDto.ID.String(), importedDto)
		if updateErr != nil {
			return dto, fmt.Errorf("failed to update import record status after NVR insert failure: %w", updateErr)
		}
		return dto, fmt.Errorf("failed to create NVR: %w", err)
	}

	// Update the import record to reflect the success
	dto.InsertStatus = true
	_, updateErr := reposity.UpdateItemByIDFromDTO[models.DTONVRImport, models.NVRImport](dto.ID.String(), dto)
	if updateErr != nil {
		return dto, fmt.Errorf("failed to update import record status after successful NVR insert: %w", updateErr)
	}

	return processedDto, nil
}

func parseNVRRecord(headers, record []string) (models.DTONVRImport, error) {
	var dto models.DTONVRImport

	for i, header := range headers {
		switch header {
		case "name":
			dto.Name = record[i]
		case "ipAddress":
			dto.IPAddress = record[i]
		case "macAddress":
			dto.MACAddress = record[i]
		case "httpPort":
			dto.HttpPort = record[i]
		case "onvifPort":
			dto.OnvifPort = record[i]
		case "managementPort":
			dto.ManagementPort = record[i]
		case "username":
			dto.Username = record[i]
		case "password":
			dto.Password = record[i]
		case "protocol":
			dto.Protocol = record[i]
		default:
			return dto, fmt.Errorf("unexpected header: %s", header)
		}
	}

	return dto, nil
}

// GenerateSampleNVRData godoc
// @Summary       Generate and download sample NVR data
// @Description   Provides a sample CSV file for users to know what data to input.
// @Tags          nvrs
// @Produce       text/csv
// @Success       200   "Sample data file"
// @Router        /nvrs/sample-data [get]
// @Security      BearerAuth
func GenerateSampleNVRData(c *gin.Context) {
	sampleData := [][]string{
		{"ipAddress", "httpPort", "username", "password"},
		{"192.168.1.2", "80", "admin", "password123"},
		{"192.168.1.3", "8080", "user", "userpass"},
		{"192.168.1.4", "8081", "admin", "adminpass"},
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=sample_nvr_data.csv")
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	for _, record := range sampleData {
		if err := writer.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV data"})
			return
		}
	}
}

// DownloadImportedNVRs godoc
// @Summary       Download the imported NVR list as a CSV file
// @Description   Provides the imported NVR list in CSV format.
// @Tags          nvrs
// @Produce       text/csv
// @Success       200   "CSV file"
// @Router        /nvrs/imported/download [get]
// @Security      BearerAuth
func DownloadImportedNVRs(c *gin.Context) {
	// Fetch imported NVRs from the database
	importedNVRs, err := fetchImportedNVRs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch imported NVRs"})
		return
	}

	// Create a new CSV file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")

	// Get filename from query or use default
	date := time.Now().Format("2006-01-02")
	defaultFilename := fmt.Sprintf("NVRs_%s.csv", date)
	filename := c.DefaultQuery("filename", defaultFilename)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv")

	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write headers
	headers := []string{"id", "name", "type", "protocol", "model", "serial", "firmwareVersion", "ipAddress", "macAddress", "httpPort", "onvifPort", "managementPort", "username", "password", "useTLS", "isOfflineSetting", "iframeURL", "lat", "long", "insertstatus", "location", "coordinate", "position", "faceRecognition", "licensePlateRecognition", "configID"}
	if err := writer.Write(headers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV headers"})
		return
	}

	// Write data
	for _, nvr := range importedNVRs {
		record := []string{
			nvr.ID.String(),
			nvr.Name,
			nvr.Type.ID,
			nvr.Protocol,
			nvr.Model,
			nvr.SerialNumber,
			nvr.FirmwareVersion,
			nvr.IPAddress,
			nvr.MACAddress,
			nvr.HttpPort,
			nvr.OnvifPort,
			nvr.ManagementPort,
			nvr.Username,
			nvr.Password,
			strconv.FormatBool(nvr.UseTLS),
			strconv.FormatBool(nvr.IsOfflineSetting),
			nvr.IFrameURL,
			nvr.Lat,
			nvr.Long,
			strconv.FormatBool(nvr.InsertStatus),
			nvr.Location,
			nvr.Coordinate,
			nvr.Position,
			strconv.FormatBool(nvr.FaceRecognition),
			strconv.FormatBool(nvr.LicensePlateRecognition),
			nvr.ConfigID.String(),
		}
		if err := writer.Write(record); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write CSV data"})
			return
		}
	}

	// Update export status to true for all exported NVRs
	for _, nvr := range importedNVRs {
		nvr.InsertStatus = true
		_, updateErr := reposity.UpdateItemByIDFromDTO[models.DTONVRImport, models.NVRImport](nvr.ID.String(), nvr)
		if updateErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update export status for NVR %s: %v", nvr.ID.String(), updateErr)})
			return
		}
	}
}

func fetchImportedNVRs() ([]models.DTONVRImport, error) {
	query := reposity.NewQuery[models.DTONVRImport, models.NVRImport]()
	query.AddConditionOfTextField("AND", "export_status", "=", false)
	nvrs, _, err := query.ExecWithPaging("+created_at", 9999, 1)
	if err != nil {
		return nil, err
	}

	// Update the exportstatus to true for all fetched NVRs
	for _, nvr := range nvrs {
		nvr.ExportStatus = true
		_, updateErr := reposity.UpdateItemByIDFromDTO[models.DTONVRImport, models.NVRImport](nvr.ID.String(), nvr)
		if updateErr != nil {
			return nil, fmt.Errorf("failed to update export status for NVR %s: %w", nvr.ID, updateErr)
		}
	}

	return nvrs, nil
}

// GetAttachedCamera godoc
// @Summary       Get all attached cameras for an NVR
// @Description   Fetches all cameras attached to a specific NVR by sending a command to the device.
// @Tags          nvrs
// @Produce       json
// @Param         id  path  string  true  "NVR ID"
// @Success       200
// @Failure       400  "Bad Request - Invalid input"
// @Failure       404  "Not Found - NVR not found"
// @Failure       500  "Internal Server Error - Unable to retrieve attached cameras"
// @Router        /nvrs/{id}/cameras [get]
// @Security      BearerAuth
func GetAttachedCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTOListRsp[models.InputProxyChannel]()

	// Get NVR ID from the request parameters
	nvrID := c.Param("id")

	// Fetch the NVR details
	dtoNVR, err := reposity.ReadItemByIDIntoDTO[models.DTO_NVR, models.NVR](nvrID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "NVR not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	// Fetch the device details associated with the NVR
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoNVR.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	query := reposity.NewQuery[models.DTOCamera, models.Camera]()
	sort := "-created_at"

	// Execute Query
	dbCameras, _, err := query.ExecNoPaging(sort)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Failed to retrieve cameras: " + err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Prepare the command to get all attached cameras
	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		RequestUUID:  RequestUUID,
		Cmd:          cmd_GetAttachedCameras,
		EventTime:    time.Now().Format(time.RFC3339),
		IPAddress:    dtoNVR.IPAddress,
		HttpPort:     dtoNVR.HttpPort,
		ProtocolType: dtoNVR.Protocol,
		OnvifPort:    dtoNVR.OnvifPort,
		UserName:     dtoNVR.Username,
		Password:     dtoNVR.Password,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	// Set timeout and ticker for waiting for a response
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetAttachedCameras)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_GetAttachedCameras
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if len(msg.PayLoad.InputProxyChannels) > 0 {
					// Compare the two lists and update with database data
					for i, returnedCamera := range msg.PayLoad.InputProxyChannels {
						for _, dbCamera := range dbCameras {
							httpPort, _ := strconv.Atoi(dbCamera.HttpPort)
							if returnedCamera.SourceInputPort.IPAddress == dbCamera.IPAddress && returnedCamera.SourceInputPort.ManagePortNo == httpPort {
								// Update the returned camera with data from the database
								msg.PayLoad.InputProxyChannels[i].SourceInputPort.NVR = dbCamera.NVR
								msg.PayLoad.InputProxyChannels[i].SourceInputPort.Box = dbCamera.Box
								msg.PayLoad.InputProxyChannels[i].SourceInputPort.MacAddress = dbCamera.MACAddress
							}
						}
					}

					jsonRsp.Data = msg.PayLoad.InputProxyChannels
					jsonRsp.Message = "Successfully retrieved attached cameras"
					c.JSON(http.StatusOK, &jsonRsp)
					return
				} else {
					jsonRsp.Code = http.StatusNotFound
					jsonRsp.Message = "No cameras returned"
					c.JSON(http.StatusNotFound, &jsonRsp)
					return
				}
			}
		}
	}
}
