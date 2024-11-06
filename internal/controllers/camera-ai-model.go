package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"vms/internal/models"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var countModelAI int64

// CreateCameraModelAI		godoc
// @Summary      	Create a new CameraModelAI
// @Description  	Takes a CameraModelAI JSON and store in DB. Return saved JSON.
// @Tags         	CameraModelAI
// @Produce			json
// @Param        	CameraModelAI  body   models.DTO_CameraModelAI  true  "CameraModelAI JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_CameraModelAI]
// @Router       	/camera-model-ai [post]
// @Security		BearerAuth
func CreateCameraModelAI(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraModelAI]()

	var dtoResp models.DTO_CameraModelAI
	if err := c.BindJSON(&dtoResp); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	countModelAI++
	dtoResp.Count = countModelAI

	dto, err := reposity.CreateItemFromDTO[models.DTO_CameraModelAI, models.CameraModelAI](dtoResp)
	if err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	jsonRsp.Data = dto
	c.JSON(http.StatusCreated, &jsonRsp)
}

// ReadCameraModelAI		godoc
// @Summary      	Create a new CameraModelAI
// @Description  	Takes a CameraModelAI JSON and store in DB. Return saved JSON.
// @Tags         	CameraModelAI
// @Produce			json
// @Param        id  path  string  true  "Search CameraModelAI by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_CameraModelAI]
// @Router       /camera-model-ai/{id} [get]
// @Security		BearerAuth
func ReadCameraModelAI(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraModelAI]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_CameraModelAI, models.CameraModelAI](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// DeleteCameraModelAI		 	godoc
// @Summary      	Delete single CameraModelAI by id
// @Description  	Delete and returns a single event whose ID value matches the id. New data must be passed in the body.
// @Tags         	CameraModelAI
// @Produce      	json
// @Param        	id  path  string  true  "Update CameraModelAI by id"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_CameraModelAI]
// @Router       	/camera-model-ai/{id} [delete]
// @Security		BearerAuth
func DeleteCameraModelAI(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraModelAI]()
	//TODO: Delete Loi
	err := reposity.DeleteItemByID[models.CameraModelAI](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	c.JSON(http.StatusNoContent, &jsonRsp)
}

// UpdateCameraModelAI		 	godoc
// @Summary      	Update single CameraModelAI by id
// @Description  	Updates and returns a single CameraModelAI whose ID value matches the id. New data must be passed in the body.
// @Tags         	CameraModelAI
// @Produce      	json
// @Param        	id  path  string  true  "Update CameraModelAI by id"
// @Param        	CameraModelAI  body      models.DTO_CameraModelAI  true  "Camera JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_CameraModelAI]
// @Router       	/camera-model-ai/{id} [put]
// @Security		BearerAuth
func UpdateCameraModelAI(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraModelAI]()

	// Get new data from body
	var dto models.DTO_CameraModelAI
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_CameraModelAI, models.CameraModelAI](c.Param("id"), dto)
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
	c.JSON(http.StatusOK, &jsonRsp)
}

// GetCameraModelAIs		godoc
// @Summary      	Get all CameraModelAI groups with query filter
// @Description  	Responds with the list of all CameraModelAI as JSON.
// @Tags         	CameraModelAI
// @Param   		keyword			query	string	false	"CameraModelAI name keyword"		minlength(1)  	maxlength(100)
// @Param   		sort  			query	string	false	"sort"						default(+created_at)
// @Param   		limit			query	int     false  	"limit"          			minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          			minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_CameraModelAI]
// @Router       	/camera-model-ai [get]
// @Security		BearerAuth
func GetCameraModelAIs(c *gin.Context) {
	jsonRspDTOCamerasBasicInfos := models.NewJsonDTOListRsp[models.DTO_CameraModelAI]()

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
	query := reposity.NewQuery[models.DTO_CameraModelAI, models.CameraModelAI]()

	// Search for keyword in name
	if keyword != "" {
		query.AddConditionOfTextField("AND", "model_name", "LIKE", keyword)
	}

	// Exec query
	dtoCameraBasics, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTOCamerasBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCamerasBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCamerasBasicInfos)
		return
	}

	jsonRspDTOCamerasBasicInfos.Count = count
	jsonRspDTOCamerasBasicInfos.Data = dtoCameraBasics
	jsonRspDTOCamerasBasicInfos.Page = int64(page)
	jsonRspDTOCamerasBasicInfos.Size = int64(len(dtoCameraBasics))
	c.JSON(http.StatusOK, &jsonRspDTOCamerasBasicInfos)
}

// GetCameraModelAIs		godoc
// @Summary      	Get all CameraModelAI groups with query filter
// @Description  	Responds with the list of all CameraModelAI as JSON.
// @Tags         	CameraModelAI
// @Param   		keyword			query	string	false	"CameraModelAI name keyword"		minlength(1)  	maxlength(100)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_CameraModelAI]
// @Router       	/camera-model-ai/count-camera-of-model-ai [get]
// @Security		BearerAuth
func GetCountOfModelAI(c *gin.Context) {
	// Gửi phản hồi JSON
	jsonRspDTOEventsBasicInfos := models.NewJsonDTOListRsp[models.DTO_CameraModelAI]()

	// Get param
	keyword := c.Query("keyword")

	// Build query
	query := reposity.NewQuery[models.DTO_CameraModelAI, models.CameraModelAI]()

	// Search for keyword in name
	if keyword != "" {
		query.AddConditionOfTextField("AND", "model_name", "LIKE", keyword)
	}

	// Exec query
	dtoEventBasics, _, err := query.ExecWithPaging("-created_at", 1, 1)
	if err != nil {
		jsonRspDTOEventsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOEventsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOEventsBasicInfos)
		return
	}
	// Calculate the count
	modelCount := int64(len(dtoEventBasics))

	// Prepare the response data
	jsonRspDTOEventsBasicInfos.Data = append(jsonRspDTOEventsBasicInfos.Data, models.DTO_CameraModelAI{Count: modelCount})

	// Send the response
	c.JSON(http.StatusOK, &jsonRspDTOEventsBasicInfos)
}
