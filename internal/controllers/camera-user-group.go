package controllers

import (
	"fmt"
	"net/http"

	"strconv"

	"vms/internal/models"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
)

// CreateUserCameraGroup godoc
// @Summary      Create a new user camera group
// @Description  Accepts a JSON payload representing a user camera group, stores it in the database, and returns the saved entity.
// @Tags         cameras
// @Accept       json
// @Produce      json
// @Param        body  body   models.DTO_User_Camera_Group_BasicInfo  true  "User Camera Group JSON"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_User_Camera_Group_BasicInfo]  "User Camera Group created successfully."
// @Failure      400  "Bad Request - Invalid JSON payload."
// @Failure      500  "Internal Server Error - Unable to create the user camera group."
// @Router       /cameras/user [post]
// @Security     BearerAuth
func CreateUserCameraGroup(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_User_Camera_Group_BasicInfo]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_User_Camera_Group_BasicInfo
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Create User Camera Group
	dto, err := reposity.CreateItemFromDTO[models.DTO_User_Camera_Group_BasicInfo, models.UserCameraGroup](dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Data = dto
	c.JSON(http.StatusCreated, &jsonRsp)
}

// ReadUserCameraGroup godoc
// @Summary      Retrieve a user camera group
// @Description  Fetches a user camera group by its unique identifier and returns it.
// @Tags         cameras
// @Produce      json
// @Param        id  path  string  true  "User Camera Group ID"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_User_Camera_Group_BasicInfo]  "User Camera Group found and returned successfully."
// @Failure      404  "Not Found - User Camera Group does not exist or could not be found."
// @Router       /cameras/user/{id} [get]
// @Security     BearerAuth
func ReadUserCameraGroup(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_User_Camera_Group_BasicInfo]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_User_Camera_Group_BasicInfo, models.UserCameraGroup](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateUserCameraGroup updates an existing user camera group.
// @Summary      Update a user camera group
// @Description  Receives a JSON payload with the updated data for a user camera group and applies the changes to the specified group.
// @Tags         cameras
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "User Camera Group ID to update"
// @Param        UserCameraGroup  body   models.DTO_User_Camera_Group_BasicInfo  true  "Updated User Camera Group JSON"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_User_Camera_Group_BasicInfo]  "Successfully updated the user camera group."
// @Failure      400  "Bad Request - JSON payload is malformed or invalid."
// @Failure      404  "Not Found - Specified user camera group does not exist."
// @Failure      500  "Internal Server Error - Unable to update the user camera group."
// @Router       /cameras/user/{id} [put]
// @Security     BearerAuth
func UpdateUserCameraGroup(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_User_Camera_Group_BasicInfo]()

	// Bind the received JSON to the DTO
	var dto models.DTO_User_Camera_Group_BasicInfo
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Bad Request - JSON payload is malformed or invalid."
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Attempt to update the user camera group using the provided DTO
	updatedDto, err := reposity.UpdateItemByIDFromDTO[models.DTO_User_Camera_Group_BasicInfo, models.UserCameraGroup](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Internal Server Error - Unable to update the user camera group."
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Data = updatedDto
	c.JSON(http.StatusOK, jsonRsp)
}

// GetUserCameras godoc
// @Summary      Get all camera groups with query filter
// @Description  Retrieves a paginated list of camera groups filtered by the provided query parameters.
// @Tags         cameras
// @Accept       json
// @Produce      json
// @Param        keyword query string false "Filter by camera name keyword" minlength(1) maxlength(100)
// @Param        sort    query string false "Sort by field and order, prefix with + for asc, - for desc" default(+created_at)
// @Param        limit   query int    false "Limit the number of items per page" minimum(1) maximum(100) default(10)
// @Param        page    query int    false "Page number for pagination" minimum(1) default(1)
// @Success      200     {object}   models.JsonDTOListRsp[models.DTO_Camera_Read_BasicInfo] "A list of camera groups"
// @Failure      500     "Internal Server Error - When the query execution fails"
// @Router       /cameras/user [get]
// @Security     BearerAuth
func GetUserCameras(c *gin.Context) {
	jsonRspDTOCamerasBasicInfos := models.NewJsonDTOListRsp[models.DTO_User_Camera_Group_BasicInfo]()

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
	query := reposity.NewQuery[models.DTO_User_Camera_Group_BasicInfo, models.UserCameraGroup]()

	// Search for keyword in name
	if keyword != "" {
		query.AddConditionOfTextField("AND", "name", "LIKE", keyword)
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
