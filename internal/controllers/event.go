package controllers

import (
	"net/http"
	"strconv"

	"vms/internal/models"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
)

// CreateEvent		godoc
// @Summary      	Create a new event
// @Description  	Takes a event JSON and store in DB. Return saved JSON.
// @Tags         	events
// @Produce			json
// @Param        	Event  body   models.DTO_Event_Created  true  "Event JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_Event_Created]
// @Router       	/events [post]
// @Security		BearerAuth
func CreateEvent(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_Event_Created]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_Event_Created
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Create Event
	dto, err := reposity.CreateItemFromDTO[models.DTO_Event_Created, models.Event](dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusCreated, &jsonRsp)
}

// ReadEvent		 godoc
// @Summary      Get single event by id
// @Description  Returns the event whose ID value matches the id.
// @Tags         events
// @Produce      json
// @Param        id  path  string  true  "Search event by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_Event]
// @Router       /events/{id} [get]
// @Security		BearerAuth
func ReadEvent(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_Event]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_Event, models.Event](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// GetEvents		godoc
// @Summary      	Get all Event groups with query filter
// @Description  	Responds with the list of all Event as JSON.
// @Tags         	events
// @Param   		keyword			query	string	false	"Event name keyword"		minlength(1)  	maxlength(100)
// @Param   		limit			query	int     false  	"limit"          			minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          			minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_Event_Read_BasicInfo]
// @Router       	/events [get]
// @Security		BearerAuth
func GetEvents(c *gin.Context) {
	jsonRspDTOEventsBasicInfos := models.NewJsonDTOListRsp[models.DTO_Event_Read_BasicInfo]()

	// Get param
	keyword := c.Query("keyword")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))

	// Build query
	query := reposity.NewQuery[models.DTO_Event_Read_BasicInfo, models.Event]()

	// Search for keyword in name
	if keyword != "" {
		query.AddConditionOfTextField("AND", "eventName", "LIKE", keyword)
	}

	// Exec query
	dtoEventBasics, count, err := query.ExecWithPaging("+created_at", limit, page)
	if err != nil {
		jsonRspDTOEventsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOEventsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOEventsBasicInfos)
		return
	}

	jsonRspDTOEventsBasicInfos.Count = count
	jsonRspDTOEventsBasicInfos.Data = dtoEventBasics
	jsonRspDTOEventsBasicInfos.Page = int64(page)
	jsonRspDTOEventsBasicInfos.Size = int64(len(dtoEventBasics))
	c.JSON(http.StatusOK, &jsonRspDTOEventsBasicInfos)
}

// GetEvents		godoc
// @Summary      	Get all Event groups with query filter
// @Description  	Responds with the list of all Event as JSON.
// @Tags         	events
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_Event_Read_BasicInfo]
// @Router       	/events/list-event-types [get]
// @Security		BearerAuth
func GetListTypeEvent(c *gin.Context) {
	jsonConfigTypeRsp := models.NewJsonDTORsp[[]models.KeyValue]()
	jsonConfigTypeRsp.Data = make([]models.KeyValue, 0)

	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "BLACK_LIST",
		Name: "Đối tượng trong danh sách đen",
	})
	jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
		ID:   "ENTER_FORBIDDEN_AREA",
		Name: "Người lạ xâm nhập khu vực cấm",
	})

	jsonConfigTypeRsp.Code = 0
	c.JSON(http.StatusOK, &jsonConfigTypeRsp)
}

