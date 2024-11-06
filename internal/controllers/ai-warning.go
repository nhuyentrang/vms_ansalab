package controllers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"vms/statuscode"
	"vms/wssignaling"

	"vms/comongo/minioclient"

	"vms/internal/models"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Param   		keyword				query	string	false	"aievent keyword"		minlength(1)  	maxlength(100)
// @Param   		sort  				query	string	false	"sort"					default(-created_at)
// @Param   		limit				query	int     false  	"limit"          		minimum(1)    	maximum(100)
// @Param   		page				query	int     false  	"page"          		minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AI_Event]
// @Router       	/aievents/searching [get]
// @Security		BearerAuth
func SearchAIEvent(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_AI_Event]()

	// Get param
	keywords := strings.Split(c.Query("keyword"), ",")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	if page < 1 {
		page = 1
	}
	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	if len(keywords) > 0 {
		for _, keyword := range keywords {
			query.AddConditionOfTextField("OR", "event_type_string", "LIKE", keyword)
			query.AddConditionOfTextField("OR", "description", "LIKE", keyword)
			query.AddConditionOfTextField("OR", "cam_name", "LIKE", keyword)
			query.AddConditionOfTextField("OR", "location", "LIKE", keyword)
			query.AddConditionOfTextField("OR", "result", "LIKE", keyword)
		}
	}
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)

	// Exec query
	dtoCabinBasics, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}

	if count == 0 {
		jsonRspDTOCabinsBasicInfos.Data = []models.DTO_AI_Event{}
	} else {
		jsonRspDTOCabinsBasicInfos.Data = dtoCabinBasics
	}
	jsonRspDTOCabinsBasicInfos.Size = int64(limit)
	jsonRspDTOCabinsBasicInfos.Page = int64(page)
	jsonRspDTOCabinsBasicInfos.Count = int64(count)

	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)
}

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Param   		location			query	string	false	"aievent location"		minlength(1)  	maxlength(100)
// @Param   		typeOfAIEvent		query	string	false	"aievent typeOfAIEvent"		minlength(1)  	maxlength(100)
// @Param   		camName			    query	string	false	"aievent cameraName"	minlength(1)  	maxlength(100)
// @Param   		startTime			query	int		false	"aievent startTime"		minimum(1)  	maximum(100)
// @Param   		endTime				query	int		false	"aievent endTime"		minimum(1)  	maximum(100)
// @Param   		sort  				query	string	false	"sort"					default(-created_at)
// @Param   		limit				query	int     false  	"limit"          		minimum(1)    	maximum(100)
// @Param   		page				query	int     false  	"page"          		minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AI_Event]
// @Router       	/aievents [get]
// @Security		BearerAuth
func GetAIEvent(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_AI_Event]()

	// Get param
	typeOfAIEvent := c.Query("typeOfAIEvent")
	// keyword := c.Query("keyword")
	location := c.Query("location")
	camName := c.Query("camName")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	if page < 1 {
		page = 1
	}
	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	if typeOfAIEvent != "" {
		query.AddConditionOfTextField("AND", "typeOfAIEvent", "=", typeOfAIEvent)
	}
	if location != "" {
		query.AddConditionOfTextField("AND", "location", "=", location)
	}
	if camName != "" {
		query.AddConditionOfTextField("AND", "cam_name", "LIKE", camName)
	}
	if startTime != "" {
		query.AddConditionOfTextField("AND", "timestamp", ">=", startTime)
	}
	if endTime != "" {
		query.AddConditionOfTextField("AND", "timestamp", "<=", endTime)
	}
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)

	// Exec query
	dtoCabinBasics, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}
	if count == 0 {
		jsonRspDTOCabinsBasicInfos.Data = []models.DTO_AI_Event{}
	} else {
		jsonRspDTOCabinsBasicInfos.Data = dtoCabinBasics
	}
	jsonRspDTOCabinsBasicInfos.Size = int64(limit)
	jsonRspDTOCabinsBasicInfos.Page = int64(page)
	jsonRspDTOCabinsBasicInfos.Count = int64(count)

	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)

}

