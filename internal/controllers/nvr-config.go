package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"vms/internal/models"

	"vms/comongo/kafkaclient"
	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdateNVR		godoc
// @Summary      	Config all nvr
// @Description  	Config and returns number of nvrs. Configs key/values must be passed in the body.
// @Tags         	nvrs
// @Produce      	json
// @Param        	nvr  body   []models.DTO_NVR_DeviceConfig{}  true  "NVR Config JSON"
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/nvrs/set-config/all [post]
// @Security		BearerAuth
func SetNVRConfigAll(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[int]()
	/*
		// Call BindJSON to bind the received JSON to
		dto := models.NewJsonDTOListReq[models.DTO_NVR_DeviceConfig]()
		if err := c.BindJSON(&dto); err != nil {
			jsonRsp.Code = http.StatusBadRequest
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusBadRequest, &jsonRsp)
			return
		}
	*/
	// Call BindJSON to bind the received JSON to
	dto := []models.DTO_NVR_DeviceConfig{}
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	// Todo: config nvr ...
	jsonRsp.Data = 1
	c.JSON(http.StatusCreated, &jsonRsp)
}

// UpdateNVR		godoc
// @Summary      	Config multiple nvr by id
// @Description  	Config and returns nvrIDs. Configs key/values must be passed in the body.
// @Tags         	nvrs
// @Produce      	json
// @Param        	nvr  body 	models.DTO_NVR_DeviceConfigBatch true  "NVR Config Batch JSON"
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/nvrs/set-config/batch [post]
// @Security		BearerAuth
func SetNVRConfigBatch(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[int]()
	/*
		// Call BindJSON to bind the received JSON to
		dto := models.NewJsonDTOListReq[models.DTO_NVR_DeviceConfig]()
		if err := c.BindJSON(&dto); err != nil {
			jsonRsp.Code = http.StatusBadRequest
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusBadRequest, &jsonRsp)
			return
		}
	*/
	// Call BindJSON to bind the received JSON to
	dto := models.DTO_NVR_DeviceConfigBatch{}
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	// Todo: config nvr ...
	jsonRsp.Data = len(dto.NVRIDs)
	c.JSON(http.StatusCreated, &jsonRsp)
}

// UpdateNVR		godoc
// @Summary      	Config multiple nvr by id
// @Description  	Config and returns cameraIDs. Configs key/values must be passed in the body.
// @Tags         	nvrs
// @Produce      	json
// @Param        	id  path  string  true  "Config nvr by id"
// @Param        	nvr  body     []models.DTO_NVR_DeviceConfig{} true  "NVR JSON"
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/nvrs/set-config/{id} [post]
// @Security		BearerAuth
func SetNVRConfig(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[int]()

	// Call BindJSON to bind the received JSON to
	dto := models.DTO_NVR_DeviceConfig{}
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Todo: config nvr ...
	jsonRsp.Data = 1
	c.JSON(http.StatusCreated, &jsonRsp)
}

// UpdateNVR		godoc
// @Summary      	Get config from all nvr
// @Description  	Config and returns number of nvr. Configs key/values must be passed in the body.
// @Tags         	nvrs
// @Produce      	json
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/nvrs/get-config/all [get]
// @Security		BearerAuth
func GetNVRConfigAll(c *gin.Context) {

	jsonRsp := models.NewJsonDTOListRsp[models.DTO_NVR_DeviceConfig]()

	// Todo: get config from all nvr ...

	c.JSON(http.StatusCreated, &jsonRsp)
}

