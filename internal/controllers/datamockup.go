package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"vms/internal/models"

	"github.com/gin-gonic/gin"
)

// GetDataMockup		godoc
// @Summary      	Get data mockup
// @Description  	Responds with the list of all camera as JSON.
// @Tags         	data mockup
// @Produce      	json
// @Success      	200  {object} models.JsonDTORsp[models.NetworkInterface]
// @Router       	/datamockup/tcpip [get]
// @Security		BearerAuth
func GetTCPIP(c *gin.Context) {
	jsonRspBasicInfos := models.NewJsonDTOListRsp[models.NetworkInterface]()

	networkInterface := models.NetworkInterface{
		Token: "etho",
		Link: models.Link{
			AutoNegotiation: true,
			Speed:           100,
			Duplex:          "Full",
		},
		DHCP:         true,
		IPv4Address:  "192.168.1.100",
		IPv4Branch:   "255.255.255.0",
		IPv4Default:  "192.168.1.1",
		IPv6Address:  "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		Length:       "1500",
		IPv6Default:  "::",
		MAC:          "00:11:22:33:44:55",
		MUT:          1500,
		DNS:          true,
		PreferredDNS: "8.8.8.8",
		AlternateDNS: "8.8.4.4",
	}

	jsonRspBasicInfos.Count = 1
	jsonRspBasicInfos.Data = []models.NetworkInterface{networkInterface}
	jsonRspBasicInfos.Page = 1
	jsonRspBasicInfos.Size = 1
	c.JSON(http.StatusOK, &jsonRspBasicInfos)
}

// GetDataMockup		godoc
// @Summary      	Get data mockup
// @Description  	Responds with the list of all camera as JSON.
// @Tags         	data mockup
// @Produce      	json
// @Success      	200  {object} models.JsonDTORsp[models.DDNS]
// @Router       	/datamockup/ddns [get]
// @Security		BearerAuth
func DDNS(c *gin.Context) {
	jsonRspBasicInfos := models.NewJsonDTOListRsp[models.DDNS]()

	ddns := models.DDNS{
		Token:         "token123",
		Enabled:       true,
		ServerAddress: "ddns.example.com",
		Domain:        "example.com",
		UserName:      "user123",
		Password:      123456,
		Confirm:       "confirm123",
	}

	jsonRspBasicInfos.Count = 1
	jsonRspBasicInfos.Data = []models.DDNS{ddns}
	jsonRspBasicInfos.Page = 1
	jsonRspBasicInfos.Size = 1
	c.JSON(http.StatusOK, &jsonRspBasicInfos)
}

// GetDataMockup		godoc
// @Summary      	Get data mockup
// @Description  	Responds with the list of all camera as JSON.
// @Tags         	data mockup
// @Produce      	json
// @Success      	200  {object} models.JsonDTORsp[models.AdminAccessProtocolList]
// @Router       	/datamockup/port [get]
// @Security		BearerAuth
func Port(c *gin.Context) {
	jsonRspBasicInfos := models.NewJsonDTOListRsp[models.AdminAccessProtocolList]()

	protocols := models.AdminAccessProtocolList{
		Protocols: []models.AdminAccessProtocol{
			{
				Id:       1,
				Enabled:  true,
				Protocol: "HTTP",
				PortNo:   8088,
			},
			{
				Id:       2,
				Enabled:  true,
				Protocol: "RTSP",
				PortNo:   554,
			},
			{
				Id:       3,
				Enabled:  false,
				Protocol: "HTTPS",
				PortNo:   443,
			},
			{
				Id:       4,
				Protocol: "DEV_MANAGE",
				PortNo:   8000,
			},
		},
	}

	jsonRspBasicInfos.Count = 1
	jsonRspBasicInfos.Data = []models.AdminAccessProtocolList{protocols}
	jsonRspBasicInfos.Page = 1
	jsonRspBasicInfos.Size = 1
	c.JSON(http.StatusOK, &jsonRspBasicInfos)
}

// GetDataMockup		godoc
// @Summary      	Get data mockup
// @Description  	Responds with the list of all camera as JSON.
// @Tags         	data mockup
// @Produce      	json
// @Success      	200  {object} models.JsonDTORsp[models.NTPList]
// @Router       	/datamockup/ntp [get]
// @Security		BearerAuth
func NTP(c *gin.Context) {

	jsonRspBasicInfos := models.NewJsonDTOListRsp[models.NTPList]()

	ntpList := models.NTPList{
		NTP: models.NTP{
			Token:         "token123",
			NTP:           true,
			DateTimeType:  "NTP",
			ServerAddress: "ntp.example.com",
			NTPPort:       "123",
			Period:        "3600",
		},
		SystemDateAndTime: models.SystemDateAndTime{
			DateTimeType:  "Manual",
			TimeZone:      "ICT-7",
			UTCDateTime:   time.Date(2023, 8, 29, 8, 30, 0, 0, time.UTC),
			LocalDateTime: time.Date(2023, 8, 29, 15, 30, 0, 0, time.FixedZone("ICT", 7*60*60)),
		},
	}

	jsonRspBasicInfos.Count = 1
	jsonRspBasicInfos.Data = []models.NTPList{ntpList}
	jsonRspBasicInfos.Page = 1
	jsonRspBasicInfos.Size = 1
	c.JSON(http.StatusOK, &jsonRspBasicInfos)

}

// GetDataMockup		godoc
// @Summary      	Get data mockup
// @Description  	Responds with the list of all camera as JSON.
// @Tags         	data mockup
// @Produce      	json
// @Success      	200  {object} models.JsonDTORsp[models.VideoOverlay]
// @Router       	/datamockup/osd [get]
// @Security		BearerAuth
func OSD(c *gin.Context) {
	vo := models.VideoOverlay{
		NormalizedScreenSize: []models.NormalizedScreenSize{
			{
				NormalizedScreenWidth:  1920,
				NormalizedScreenHeight: 1080,
			},
		},
		Attribute: []models.Attribute{
			{
				Transparent: true,
				Flashing:    false,
			},
		},
		TextOverlayList: []models.TextOverlayList{
			{
				TextOverlays: []models.TextOverlays{
					{
						ID:          1,
						Enabled:     true,
						PositionX:   100,
						PositionY:   100,
						DisplayText: "Hello, world!",
					},
				},
			},
		},
		DateTimeOverlay: []models.DateTimeOverlay{
			{
				Enabled:     true,
				PositionX:   200,
				PositionY:   200,
				DateStyle:   "yyyy-MM-dd",
				TimeStyle:   "HH:mm:ss",
				DisplayWeek: false,
			},
		},
		ChannelNameOverlay: []models.ChannelNameOverlay{
			{
				Enabled:   false,
				PositionX: 0,
				PositionY: 0,
			},
		},
	}

	jsonBytes, err := json.Marshal(vo)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, "application/json", jsonBytes)
}