// GetRoutine		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Param   		keyword				query	string	false	"aievent keyword"		minlength(1)  	maxlength(100)
// @Param   		eventType      		query   string  false   "event type: AI_EVENT_BLACKLIST_FACE_RECOGNITION, LICENSE_PLATE_RECOGNITION"
// @Param   		eventType1      		query   string  false   "event type 1: AI_EVENT_UNKNOWN_FACE_RECOGNITION"
// @Param   		eventType2      		query   string  false   "event type 2: AI_EVENT_PERSON_RECOGNITION"
// @Param   		startTime			query	int		false	"aievent startTime"		minimum(1)  	maximum(100)
// @Param   		endTime				query	int		false	"aievent endTime"		minimum(1)  	maximum(100)
// @Param   		sort  				query	string	false	"sort"					default(-created_at)
// @Param   		cameraName  				query	string	false	"Camera Name"
// @Param   		location  				query	string	false	"Camera Name"
// @Param   		limit				query	int     false  	"limit"          		minimum(1)    	maximum(100)
// @Param   		page				query	int     false  	"page"          		minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AIWaring]
// @Router       	/aievents/routine [get]
// @Security		BearerAuth
func AIEventRoutineHandler(c *gin.Context) {

	jsonRspDTOBasicInfos := models.NewJsonDTOListRsp[models.DTO_AIWaring]()

	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	keywords := strings.Split(c.Query("keyword"), ",")
	eventType := c.Query("eventType")
	eventType1 := c.Query("eventType1")
	eventType2 := c.Query("eventType2")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	cameraNames := strings.Split(c.Query("cameraName"), ",")

	location := c.Query("location")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}

	query := reposity.NewQuery[models.DTO_AIWaring, models.AIWaring]()
	if len(keywords) > 0 {
		for _, keyword := range keywords {
			if keyword != "" {
				query.AddConditionOfTextField("OR", "result", "LIKE", keyword)
				query.AddConditionOfTextField("AND", "deleted_mark", "=", false)
			}
		}
	}

	if len(cameraNames) > 0 {
		for _, cameraName := range cameraNames {
			if cameraName != "" {
				query.AddConditionOfTextField("OR", "cam_name", "LIKE", cameraName)
				query.AddConditionOfTextField("AND", "deleted_mark", "=", false)
			}

		}
	}

	if location != "" {
		query.AddConditionOfTextField("AND", "location", "LIKE", location)
	}
	eventTypes := []string{}
	if eventType != "" {
		eventTypes = append(eventTypes, eventType)
	}
	if eventType1 != "" {
		eventTypes = append(eventTypes, eventType1)
	}
	if eventType2 != "" {
		eventTypes = append(eventTypes, eventType2)
	}

	// Thêm điều kiện vào câu truy vấn
	if len(eventTypes) == 1 {
		// Chỉ có một giá trị, sử dụng "=" cho event_type
		query.AddConditionOfTextField("AND", "event_type", "=", eventTypes[0])
	} else if len(eventTypes) > 1 {
		// Có nhiều giá trị, sử dụng "IN" cho event_type
		query.AddConditionOfTextField("AND", "event_type", "IN", eventTypes)
	}

	// if eventType != "" {
	// 	query.AddConditionOfTextField("AND", "event_type", "=", eventType)
	// }

	// if eventType != "" {
	// 	query.AddConditionOfTextField("AND", "event_type", "=", eventType1)
	// }
	// if eventType != "" {
	// 	query.AddConditionOfTextField("AND", "event_type", "=", eventType2)
	// }

	if startTime != "" {
		query.AddConditionOfTextField("AND", "timestamp", ">=", startTime)
	}
	if endTime != "" {
		query.AddConditionOfTextField("AND", "timestamp", "<=", endTime)
	}

	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)

	// Execute the query with paging
	sort := "-timestamp"
	dtoBasics, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTOBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOBasicInfos)
		return
	}

	for key, value := range dtoBasics {
		//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
		image, err := minioclient.GetPresignedURL(value.StorageBucket, value.Image)
		if err != nil {
			// panic("Failed to GetURLFileToken, err: " + err.Error())

		}
		imageObject, err2 := minioclient.GetPresignedURL(value.StorageBucket, value.ImageObject)
		if err2 != nil {
			// panic("Failed to GetURLFileToken, err: " + err2.Error())

		}
		imageResult, err3 := minioclient.GetPresignedURL(value.StorageBucket, value.ImageResult)
		if err3 != nil {
			// panic("Failed to GetURLFileToken, err: " + err3.Error())

		}
		//fmt.Println("ImageURL: ", image, imageObject, imageResult)

		dtoBasics[key].Image = image
		dtoBasics[key].ImageObject = imageObject
		dtoBasics[key].ImageResult = imageResult
	}

	if count == 0 {
		jsonRspDTOBasicInfos.Data = []models.DTO_AIWaring{}
	} else {
		jsonRspDTOBasicInfos.Data = dtoBasics
	}
	jsonRspDTOBasicInfos.Size = int64(limit)
	jsonRspDTOBasicInfos.Page = int64(page)
	//jsonRspDTOCabinsBasicInfos.Count = int64(math.Ceil(float64(count) / float64(limit))) Total Pages
	jsonRspDTOBasicInfos.Count = int64(count) // Total Entries

	// Send response
	c.JSON(http.StatusOK, &jsonRspDTOBasicInfos)
}

