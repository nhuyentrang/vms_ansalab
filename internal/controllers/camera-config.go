package controllers

import (
	"net/http"

	"vms/internal/models"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateCameraConfig	godoc
// @Summary      	Get a cameraConfig by ID
// @Description  	Fetches a cameraConfig by its ID and returns it as JSON.
// @Tags         	camera-configs
// @Produce      	json
// @Param        	id   path     string  true  "CameraConfig ID"
// @Success      	200  {object} models.DTO_CameraConfig
// @Failure      	404  {object} models.DTO_CameraConfig
// @Router       	/camera-configs/{id} [post]
// @Security      BearerAuth
func CreateCameraConfig(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraConfig]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_CameraConfig, models.CameraConfig](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}
	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// UpdateCameraConfig          godoc
// @Summary       Update a single CameraConfig by ID
// @Description   Updates and returns a single CameraConfig whose ID matches the parameter id.
// @Tags          CameraConfig
// @Produce       json
// @Param         id path string true "ID of the CameraConfig to update"
// @Param         CameraConfig body models.DTO_CameraConfig true "Updated CameraConfig object"
// @Success       200 {object} models.DTO_CameraConfig
// @Router        /camera-config/{id} [put]
// @Security      BearerAuth
func UpdateCameraConfig(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraConfig]()

	// Get new data from body
	var dto models.DTO_CameraConfig
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_CameraConfig, models.CameraConfig](c.Param("id"), dto)
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

// DeleteCameraConfig          godoc
// @Summary       Delete a single CameraConfig by ID
// @Description   Deletes and returns a single CameraConfig whose ID matches the parameter id.
// @Tags          CameraConfig
// @Produce       json
// @Param         id path string true "ID of the CameraConfig to delete"
// @Success       204 "No content"
// @Router        /camera-config/{id} [delete]
// @Security      BearerAuth
func DeleteCameraConfig(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_CameraConfig]()

	// Attempt to delete the CameraConfig by ID
	err := reposity.DeleteItemByID[models.CameraConfig](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	// Return success with no content\
	c.JSON(http.StatusNoContent, &jsonRsp)
}

// SetCameraConfigAll	godoc
// @Summary      	Set configurations for all cameras
// @Description  	Takes a list of camera configurations JSON and apply them. Returns the number of configurations applied.
// @Tags         	camera-configs
// @Produce      	json
// @Param        	cameraConfigs  body   []models.DTOCameraDeviceConfig  true  "List of CameraConfig JSON"
// @Success      	201  {object}  models.JsonDTORsp[int]
// @Failure      	400  {object}  models.JsonDTORsp[int]
// @Router       	/camera-configs [post]
// @Security      BearerAuth
func SetCameraConfigAll(c *gin.Context) {

	jsonRspDTOCameraConfigList := models.NewJsonDTORsp[int]()
	/*
		// Call BindJSON to bind the received JSON to
		dtoCameraConfigs := models.NewJsonDTOListReq[models.DTOCameraDeviceConfig]()
		if err := c.BindJSON(&dtoCameraConfigs); err != nil {
			jsonRspDTOCameraConfigList.Code = http.StatusBadRequest
			jsonRspDTOCameraConfigList.Message = err.Error()
			c.JSON(http.StatusBadRequest, &jsonRspDTOCameraConfigList)
			return
		}
	*/
	// Call BindJSON to bind the received JSON to
	dtoCameraConfigs := []models.DTOCameraDeviceConfig{}
	if err := c.BindJSON(&dtoCameraConfigs); err != nil {
		jsonRspDTOCameraConfigList.Code = http.StatusBadRequest
		jsonRspDTOCameraConfigList.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRspDTOCameraConfigList)
		return
	}
	// Todo: config camera ...
	jsonRspDTOCameraConfigList.Data = 1
	c.JSON(http.StatusCreated, &jsonRspDTOCameraConfigList)
}

