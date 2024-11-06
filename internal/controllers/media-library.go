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

// CreateMediaLibrary		godoc
// @Summary      	Create a new mediaLibrary
// @Description  	Takes a mediaLibrary JSON and store in DB. Return saved JSON.
// @Tags         	media-libraries
// @Produce			json
// @Param        	mediaLibrary  body   models.DTO_MediaLibrary_Create  true  "MediaLibrary JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_MediaLibrary_Create]
// @Router       	/media-libraries [post]
// @Security		BearerAuth
func CreateMediaLibrary(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_MediaLibrary_Create]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_MediaLibrary_Create
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = statuscode.StatusBindingInputJsonFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Create
	dto, err := reposity.CreateItemFromDTO[models.DTO_MediaLibrary_Create, models.MediaLibrary](dto)
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

// ReadMediaLibrary		 godoc
// @Summary      Get single mediaLibrary by id
// @Description  Returns the mediaLibrary whose ID value matches the id.
// @Tags         media-libraries
// @Produce      json
// @Param        id  path  string  true  "Read mediaLibrary by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_MediaLibrary]
// @Router       /media-libraries/{id} [get]
// @Security		BearerAuth
func ReadMediaLibrary(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_MediaLibrary]()
	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_MediaLibrary, models.MediaLibrary](c.Param("id"))
	if err != nil {
		jsonRsp.Code = statuscode.StatusUpdateItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateMediaLibrary		 	godoc
// @Summary      	Update single mediaLibrary by id
// @Description  	Updates and returns a single mediaLibrary whose ID value matches the id. New data must be passed in the body.
// @Tags         	media-libraries
// @Produce      	json
// @Param        	id  path  string  true  "Update mediaLibrary by id"
// @Param        	mediaLibrary  body      models.DTO_MediaLibrary_Create  true  "MediaLibrary JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_MediaLibrary_Create]
// @Router       	/media-libraries/{id} [put]
// @Security		BearerAuth
func UpdateMediaLibrary(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_MediaLibrary_Create]()

	// Get new data from body
	var dto models.DTO_MediaLibrary_Create
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_MediaLibrary_Create, models.MediaLibrary](c.Param("id"), dto)
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

// DeleteMediaLibrary	 godoc
// @Summary      Remove single mediaLibrary by id
// @Description  Delete a single entry from the reposity based on id.
// @Tags         media-libraries
// @Produce      json
// @Param        id  path  string  true  "Delete mediaLibrary by id"
// @Success      204
// @Router       /media-libraries/{id} [delete]
// @Security		BearerAuth
func DeleteMediaLibrary(c *gin.Context) {
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

// GetMediaLibraries		godoc
// @Summary      	Get all camera groups with query filter
// @Description  	Responds with the list of all mediaLibrary as JSON.
// @Tags         	media-libraries
// @Param   		keyword			query	string	false	"media library name keyword"	minlength(1)  	maxlength(100)
// @Param   		type			query	string	false	"media library type"			minlength(1)  	maxlength(100)
// @Param   		fromDate		query	int		false	"media library start date, timestamp millisecond"
// @Param   		toDate			query	int		false	"media library end date, timestamp millisecond"
// @Param   		sort  			query	string	false	"sort"							default(-created_at)
// @Param   		limit			query	int     false  	"limit"          				minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          				minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}  models.JsonDTOListRsp[models.DTO_MediaLibrary_Read_BasicInfo]
// @Router       	/media-libraries [get]
// @Security		BearerAuth
func GetMediaLibraries(c *gin.Context) {
	jsonRsp := models.NewJsonDTOListRsp[models.DTO_MediaLibrary_Read_BasicInfo]()

	// Get param
	keyword := c.Query("keyword")
	fileType := c.Query("type")
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
		" - type: ", fileType,
		" - fromDate: ", fromDate,
		" - toDate: ", toDate,
		" - sort: ", sort,
		" - limit: ", limit,
		" - page: ", page)

	// Build query
	query := reposity.NewQuery[models.DTO_MediaLibrary_Read_BasicInfo, models.MediaLibrary]()
	if keyword != "" {
		//query.AddTwoConditionOfTextField("AND", "name", "LIKE", keyword, "OR", "code", "LIKE", keyword)
		query.AddConditionOfTextField("AND", "name", "LIKE", keyword)
	}
	if fileType != "" {
		query.AddConditionOfJsonbField("AND", "type", "id", "=", fileType)
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

// GetMediaFileTypes		godoc
// @Summary      	Get types of media file
// @Description  	Responds with the list of all type of media file as JSON.
// @Tags         	media-libraries
// @Produce      	json
// @Success      	200  {object}  models.JsonMediaFileTypeRsp
// @Router       	/media-libraries/options/file-types [get]
// @Security		BearerAuth
func GetMediaFileTypes(c *gin.Context) {

	var jsonMediaFileTypeRsp models.JsonMediaFileTypeRsp

	jsonMediaFileTypeRsp.Data = make([]models.KeyValue, 0)

	jsonMediaFileTypeRsp.Data = append(jsonMediaFileTypeRsp.Data, models.KeyValue{
		ID:   "image",
		Name: "Ảnh",
	})

	jsonMediaFileTypeRsp.Data = append(jsonMediaFileTypeRsp.Data, models.KeyValue{
		ID:   "video",
		Name: "Video",
	})

	jsonMediaFileTypeRsp.Code = 0
	jsonMediaFileTypeRsp.Page = 1
	jsonMediaFileTypeRsp.Count = len(jsonMediaFileTypeRsp.Data)
	jsonMediaFileTypeRsp.Size = jsonMediaFileTypeRsp.Count
	c.JSON(http.StatusOK, &jsonMediaFileTypeRsp)
}