// ReadCabin		 godoc
// @Summary      Get single cabin by id
// @Description  Returns the cabin whose ID value matches the id.
// @Tags         aiwarnings
// @Produce      json
// @Param        id  path  string  true  "Search aiwarnings by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_AIEvent]
// @Router       /aievents/{id} [get]
// @Security		BearerAuth
func ReadCabinEventAI(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_AIEvent]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_AIEvent, models.AIWaring](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Produce      	json
// @Param   		startTime			query	int     false  	"startTime"          			minimum(1)    	maximum(1000000000000000000000000000000)
// @Param   		endTime				query	int     false  	"endTime"          				minimum(1)    	maximum(1000000000000000000000000000000)
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AIEvent]
// @Router       	/aievents/columnchartday [get]
// @Security		BearerAuth
func ColumnChartDay(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_AIEvent_Count]()

	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	startTimes, err := strconv.ParseInt(startTime, 10, 64)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCabinsBasicInfos.Message = "Invalid startTime"
		c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
		return
	}

	endTimes, err := strconv.ParseInt(endTime, 10, 64)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCabinsBasicInfos.Message = "Invalid endTime"
		c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
		return
	}
	days := GetDaysInRange(startTimes, endTimes)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
		return
	}
	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)
	query.AddConditionOfTextField("AND", "timestamp", ">=", startTimes)
	query.AddConditionOfTextField("AND", "timestamp", "<=", endTimes)

	dtoCabinBasics, _, err := query.ExecNoPaging("+created_at")
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}
	dataByDate := make(map[string]map[string]int)

	for _, day := range days {
		for _, dto := range dtoCabinBasics {
			eventType := dto.EventType
			dateStr := day.Format("2006-01-02")
			dateStrData := dto.ConverTimestamp.Format("2006-01-02")
			if dateStr == dateStrData {
				if _, ok := dataByDate[dateStr]; !ok {
					dataByDate[dateStr] = make(map[string]int)
				}
				dataByDate[dateStr][eventType]++
			}
		}

	}

	dataAIEvent := make([]models.DTO_AIEvent_Count, 0)
	for dateStr, eventCounts := range dataByDate {
		timestamp, _ := time.Parse("2006-01-02", dateStr)
		dataAIEvent = append(dataAIEvent, models.DTO_AIEvent_Count{
			Date:      timestamp,
			SABOTAGE:  eventCounts["AI_EVENT_SABOTAGE_DETECTION"],
			DANGEROUS: eventCounts["AI_EVENT_DANGEROUS_OBJECT_DETECTION"],
			ABNORMAL:  eventCounts["AI_EVENT_ABNORMAL_BEHAVIOR_DETECTION"],
		})
	}
	jsonRspDTOCabinsBasicInfos.Data = dataAIEvent
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)

}

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Produce      	json
// @Param   		startTime			query	int     false  	"startTime"          			minimum(1)    	maximum(1000000000000000000000000000000)
// @Param   		endTime				query	int     false  	"endTime"          				minimum(1)    	maximum(1000000000000000000000000000000)
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AIEvent]
// @Router       	/aievents/columncharthour [get]
// @Security		BearerAuth
func ColumnChartHour(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_AIEvent_Count]()

	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	startTimes, err := strconv.ParseInt(startTime, 10, 64)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCabinsBasicInfos.Message = "Invalid startTime"
		c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
		return
	}

	endTimes, err := strconv.ParseInt(endTime, 10, 64)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCabinsBasicInfos.Message = "Invalid endTime"
		c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
		return
	}
	days := GetHoursInRange(startTimes, endTimes)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
		return
	}
	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)
	query.AddConditionOfTextField("AND", "timestamp", ">=", startTimes)
	query.AddConditionOfTextField("AND", "timestamp", "<=", endTimes)

	dtoCabinBasics, count, err := query.ExecNoPaging("+created_at")
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}

	fmt.Println("======> Count: ", count)
	dataByDate := make(map[string]map[string]int)

	for _, day := range days {
		for _, dto := range dtoCabinBasics {
			eventType := dto.EventType
			dateStr := day.Format("2006-01-02 15")
			dateStrData := dto.ConverTimestamp.Format("2006-01-02 15")
			if dateStr == dateStrData {
				if _, ok := dataByDate[dateStr]; !ok {
					dataByDate[dateStr] = make(map[string]int)
				}
				dataByDate[dateStr][eventType]++
			}
		}

	}

	dataAIEvent := make([]models.DTO_AIEvent_Count, 0)
	for dateStr, eventCounts := range dataByDate {
		timestamp, _ := time.Parse("2006-01-02 15", dateStr)
		dataAIEvent = append(dataAIEvent, models.DTO_AIEvent_Count{
			Date:      timestamp,
			SABOTAGE:  eventCounts["AI_EVENT_SABOTAGE_DETECTION"],
			DANGEROUS: eventCounts["AI_EVENT_DANGEROUS_OBJECT_DETECTION"],
			ABNORMAL:  eventCounts["AI_EVENT_ABNORMAL_BEHAVIOR_DETECTION"],
		})
	}
	jsonRspDTOCabinsBasicInfos.Data = dataAIEvent
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)

}

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Produce      	json
// @Param   		startTime			query	int     false  	"startTime"          			minimum(1)    	maximum(1000000000000000000000000000000)
// @Param   		endTime				query	int     false  	"endTime"          				minimum(1)    	maximum(1000000000000000000000000000000)
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AIEvent]
// @Router       	/aievents/columnchartweek [get]
// @Security		BearerAuth
func ColumnChartWeek(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_AIEvent_Count]()

	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	startTimes, err := strconv.ParseInt(startTime, 10, 64)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCabinsBasicInfos.Message = "Invalid startTime"
		c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
		return
	}

	endTimes, err := strconv.ParseInt(endTime, 10, 64)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCabinsBasicInfos.Message = "Invalid endTime"
		c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
		return
	}
	days := GetWeeksInRange(startTimes, endTimes)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
		return
	}
	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)
	query.AddConditionOfTextField("AND", "timestamp", ">=", startTimes)
	query.AddConditionOfTextField("AND", "timestamp", "<=", endTimes)

	dtoCabinBasics, count, err := query.ExecNoPaging("+created_at")
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}

	fmt.Println("======> Count: ", count)
	dataByDate := make(map[string]map[string]int)

	for _, day := range days {
		for _, dto := range dtoCabinBasics {
			eventType := dto.EventType
			dateStr := day.Format("2006-01-02")
			dateStrData := dto.ConverTimestamp.Format("2006-01-02")
			if dateStr <= dateStrData {
				if _, ok := dataByDate[dateStr]; !ok {
					dataByDate[dateStr] = make(map[string]int)
				}
				dataByDate[dateStr][eventType]++
			}
		}

	}

	dataAIEvent := make([]models.DTO_AIEvent_Count, 0)
	for dateStr, eventCounts := range dataByDate {
		timestamp, _ := time.Parse("2006-01-02", dateStr)
		dataAIEvent = append(dataAIEvent, models.DTO_AIEvent_Count{
			Date:      timestamp,
			SABOTAGE:  eventCounts["AI_EVENT_SABOTAGE_DETECTION"],
			DANGEROUS: eventCounts["AI_EVENT_DANGEROUS_OBJECT_DETECTION"],
			ABNORMAL:  eventCounts["AI_EVENT_ABNORMAL_BEHAVIOR_DETECTION"],
		})
	}
	jsonRspDTOCabinsBasicInfos.Data = dataAIEvent
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)

}

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Produce      	json
// @Param   		startTime			query	int     false  	"startTime"          			minimum(1)    	maximum(1000000000000000000000000000000)
// @Param   		endTime				query	int     false  	"endTime"          				minimum(1)    	maximum(1000000000000000000000000000000)
// @Param   		time      			query   string  true    "Time value:day ,month ,or year"  enum(day,month,year)
// @Param   		type      			query   string  true    "Type value:image ,video, all"  enum(image,video,all)
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AIEvent]
// @Router       	/aievents/imageai [get]
// @Security		BearerAuth
func ImageAI(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_AIEvent_Image]()

	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	times := c.Query("time")
	types := c.Query("type")

	var TimeAPI string = ""
	var TypeAPI string = ""

	switch types {
	case "video":
		TypeAPI = "video"
	case "image":
		TypeAPI = "image"
	case "all":
		TypeAPI = "all"
	default:
	}

	switch times {
	case "day":
		TimeAPI = "2006-01-02"
	case "month":
		TimeAPI = "2006-01"
	case "year":
		TimeAPI = "2006"
	default:
	}

	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)
	if startTime != "" {
		startTimes, err := strconv.ParseInt(startTime, 10, 64)
		if err != nil {
			jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
			jsonRspDTOCabinsBasicInfos.Message = "Invalid startTime"
			c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
			return
		}
		query.AddConditionOfTextField("AND", "timestamp", ">=", startTimes)
	}
	if endTime != "" {
		endTimes, err := strconv.ParseInt(endTime, 10, 64)
		if err != nil {
			jsonRspDTOCabinsBasicInfos.Code = http.StatusBadRequest
			jsonRspDTOCabinsBasicInfos.Message = "Invalid endTime"
			c.JSON(http.StatusBadRequest, &jsonRspDTOCabinsBasicInfos)
			return
		}
		query.AddConditionOfTextField("AND", "timestamp", "<=", endTimes)
	}
	dtoAiEventBasics, _, err := query.ExecNoPaging("+created_at")
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}

	// Map to store DTO_AIEvent_ImageItem by date
	imageItemsByDate := make(map[string][]models.DTO_AIEvent_ImageItem)

	// Loop through DTO_AIEvent list and group DTO_AIEvent_ImageItem by date
	for _, dtoAiEvent := range dtoAiEventBasics {
		eventTime := time.Unix(int64(dtoAiEvent.Timestamp)/1000, (int64(dtoAiEvent.Timestamp)%1000)*1000000)

		date := eventTime.Format(TimeAPI)
		if TypeAPI == "all" {
			imageItem := models.DTO_AIEvent_ImageItem{
				Image:         dtoAiEvent.Image,
				ImageResult:   dtoAiEvent.ImageResult,
				ImageObject:   dtoAiEvent.ImageObject,
				Video:         dtoAiEvent.Video,
				StorageBucket: dtoAiEvent.StorageBucket,
			}
			imageItemsByDate[date] = append(imageItemsByDate[date], imageItem)
		}

		if TypeAPI == "video" {
			imageItem := models.DTO_AIEvent_ImageItem{
				Video:         dtoAiEvent.Video,
				StorageBucket: dtoAiEvent.StorageBucket,
			}
			imageItemsByDate[date] = append(imageItemsByDate[date], imageItem)
		}

		if TypeAPI == "image" {
			imageItem := models.DTO_AIEvent_ImageItem{
				Image:         dtoAiEvent.Image,
				ImageResult:   dtoAiEvent.ImageResult,
				ImageObject:   dtoAiEvent.ImageObject,
				StorageBucket: dtoAiEvent.StorageBucket,
			}
			imageItemsByDate[date] = append(imageItemsByDate[date], imageItem)
		}
	}

	// Initialize slice to store DTO_AIEvent_Image
	var dtoAiEventImages []models.DTO_AIEvent_Image

	// Loop through map and create DTO_AIEvent_Image for each date
	for date, imageItems := range imageItemsByDate {
		dtoAiEventImage := models.DTO_AIEvent_Image{
			Time:        time.Time{},
			DataAIEvent: imageItems,
		}
		dtoAiEventImage.Time, _ = time.Parse(TimeAPI, date)
		dtoAiEventImages = append(dtoAiEventImages, dtoAiEventImage)
	}

	jsonRspDTOCabinsBasicInfos.Data = dtoAiEventImages
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)
}