// UpdateCamera		godoc
// @Summary      	Config multiple camera by id
// @Description  	Config and returns cameraIDs. Configs key/values must be passed in the body.
// @Tags         	cameras
// @Produce      	json
// @Param        	camera  body 	models.DTOCameraDeviceConfigBatch true  "Camera Config Batch JSON"
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/cameras/set-config/batch [post]
// @Security		BearerAuth
func SetCameraConfigBatch(c *gin.Context) {

	jsonRspDTOCameraConfigList := models.NewJsonDTORsp[int]()
	/*
		// Call BindJSON to bind the received JSON to
		dtoCameraConfigs := models.NewJsonDTOListReq[models.DTOCameraDeviceConfig]()
		if err := c.BindJSON(&dtoCameraConfigs); err != nil {
			jsonRspDTOCameraConfigList.Code = http.StatusBadRequest
			jsonRspDTOCameraConfigList.Message = err.Error()
			c.JSON(http.StatusBadRequest, &jsonRspDTOCameraConfigList)
			return
		}
	*/
	// Call BindJSON to bind the received JSON to
	dtoCameraConfigs := models.DTOCameraDeviceConfigBatch{}
	if err := c.BindJSON(&dtoCameraConfigs); err != nil {
		jsonRspDTOCameraConfigList.Code = http.StatusBadRequest
		jsonRspDTOCameraConfigList.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRspDTOCameraConfigList)
		return
	}
	// Todo: config camera ...
	jsonRspDTOCameraConfigList.Data = len(dtoCameraConfigs.CameraIDs)
	c.JSON(http.StatusCreated, &jsonRspDTOCameraConfigList)
}

// UpdateCamera		godoc
// @Summary      	Config multiple camera by id
// @Description  	Config and returns cameraIDs. Configs key/values must be passed in the body.
// @Tags         	cameras
// @Produce      	json
// @Param        	id  path  string  true  "Config camera by id"
// @Param        	camera  body     []models.DTOCameraDeviceConfig{} true  "Camera JSON"
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/cameras/set-config/{id} [post]
// @Security		BearerAuth
func SetCameraConfig(c *gin.Context) {

	jsonRspDTOCameraConfigList := models.NewJsonDTORsp[int]()

	// Call BindJSON to bind the received JSON to
	dtoCameraConfigs := models.DTOCameraDeviceConfig{}
	if err := c.BindJSON(&dtoCameraConfigs); err != nil {
		jsonRspDTOCameraConfigList.Code = http.StatusBadRequest
		jsonRspDTOCameraConfigList.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRspDTOCameraConfigList)
		return
	}

	// Todo: config camera ...
	jsonRspDTOCameraConfigList.Data = 1
	c.JSON(http.StatusCreated, &jsonRspDTOCameraConfigList)
}

// UpdateCamera		godoc
// @Summary      	Get config from all camera
// @Description  	Config and returns number of camera. Configs key/values must be passed in the body.
// @Tags         	cameras
// @Produce      	json
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/cameras/get-config/all [get]
// @Security		BearerAuth
func GetCameraConfigAll(c *gin.Context) {

	jsonRspDTOCameraConfigList := models.NewJsonDTOListRsp[models.DTOCameraDeviceConfig]()

	// Todo: get config from all camera ...

	c.JSON(http.StatusCreated, &jsonRspDTOCameraConfigList)
}

// UpdateCamera		godoc
// @Summary      	Config multiple camera by ids
// @Description  	Config and returns number of camera. Configs key/values must be passed in the body.
// @Tags         	cameras
// @Produce      	json
// @Param        	ids  body   []string{}  true  "camera IDs JSON"
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/cameras/get-config/batch [POST]
// @Security		BearerAuth
func GetCameraConfigBatch(c *gin.Context) {

	jsonRspDTOCameraConfigList := models.NewJsonDTOListRsp[models.DTOCameraDeviceConfig]()

	/*
		// Call BindJSON to bind the received JSON to
		dtoCameraConfigs := models.NewJsonDTOListReq[string]()
		if err := c.BindJSON(&dtoCameraConfigs); err != nil {
			jsonRspDTOCameraConfigList.Code = http.StatusBadRequest
			jsonRspDTOCameraConfigList.Message = err.Error()
			c.JSON(http.StatusBadRequest, &jsonRspDTOCameraConfigList)
			return
		}
	*/

	// Call BindJSON to bind the received JSON to
	dtoCameraConfigs := []string{}
	if err := c.BindJSON(&dtoCameraConfigs); err != nil {
		jsonRspDTOCameraConfigList.Code = http.StatusBadRequest
		jsonRspDTOCameraConfigList.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRspDTOCameraConfigList)
		return
	}

	cam1Config := models.DTOCameraDeviceConfig{}

	if len(dtoCameraConfigs) > 0 {
		cam1Config.ID = uuid.MustParse(dtoCameraConfigs[0])

		// Add config TCP/IP
		cam1Config.Configs = make([]models.KeyValue, 0)
		cam1Config.Configs = append(cam1Config.Configs, models.KeyValue{
			ID:   "tcpip-nicType-auto",
			Name: "Auto",
		})
		cam1Config.Configs = append(cam1Config.Configs, models.KeyValue{
			ID:   "tcpip-dhcp-on",
			Name: "ON",
		})
		cam1Config.Configs = append(cam1Config.Configs, models.KeyValue{
			ID:   "tcpip-multicastDiscovery-enable",
			Name: "ENABLE",
		})
		cam1Config.Configs = append(cam1Config.Configs, models.KeyValue{
			ID:   "ddns-type-auto",
			Name: "Auto",
		})
		// Todo: get config of all camera IDs ...
		jsonRspDTOCameraConfigList.Data = append(jsonRspDTOCameraConfigList.Data, cam1Config)
	}
	c.JSON(http.StatusCreated, &jsonRspDTOCameraConfigList)
}

