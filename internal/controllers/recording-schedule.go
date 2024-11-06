package controllers

import (
	"net/http"
	"strconv"

	"vms/internal/models"
	"vms/statuscode"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
)

// CreateRecordingSchedule		godoc
// @Summary      	Create a new recordingSchedule
// @Description  	Takes a recordingSchedule JSON and store in DB. Return saved JSON.
// @Tags         	recording-schedules
// @Produce			json
// @Param        	recordingSchedule  body   models.DTO_RecordingSchedule_Create  true  "RecordingSchedule JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_RecordingSchedule]
// @Router       	/recording-schedules [post]
// @Security		BearerAuth
func CreateRecordingSchedule(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_RecordingSchedule_Create]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_RecordingSchedule_Create
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = statuscode.StatusBindingInputJsonFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Create new block
	dto, err := reposity.CreateItemFromDTO[models.DTO_RecordingSchedule_Create, models.RecordingSchedule](dto)
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

// ReadRecordingSchedule		 godoc
// @Summary      Get single recordingSchedule by id
// @Description  Returns the recordingSchedule whose ID value matches the id.
// @Tags         recording-schedules
// @Produce      json
// @Param        id  path  string  true  "Read recordingSchedule by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_RecordingSchedule]
// @Router       /recording-schedules/{id} [get]
// @Security		BearerAuth
func ReadRecordingSchedule(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_RecordingSchedule]()
	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_RecordingSchedule, models.RecordingSchedule](c.Param("id"))
	if err != nil {
		jsonRsp.Code = statuscode.StatusUpdateItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateRecordingSchedule		 	godoc
// @Summary      	Update single recordingSchedule by id
// @Description  	Updates and returns a single recordingSchedule whose ID value matches the id. New data must be passed in the body.
// @Tags         	recording-schedules
// @Produce      	json
// @Param        	id  path  string  true  "Update recordingSchedule by id"
// @Param        	recordingSchedule  body      models.DTO_RecordingSchedule_Create  true  "RecordingSchedule JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_RecordingSchedule_Create]
// @Router       	/recording-schedules/{id} [put]
// @Security		BearerAuth
func UpdateRecordingSchedule(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_RecordingSchedule_Create]()

	// Get new data from body
	var dto models.DTO_RecordingSchedule_Create
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_RecordingSchedule_Create, models.RecordingSchedule](c.Param("id"), dto)
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

// DeleteRecordingSchedule	 godoc
// @Summary      Remove single recordingSchedule by id
// @Description  Delete a single entry from the reposity based on id.
// @Tags         recording-schedules
// @Produce      json
// @Param        id  path  string  true  "Delete recordingSchedule by id"
// @Success      204
// @Router       /recording-schedules/{id} [delete]
// @Security		BearerAuth
func DeleteRecordingSchedule(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_RecordingSchedule]()

	err := reposity.DeleteItemByID[models.RecordingSchedule](c.Param("id"))
	if err != nil {
		jsonRsp.Code = statuscode.StatusDeleteItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	c.JSON(http.StatusNoContent, &jsonRsp)
}

// GetRecordingSchedules		godoc
// @Summary      	Get all recording schedules with query filter
// @Description  	Responds with the list of all recordingSchedule as JSON.
// @Tags         	recording-schedules
// @Param   		keyword			query	string	false	"recordingSchedule name keyword"		minlength(1)  	maxlength(100)
// @Param   		sort  			query	string	false	"sort"							default(-created_at)
// @Param   		limit			query	int     false  	"limit"          				minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          				minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}  models.JsonDTOListRsp[models.DTO_RecordingSchedule_Read_BasicInfo]
// @Router       	/recording-schedules [get]
// @Security		BearerAuth
func GetRecordingSchedules(c *gin.Context) {

	jsonRsp := models.NewJsonDTOListRsp[models.DTO_RecordingSchedule_Read_BasicInfo]()

	// Get param
	keyword := c.Query("keyword")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))

	// Build query
	query := reposity.NewQuery[models.DTO_RecordingSchedule_Read_BasicInfo, models.RecordingSchedule]()

	if keyword != "" {
		query.AddConditionOfTextField("AND", "name", "LIKE", keyword)
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