// CreateCabin		godoc
// @Summary      	Create a new cabin
// @Description  	Takes a event JSON and store in DB. Return saved JSON.
// @Tags         	aiwarnings
// @Produce			json
// @Param        	cabin  body   models.DTO_CameraStatus  true  "CameraStatus JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_CameraStatus]
// @Router       	/aievents/camerastatus [post]
// @Security		BearerAuth
func CameraStatus(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraStatus]()
	// Call BindJSON to bind the received JSON to
	var dtoReq models.DTO_CameraStatus
	if err := c.BindJSON(&dtoReq); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	cabinData, err := reposity.ReadItemByIDIntoDTO[models.DTO_Cabin, models.Cabin](dtoReq.CabinID.String())
	if err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = SaveDevicelog(dtoReq.CameraName, "CAMERA", dtoReq.CameraStatusConnect, cabinData.CabinCode, "MQTT")
	if err != nil {
		fmt.Printf("\t\t> Failed to SaveDevicelog, error: %s\n", err)
	}

	err = SaveNewAlertReport(cabinData.ID, dtoReq.CameraID, ConvertConnectToAlertTypeString(dtoReq.CameraStatusConnect), cabinData.Name, "CAMERA", cabinData.Location, true)
	if err != nil {
		fmt.Printf("\t\t> Failed to SaveNewAlertReport, error: %s\n", err)
	}

	err = SaveNewSystemWaring(cabinData.ID, dtoReq.CameraID, ConvertConnectToAlertTypeString(dtoReq.CameraStatusConnect), cabinData.Name, "CAMREA", cabinData.Location, true)
	if err != nil {
		fmt.Printf("\t\t> Failed to SaveNewAlertReport, error: %s\n", err)
	}

	timestampMs := int64(dtoReq.EpochTime)
	seconds := timestampMs / 1000
	nanoseconds := (timestampMs % 1000) * 1000000
	eventTime := time.Unix(seconds, nanoseconds)
	report := models.DTO_Report{
		ID:         uuid.New(),
		CabinID:    cabinData.ID,
		DeviceName: "CAMERA",
		Status:     "NEW", //sensorStatus
		CabinName:  cabinData.Name,
		SensorID:   dtoReq.CameraID,
		Deleted:    false,
		AlertDate:  eventTime,
		AlertType:  ConvertConnectToAlertTypeString(dtoReq.CameraStatusConnect),
		Location:   cabinData.Location,
	}

	wssignaling.SendNotifyMessage("reports", "reports", report)
	wssignaling.SendNotifyMessage("system_waring", "system_waring", report)

	c.JSON(http.StatusOK, &jsonRsp)
}