// UpdateCamera		godoc
// @Summary      	Config multiple camera by ids
// @Description  	Config and returns number of camera. Configs key/values must be passed in the body.
// @Tags         	cameras
// @Produce      	json
// @Param        	id  path  string  true  "Config camera by id"
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/cameras/get-config/{id} [get]
// @Security		BearerAuth
func GetCameraConfig(c *gin.Context) {

	jsonRspDTOCameraConfigList := models.NewJsonDTORsp[models.DTOCameraDeviceConfig]()

	// Get config for camera with ID = (c.Param("id"))
	// Todo: get camera's config ...

	//jsonRspDTOCameraConfigList.Data = 1
	c.JSON(http.StatusNotFound, &jsonRspDTOCameraConfigList)
}

// GetCameraConfigOptions		godoc
// @Summary      	Get items of camera configs for select box
// @Description  	Responds with the list of item for camera config select box depend on query params.
// @Tags         	cameras
// @Produce      	json
// @Param   		tcpip		query	string	false	"items of tcpip config subgroup"	Enums(nicType,dhcp,mtuPlaceHolder,multicastDiscovery)
// @Param   		ddns		query	string	false	"items of ddns config subgroup"		Enums(dnsType, portPlaceHolder)
// @Param   		port		query	string	false	"items of port config subgroup"		Enums(httpPortPlaceHolder, httpsPortPlaceHolder, rtspPortPlaceHolder, servicePortPlaceHolder)
// @Param   		nat			query	string	false	"items of nat config subgroup"		Enums(portMappingMode,portType,StatusType,httpPortPlaceHolder,httpsPortPlaceHolder,rtspPortPlaceHolder, servicePortPlaceHolder)
// @Param   		video		query	string	false	"items of video config subgroup"	Enums(streamType,videoType,resolution,bitrateType,videoQuality,frameRatePlaceHolder,maxBitRatePlaceHolder,videoEncoding,h264Plus,profile,iframeIntervalPlaceHolder,smoothing)
// @Param   		audio		query	string	false	"items of audio config subgroup"	Enums(audioEncoding,audioInput,inputVolumePlaceHolder,environmentNoiseFilter)
// @Param   		image		query	string	false	"items of image config subgroup"	Enums(dayNightMode,exposure,areaBLC,whiteBalance,gain,shutter,wdr,noiseReduction,hlc,mirror)
// @Param   		osd			query	string	false	"items of osd config subgroup"		Enums(displayMode,osdSize,fontColor,alignment,timeFormat,dateFormat)
// @Success      	200  		{object}  		models.JsonDTORsp[[]models.KeyValue]
// @Router       	/cameras/options/config-types [get]
// @Security		BearerAuth
func GetCameraConfigTypeOptions(c *gin.Context) {

	tcpip := c.Query("tcpip")
	//ddns := c.Query("ddns")
	//port := c.Query("port")
	//nat := c.Query("nat")
	//video := c.Query("video")
	//audio := c.Query("audio")
	//image := c.Query("image")
	//osd := c.Query("osd")

	jsonConfigTypeRsp := models.NewJsonDTORsp[[]models.KeyValue]()
	jsonConfigTypeRsp.Data = make([]models.KeyValue, 0)

	// Search for keyword in name
	if tcpip == "nicType" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "tcpip-nicType-auto",
			Name: "Auto",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "tcpip-nicType-manual",
			Name: "Manual",
		})
	} else if tcpip == "dhcp" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "tcpip-dhcp-on",
			Name: "ON",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "tcpip-dhcp-off",
			Name: "OFF",
		})
	} else if tcpip == "mtuPlaceHolder" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "tcpip-mtuPlaceHolder",
			Name: "0-1500",
		})
	} else if tcpip == "multicastDiscovery" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "tcpip-multicastDiscovery-enable",
			Name: "ENABLE",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "tcpip-multicastDiscovery-disable",
			Name: "DISABLE",
		})
	}

	jsonConfigTypeRsp.Code = 0
	c.JSON(http.StatusOK, &jsonConfigTypeRsp)
}
