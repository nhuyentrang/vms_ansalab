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

// ReadCamera	 godoc
// @Summary      Get single camera by id
// @Description  Returns the camera whose ID value matches the id.
// @Tags         nvrs
// @Produce      json
// @Param   	 idNVR			query	string	true	"id NVR keyword"		minlength(1)  	maxlength(100)
// @Param        Event  body   models.DTO_ChangePassword  true  "DTO_ChangePassword JSON"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_ChangePassword]
// @Router       /nvrs/config/changepassword [put]
// @Security	 BearerAuth
func ChangePassword(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_ChangePassword]()
	idNVR := c.Query("idNVR")

	// Get new data from body
	var dto models.DTO_ChangePassword
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dto.ID = idNVR

	dtoNVR, err := reposity.ReadItemByIDIntoDTO[models.DTO_NVR, models.NVR](idNVR)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
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
		CommandID:      dtoDevice.ModelID,
		RequestUUID:    RequestUUID,
		Cmd:            cmd_ChangePassWord,
		EventTime:      time.Now().Format(time.RFC3339),
		ChangePassword: dto,
		ProtocolType:   dtoNVR.Protocol,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

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

				err := reposity.UpdateSingleColumn[models.NVR](idNVR, "password", dto.PasswordNew)
				if err != nil {
					jsonRsp.Code = http.StatusInternalServerError
					jsonRsp.Message = err.Error()
					c.JSON(http.StatusInternalServerError, &jsonRsp)
					return
				}

				jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
				jsonRsp.Data = dto
				c.JSON(http.StatusOK, &jsonRsp)
				return
			}
		}
	}
}