// DeleteCabin	 godoc
// @Summary      Remove single cabin by id
// @Description  Delete a single entry from the reposity based on id.
// @Tags         aiwarnings
// @Produce      json
// @Param        id  path  string  true  "Delete cabin by id"
// @Success      200
// @Router       /aievents/delete/{id} [delete]
// @Security		BearerAuth
func DeleteAievent(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_AI_Event]()

	err := reposity.DeleteItemByID[models.AIWaring](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusNoContent, &jsonRsp)
}

// UpdateCabin		 	godoc
// @Summary      	Update single cabin by id
// @Description  	Updates and returns a single cabin whose ID value matches the id. New data must be passed in the body.
// @Tags         	aiwarnings
// @Produce      	json
// @Param        	id  path  string  true  "Update cabin by id"
// @Param        	aievent  body      models.DTO_AI_Event  true  "Cabin JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_AI_Event]
// @Router       	/aievents/update/{id} [put]
// @Security		BearerAuth
func Updateaievent(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_AI_Event]()

	// Get new data from body
	var dto models.DTO_AI_Event
	if err := c.ShouldBindJSON(&dto); err != nil {
		fmt.Println(err)
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_AI_Event, models.AIWaring](c.Param("id"), dto)
	if err != nil {
		fmt.Println(err)
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Override ID
	dto.ID, _ = uuid.Parse(c.Param("id"))

	// Return
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Produce      	json
// @Param   		startTime			query	int     false  	"startTime"          			minimum(1)    	maximum(1000000000000000000000000000000)
// @Param   		endTime				query	int     false  	"endTime"          				minimum(1)    	maximum(1000000000000000000000000000000)
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AI_Event]
// @Router       	/aievents/columnchart/week [get]
// @Security		BearerAuth
func ColumnChartWeeks(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_AIEvent_Count]()

	// startTime := c.Query("startTime")
	// endTime := c.Query("endTime")

	currentTime := time.Now().Unix()
	startTime := currentTime - (7 * 24 * 60 * 60)
	endTime := currentTime

	days := GetDaysInRange(startTime, endTime)

	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)
	query.AddConditionOfTextField("AND", "timestamp", ">=", startTime)
	query.AddConditionOfTextField("AND", "timestamp", "<=", endTime)

	dtoCabinBasics, _, err := query.ExecNoPaging("+created_at")
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}
	dataByDate := make(map[string]map[string]int)

	for _, day := range days {
		for _, dto := range dtoCabinBasics {
			eventType := dto.EventType
			dateStr := day.Format("2006-01-02")
			dateStrData := dto.ConverTimestamp.Format("2006-01-02")
			if dateStr == dateStrData {
				if _, ok := dataByDate[dateStr]; !ok {
					dataByDate[dateStr] = make(map[string]int)
				}
				dataByDate[dateStr][eventType]++
			}
		}

	}

	dataAIEvent := make([]models.DTO_AIEvent_Count, 0)
	for dateStr, eventCounts := range dataByDate {
		timestamp, _ := time.Parse("2006-01-02", dateStr)
		dataAIEvent = append(dataAIEvent, models.DTO_AIEvent_Count{
			Date:      timestamp,
			SABOTAGE:  eventCounts["AI_EVENT_SABOTAGE_DETECTION"],
			DANGEROUS: eventCounts["AI_EVENT_DANGEROUS_OBJECT_DETECTION"],
			ABNORMAL:  eventCounts["AI_EVENT_ABNORMAL_BEHAVIOR_DETECTION"],
		})
	}
	jsonRspDTOCabinsBasicInfos.Data = dataAIEvent
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)
}

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AI_Event]
// @Router       	/aievents/columnchart/day [get]
// @Security		BearerAuth
func ColumnChartDays(c *gin.Context) {
	// Khởi tạo biến để lưu trữ kết quả
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_AIEvent_Count]()

	// currentTime := time.Now().Unix()
	endTime := 1705377105000
	startTime := 1705117905000

	hours := Get4HoursIntervalsInRangeStamp(1705117905000, 1705377105000)
	fmt.Println("======> hours: ", hours)

	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)
	query.AddConditionOfTextField("AND", "timestamp", ">=", startTime)
	query.AddConditionOfTextField("AND", "timestamp", "<=", endTime)

	dtoCabinBasics, count, err := query.ExecNoPaging("+created_at")
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}
	fmt.Println("======> count: ", count)

	dataByDate := make(map[string]map[string]int)

	for _, day := range hours {
		for _, dto := range dtoCabinBasics {
			time := time.Unix(day/1000, 0)
			dateStr := time.Format("2006-01-02 15")
			previous4Hours := day - (4 * 60 * 60 * 1000)

			if _, ok := dataByDate[dateStr]; !ok {
				dataByDate[dateStr] = make(map[string]int)
			}

			if previous4Hours <= dto.Timestamp && dto.Timestamp <= day {
				dataByDate[dateStr][dto.EventType]++
			}
		}
	}
	fmt.Println(len(dataByDate))

	dataAIEvent := make([]models.DTO_AIEvent_Count, 0)
	for dateStr, eventCounts := range dataByDate {
		timestamp, _ := time.Parse("2006-01-02 15", dateStr)
		dataAIEvent = append(dataAIEvent, models.DTO_AIEvent_Count{
			Date:      timestamp,
			SABOTAGE:  eventCounts["AI_EVENT_SABOTAGE_DETECTION"],
			DANGEROUS: eventCounts["AI_EVENT_DANGEROUS_OBJECT_DETECTION"],
			ABNORMAL:  eventCounts["AI_EVENT_ABNORMAL_BEHAVIOR_DETECTION"],
		})
	}
	jsonRspDTOCabinsBasicInfos.Data = dataAIEvent
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)
}

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Produce      	json
// @Param   		startTime			query	int     false  	"startTime"          			minimum(1)    	maximum(1000000000000000000000000000000)
// @Param   		endTime				query	int     false  	"endTime"          				minimum(1)    	maximum(1000000000000000000000000000000)
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AI_Event]
// @Router       	/aievents/columnchart/hour [get]
// @Security		BearerAuth
func ColumnChartHours(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_AIEvent_Count]()

	// currentTime := time.Now().Unix()
	// endTime := currentTime
	// startTime := endTime - (24 * 60 * 60)

	endTime := 1705230000000
	startTime := 1705226400000

	minutes := Get10MinuteIntervalsInRangeStamp(1705226400000, 1705230000000)

	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)
	query.AddConditionOfTextField("AND", "timestamp", ">=", startTime)
	query.AddConditionOfTextField("AND", "timestamp", "<=", endTime)

	dtoCabinBasics, _, err := query.ExecNoPaging("+created_at")
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}
	dataByDate := make(map[string]map[string]int)

	for _, day := range minutes {
		date := time.Unix(day/1000, 0)
		fmt.Println("===> time", date)
		layout := "2006-01-02 15:04"

		dateStr := date.Format(layout)
		previous4Hours := day - (10 * 60 * 1000)
		fmt.Println("===> dateStr: ", dateStr)
		if _, ok := dataByDate[dateStr]; !ok {
			dataByDate[dateStr] = make(map[string]int)
		}
		for _, dto := range dtoCabinBasics {

			if previous4Hours <= dto.Timestamp && dto.Timestamp <= day {
				dataByDate[dateStr][dto.EventType]++
			}
		}
	}

	dataAIEvent := make([]models.DTO_AIEvent_Count, 0)
	for dateStr, eventCounts := range dataByDate {
		layout := "2006-01-02 15:04"
		timestamp, _ := time.Parse(layout, dateStr)
		fmt.Println("====>: ", timestamp, dateStr)
		dataAIEvent = append(dataAIEvent, models.DTO_AIEvent_Count{
			Date:      timestamp,
			SABOTAGE:  eventCounts["AI_EVENT_SABOTAGE_DETECTION"],
			DANGEROUS: eventCounts["AI_EVENT_DANGEROUS_OBJECT_DETECTION"],
			ABNORMAL:  eventCounts["AI_EVENT_ABNORMAL_BEHAVIOR_DETECTION"],
		})
	}
	jsonRspDTOCabinsBasicInfos.Data = dataAIEvent
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)
}

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Produce      	json
// @Param   		type			query	string	false	"type aievent location, aievent"		minlength(1)  	maxlength(100)
// @Param   		format			query	string	false	"date format day, hour, week"		minlength(1)  	maxlength(100)
// @Success      	200  {object}   models.JsonDTOListRsp[models.Top5Event]
// @Router       	/aievents/columnchart/top5 [get]
// @Security		BearerAuth
func Top5(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.Top5Event]()

	typeAI := c.Query("type")
	format := c.Query("keyword")

	startTime := int64(0)
	endTime := time.Now().UnixNano()
	// endTime := 1705377105000
	// startTime := 1705117905000

	if format == "hour" {
		currentTimeInSeconds := endTime / int64(time.Second)
		startTime = (currentTimeInSeconds / 3600) * 3600
		startTime += 3600
	}

	if format == "day" {
		currentTimeInSeconds := endTime / int64(time.Second)
		startTime = (currentTimeInSeconds / 86400) * 86400
		startTime += 86400
	}

	if format == "week" {
		currentTimeInSeconds := endTime / int64(time.Second)
		startTime = (currentTimeInSeconds / 86400) * 86400
	}

	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)
	query.AddConditionOfTextField("AND", "timestamp", ">=", startTime)
	query.AddConditionOfTextField("AND", "timestamp", "<=", endTime)

	dtoCabinBasics, _, err := query.ExecNoPaging("+created_at")
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}

	var top5Events []models.Top5Event
	dataByDate := make(map[string]int)

	for _, data := range dtoCabinBasics {
		fmt.Println("======> data.Location", data.Location)
		if typeAI == "aievent" {
			dataByDate[data.EventType]++
		}

		if typeAI == "location" {
			dataByDate[data.Location]++
		}
	}
	for event, count := range dataByDate {
		top5Events = append(top5Events, models.Top5Event{Name: event, Count: strconv.Itoa(count)})
	}

	sort.Slice(top5Events, func(i, j int) bool {
		countI, _ := strconv.Atoi(top5Events[i].Count)
		countJ, _ := strconv.Atoi(top5Events[j].Count)
		return countI > countJ
	})

	if len(top5Events) > 5 {
		top5Events = top5Events[:5]
	}

	jsonRspDTOCabinsBasicInfos.Data = top5Events
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)
}