// UpdateEvent		 	godoc
// @Summary      	Update single event by id
// @Description  	Updates and returns a single event whose ID value matches the id. New data must be passed in the body.
// @Tags         	events
// @Produce      	json
// @Param        	id  path  string  true  "Update event by id"
// @Param        	event  body      models.DTO_Event_UpdateDescription  true  "Event JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_Event_UpdateDescription]
// @Router       	/events/{id}/description [put]
// @Security		BearerAuth
func UpdateDescription(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_Event_UpdateDescription]()

	var dto models.DTO_Event_UpdateDescription
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_Event_UpdateDescription, models.Event](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateEvent		 	godoc
// @Summary      	Update single event by id
// @Description  	Updates and returns a single event whose ID value matches the id. New data must be passed in the body.
// @Tags         	events
// @Produce      	json
// @Param        	id  path  string  true  "Update event by id"
// @Param        	event  body      models.DTO_Event_UpdateStatus  true  "Event JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_Event_UpdateStatus]
// @Router       	/events/{id}/update-event-status [put]
// @Security		BearerAuth
func UpdateStatus(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_Event_UpdateStatus]()

	var dto models.DTO_Event_UpdateStatus
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_Event_UpdateStatus, models.Event](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// GetEvents		godoc
// @Summary      	Get all Event groups with query filter
// @Description  	Responds with the list of all Event as JSON.
// @Tags         	events
// @Param   		areaId			query	string	false	"Event name areaId"		minlength(1)  	maxlength(100)
// @Param   		deviceId		query	string	false	"Event name deviceId"		minlength(1)  	maxlength(100)
// @Param   		startTime		query	string	false	"Event name startTime"		minlength(1)  	maxlength(100)
// @Param   		endTime			query	string	false	"Event name endTime"		minlength(1)  	maxlength(100)
// @Param   		keyword			query	string	false	"Event name keyword"		minlength(1)  	maxlength(100)
// @Param   		eventType		query	string	false	"Event name eventType"		minlength(1)  	maxlength(100)
// @Param   		limit			query	int     false  	"limit"          			minimum(1)    	default(50)
// @Param   		page			query	int     false  	"page"          			minimum(1) 		default(1)
// @Param   		status			query	string     false  	"status"          		minlength(1)  	maxlength(100)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_Event_Read_BasicInfo]
// @Router       	/events/filter [get]
// @Security		BearerAuth
func GetEventFilter(c *gin.Context) {
	jsonRspDTOEventsBasicInfos := models.NewJsonDTOListRsp[models.DTO_Event_Read_BasicInfo]()

	// Get param
	keyword := c.Query("keyword")
	areaId := c.Query("areaId")
	deviceId := c.Query("deviceId")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	eventType := c.Query("eventType")
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))

	// Build query
	query := reposity.NewQuery[models.DTO_Event_Read_BasicInfo, models.Event]()

	// Search for keyword in name
	if keyword != "" {
		query.AddConditionOfTextField("AND", "eventName", "LIKE", keyword)
	}

	if areaId != "" {
		query.AddConditionOfTextField("AND", "areaId", "=", areaId)
	}
	if deviceId != "" {
		query.AddConditionOfTextField("AND", "deviceId", "=", deviceId)
	}
	if eventType != "" {
		query.AddConditionOfTextField("AND", "eventType", "=", eventType)
	}
	if status != "" {
		query.AddConditionOfTextField("AND", "status", "=", status)
	}

	if endTime != "" {
	}

	if startTime != "" {
	}

	// Exec query
	dtoEventBasics, count, err := query.ExecWithPaging("+created_at", limit, page)
	if err != nil {
		jsonRspDTOEventsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOEventsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOEventsBasicInfos)
		return
	}

	jsonRspDTOEventsBasicInfos.Count = count
	jsonRspDTOEventsBasicInfos.Data = dtoEventBasics
	jsonRspDTOEventsBasicInfos.Page = int64(page)
	jsonRspDTOEventsBasicInfos.Size = int64(len(dtoEventBasics))
	c.JSON(http.StatusOK, &jsonRspDTOEventsBasicInfos)
}

// UpdateEvent		 	godoc
// @Summary      	Update single event by id
// @Description  	Updates and returns a single event whose ID value matches the id. New data must be passed in the body.
// @Tags         	events
// @Produce      	json
// @Param        	id  path  string  true  "Update event by id"
// @Param        	event  body      models.DTO_Event_ImageURL  true  "Event JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_Event_ImageURL]
// @Router       	/events/{id}/update-image [put]
// @Security		BearerAuth
func UpdateImageURL(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_Event_ImageURL]()

	var dto models.DTO_Event_ImageURL
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_Event_ImageURL, models.Event](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}
