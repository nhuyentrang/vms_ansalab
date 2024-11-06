package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"vms/internal/models"
	"vms/statuscode"

	"vms/comongo/kafkaclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateCameraGroup		godoc
// @Summary      	Create a new cameraGroup
// @Description  	Takes a cameraGroup JSON and store in DB. Return saved JSON.
// @Tags         	camera-groups
// @Produce			json
// @Param        	cameraGroup  body   models.DTO_CameraGroup_Create  true  "CameraGroup JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_CameraGroup]
// @Router       	/camera-groups [post]
// @Security		BearerAuth
func CreateCameraGroup(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraGroup_Create]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_CameraGroup_Create
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = statuscode.StatusBindingInputJsonFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Check if there is at least one camera in the list then check for duplication
	// if len(dto.Cameras) > 0 {
	// 	jsonRspDupl := models.NewJsonDTOListRsp[models.DTO_CameraGroup_Read_BasicInfo]()
	// 	firstCameraID := dto.Cameras[0].ID
	// 	query := reposity.NewQuery[models.DTO_CameraGroup_Read_BasicInfo, models.CameraGroup]()
	// 	sort := "-created_at"
	// 	jsonCondition := fmt.Sprintf("[{\"id\":\"%s\"}]", firstCameraID)
	// 	query.AddConditionOfTextField("AND", "cameras", "@>", jsonCondition)

	// 	//Exec Query
	// 	_, count, err := query.ExecNoPaging(sort)
	// 	if err != nil {
	// 		jsonRsp.Code = statuscode.StatusSearchItemFailed
	// 		jsonRsp.Message = err.Error()
	// 		c.JSON(http.StatusInternalServerError, &jsonRsp)
	// 		return
	// 	}
	// 	jsonRspDupl.Count = count

	// 	// If duplication is found
	// 	if jsonRspDupl.Count > 0 {
	// 		jsonRsp.Code = statuscode.StatusCreateItemFailed
	// 		jsonRsp.Message = "Camera already exists in other group"
	// 		c.JSON(http.StatusInternalServerError, &jsonRsp)
	// 		return
	// 	}
	// }

	// Create new block
	dto, err := reposity.CreateItemFromDTO[models.DTO_CameraGroup_Create, models.CameraGroup](dto)
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

// ReadCameraGroup		 godoc
// @Summary      Get single cameraGroup by id
// @Description  Returns the cameraGroup whose ID value matches the id.
// @Tags         camera-groups
// @Produce      json
// @Param        id  path  string  true  "Read cameraGroup by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_CameraGroup]
// @Router       /camera-groups/{id} [get]
// @Security		BearerAuth
func ReadCameraGroup(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraGroup]()
	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_CameraGroup, models.CameraGroup](c.Param("id"))
	if err != nil {
		jsonRsp.Code = statuscode.StatusUpdateItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateCameraGroup		 	godoc
// @Summary      	Update single cameraGroup by id
// @Description  	Updates and returns a single cameraGroup whose ID value matches the id. New data must be passed in the body.
// @Tags         	camera-groups
// @Produce      	json
// @Param        	id  path  string  true  "Update cameraGroup by id"
// @Param        	cameraGroup  body      models.DTO_CameraGroup_Create  true  "CameraGroup JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_CameraGroup_Create]
// @Router       	/camera-groups/{id} [put]
// @Security		BearerAuth
func UpdateCameraGroup(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraGroup_Create]()
	jsonRspDtb := models.NewJsonDTORsp[models.DTO_CameraGroup_Create]()

	// Receive current record from database to check duplicate
	dto1, errdtb := reposity.ReadItemByIDIntoDTO[models.DTO_CameraGroup_Create, models.CameraGroup](c.Param("id"))
	if errdtb != nil {
		jsonRspDtb.Code = statuscode.StatusUpdateItemFailed
		jsonRspDtb.Data = dto1
		jsonRspDtb.Message = errdtb.Error()
		c.JSON(http.StatusNotFound, &jsonRspDtb)
		return
	}
	jsonRspDtb.Data = dto1

	// Get new data from body
	var dto models.DTO_CameraGroup_Create
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	//#region Only add 1 camera in 1 group
	// currentIDs := extractCameraIDs(dto1)
	// pendingIDs := extractCameraIDs(dto)

	// // If user is removing camera from list
	// if len(currentIDs) < len(pendingIDs) {

	// 	// Find unique IDs in the pending record
	// 	uniqueIDs := findUniqueIDs(pendingIDs, currentIDs)

	// 	// If uniqueIDs contains IDs, those are new cameras not present in the current record
	// 	if len(uniqueIDs) > 0 {
	// 		jsonRspDupl := models.NewJsonDTOListRsp[models.DTO_CameraGroup_Read_BasicInfo]()
	// 		firstCameraID := uniqueIDs[0]
	// 		query := reposity.NewQuery[models.DTO_CameraGroup_Read_BasicInfo, models.CameraGroup]()
	// 		sort := "-created_at"
	// 		jsonCondition := fmt.Sprintf("[{\"id\":\"%s\"}]", firstCameraID)
	// 		query.AddConditionOfTextField("AND", "cameras", "@>", jsonCondition)

	// 		//Exec Query
	// 		dto2, count, err := query.ExecNoPaging(sort)
	// 		if err != nil {
	// 			jsonRspDupl.Code = statuscode.StatusUpdateItemFailed
	// 			jsonRspDupl.Message = err.Error()
	// 			c.JSON(http.StatusInternalServerError, &jsonRspDupl)
	// 			return
	// 		}
	// 		jsonRspDupl.Data = dto2
	// 		jsonRspDupl.Count = count

	// 		// If duplication is found
	// 		if jsonRspDupl.Count > 0 {
	// 			jsonRspDupl.Code = statuscode.StatusUpdateItemFailed
	// 			jsonRspDupl.Message = "Camera already exists in other group"
	// 			c.JSON(http.StatusInternalServerError, &jsonRspDupl)
	// 			return
	// 		}
	// 	}
	// }
	// #endregion
	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_CameraGroup_Create, models.CameraGroup](c.Param("id"), dto)
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