// GetCabins		godoc
// @Summary      	Get all aievent groups with query filter
// @Description  	Responds with the list of all aievent as JSON.
// @Tags         	aiwarnings
// @Param   		location			query	string	false	"aievent location"		minlength(1)  	maxlength(100)
// @Param   		eventType		query	string	false	"aievent eventType"		minlength(1)  	maxlength(100)
// @Param   		camName			    query	string	false	"aievent cameraName"	minlength(1)  	maxlength(100)
// @Param   		startTime			query	int		false	"aievent startTime"		minimum(1)  	maximum(100)
// @Param   		endTime				query	int		false	"aievent endTime"		minimum(1)  	maximum(100)
// @Param   		sort  				query	string	false	"sort"					default(-created_at)
// @Param   		limit				query	int     false  	"limit"          		minimum(1)    	maximum(100)
// @Param   		page				query	int     false  	"page"          		minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_AI_Event]
// @Router       	/aievents/genimage [get]
// @Security		BearerAuth
func GetAIEventGenImage(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_AI_Event]()

	// Get param
	typeOfAIEvent := c.Query("eventType")
	location := c.Query("location")
	camName := c.Query("camName")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 25
	}
	if page < 1 {
		page = 1
	}
	query := reposity.NewQuery[models.DTO_AI_Event, models.AIWaring]()
	if typeOfAIEvent != "" {
		query.AddConditionOfTextField("AND", "event_type", "=", typeOfAIEvent)
	}
	if location != "" {
		query.AddConditionOfTextField("AND", "location", "=", location)
	}
	if camName != "" {
		query.AddConditionOfTextField("AND", "cam_name", "LIKE", camName)
	}
	if startTime != "" {
		query.AddConditionOfTextField("AND", "timestamp", ">=", startTime)
	}
	if endTime != "" {
		query.AddConditionOfTextField("AND", "timestamp", "<=", endTime)
	}
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)

	// Exec query
	dtoCabinBasics, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}

	for key, value := range dtoCabinBasics {
		//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
		image, err := minioclient.GetPresignedURL(value.StorageBucket, value.Image)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		imageObject, err := minioclient.GetPresignedURL(value.StorageBucket, value.ImageObject)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		imageResult, err := minioclient.GetPresignedURL(value.StorageBucket, value.ImageResult)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		//fmt.Println("ImageURL: ", image, imageObject, imageResult)

		dtoCabinBasics[key].Image = image
		dtoCabinBasics[key].ImageObject = imageObject
		dtoCabinBasics[key].ImageResult = imageResult
	}

	if count == 0 {
		jsonRspDTOCabinsBasicInfos.Data = []models.DTO_AI_Event{}
	} else {
		jsonRspDTOCabinsBasicInfos.Data = dtoCabinBasics
	}
	jsonRspDTOCabinsBasicInfos.Size = int64(limit)
	jsonRspDTOCabinsBasicInfos.Page = int64(page)
	jsonRspDTOCabinsBasicInfos.Count = int64(count)
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)

}

