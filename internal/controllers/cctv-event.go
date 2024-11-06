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

// CreateCCTVEvent		godoc
// @Summary      	Create a new cctvEvent
// @Description  	Takes a cctvEvent JSON and store in DB. Return saved JSON.
// @Tags         	cctv-events
// @Produce			json
// @Param        	cctvEvent  body   models.DTO_CCTVEvent_Create  true  "CCTVEvent JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_CCTVEvent_Create]
// @Router       	/cctv-events [post]
// @Security		BearerAuth
func CreateCCTVEvent(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CCTVEvent_Create]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_CCTVEvent_Create
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = statuscode.StatusBindingInputJsonFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Create
	dto, err := reposity.CreateItemFromDTO[models.DTO_CCTVEvent_Create, models.CCTVEvent](dto)
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

// ReadCCTVEvent		 godoc
// @Summary      Get single cctvEvent by id
// @Description  Returns the cctvEvent whose ID value matches the id.
// @Tags         cctv-events
// @Produce      json
// @Param        id  path  string  true  "Read cctvEvent by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_CCTVEvent]
// @Router       /cctv-events/{id} [get]
// @Security		BearerAuth
func ReadCCTVEvent(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CCTVEvent]()
	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_CCTVEvent, models.CCTVEvent](c.Param("id"))
	if err != nil {
		jsonRsp.Code = statuscode.StatusUpdateItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateCCTVEvent		 	godoc
// @Summary      	Update single cctvEvent by id
// @Description  	Updates and returns a single cctvEvent whose ID value matches the id. New data must be passed in the body.
// @Tags         	cctv-events
// @Produce      	json
// @Param        	id  path  string  true  "Update cctvEvent by id"
// @Param        	cctvEvent  body      models.DTO_CCTVEvent_Create  true  "CCTVEvent JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_CCTVEvent_Create]
// @Router       	/cctv-events/{id} [put]
// @Security		BearerAuth
func UpdateCCTVEvent(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CCTVEvent_Create]()

	// Get new data from body
	var dto models.DTO_CCTVEvent_Create
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_CCTVEvent_Create, models.CCTVEvent](c.Param("id"), dto)
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

// DeleteCCTVEvent	 godoc
// @Summary      Remove single cctvEvent by id
// @Description  Delete a single entry from the reposity based on id.
// @Tags         cctv-events
// @Produce      json
// @Param        id  path  string  true  "Delete cctvEvent by id"
// @Success      204
// @Router       /cctv-events/{id} [delete]
// @Security		BearerAuth
func DeleteCCTVEvent(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CCTVEvent]()

	err := reposity.DeleteItemByID[models.CCTVEvent](c.Param("id"))
	if err != nil {
		jsonRsp.Code = statuscode.StatusDeleteItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	c.JSON(http.StatusNoContent, &jsonRsp)
}

