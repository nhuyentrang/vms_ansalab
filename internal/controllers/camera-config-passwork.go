package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"vms/internal/models"

	"vms/comongo/kafkaclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ReadCamera	 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         cameras
// @Produce      json
// @Param   	 idCamera			query	string	true	"id NVR keyword"		minlength(1)  	maxlength(100)
// @Param        Event  body   models.DTO_ChangePassword  true  "DTO_ChangePassword JSON"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_ChangePassword]
// @Router       /cameras/config/changepassword [put]
// @Security	 BearerAuth
func ChangePasswordCamera(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_ChangePassword]()
	idCamera := c.Query("idCamera")

	// Bind JSON request to DTO
	var dto models.DTO_ChangePassword
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Invalid request body: " + err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dto.ID = idCamera
	// Sanitize the password
	dto.PasswordNew = strings.TrimSpace(dto.PasswordNew)

	// Fetch Camera details
	dtoCamera, err := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](idCamera)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Camera not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	// Fetch Device details
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoCamera.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	dtoCameraConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraConfig, models.CameraConfig]("id = ?", dtoCamera.ConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)

	}

	dtoCameraNetworkConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NetworkConfig, models.NetworkConfig]("id = ?", dtoCameraConfig.NetworkConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)

	}

	// dtoCameraVideoConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_VideoConfig, models.VideoConfig]("id = ?", dtoCameraConfig.VideoConfigID)
	// if err != nil {
	// 	jsonRsp.Code = http.StatusInternalServerError
	// 	jsonRsp.Message = err.Error()
	// 	c.JSON(http.StatusInternalServerError, &jsonRsp)

	// }

	// Prepare the command
	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:      dtoDevice.ModelID,
		RequestUUID:    RequestUUID,
		Cmd:            cmd_ChangePassWord,
		EventTime:      time.Now().Format(time.RFC3339),
		ChangePassword: dto,
		ProtocolType:   dtoCamera.Protocol,
		UserName:       dtoCamera.Username,
		Password:       strings.TrimSpace(dto.PasswordOld),
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
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_ChangePassWord)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_ChangePassWord
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.Status != "FAILURE" {
					// Update camera password
					err := reposity.UpdateSingleColumn[models.Camera](dtoCamera.ID.String(), "password", dto.PasswordNew)
					if err != nil {
						jsonRsp.Code = http.StatusInternalServerError
						jsonRsp.Message = "Failed to update camera password: " + err.Error()
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						return
					}
					// Call UpdatePasswordCamera function to update RTSP URL
					dtoCamera.Password = dto.PasswordNew
					videoStream := UpdatePasswordCamera(dtoCamera, dtoCameraNetworkConfig, dtoCamera.Streams)
					// Update camera streams
					err = reposity.UpdateSingleColumn[models.Camera](dtoCamera.ID.String(), "streams", videoStream)
					if err != nil {
						jsonRsp.Code = http.StatusInternalServerError
						jsonRsp.Message = "Failed to update video stream: " + err.Error()
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						return
					}

					jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
					jsonRsp.Data = dto
					c.JSON(http.StatusOK, &jsonRsp)
					return
				} else {
					jsonRsp.Message = "Failed to change device's password"
					jsonRsp.Data = dto
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return
				}
			}
		}
	}
}

// ChangePasswordCameraSeries godoc
// @Summary      Change passwords for multiple cameras
// @Description  Changes the passwords for a list of cameras.
// @Tags         cameras
// @Produce      json
// @Param        Event  body   []models.DTO_ChangePassword  true  "DTO_ChangePassword JSON array"
// @Success      200   {object}  models.JsonDTORsp[[]models.DTO_ChangePassword]
// @Router       /cameras/config/changepasswordseries [put]
// @Security     BearerAuth
func ChangePasswordCameraSeries(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.DTO_ChangePassword]()

	// Get new data from body
	var dtos []models.DTO_ChangePassword
	if err := c.ShouldBindJSON(&dtos); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Invalid request body: " + err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Fetch the details of the first camera to get the device information
	dtoCamera, err := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](dtos[0].ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Camera not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	// Fetch Device details
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoCamera.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	//TODO: Each camera have a different protocol to edit password
	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:            dtoDevice.ModelID,
		RequestUUID:          RequestUUID,
		Cmd:                  cmd_ChangePassWordSeries,
		EventTime:            time.Now().Format(time.RFC3339),
		ChangePasswordSeries: dtos,
		ProtocolType:         dtoCamera.Protocol,
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
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_ChangePassWordSeries)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_ChangePassWordSeries
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)
				if msg.PayLoad.Status != "FAILURE" {

				}
				for _, dto := range msg.PayLoad.ChangePasswordSeries {
					if dto.Status == true {
						// Update camera password
						err := reposity.UpdateSingleColumn[models.Camera](dto.ID, "password", dto.PasswordNew)
						if err != nil {
							jsonRsp.Code = http.StatusInternalServerError
							jsonRsp.Message = "Failed to update camera password: " + err.Error()
							c.JSON(http.StatusInternalServerError, &jsonRsp)
							return
						}

						// Fetch updated camera details
						dtoCamera, err := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](dto.ID)
						if err != nil {
							jsonRsp.Code = http.StatusNotFound
							jsonRsp.Message = "Camera not found: " + err.Error()
							c.JSON(http.StatusNotFound, &jsonRsp)
							return
						}

						// Fetch related configurations
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

						// dtoCameraVideoConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_VideoConfig, models.VideoConfig]("id = ?", dtoCameraConfig.VideoConfigID)
						// if err != nil {
						// 	jsonRsp.Code = http.StatusInternalServerError
						// 	jsonRsp.Message = err.Error()
						// 	c.JSON(http.StatusInternalServerError, &jsonRsp)

						// }

						// Call UpdatePasswordCamera function to update RTSP URL
						dtoCamera.Password = dto.PasswordNew
						videoStream := UpdatePasswordCamera(dtoCamera, dtoCameraNetworkConfig, dtoCamera.Streams)
						// Update camera streams
						err = reposity.UpdateSingleColumn[models.Camera](dtoCamera.ID.String(), "streams", videoStream)
						if err != nil {
							jsonRsp.Code = http.StatusInternalServerError
							jsonRsp.Message = "Failed to update video stream: " + err.Error()
							c.JSON(http.StatusInternalServerError, &jsonRsp)
							return
						}
					}
				}
				jsonEditResp := models.NewJsonDTORsp[[]models.ChangePasswordStatuses]()
				jsonEditResp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
				jsonEditResp.Data = msg.PayLoad.ChangePasswordsStatuses
				jsonEditResp.Code = http.StatusInternalServerError
				c.JSON(http.StatusOK, &jsonEditResp)
				return
			}
		}
	}
}