// CreateAIWarning		godoc
// @Summary      	Create a new aiEvent
// @Description  	Takes a aiEvent JSON and store in DB. Return saved JSON.
// @Tags         	ai-events
// @Produce			json
// @Param        	aiEvent  body   models.DTO_AI_Event  true  "AIEvent JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_AI_Event]
// @Router       	/aievents [post]
// @Security		BearerAuth
func CreateAIWarning(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_AI_Event]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_AI_Event
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = statuscode.StatusBindingInputJsonFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Create
	dto, err := reposity.CreateItemFromDTO[models.DTO_AI_Event, models.AIWaring](dto)
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

// ReadAIWarning		 godoc
// @Summary      Get single aiEvent by id
// @Description  Returns the aiEvent whose ID value matches the id.
// @Tags         ai-events
// @Produce      json
// @Param        id  path  string  true  "Read aiEvent by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_AI_Event]
// @Router       /aievents/{id} [get]
// @Security		BearerAuth
func ReadAIWarning(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_AI_Event]()
	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_AI_Event, models.AIEvent](c.Param("id"))
	if err != nil {
		jsonRsp.Code = statuscode.StatusUpdateItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateAIWarning		 	godoc
// @Summary      	Update single aiEvent by id
// @Description  	Updates and returns a single aiEvent whose ID value matches the id. New data must be passed in the body.
// @Tags         	ai-events
// @Produce      	json
// @Param        	id  path  string  true  "Update aiEvent by id"
// @Param        	aiEvent  body      models.DTO_AIWaring_Update  true  "AIEvent JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_AIWaring_Update]
// @Router       	/aievents/{id} [put]
// @Security		BearerAuth
func UpdateAIWarning(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_AIWaring_Update]()

	// Get new data from body
	var dto models.DTO_AIWaring_Update
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Update entity from DTO
	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_AIWaring_Update, models.AIWaring](c.Param("id"), dto)
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

// DeleteAIWarning	 godoc
// @Summary      Remove single aiEvent by id
// @Description  Delete a single entry from the reposity based on id.
// @Tags         ai-events
// @Produce      json
// @Param        id  path  string  true  "Delete aiEvent by id"
// @Success      204
// @Router       /aievents/{id} [delete]
// @Security		BearerAuth
func DeleteAIWarning(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_AI_Event]()

	// Soft delete aievent
	err := reposity.UpdateSingleColumn[models.AIWaring](c.Param("id"), "deleted_mark", true)
	if err != nil {
		jsonRsp.Code = statuscode.StatusDeleteItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	/*
		// Hard delete aievent
			err := reposity.DeleteItemByID[models.AIWaring](c.Param("id"))
			if err != nil {
				jsonRsp.Code = statuscode.StatusDeleteItemFailed
				jsonRsp.Message = err.Error()
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}
	*/
	c.JSON(http.StatusOK, &jsonRsp)
}

// // GetAIDeviceTypes		godoc
// // @Summary      	Get types of ai device
// // @Description  	Responds with the list of all ai device as JSON.
// // @Tags         	ai-events
// // @Produce      	json
// // @Success      	200  {object}  models.JsonAIDeviceTypeRsp
// // @Router       	/ai-events/options/device-types [get]
// // @Security		BearerAuth
// func GetAIDeviceTypes(c *gin.Context) {

// 	var jsonAIDeviceTypeRsp models.JsonAIDeviceTypeRsp

// 	jsonAIDeviceTypeRsp.Data = make([]models.KeyValue, 0)

// 	jsonAIDeviceTypeRsp.Data = append(jsonAIDeviceTypeRsp.Data, models.KeyValue{
// 		ID:   "ipCamera",
// 		Name: "IP Camera",
// 	})

// 	jsonAIDeviceTypeRsp.Data = append(jsonAIDeviceTypeRsp.Data, models.KeyValue{
// 		ID:   "smartCamera",
// 		Name: "Smart Camera",
// 	})