// UpdateNVR		godoc
// @Summary      	Config multiple nvr by ids
// @Description  	Config and returns number of nvr. Configs key/values must be passed in the body.
// @Tags         	nvrs
// @Produce      	json
// @Param        	ids  body   []string{}  true  "nvr IDs JSON"
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/nvrs/get-config/batch [POST]
// @Security		BearerAuth
func GetNVRConfigBatch(c *gin.Context) {

	jsonRsp := models.NewJsonDTOListRsp[models.DTO_NVR_DeviceConfig]()

	/*
		// Call BindJSON to bind the received JSON to
		dto := models.NewJsonDTOListReq[string]()
		if err := c.BindJSON(&dto); err != nil {
			jsonRsp.Code = http.StatusBadRequest
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusBadRequest, &jsonRsp)
			return
		}
	*/
	// Call BindJSON to bind the received JSON to
	dto := []string{}
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	cam1Config := models.DTO_NVR_DeviceConfig{}

	if len(dto) > 0 {
		cam1Config.ID = uuid.MustParse(dto[0])

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
		// Todo: get config of all nvr IDs ...
		jsonRsp.Data = append(jsonRsp.Data, cam1Config)
	}
	c.JSON(http.StatusCreated, &jsonRsp)
}

// UpdateNVR		godoc
// @Summary      	Config multiple nvr by ids
// @Description  	Config and returns number of nvr. Configs key/values must be passed in the body.
// @Tags         	nvrs
// @Produce      	json
// @Param        	id  path  string  true  "Config nvr by id"
// @Success      	200  {object}  models.JsonDTORsp[int]
// @Router       	/nvrs/get-config/{id} [get]
// @Security		BearerAuth
func GetNVRConfig(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_NVR_DeviceConfig]()

	// Get config for nvr with ID = (c.Param("id"))
	// Todo: get nvr's config ...

	//jsonRsp.Data = 1
	c.JSON(http.StatusNotFound, &jsonRsp)
}

// GetNVRConfigOptions		godoc
// @Summary      	Get items of nvr configs for select box
// @Description  	Responds with the list of item for nvr config select box depend on query params.
// @Tags         	nvrs
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
// @Router       	/nvrs/options/config-types [get]
// @Security		BearerAuth
func GetNVRConfigTypeOptions(c *gin.Context) {

	tcpip := c.Query("tcpip")
	// ddns := c.Query("ddns")
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
			ID:   "Auto",
			Name: "Auto",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "10M_Half_Dup",
			Name: "10M Half-Dup",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "10M_Full_Dup",
			Name: "10M Full-Dup",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "100M_Half_Dup",
			Name: "100M Half-Dup",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "100M_Full_Dup",
			Name: "100M Full-Dup",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "10M_Half_Dup",
			Name: "10M Half-Dup",
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

// UpdateVideoConfigNVR godoc
// @Summary      Update video configuration for a specific NVR
// @Description  Updates the video configuration of the camera identified by the provided ID.
// @Tags         nvrs
// @Produce      json
// @Param        id           path     string  true  "NVR ID"
// @Param        videoConfig  body   models.DTO_VideoConfig  true  "DTO_VideoConfig JSON"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_VideoConfig]
// @Failure      400   {object}  models.JsonDTORsp[models.DTO_VideoConfig]  "Bad Request"
// @Failure      404   {object}  models.JsonDTORsp[models.DTO_VideoConfig]  "Device Not Found"
// @Failure      408   {object}  models.JsonDTORsp[models.DTO_VideoConfig]  "Request Timeout"
// @Failure      500   {object}  models.JsonDTORsp[models.DTO_VideoConfig]  "Internal Server Error"
// @Router       /nvrs/config/videoconfig/{id} [put]
// @Security     BearerAuth
func UpdateVideoConfigNVR(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_VideoConfig]()
	idNVR := c.Param("id")

	// Get new data from body
	var dto models.DTO_VideoConfig
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	cameraID := dto.CameraID
	dtoCamera, errCam := reposity.ReadItemByIDIntoDTO[models.DTOCamera, models.Camera](cameraID.String())
	if errCam != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = errCam.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	dtoNVR, err := reposity.ReadItemByIDIntoDTO[models.DTO_NVR, models.NVR](idNVR)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	dtoNVRConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVRConfig, models.NVRConfig]("id = ?", dtoNVR.ConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Update entity from DTO
	UpdatedDTO, err := reposity.UpdateItemByIDFromDTO[models.DTO_VideoConfig, models.VideoConfig](dtoNVRConfig.VideoConfigID.String(), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Map VideoConfigInfo to StreamingChannelList
	streamingChannelList := mapVideoConfigToStreamingChannelsNVR(UpdatedDTO, dtoCamera.NVR.Channel)

	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoNVR.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	RequestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:               dtoDevice.ModelID,
		Cmd:                     cmd_SetVideoConfigOfNVR,
		EventTime:               time.Now().Format(time.RFC3339),
		ProtocolType:            dtoNVR.Protocol,
		RequestUUID:             RequestUUID,
		IPAddress:               dtoNVR.IPAddress,
		UserName:                dtoNVR.Username,
		Password:                dtoNVR.Password,
		OnvifPort:               dtoNVR.OnvifPort,
		HttpPort:                dtoNVR.HttpPort,
		Channel:                 dtoCamera.NVR.Channel,
		StreamingChannelListNVR: streamingChannelList,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	cmsSStr, err := json.MarshalIndent(cmd, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	} else {
		fmt.Println(string(cmsSStr))
	}
	// Set timeout and ticker
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	// Wait for response
	for {
		select {
		case <-timeout:
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_SetVideoConfigOfNVR)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_SetVideoConfigOfNVR
			jsonRsp.Code = http.StatusRequestTimeout
			c.JSON(http.StatusRequestTimeout, &jsonRsp)
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(RequestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				jsonRsp.Message = "Command sent successfully, please check the table to grasp specific information: " + msg.PayLoad.Cmd
				jsonRsp.Data = UpdatedDTO
				jsonRsp.Code = http.StatusOK
				c.JSON(http.StatusOK, &jsonRsp)
				return
			}
		}
	}
}

