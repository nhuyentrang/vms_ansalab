package controllers

import (
	"net/http"

	"strconv"

	"vms/internal/models"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
)

// CreateSystemIncidentLog godoc
// @Summary      Create a new system incident log
// @Description  Accepts a JSON payload representing a system incident log, stores it in the database, and returns the saved entity.
// @Tags         incidents
// @Accept       json
// @Produce      json
// @Param        body  body   models.DTO_System_Incident_BasicInfo  true  "System Incident Log JSON"
// @Success      201  {object}  models.JsonDTORsp[models.DTO_System_Incident_BasicInfo]  "System Incident Log created successfully."
// @Failure      400  "Bad Request - Invalid JSON payload."
// @Failure      500  "Internal Server Error - Unable to create the system incident log."
// @Router       /incidents/logs [post]
// @Security     BearerAuth
func CreateSystemIncidentLog(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_System_Incident_BasicInfo]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_System_Incident_BasicInfo
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Create Incident in system
	dto, err := reposity.CreateItemFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Data = dto
	c.JSON(http.StatusCreated, &jsonRsp)
}

// ReadSystemIncidentLog godoc
// @Summary      Retrieve a system incident log
// @Description  Fetches a system incident log by its unique identifier and returns it.
// @Tags         incidents
// @Produce      json
// @Param        id  path  string  true  "System Incident Log ID"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_System_Incident_BasicInfo]  "System Incident Log found and returned successfully."
// @Failure      404  "Not Found - System Incident Log does not exist or could not be found."
// @Router       /incidents/logs/{id} [get]
// @Security     BearerAuth
func ReadSystemIncidentLog(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_System_Incident_BasicInfo]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateSystemIncidentLog godoc
// @Summary      Update a system incident log
// @Description  Receives a JSON payload with the updated data for a system incident log and applies the changes to the specified log.
// @Tags         incidents
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "System Incident Log ID to update"
// @Param        IncidentLog  body   models.DTO_System_Incident_BasicInfo  true  "Updated System Incident Log JSON"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_System_Incident_BasicInfo]  "Successfully updated the system incident log."
// @Failure      400  "Bad Request - JSON payload is malformed or invalid."
// @Failure      404  "Not Found - Specified system incident log does not exist."
// @Failure      500  "Internal Server Error - Unable to update the system incident log."
// @Router       /incidents/logs/{id} [put]
// @Security     BearerAuth
func UpdateSystemIncidentLog(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_System_Incident_BasicInfo]()

	// Bind the received JSON to the DTO
	var dto models.DTO_System_Incident_BasicInfo
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Bad Request - JSON payload is malformed or invalid."
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Attempt to update the user camera group using the provided DTO
	updatedDto, err := reposity.UpdateItemByIDFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](c.Param("id"), dto)
	if err != nil {
		// Here you may want to differentiate between not found and other server errors for accurate status codes
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Internal Server Error - Unable to update the user camera group."
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Data = updatedDto
	c.JSON(http.StatusOK, jsonRsp)
}