// 	jsonAIDeviceTypeRsp.Data = append(jsonAIDeviceTypeRsp.Data, models.KeyValue{
// 		ID:   "nvr",
// 		Name: "NVR",
// 	})

// 	jsonAIDeviceTypeRsp.Data = append(jsonAIDeviceTypeRsp.Data, models.KeyValue{
// 		ID:   "smartNVR",
// 		Name: "Smart NVR",
// 	})

// 	jsonAIDeviceTypeRsp.Code = 0
// 	jsonAIDeviceTypeRsp.Page = 1
// 	jsonAIDeviceTypeRsp.Count = len(jsonAIDeviceTypeRsp.Data)
// 	jsonAIDeviceTypeRsp.Size = jsonAIDeviceTypeRsp.Count
// 	c.JSON(http.StatusOK, &jsonAIDeviceTypeRsp)
// }

// // GetAIEventTypes		godoc
// // @Summary      	Get types of ai event
// // @Description  	Responds with the list of all ai event type as JSON.
// // @Tags         	ai-events
// // @Produce      	json
// // @Success      	200  {object}  models.JsonAIEventTypeRsp
// // @Router       	/ai-events/options/types [get]
// // @Security		BearerAuth
// func GetAIEventTypes(c *gin.Context) {

// 	var jsonAIEventTypeRsp models.JsonAIEventTypeRsp

// 	jsonAIEventTypeRsp.Data = make([]models.KeyValue, 0)

// 	jsonAIEventTypeRsp.Data = append(jsonAIEventTypeRsp.Data, models.KeyValue{
// 		ID:   "faceDetection",
// 		Name: "Phát hiện mặt",
// 	})

// 	jsonAIEventTypeRsp.Data = append(jsonAIEventTypeRsp.Data, models.KeyValue{
// 		ID:   "instructionDetection",
// 		Name: "Phát hiện vượt rào",
// 	})

// 	jsonAIEventTypeRsp.Code = 0
// 	jsonAIEventTypeRsp.Page = 1
// 	jsonAIEventTypeRsp.Count = len(jsonAIEventTypeRsp.Data)
// 	jsonAIEventTypeRsp.Size = jsonAIEventTypeRsp.Count
// 	c.JSON(http.StatusOK, &jsonAIEventTypeRsp)
// }

// // GetAIEventStatusTypes		godoc
// // @Summary      	Get types of ai event status
// // @Description  	Responds with the list of all ai event status type as JSON.
// // @Tags         	ai-events
// // @Produce      	json
// // @Success      	200  {object}  models.JsonAIEventStatusTypeRsp
// // @Router       	/ai-events/options/status-types [get]
// // @Security		BearerAuth
// func GetAIEventStatusTypes(c *gin.Context) {

// 	var jsonAIEventStatusTypeRsp models.JsonAIEventStatusTypeRsp

// 	jsonAIEventStatusTypeRsp.Data = make([]models.KeyValue, 0)

// 	jsonAIEventStatusTypeRsp.Data = append(jsonAIEventStatusTypeRsp.Data, models.KeyValue{
// 		ID:   "new",
// 		Name: "Mới",
// 	})

// 	jsonAIEventStatusTypeRsp.Data = append(jsonAIEventStatusTypeRsp.Data, models.KeyValue{
// 		ID:   "inprogress",
// 		Name: "Đang xử lý",
// 	})

// 	jsonAIEventStatusTypeRsp.Data = append(jsonAIEventStatusTypeRsp.Data, models.KeyValue{
// 		ID:   "resolved",
// 		Name: "Đã xử lý",
// 	})

// 	jsonAIEventStatusTypeRsp.Data = append(jsonAIEventStatusTypeRsp.Data, models.KeyValue{
// 		ID:   "closed",
// 		Name: "Đóng",
// 	})

// 	jsonAIEventStatusTypeRsp.Code = 0
// 	jsonAIEventStatusTypeRsp.Page = 1
// 	jsonAIEventStatusTypeRsp.Count = len(jsonAIEventStatusTypeRsp.Data)
// 	jsonAIEventStatusTypeRsp.Size = jsonAIEventStatusTypeRsp.Count
// 	c.JSON(http.StatusOK, &jsonAIEventStatusTypeRsp)
// }

// // GetAIEventLevelTypes		godoc
// // @Summary      	Get types of ai event level
// // @Description  	Responds with the list of all ai event level type as JSON.
// // @Tags         	ai-events
// // @Produce      	json
// // @Success      	200  {object}  models.JsonAIEventLevelTypeRsp
// // @Router       	/ai-events/options/level-types [get]
// // @Security		BearerAuth
// func GetAIEventLevelTypes(c *gin.Context) {

// 	var jsonAIEventLevelTypeRsp models.JsonAIEventLevelTypeRsp

// 	jsonAIEventLevelTypeRsp.Data = make([]models.KeyValue, 0)

// 	jsonAIEventLevelTypeRsp.Data = append(jsonAIEventLevelTypeRsp.Data, models.KeyValue{
// 		ID:   "low",
// 		Name: "Thấp",
// 	})

// 	jsonAIEventLevelTypeRsp.Data = append(jsonAIEventLevelTypeRsp.Data, models.KeyValue{
// 		ID:   "mid",
// 		Name: "Trung bình",
// 	})

// 	jsonAIEventLevelTypeRsp.Data = append(jsonAIEventLevelTypeRsp.Data, models.KeyValue{
// 		ID:   "high",
// 		Name: "Cao",
// 	})

// 	jsonAIEventLevelTypeRsp.Code = 0
// 	jsonAIEventLevelTypeRsp.Page = 1
// 	jsonAIEventLevelTypeRsp.Count = len(jsonAIEventLevelTypeRsp.Data)
// 	jsonAIEventLevelTypeRsp.Size = jsonAIEventLevelTypeRsp.Count
// 	c.JSON(http.StatusOK, &jsonAIEventLevelTypeRsp)
// }
