package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"strconv"

	"vms/internal/models"
	services "vms/internal/services/ws"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateUserCameraAIProperty godoc
// @Summary      Create AI properties for a user camera
// @Description  Accepts a JSON payload representing the basic info for AI properties of a user camera and stores it in the database.
// @Tags         cameras
// @Accept       json
// @Produce      json
// @Param        body  body   models.DTO_Camera_AI_Property_BasicInfo  true  "Camera AI Property Basic Info JSON"
// @Success      201  {object}  models.JsonDTORsp[models.DTO_Camera_AI_Property_BasicInfo]  "Camera AI properties created successfully."
// @Failure      400  "Bad Request - Invalid JSON payload."
// @Failure      500  "Internal Server Error - Unable to create camera AI properties."
// @Router       /cameras/user/ai-properties [post]
// @Security     BearerAuth
func CreateUserCameraAIProperty(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_Camera_AI_Property_BasicInfo]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_Camera_AI_Property_BasicInfo
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Send message to aieagent over websocket
	err := services.WSNotifyAIEngineUpdateCameraAIProperty(dto.CameraID, dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Get the ID of the camera model AI from request
	var modelAIID string = ""
	for _, aiCamProp := range dto.CameraAIPropertyList {
		if aiCamProp.CameraModelAI.ID != uuid.Nil {
			modelAIID = aiCamProp.CameraModelAI.ID.String()
			break
		}
	}

	if modelAIID == "" {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Bad Request - Camera Model AI ID is required."
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dtoCameraAI, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraModelAI, models.CameraModelAI]("id = ?", modelAIID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Todo: fix this code, currently it increase count every time toggle button is clicked
	dtoCameraAI.Count++
	reposity.UpdateItemByIDFromDTO[models.DTO_CameraModelAI, models.CameraModelAI](dtoCameraAI.ID.String(), dtoCameraAI)

	// Check if the camera AI property for this camera is already exists
	cameraID := dto.CameraID.String()

	dtoAICamProps, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty]("camera_id = ?", cameraID)
	if err != nil {
		// Create new camera AI property for this camera id
		dtoAICamProps, err = reposity.CreateItemFromDTO[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty](dto)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

	} else {
		// Update existing camera AI property for this camera id
		dtoAICamProps, err = reposity.UpdateItemByIDFromDTO[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty](dtoAICamProps.ID.String(), dto)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

	}

	jsonRsp.Data = dtoAICamProps
	c.JSON(http.StatusCreated, &jsonRsp)

}

// ReadUserCameraAIProperty godoc
// @Summary      Retrieve AI property information for a user camera by ID
// @Description  Fetches AI property information for a user camera using its unique identifier and returns it.
// @Tags         cameras
// @Produce      json
// @Param        id  path  string  true  "Camera AI Property ID"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_Camera_AI_Property_BasicInfo]  "Camera AI property information found and returned successfully."
// @Failure      404  "Not Found - Camera AI property does not exist or could not be found."
// @Router       /cameras/user/ai-properties/{id} [get]
// @Security     BearerAuth
func ReadUserCameraAIProperty(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_Camera_AI_Property_BasicInfo]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// ReadUserCameraAIPropertybyCameraID godoc
// @Summary      Retrieve AI property information for a user camera by CameraID
// @Description  Fetches AI property information for a user camera using its unique identifier and returns it.
// @Tags         cameras
// @Produce      json
// @Param        camera_id  path  string  true  "Camera AI Property ID"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_Camera_AI_Property_BasicInfo]  "Camera AI property information found and returned successfully."
// @Failure      404  "Not Found - Camera AI property does not exist or could not be found."
// @Router       /cameras/user/ai-properties/camera/{camera_id} [get]
// @Security     BearerAuth
func ReadUserCameraAIPropertybyCameraID(c *gin.Context) {

	jsonRspDTOCamerasBasicInfos := models.NewJsonDTOListRsp[models.DTO_Camera_AI_Property_BasicInfo]()

	// Get param
	keyword := c.Param("camera_id")
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}

	// Build query
	query := reposity.NewQuery[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty]()

	// Search for keyword in name
	if keyword != "" {
		query.AddConditionOfTextField("AND", "camera_id", "=", keyword)
	}

	// Exec query
	dtoCameraBasics, count, err := query.ExecWithPaging("+created_at", 10, 1)
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