func mapVideoConfigToStreamingChannelsNVR(dtoConfig models.DTO_VideoConfig, channelCamera string) models.StreamingChannelListNVR {
	streamingChannels := make([]models.StreamingChannelNVRUpdate, 0)

	for i, vci := range dtoConfig.VideoConfigInfo {
		width, height := parseResolution(vci.Resolution)
		numchannelCamera, err := strconv.Atoi(channelCamera)
		if err != nil {
			fmt.Println("Error Converted numchannelCamera to number:", err)
		}
		id := 1 + i + numchannelCamera*100
		streamingChannel := models.StreamingChannelNVRUpdate{
			ID: id,
			TransportNVR: models.TransportNVR{
				RtspPortNo:    554,
				MaxPacketSize: 1446,
				ControlProtocolList: models.ControlProtocolListNVR{
					ControlProtocols: []models.ControlProtocol{
						{StreamingTransport: "RTSP"},
					},
				},
			},
			//ChannelName: fmt.Sprintf("%d", id),
			Enabled: true,
			Video: models.Video{
				VideoCodecType:          vci.VideoEncoding,
				VideoResolutionWidth:    width,
				VideoResolutionHeight:   height,
				VideoQualityControlType: vci.BitrateType,
				FixedQuality:            stringToInt(vci.VideoQuality),
				VbrUpperCap:             stringToInt(vci.MaxBitrate),
				MaxFrameRate:            stringToInt(vci.FrameRate),
				H265Profile:             vci.H265,
				VideoScanType:           "progressive",
			},
		}

		streamingChannels = append(streamingChannels, streamingChannel)
	}

	return models.StreamingChannelListNVR{
		StreamingChannels: streamingChannels,
	}
}