// GetSystemIncidentLogs godoc
// @Summary      Get all system incident logs with query filter
// @Description  Retrieves a paginated list of system incident logs filtered by the provided query parameters.
// @Tags         incidents
// @Accept       json
// @Produce      json
// @Param        event_name query string false "Filter by incident name event_name" Enums(Đã kết nối, Mất kết nối, Đã xóa camera)
// @Param   	 device_type	query	string	false	"Incident Device Type"	Enums(nvr, camera)
// @Param        device_name query string false "Filter by device name device_name" minlength(1) maxlength(100)
// @Param        type	query	string	false	"Incident status"	Enums(Active, Deactive)
// @Param        sort    query string false "Sort by field and order, prefix with + for asc, - for desc" default(+status)
// @Param        limit   query int    false "Limit the number of items per page" minimum(1) maximum(100) default(10)
// @Param        page    query int    false "Page number for pagination" minimum(1) default(1)
// @Success      200     {object}   models.JsonDTOListRsp[models.DTO_System_Incident_BasicInfo] "A list of system incident logs"
// @Failure      500     "Internal Server Error - When the query execution fails"
// @Router       /incidents/logs [get]
// @Security     BearerAuth
func GetSystemIncidentLogs(c *gin.Context) {
	jsonRspDTOSystemIncidentBasicInfos := models.NewJsonDTOListRsp[models.DTO_System_Incident_BasicInfo]()

	// Get param
	event_name := c.Query("event_name")
	device_name := c.Query("device_name")
	device_type := c.Query("device_type")
	Type := c.Query("type")
	sort := c.Query("sort")
	keyword := c.Query("keyword")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	// Build query
	query := reposity.NewQuery[models.DTO_System_Incident_BasicInfo, models.SystemIncident]()

	// Search for keyword in name
	if keyword != "" {
		query.AddTwoConditionOfTextField("AND", "event_name", "=", keyword, "OR", "device_type", "=", keyword)
		// query.AddConditionOfJsonbField("AND", "event_name", "device_type", "LIKE", keyword)
	}
	if event_name != "" {
		query.AddConditionOfTextField("AND", "event_name", "=", event_name)
	}
	if device_name != "" {
		query.AddConditionOfTextField("AND", "source", "LIKE", device_name)
	}
	if device_type != "" {
		query.AddConditionOfTextField("AND", "device_type", "=", device_type)
	}
	if Type != "" {
		query.AddConditionOfTextField("AND", "type", "=", Type)
	}

	// Exec query
	dtoCameraBasics, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTOSystemIncidentBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOSystemIncidentBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOSystemIncidentBasicInfos)
		return
	}

	jsonRspDTOSystemIncidentBasicInfos.Count = count
	jsonRspDTOSystemIncidentBasicInfos.Data = dtoCameraBasics
	jsonRspDTOSystemIncidentBasicInfos.Page = int64(page)
	jsonRspDTOSystemIncidentBasicInfos.Size = int64(len(dtoCameraBasics))
	c.JSON(http.StatusOK, &jsonRspDTOSystemIncidentBasicInfos)
}

// UpdateSystemIncidentLogStatus godoc
// @Summary      Update a system incident log
// @Description  Receives a JSON payload with the updated data for a system incident log and applies the changes to the specified log.
// @Tags         incidents
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "System Incident Log ID to update"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_System_Incident_BasicInfo]  "Successfully updated the system incident log."
// @Failure      400  "Bad Request - JSON payload is malformed or invalid."
// @Failure      404  "Not Found - Specified system incident log does not exist."
// @Failure      500  "Internal Server Error - Unable to update the system incident log."
// @Router       /incidents/logs/status/{id} [put]
// @Security     BearerAuth
func UpdateSystemIncidentLogStatus(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_System_Incident_BasicInfo]()

	// Retrieve the existing record from the database
	existingDto, err := reposity.ReadItemByIDIntoDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	// Update only the status and type fields
	existingDto.Type = "Active"
	existingDto.Status = "Đã xử lý"
	existingDto.EventName = "Đã kết nối"

	// Save the updated record back to the database
	updatedDto, err := reposity.UpdateItemByIDFromDTO[models.DTO_System_Incident_BasicInfo, models.SystemIncident](c.Param("id"), existingDto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Internal Server Error - Unable to update the item."
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Data = updatedDto
	c.JSON(http.StatusOK, jsonRsp)
}

// DeleteSystemIncident	 godoc
// @Summary      Remove single Incident by id
// @Description  Delete a single entry from the reposity based on id.
// @Tags         incidents
// @Produce      json
// @Param        id  path  string  true  "Delete System Incident by id"
// @Success      204
// @Router       /incidents/{id} [delete]
// @Security		BearerAuth
func DeleteSystemIncident(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_System_Incident_BasicInfo]()

	err := reposity.DeleteItemByID[models.SystemIncident](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	c.JSON(http.StatusOK, &jsonRsp)
}
