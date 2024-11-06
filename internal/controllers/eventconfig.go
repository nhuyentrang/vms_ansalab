package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"vms/internal/models"

	"vms/comongo/kafkaclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateEventconfig		godoc
// @Summary      	Create a new event config camera
// @Description  	Takes a eventConfig JSON and store in DB. Return saved JSON.
// @Tags         	event-config
// @Produce			json
// @Param        	eventConfig  body   models.DTO_EventConfig  true  "EventConfig JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_EventConfig]
// @Router       	/event-config [post]
// @Security		BearerAuth
func CreateEventconfig(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_EventConfig]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_EventConfig
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Create EventConfig
	dto, err := reposity.CreateItemFromDTO[models.DTO_EventConfig, models.EventConfig](dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Data = dto
	c.JSON(http.StatusCreated, &jsonRsp)
}

// ReadEventConfig		 godoc
// @Summary      Get single EventConfig by id
// @Description  Returns the event-config whose ID value matches the id.
// @Tags         event-config
// @Produce      json
// @Param        id  path  string  true  "Search event-config by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_EventConfig]
// @Router       /event-config/{id} [get]
// @Security		BearerAuth
func ReadEventconfig(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_EventConfig]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_EventConfig, models.EventConfig](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateEventConfig		 	godoc
// @Summary      	Update single event config by id
// @Description  	Updates and returns a single event config whose ID value matches the id. New data must be passed in the body.
// @Tags         	event-config
// @Produce      	json
// @Param        	id  path  string  true  "Update camera by id"
// @Param        	eventconfig  body      models.DTO_EventConfig  true  "EventConfig JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_EventConfig]
// @Router       	/event-config/{id} [put]
// @Security		BearerAuth
func UpdateEventconfig(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_EventConfig]()

	// Get new data from body
	var dto models.DTO_EventConfig
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_EventConfig, models.EventConfig](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	dtoCamera, err := reposity.ReadItemWithFilterIntoDTO[models.DTOCamera, models.Camera]("id = ?", dto.CameraID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Override ID
	dto.ID, _ = uuid.Parse(c.Param("id"))

	cmd := models.DeviceCommand{
		CommandID: uuid.New().String(),
		Cmd:       "",
		EventTime: time.Now().Format(time.RFC3339),
		IPAddress: dtoCamera.IPAddress,
		HttpPort:  dtoCamera.HttpPort,
		UserName:  dtoCamera.Username,
		Password:  dtoCamera.Password,
		IndexNVR:  dtoCamera.IndexNVR,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	// Return
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// DeleteEventConfig	 godoc
// @Summary      Remove single camera by id
// @Description  Delete a single entry from the reposity based on id.
// @Tags         event-config
// @Produce      json
// @Param        id  path  string  true  "Delete event config by id"
// @Success      204
// @Router       /event-config/{id} [delete]
// @Security		BearerAuth
func DeleteEventconfig(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_EventConfig]()

	err := reposity.DeleteItemByID[models.EventConfig](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	c.JSON(http.StatusNoContent, &jsonRsp)
}

// ScanEventConfig	 godoc
// @Summary      Remove single camera by id
// @Description  Scan for supported models
// @Tags         event-config
// @Produce      json
// @Success      204
// @Router       /event-config/scan [post]
// @Security		BearerAuth
func PostScanEventModel(c *gin.Context) {
	jsonRsp := models.NewJsonDTOListRsp[models.DTO_CmdEvent]()

	cmd := models.DTO_CmdEvent{
		CommandID: uuid.New().String(),
		Cmd:       "EventScan",
		Code:      "EventScan",
		EventTime: time.Now().Format(time.RFC3339),
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), "DEV_EVENT_CONFIG")
	jsonRsp.Message = "Call successfully"
	c.JSON(http.StatusOK, &jsonRsp)
}