// ReadNVRConfig godoc
// @Summary      Get single nvr by id
// @Description  Returns the nvr whose ID value matches the id.
// @Tags         cameras
// @Produce      json
// @Param   	 id			path	string	true	"id nvr keyword"		minlength(1)  	maxlength(100)
// @Success      200   {object}  models.JsonDTORsp[models.DTO_VideoConfig]
// @Router       /nvrs/config/videoconfigNVR/{id} [get]
// @Security	 BearerAuth
func GetVideoConfigNVR(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_VideoConfig]()
	nvrID := c.Param("id")

	dtoNVR, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVR, models.NVR]("id = ?", nvrID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	dtoNVRConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVRConfig, models.NVRConfig]("id = ?", dtoNVR.ConfigID)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	dtoDevice, err := reposity.ReadItemByIDIntoDTO[models.DTO_Device, models.Device](dtoNVR.Box.ID)
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = "Device not found: " + err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	requestUUID := uuid.New()
	cmd := models.DeviceCommand{
		CommandID:    dtoDevice.ModelID,
		Cmd:          cmd_GetVideoConfigNVR,
		EventTime:    time.Now().Format(time.RFC3339),
		ProtocolType: dtoNVR.Protocol,
		RequestUUID:  requestUUID,
		IPAddress:    dtoNVR.IPAddress,
		UserName:     dtoNVR.Username,
		Password:     dtoNVR.Password,
		OnvifPort:    dtoNVR.OnvifPort,
		HttpPort:     dtoNVR.HttpPort,
	}
	cmsStr, err := json.Marshal(cmd)
	if err != nil {
		fmt.Println(err)
	}
	cmsSStr, err := json.MarshalIndent(cmd, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
	} else {
		fmt.Println(string(cmsSStr))
	}
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	// Set timeout and ticker
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			dtoNVRConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_NVRConfig, models.NVRConfig]("id = ?", dtoNVR.ConfigID)
			if err != nil {
				jsonRsp.Code = http.StatusInternalServerError
				jsonRsp.Message = err.Error()
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}

			dtoNVRVideoConfig, err := reposity.ReadItemWithFilterIntoDTO[models.DTO_VideoConfig, models.VideoConfig]("id = ?", dtoNVRConfig.VideoConfigID)
			if err != nil {
				jsonRsp.Code = http.StatusInternalServerError
				jsonRsp.Message = err.Error()
				c.JSON(http.StatusInternalServerError, &jsonRsp)
				return
			}
			jsonRsp.Data = dtoNVRVideoConfig
			c.JSON(http.StatusOK, &jsonRsp)
			fmt.Println("\t\t> error, waiting for command result timed out: ", cmd_GetOnSiteVideoConfig)
			jsonRsp.Message = "No response from the device, timeout: " + cmd_GetOnSiteVideoConfig
			jsonRsp.Code = http.StatusRequestTimeout
			return

		case <-ticker.C:
			if storedMsg, ok := messageMap.Load(requestUUID); ok {
				msg := storedMsg.(*models.KafkaJsonVMSMessage)

				if msg.PayLoad.VideoConfigNVR != nil {

					for i := range msg.PayLoad.VideoConfigNVR.VideoConfigInfo {
						for _, camera := range *dtoNVR.Cameras {
							if msg.PayLoad.VideoConfigNVR.VideoConfigInfo[i].CameraChannel == camera.Channel {
								msg.PayLoad.VideoConfigNVR.VideoConfigInfo[i].CameraName = camera.Name
							}
						}
					}
					// Update entity from DTO
					dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_VideoConfig, models.VideoConfig](dtoNVRConfig.VideoConfigID.String(), *msg.PayLoad.VideoConfigNVR)
					if err != nil {
						jsonRsp.Code = http.StatusInternalServerError
						jsonRsp.Message = err.Error()
						c.JSON(http.StatusInternalServerError, &jsonRsp)
						return
					}

					jsonRsp.Data = dto
					c.JSON(http.StatusOK, &jsonRsp)
					return
				} else {
					jsonRsp.Message = "No video config data received"
					jsonRsp.Code = http.StatusNoContent
					c.JSON(http.StatusNoContent, &jsonRsp)

				}
			}
		}
	}
}