// GetCCTVEvents		godoc
// @Summary      	Get all cctv events with query filter
// @Description  	Responds with the list of all cctvEvent as JSON.
// @Tags         	cctv-events
// @Param   		keyword			query	string	false	"event name keyword"			minlength(1)  	maxlength(100)
// @Param   		type			query	string	false	"cctv event type"				Enums(statusOfDevice, motionDetection)
// @Param   		deviceType		query	string	false	"cctv device type"				Enums(ipCamera, smartCamera, nvr, smartNVR)
// @Param   		status			query	string	false	"cctv event status"				Enums(new, inprogress, resolved, closed)
// @Param   		fromDate		query	int		false	"cctv event start date, timestamp millisecond"
// @Param   		toDate			query	int		false	"cctv event end date, timestamp millisecond"
// @Param   		siteID  		query	string	false	"site uuid"						minlength(1)  	maxlength(100)
// @Param   		sort  			query	string	false	"sort"							default(-created_at)
// @Param   		limit			query	int     false  	"limit"          				minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          				minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}  models.JsonDTOListRsp[models.DTO_CCTVEvent_Read_BasicInfo]
// @Router       	/cctv-events [get]
// @Security		BearerAuth
func GetCCTVEvents(c *gin.Context) {
	jsonRsp := models.NewJsonDTOListRsp[models.DTO_CCTVEvent_Read_BasicInfo]()

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
	query := reposity.NewQuery[models.DTO_CCTVEvent_Read_BasicInfo, models.CCTVEvent]()
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

// GetCCTVDeviceTypes		godoc
// @Summary      	Get types of cctv device
// @Description  	Responds with the list of all cctv device as JSON.
// @Tags         	cctv-events
// @Produce      	json
// @Success      	200  {object}  models.JsonCCTVDeviceTypeRsp
// @Router       	/cctv-events/options/device-types [get]
// @Security		BearerAuth
func GetCCTVDeviceTypes(c *gin.Context) {

	var jsonCCTVDeviceTypeRsp models.JsonCCTVDeviceTypeRsp

	jsonCCTVDeviceTypeRsp.Data = make([]models.KeyValue, 0)

	jsonCCTVDeviceTypeRsp.Data = append(jsonCCTVDeviceTypeRsp.Data, models.KeyValue{
		ID:   "ipCamera",
		Name: "IP Camera",
	})

	jsonCCTVDeviceTypeRsp.Data = append(jsonCCTVDeviceTypeRsp.Data, models.KeyValue{
		ID:   "smartCamera",
		Name: "Smart Camera",
	})

	jsonCCTVDeviceTypeRsp.Data = append(jsonCCTVDeviceTypeRsp.Data, models.KeyValue{
		ID:   "nvr",
		Name: "NVR",
	})

	jsonCCTVDeviceTypeRsp.Data = append(jsonCCTVDeviceTypeRsp.Data, models.KeyValue{
		ID:   "smartNVR",
		Name: "Smart NVR",
	})

	jsonCCTVDeviceTypeRsp.Code = 0
	jsonCCTVDeviceTypeRsp.Page = 1
	jsonCCTVDeviceTypeRsp.Count = len(jsonCCTVDeviceTypeRsp.Data)
	jsonCCTVDeviceTypeRsp.Size = jsonCCTVDeviceTypeRsp.Count
	c.JSON(http.StatusOK, &jsonCCTVDeviceTypeRsp)
}

// GetCCTVEventTypes		godoc
// @Summary      	Get types of cctv event
// @Description  	Responds with the list of all cctv event type as JSON.
// @Tags         	cctv-events
// @Produce      	json
// @Success      	200  {object}  models.JsonCCTVEventTypeRsp
// @Router       	/cctv-events/options/types [get]
// @Security		BearerAuth
func GetCCTVEventTypes(c *gin.Context) {

	var jsonCCTVEventTypeRsp models.JsonCCTVEventTypeRsp

	jsonCCTVEventTypeRsp.Data = make([]models.KeyValue, 0)

	jsonCCTVEventTypeRsp.Data = append(jsonCCTVEventTypeRsp.Data, models.KeyValue{
		ID:   "statusOfDevice",
		Name: "Cảnh báo trạng thái của thiết bị",
	})

	jsonCCTVEventTypeRsp.Data = append(jsonCCTVEventTypeRsp.Data, models.KeyValue{
		ID:   "motionDetection",
		Name: "Cảnh báo phát hiện chuyển động",
	})

	jsonCCTVEventTypeRsp.Code = 0
	jsonCCTVEventTypeRsp.Page = 1
	jsonCCTVEventTypeRsp.Count = len(jsonCCTVEventTypeRsp.Data)
	jsonCCTVEventTypeRsp.Size = jsonCCTVEventTypeRsp.Count
	c.JSON(http.StatusOK, &jsonCCTVEventTypeRsp)
}

// GetCCTVEventStatusTypes		godoc
// @Summary      	Get types of cctv event status
// @Description  	Responds with the list of all cctv event status type as JSON.
// @Tags         	cctv-events
// @Produce      	json
// @Success      	200  {object}  models.JsonCCTVEventStatusTypeRsp
// @Router       	/cctv-events/options/status-types [get]
// @Security		BearerAuth
func GetCCTVEventStatusTypes(c *gin.Context) {

	var jsonCCTVEventStatusTypeRsp models.JsonCCTVEventStatusTypeRsp

	jsonCCTVEventStatusTypeRsp.Data = make([]models.KeyValue, 0)

	jsonCCTVEventStatusTypeRsp.Data = append(jsonCCTVEventStatusTypeRsp.Data, models.KeyValue{
		ID:   "new",
		Name: "Mới",
	})

	jsonCCTVEventStatusTypeRsp.Data = append(jsonCCTVEventStatusTypeRsp.Data, models.KeyValue{
		ID:   "inprogress",
		Name: "Đang xử lý",
	})

	jsonCCTVEventStatusTypeRsp.Data = append(jsonCCTVEventStatusTypeRsp.Data, models.KeyValue{
		ID:   "resolved",
		Name: "Đã xử lý",
	})

	jsonCCTVEventStatusTypeRsp.Data = append(jsonCCTVEventStatusTypeRsp.Data, models.KeyValue{
		ID:   "closed",
		Name: "Đóng",
	})

	jsonCCTVEventStatusTypeRsp.Code = 0
	jsonCCTVEventStatusTypeRsp.Page = 1
	jsonCCTVEventStatusTypeRsp.Count = len(jsonCCTVEventStatusTypeRsp.Data)
	jsonCCTVEventStatusTypeRsp.Size = jsonCCTVEventStatusTypeRsp.Count
	c.JSON(http.StatusOK, &jsonCCTVEventStatusTypeRsp)
}

// GetCCTVEventLevelTypes		godoc
// @Summary      	Get types of cctv event level
// @Description  	Responds with the list of all cctv event level type as JSON.
// @Tags         	cctv-events
// @Produce      	json
// @Success      	200  {object}  models.JsonCCTVEventLevelTypeRsp
// @Router       	/cctv-events/options/level-types [get]
// @Security		BearerAuth
func GetCCTVEventLevelTypes(c *gin.Context) {

	var jsonCCTVEventLevelTypeRsp models.JsonCCTVEventLevelTypeRsp

	jsonCCTVEventLevelTypeRsp.Data = make([]models.KeyValue, 0)

	jsonCCTVEventLevelTypeRsp.Data = append(jsonCCTVEventLevelTypeRsp.Data, models.KeyValue{
		ID:   "low",
		Name: "Thấp",
	})

	jsonCCTVEventLevelTypeRsp.Data = append(jsonCCTVEventLevelTypeRsp.Data, models.KeyValue{
		ID:   "mid",
		Name: "Trung bình",
	})

	jsonCCTVEventLevelTypeRsp.Data = append(jsonCCTVEventLevelTypeRsp.Data, models.KeyValue{
		ID:   "high",
		Name: "Cao",
	})

	jsonCCTVEventLevelTypeRsp.Code = 0
	jsonCCTVEventLevelTypeRsp.Page = 1
	jsonCCTVEventLevelTypeRsp.Count = len(jsonCCTVEventLevelTypeRsp.Data)
	jsonCCTVEventLevelTypeRsp.Size = jsonCCTVEventLevelTypeRsp.Count
	c.JSON(http.StatusOK, &jsonCCTVEventLevelTypeRsp)
}
