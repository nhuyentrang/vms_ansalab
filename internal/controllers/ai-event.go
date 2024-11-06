package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"vms/internal/models"
	"vms/statuscode"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
)

// CreateAIEvent		godoc
// @Summary      	Create a new aiEvent
// @Description  	Takes a aiEvent JSON and store in DB. Return saved JSON.
// @Tags         	ai-events
// @Produce			json
// @Param        	aiEvent  body   models.DTO_AIEvent_Create  true  "AIEvent JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_AIEvent_Create]
// @Router       	/ai-events [post]
// @Security		BearerAuth
func CreateAIEvent(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_AIEvent_Create]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_AIEvent_Create
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = statuscode.StatusBindingInputJsonFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Create
	dto, err := reposity.CreateItemFromDTO[models.DTO_AIEvent_Create, models.AIEvent](dto)
	if err != nil {
		jsonRsp.Code = statuscode.StatusCreateItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Response
	jsonRsp.Data = dto
	c.JSON(http.StatusCreated, &jsonRsp)
}

// ReadAIEvent		 godoc
// @Summary      Get single aiEvent by id
// @Description  Returns the aiEvent whose ID value matches the id.
// @Tags         ai-events
// @Produce      json
// @Param        id  path  string  true  "Read aiEvent by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_AIEvent]
// @Router       /ai-events/{id} [get]
// @Security		BearerAuth
func ReadAIEvent(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_AIEvent]()
	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_AIEvent, models.AIEvent](c.Param("id"))
	if err != nil {
		jsonRsp.Code = statuscode.StatusUpdateItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateAIEvent		 	godoc
// @Summary      	Update single aiEvent by id
// @Description  	Updates and returns a single aiEvent whose ID value matches the id. New data must be passed in the body.
// @Tags         	ai-events
// @Produce      	json
// @Param        	id  path  string  true  "Update aiEvent by id"
// @Param        	aiEvent  body      models.DTO_AIEvent_Create  true  "AIEvent JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_AIEvent_Create]
// @Router       	/ai-events/{id} [put]
// @Security		BearerAuth
func UpdateAIEvent(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_AIEvent_Create]()

	// Get new data from body
	var dto models.DTO_AIEvent_Create
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_AIEvent_Create, models.AIEvent](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Return
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// DeleteAIEvent	 godoc
// @Summary      Remove single aiEvent by id
// @Description  Delete a single entry from the reposity based on id.
// @Tags         ai-events
// @Produce      json
// @Param        id  path  string  true  "Delete aiEvent by id"
// @Success      204
// @Router       /ai-events/{id} [delete]
// @Security		BearerAuth
func DeleteAIEvent(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_AIEvent]()

	err := reposity.DeleteItemByID[models.AIEvent](c.Param("id"))
	if err != nil {
		jsonRsp.Code = statuscode.StatusDeleteItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	c.JSON(http.StatusOK, &jsonRsp)
}

// GetAIEvents		godoc
// @Summary      	Get all ai events with query filter
// @Description  	Responds with the list of all aiEvent as JSON.
// @Tags         	ai-events
// @Param   		keyword			query	string	false	"event name keyword"			minlength(1)  	maxlength(100)
// @Param   		type			query	string	false	"ai event type"				Enums(statusOfDevice, motionDetection)
// @Param   		deviceType		query	string	false	"ai device type"				Enums(ipCamera, smartCamera, nvr, smartNVR)
// @Param   		status			query	string	false	"ai event status"				Enums(new, inprogress, resolved, closed)
// @Param   		fromDate		query	int		false	"ai event start date, timestamp millisecond"
// @Param   		toDate			query	int		false	"ai event end date, timestamp millisecond"
// @Param   		siteID  		query	string	false	"site uuid"						minlength(1)  	maxlength(100)
// @Param   		sort  			query	string	false	"sort"							default(-created_at)
// @Param   		limit			query	int     false  	"limit"          				minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          				minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}  models.JsonDTOListRsp[models.DTO_AIEvent_Read_BasicInfo]
// @Router       	/ai-events [get]
// @Security		BearerAuth
func GetAIEvents(c *gin.Context) {
	jsonRsp := models.NewJsonDTOListRsp[models.DTO_AIEvent_Read_BasicInfo]()

	// Get param
	keyword := c.Query("keyword")
	eventType := c.Query("type")
	deviceType := c.Query("deviceType")
	status := c.Query("status")
	siteID := c.Query("siteID")
	fromDate, _ := strconv.Atoi(c.Query("fromDate"))
	toDate, _ := strconv.Atoi(c.Query("toDate"))
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	fmt.Println(
		"keyword: ", keyword,
		" - eventType: ", eventType,
		" - deviceType: ", deviceType,
		" - status: ", status,
		" - siteID: ", siteID,
		" - fromDate: ", fromDate,
		" - toDate: ", toDate,
		" - sort: ", sort,
		" - limit: ", limit,
		" - page: ", page)

	// Build query
	query := reposity.NewQuery[models.DTO_AIEvent_Read_BasicInfo, models.AIEvent]()
	if keyword != "" {
		//query.AddTwoConditionOfTextField("AND", "name", "LIKE", keyword, "OR", "code", "LIKE", keyword)
		query.AddConditionOfJsonbField("AND", "type", "name", "LIKE", keyword)
	}
	if eventType != "" {
		query.AddConditionOfJsonbField("AND", "type", "id", "=", eventType)
	}
	if deviceType != "" {
		query.AddConditionOfJsonbField("AND", "device_type", "id", "=", deviceType)
	}
	if status != "" {
		query.AddConditionOfJsonbField("AND", "status", "id", "=", status)
	}
	if siteID != "" {
		query.AddConditionOfJsonbField("AND", "site_info", "id", "=", siteID)
	}
	if fromDate > 0 {
		query.AddConditionOfTextField("AND", "atts", ">", fromDate)
	}
	if toDate > 0 {
		query.AddConditionOfTextField("AND", "atts", "<", toDate)
	}
	// Exec query
	dtos, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRsp.Code = statuscode.StatusSearchItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Count = count
	jsonRsp.Data = dtos
	jsonRsp.Page = int64(page)
	jsonRsp.Size = int64(len(dtos))
	c.JSON(http.StatusOK, &jsonRsp)
}

// GetAIDeviceTypes		godoc
// @Summary      	Get types of ai device
// @Description  	Responds with the list of all ai device as JSON.
// @Tags         	ai-events
// @Produce      	json
// @Success      	200  {object}  models.JsonAIDeviceTypeRsp
// @Router       	/ai-events/options/device-types [get]
// @Security		BearerAuth
func GetAIDeviceTypes(c *gin.Context) {

	var jsonAIDeviceTypeRsp models.JsonAIDeviceTypeRsp

	jsonAIDeviceTypeRsp.Data = make([]models.KeyValue, 0)

	jsonAIDeviceTypeRsp.Data = append(jsonAIDeviceTypeRsp.Data, models.KeyValue{
		ID:   "ipCamera",
		Name: "IP Camera",
	})

	jsonAIDeviceTypeRsp.Data = append(jsonAIDeviceTypeRsp.Data, models.KeyValue{
		ID:   "smartCamera",
		Name: "Smart Camera",
	})

	jsonAIDeviceTypeRsp.Data = append(jsonAIDeviceTypeRsp.Data, models.KeyValue{
		ID:   "nvr",
		Name: "NVR",
	})

	jsonAIDeviceTypeRsp.Data = append(jsonAIDeviceTypeRsp.Data, models.KeyValue{
		ID:   "smartNVR",
		Name: "Smart NVR",
	})

	jsonAIDeviceTypeRsp.Code = 0
	jsonAIDeviceTypeRsp.Page = 1
	jsonAIDeviceTypeRsp.Count = len(jsonAIDeviceTypeRsp.Data)
	jsonAIDeviceTypeRsp.Size = jsonAIDeviceTypeRsp.Count
	c.JSON(http.StatusOK, &jsonAIDeviceTypeRsp)
}

// GetAIEventTypes		godoc
// @Summary      	Get types of ai event
// @Description  	Responds with the list of all ai event type as JSON.
// @Tags         	ai-events
// @Produce      	json
// @Success      	200  {object}  models.JsonAIEventTypeRsp
// @Router       	/ai-events/options/types [get]
// @Security		BearerAuth
func GetAIEventTypes(c *gin.Context) {

	var jsonAIEventTypeRsp models.JsonAIEventTypeRsp

	jsonAIEventTypeRsp.Data = make([]models.KeyValue, 0)

	jsonAIEventTypeRsp.Data = append(jsonAIEventTypeRsp.Data, models.KeyValue{
		ID:   "faceDetection",
		Name: "Phát hiện mặt",
	})

	jsonAIEventTypeRsp.Data = append(jsonAIEventTypeRsp.Data, models.KeyValue{
		ID:   "instructionDetection",
		Name: "Phát hiện vượt rào",
	})

	jsonAIEventTypeRsp.Code = 0
	jsonAIEventTypeRsp.Page = 1
	jsonAIEventTypeRsp.Count = len(jsonAIEventTypeRsp.Data)
	jsonAIEventTypeRsp.Size = jsonAIEventTypeRsp.Count
	c.JSON(http.StatusOK, &jsonAIEventTypeRsp)
}

// GetAIEventStatusTypes		godoc
// @Summary      	Get types of ai event status
// @Description  	Responds with the list of all ai event status type as JSON.
// @Tags         	ai-events
// @Produce      	json
// @Success      	200  {object}  models.JsonAIEventStatusTypeRsp
// @Router       	/ai-events/options/status-types [get]
// @Security		BearerAuth
func GetAIEventStatusTypes(c *gin.Context) {

	var jsonAIEventStatusTypeRsp models.JsonAIEventStatusTypeRsp

	jsonAIEventStatusTypeRsp.Data = make([]models.KeyValue, 0)

	jsonAIEventStatusTypeRsp.Data = append(jsonAIEventStatusTypeRsp.Data, models.KeyValue{
		ID:   "new",
		Name: "Mới",
	})

	jsonAIEventStatusTypeRsp.Data = append(jsonAIEventStatusTypeRsp.Data, models.KeyValue{
		ID:   "inprogress",
		Name: "Đang xử lý",
	})

	jsonAIEventStatusTypeRsp.Data = append(jsonAIEventStatusTypeRsp.Data, models.KeyValue{
		ID:   "resolved",
		Name: "Đã xử lý",
	})

	jsonAIEventStatusTypeRsp.Data = append(jsonAIEventStatusTypeRsp.Data, models.KeyValue{
		ID:   "closed",
		Name: "Đóng",
	})

	jsonAIEventStatusTypeRsp.Code = 0
	jsonAIEventStatusTypeRsp.Page = 1
	jsonAIEventStatusTypeRsp.Count = len(jsonAIEventStatusTypeRsp.Data)
	jsonAIEventStatusTypeRsp.Size = jsonAIEventStatusTypeRsp.Count
	c.JSON(http.StatusOK, &jsonAIEventStatusTypeRsp)
}

// GetAIEventLevelTypes		godoc
// @Summary      	Get types of ai event level
// @Description  	Responds with the list of all ai event level type as JSON.
// @Tags         	ai-events
// @Produce      	json
// @Success      	200  {object}  models.JsonAIEventLevelTypeRsp
// @Router       	/ai-events/options/level-types [get]
// @Security		BearerAuth
func GetAIEventLevelTypes(c *gin.Context) {

	var jsonAIEventLevelTypeRsp models.JsonAIEventLevelTypeRsp

	jsonAIEventLevelTypeRsp.Data = make([]models.KeyValue, 0)

	jsonAIEventLevelTypeRsp.Data = append(jsonAIEventLevelTypeRsp.Data, models.KeyValue{
		ID:   "low",
		Name: "Thấp",
	})

	jsonAIEventLevelTypeRsp.Data = append(jsonAIEventLevelTypeRsp.Data, models.KeyValue{
		ID:   "mid",
		Name: "Trung bình",
	})

	jsonAIEventLevelTypeRsp.Data = append(jsonAIEventLevelTypeRsp.Data, models.KeyValue{
		ID:   "high",
		Name: "Cao",
	})

	jsonAIEventLevelTypeRsp.Code = 0
	jsonAIEventLevelTypeRsp.Page = 1
	jsonAIEventLevelTypeRsp.Count = len(jsonAIEventLevelTypeRsp.Data)
	jsonAIEventLevelTypeRsp.Size = jsonAIEventLevelTypeRsp.Count
	c.JSON(http.StatusOK, &jsonAIEventLevelTypeRsp)
}