// UpdateUserCameraAIProperty godoc
// @Summary      Update AI Property for a user's camera
// @Description  Receives a JSON payload with the updated data for the AI Property and applies the changes.
// @Tags         cameras
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "AI Property ID to update"
// @Param        id  body   models.DTO_Camera_AI_Property_BasicInfo  true  "JSON payload containing the updated Camera AI Property details"
// @Success      200  {object}  models.JsonDTORsp[models.DTO_Camera_AI_Property_BasicInfo]  "Successfully updated the camera AI Property."
// @Failure      400  "Bad Request - JSON payload is malformed or invalid."
// @Failure      404  "Not Found - The specified AI Property does not exist."
// @Failure      500  "Internal Server Error - Unable to update the camera AI property."
// @Router       /cameras/user/ai-properties/{id} [put]
// @Security     BearerAuth
func UpdateUserCameraAIProperty(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_Camera_AI_Property_BasicInfo]()

	// Bind the received JSON to the DTO
	var dto models.DTO_Camera_AI_Property_BasicInfo
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Bad Request - JSON payload is malformed or invalid."
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Set ID
	dto.ID = uuid.MustParse(c.Param("id"))

	// Send message to aieagent over websocket
	err := services.WSNotifyAIEngineUpdateCameraAIProperty(dto.CameraID, dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Get the ID of the camera model AI from request
	var modelAIID string = ""
	for _, aiCamProp := range dto.CameraAIPropertyList {
		if aiCamProp.CameraModelAI.ID != uuid.Nil {
			modelAIID = aiCamProp.CameraModelAI.ID.String()
			break
		}
	}

	if modelAIID == "" {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Bad Request - Camera Model AI ID is required."
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dtoCameraAI, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraModelAI, models.CameraModelAI]("id = ?", modelAIID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	dtoCameraAI.Count++

	reposity.UpdateItemByIDFromDTO[models.DTO_CameraModelAI, models.CameraModelAI](dtoCameraAI.ID.String(), dtoCameraAI)

	// Attempt to update the user camera group using the provided DTO
	updatedDto, err := reposity.UpdateItemByIDFromDTO[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = "Internal Server Error - Unable to update the camera AI Property."
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	jsonRsp.Data = updatedDto
	c.JSON(http.StatusOK, jsonRsp)
}

// UpdateUserCamerasAIProperties godoc
// @Summary      Update AI Properties for multiple user cameras
// @Description  Receives a JSON payload with the updated data for the AI Properties and applies the changes to multiple cameras.
// @Tags         cameras
// @Accept       json
// @Produce      json
// @Param        body  body   []models.DTO_Camera_AI_Property_BasicInfo  true  "JSON payload containing the updated Camera AI Properties details"
// @Success      200  {object}  models.JsonDTOListRsp[models.DTO_Camera_AI_Property_BasicInfo]  "Successfully updated the camera AI Properties."
// @Failure      400  "Bad Request - JSON payload is malformed or invalid."
// @Failure      404  "Not Found - One or more specified AI Properties do not exist."
// @Failure      500  "Internal Server Error - Unable to update the camera AI properties."
// @Router       /cameras/user/ai-properties/cameras [put]
// @Security     BearerAuth
func UpdateUserCamerasAIProperties(c *gin.Context) {
	jsonRsp := models.NewJsonDTOListRsp[models.DTO_Camera_AI_Property_BasicInfo]()

	// Read the raw JSON payload
	var rawPayload []byte
	var err error
	if rawPayload, err = io.ReadAll(c.Request.Body); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Bad Request - Unable to read JSON payload."
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Unmarshal the raw payload to validate and extract camera_id
	var rawDtos []map[string]interface{}
	if err := json.Unmarshal(rawPayload, &rawDtos); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Bad Request - JSON payload is malformed or invalid."
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	var dtos []models.DTO_Camera_AI_Property_BasicInfo
	for _, rawDto := range rawDtos {
		var cameraID uuid.UUID
		cameraIDStr, ok := rawDto["camera_id"].(string)
		if !ok || cameraIDStr == "" {
			fmt.Println("camera_id entry doesn't exist or is empty")
			cameraID = uuid.Nil
		} else {
			cameraID = uuid.MustParse(cameraIDStr)
		}

		var parsedID uuid.UUID
		idStr, ok := rawDto["id"].(string)
		if !ok || idStr == "" {
			fmt.Println("id entry doesn't exist or is empty")
			parsedID = uuid.Nil
		} else {
			parsedID = uuid.MustParse(idStr)
		}

		dto := models.DTO_Camera_AI_Property_BasicInfo{
			CameraID: cameraID,
			ID:       parsedID,
		}

		// Check if cameraaiproperty is empty, create new properties if necessary
		cameraAIPropertyList, ok := rawDto["cameraaiproperty"].([]interface{})
		if !ok || len(cameraAIPropertyList) == 0 {
			dto.CameraAIPropertyList = []models.CameraAIProperty{
				{
					CameraAIZone:  models.CameraVirtualProperty{Vzone: models.VZoneCoordinate{Point_A: models.Point{X: 0, Y: 0}, Point_B: models.Point{X: 0, Y: 0}, Point_C: models.Point{X: 0, Y: 0}, Point_D: models.Point{X: 0, Y: 0}}},
					CameraModelAI: models.CameraModelAI{ID: uuid.MustParse("56f49f4d-7850-48c9-87b7-135a1489ef10"), Type: "face_recognition", ModelName: "Nhận diện khuôn mặt"},
					CalendarDays:  []models.CalendarDays{},
					IsActive:      false,
				},
				{
					CameraAIZone:  models.CameraVirtualProperty{Vzone: models.VZoneCoordinate{Point_A: models.Point{X: 0, Y: 0}, Point_B: models.Point{X: 0, Y: 0}, Point_C: models.Point{X: 0, Y: 0}, Point_D: models.Point{X: 0, Y: 0}}},
					CameraModelAI: models.CameraModelAI{ID: uuid.MustParse("56f49f4d-7850-48c9-87b7-135a1489ef11"), Type: "license_plate_recognition", ModelName: "Nhận diện biển số"},
					CalendarDays:  []models.CalendarDays{},
					IsActive:      false,
				},
			}
		} else {
			// Handle the case where cameraaiproperty is not empty
			var properties []models.CameraAIProperty
			for _, prop := range cameraAIPropertyList {
				var property models.CameraAIProperty
				propBytes, _ := json.Marshal(prop)
				if err := json.Unmarshal(propBytes, &property); err != nil {
					continue // Skip this property if it is malformed
				}
				properties = append(properties, property)
			}

			dto.CameraAIPropertyList = properties

			// Check if both properties exist and create the missing one
			existsFaceRecognition := false
			existsLicensePlateRecognition := false
			for _, prop := range dto.CameraAIPropertyList {
				if prop.CameraModelAI.Type == "face_recognition" {
					existsFaceRecognition = true
				}
				if prop.CameraModelAI.Type == "license_plate_recognition" {
					existsLicensePlateRecognition = true
				}
			}
			if !existsFaceRecognition {
				dto.CameraAIPropertyList = append(dto.CameraAIPropertyList, models.CameraAIProperty{
					CameraAIZone:  models.CameraVirtualProperty{Vzone: models.VZoneCoordinate{Point_A: models.Point{X: 0, Y: 0}, Point_B: models.Point{X: 0, Y: 0}, Point_C: models.Point{X: 0, Y: 0}, Point_D: models.Point{X: 0, Y: 0}}},
					CameraModelAI: models.CameraModelAI{ID: uuid.MustParse("56f49f4d-7850-48c9-87b7-135a1489ef10"), Type: "face_recognition", ModelName: "Nhận diện khuôn mặt"},
					CalendarDays:  []models.CalendarDays{},
					IsActive:      false,
				})
			}
			if !existsLicensePlateRecognition {
				dto.CameraAIPropertyList = append(dto.CameraAIPropertyList, models.CameraAIProperty{
					CameraAIZone:  models.CameraVirtualProperty{Vzone: models.VZoneCoordinate{Point_A: models.Point{X: 0, Y: 0}, Point_B: models.Point{X: 0, Y: 0}, Point_C: models.Point{X: 0, Y: 0}, Point_D: models.Point{X: 0, Y: 0}}},
					CameraModelAI: models.CameraModelAI{ID: uuid.MustParse("56f49f4d-7850-48c9-87b7-135a1489ef11"), Type: "license_plate_recognition", ModelName: "Nhận diện biển số"},
					CalendarDays:  []models.CalendarDays{},
					IsActive:      false,
				})
			}
		}

		dtos = append(dtos, dto)
	}

	var updatedDtos []models.DTO_Camera_AI_Property_BasicInfo
	for _, dto := range dtos {
		// Check if the entry exists in the database
		_, err := reposity.ReadItemByIDIntoDTO[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty](dto.ID.String())
		if err != nil {
			// If it doesn't exist, create a new entry
			dto, err = reposity.CreateItemFromDTO[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty](dto)
			if err != nil {
				jsonRsp.Code = http.StatusInternalServerError
				jsonRsp.Message = "Internal Server Error - Unable to create the camera AI Property."
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}
		}

		// Send message to aieagent over websocket
		err = services.WSNotifyAIEngineUpdateCameraAIProperty(dto.CameraID, dto)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

		// If send message success, update the camera AI property
		lastElement := dto.CameraAIPropertyList[len(dto.CameraAIPropertyList)-1]
		dtoCameraAI, errAI := reposity.ReadItemWithFilterIntoDTO[models.DTO_CameraModelAI, models.CameraModelAI]("id = ?", lastElement.CameraModelAI.ID)
		if errAI != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = errAI.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		if dtoCameraAI.Count == 0 {
			dtoCameraAI.Count = 1
		} else {
			dtoCameraAI.Count++
		}
		reposity.UpdateItemByIDFromDTO[models.DTO_CameraModelAI, models.CameraModelAI](dtoCameraAI.ID.String(), dtoCameraAI)

		// Update the camera AI property
		updatedDto, err := reposity.UpdateItemByIDFromDTO[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty](dto.ID.String(), dto)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = "Internal Server Error - Unable to update the camera AI Property."
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
		updatedDtos = append(updatedDtos, updatedDto)
	}

	jsonRsp.Data = updatedDtos
	c.JSON(http.StatusOK, jsonRsp)
}

// GetUserCameraAIProperty godoc
// @Summary      Get all camera AI properties with query filter
// @Description  Retrieves a paginated list of camera AI properties filtered by the provided query parameters. This allows for the search of AI properties based on specific criteria, such as name keywords, and ordering by fields.
// @Tags         cameras
// @Accept       json
// @Produce      json
// @Param        keyword query string false "Filter by AI property name keyword" minlength(1) maxlength(100)
// @Param        sort    query string false "Sort by field and order, prefix with + for asc, - for desc" default(+created_at)
// @Param        limit   query int    false "Limit the number of items per page" minimum(1) maximum(100) default(10)
// @Param        page    query int    false "Page number for pagination" minimum(1) default(1)
// @Success      200     {object}   models.JsonDTOListRsp[models.DTO_Camera_AI_Property_BasicInfo] "A list of camera AI properties"
// @Failure      500     "Internal Server Error - When the query execution fails"
// @Router       /cameras/user/ai-properties [get]
// @Security     BearerAuth
func GetUserCameraAIProperty(c *gin.Context) {
	jsonRspDTOCamerasBasicInfos := models.NewJsonDTOListRsp[models.DTO_Camera_AI_Property_BasicInfo]()

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
	query := reposity.NewQuery[models.DTO_Camera_AI_Property_BasicInfo, models.CameraAIEventProperty]()

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