// Utility function to extract camera IDs from a DTO_CameraGroup_Create instance
func extractCameraIDs(dto models.DTO_CameraGroup_Create) []string {
	ids := make([]string, len(dto.Cameras))

	// Iterate over the Cameras slice within the DTO
	for i, camera := range dto.Cameras {
		ids[i] = camera.ID
	}

	return ids
}

// Utility function to find unique IDs in the first slice that aren't in the second
func findUniqueIDs(slice1, slice2 []string) []string {
	unique := make([]string, 0)
	set := make(map[string]struct{})
	for _, id := range slice2 {
		set[id] = struct{}{}
	}
	for _, id := range slice1 {
		if _, found := set[id]; !found {
			unique = append(unique, id)
		}
	}
	return unique
}

// DeleteCameraGroup	 godoc
// @Summary      Remove single cameraGroup by id
// @Description  Delete a single entry from the reposity based on id.
// @Tags         camera-groups
// @Produce      json
// @Param        id  path  string  true  "Delete cameraGroup by id"
// @Success      204
// @Router       /camera-groups/{id} [delete]
// @Security		BearerAuth
func DeleteCameraGroup(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraGroup]()

	err := reposity.DeleteItemByID[models.CameraGroup](c.Param("id"))
	if err != nil {
		jsonRsp.Code = statuscode.StatusDeleteItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	c.JSON(http.StatusNoContent, &jsonRsp)
}

// GetCameraGroups		godoc
// @Summary      	Get all camera groups with query filter
// @Description  	Responds with the list of all cameraGroup as JSON.
// @Tags         	camera-groups
// @Param   		keyword			query	string	false	"cameraGroup name keyword"		minlength(1)  	maxlength(100)
// @Param   		sort  			query	string	false	"sort"							default(-created_at)
// @Param   		limit			query	int     false  	"limit"          				minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          				minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}  models.JsonDTOListRsp[models.DTO_CameraGroup_Read_BasicInfo]
// @Router       	/camera-groups [get]
// @Security		BearerAuth
func GetCameraGroups(c *gin.Context) {

	jsonRsp := models.NewJsonDTOListRsp[models.DTO_CameraGroup_Read_BasicInfo]()

	// Get param
	keyword := c.Query("keyword")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))

	// Build query
	query := reposity.NewQuery[models.DTO_CameraGroup_Read_BasicInfo, models.CameraGroup]()

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

	cmd := models.DeviceCommand{
		CommandID: uuid.New().String(),
		Cmd:       "PingCamera",
		EventTime: time.Now().Format(time.RFC3339),
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)

	jsonRsp.Count = count
	jsonRsp.Data = dtos
	jsonRsp.Page = int64(page)
	jsonRsp.Size = int64(len(dtos))
	c.JSON(http.StatusOK, &jsonRsp)
}

// GetCameraFromGroups godoc
// @Summary      Get all camera groups containing a specific camera ID
// @Description  Responds with the list of all camera groups that contain the specified camera ID as JSON.
// @Tags         camera-groups
// @Param        	id  path  string  true  "Camera ID"
// @Param        sort         query   string  false "Sort"                      default(-created_at)
// @Produce      json
// @Success      200  {object}  models.JsonDTOListRsp[models.DTO_CameraGroup_Read_BasicInfo]
// @Router       /camera-groups/cameras/{id} [get]
// @Security     BearerAuth
func GetCameraFromGroups(c *gin.Context) {

	jsonRsp := models.NewJsonDTOListRsp[models.DTO_CameraGroup_Read_BasicInfo]()

	// Get param
	keyword := c.Param("id")
	sort := c.Query("sort")

	// Build query
	query := reposity.NewQuery[models.DTO_CameraGroup_Read_BasicInfo, models.CameraGroup]()

	if keyword != "" {
		jsonCondition := fmt.Sprintf("[{\"id\":\"%s\"}]", keyword)
		query.AddConditionOfTextField("AND", "cameras", "@>", jsonCondition)
	}
	fmt.Printf("query: %v\n", query)
	// Exec query
	dtos, count, err := query.ExecNoPaging(sort)
	if err != nil {
		jsonRsp.Code = statuscode.StatusSearchItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Count = count
	jsonRsp.Data = dtos
	jsonRsp.Size = int64(len(dtos))
	c.JSON(http.StatusOK, &jsonRsp)
}
