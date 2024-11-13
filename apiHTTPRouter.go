package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Message resp struct
type Message struct {
	Status  int         `json:"status"`
	Payload interface{} `json:"payload"`
}

// HTTPAPIServer start http server routes
func HTTPAPIServer() {
	//Set HTTP API mode
	log.WithFields(logrus.Fields{
		"module": "http_server",
		"func":   "RTSPServer",
		"call":   "Start",
	}).Infoln("Server HTTP start")
	var public *gin.Engine
	if !Storage.ServerHTTPDebug() {
		gin.SetMode(gin.ReleaseMode)
		public = gin.New()
	} else {
		gin.SetMode(gin.DebugMode)
		public = gin.Default()
	}

	public.Use(CrossOrigin())
	//Add private login password protect methods
	privat := public.Group("/")
	if Storage.ServerHTTPLogin() != "" && Storage.ServerHTTPPassword() != "" {
		privat.Use(gin.BasicAuth(gin.Accounts{Storage.ServerHTTPLogin(): Storage.ServerHTTPPassword()}))
	}

	/*
		Static HTML Files Demo Mode
	*/

	if Storage.ServerHTTPDemo() {
		public.LoadHTMLGlob(Storage.ServerHTTPDir() + "/templates/*")
		public.GET("/", HTTPAPIServerIndex)
		public.GET("/pages/stream/list", HTTPAPIStreamList)
		public.GET("/pages/stream/add", HTTPAPIAddStream)
		public.GET("/pages/stream/edit/:uuid", HTTPAPIEditStream)
		public.GET("/pages/player/hls/:uuid/:channel", HTTPAPIPlayHls)
		public.GET("/pages/player/mse/:uuid/:channel", HTTPAPIPlayMse)
		public.GET("/pages/player/webrtc/:uuid/:channel", HTTPAPIPlayWebrtc)
		public.GET("/pages/multiview", HTTPAPIMultiview)
		public.Any("/pages/multiview/full", HTTPAPIFullScreenMultiView)
		public.GET("/pages/documentation", HTTPAPIServerDocumentation)
		public.GET("/pages/player/all/:uuid/:channel", HTTPAPIPlayAll)
		public.GET("/pages/aiEvent", aiEventHandler)
		public.GET("/pages/playback/:id", playbackHandler)
		public.StaticFS("/static", http.Dir(Storage.ServerHTTPDir()+"/static"))
	}

	/*
		Stream Control elements
	*/

	privat.GET("/streams", HTTPAPIServerStreams)
	privat.POST("/stream/:uuid/add", HTTPAPIServerStreamAdd)
	privat.POST("/stream/:uuid/edit", HTTPAPIServerStreamEdit)
	privat.GET("/stream/:uuid/delete", HTTPAPIServerStreamDelete)
	privat.GET("/stream/:uuid/reload", HTTPAPIServerStreamReload)
	privat.GET("/stream/:uuid/info", HTTPAPIServerStreamInfo)

	/*
		Streams Multi Control elements
	*/

	privat.POST("/streams/multi/control/add", HTTPAPIServerStreamsMultiControlAdd)
	privat.POST("/streams/multi/control/delete", HTTPAPIServerStreamsMultiControlDelete)

	/*
		Stream Channel elements
	*/

	privat.POST("/stream/:uuid/channel/:channel/add", HTTPAPIServerStreamChannelAdd)
	privat.POST("/stream/:uuid/channel/:channel/edit", HTTPAPIServerStreamChannelEdit)
	privat.GET("/stream/:uuid/channel/:channel/delete", HTTPAPIServerStreamChannelDelete)
	privat.GET("/stream/:uuid/channel/:channel/codec", HTTPAPIServerStreamChannelCodec)
	privat.GET("/stream/:uuid/channel/:channel/reload", HTTPAPIServerStreamChannelReload)
	privat.GET("/stream/:uuid/channel/:channel/info", HTTPAPIServerStreamChannelInfo)

	/*
		Stream video elements
	*/
	//HLS
	public.GET("/stream/:uuid/channel/:channel/hls/live/index.m3u8", HTTPAPIServerStreamHLSM3U8)
	public.GET("/stream/:uuid/channel/:channel/hls/live/segment/:seq/file.ts", HTTPAPIServerStreamHLSTS)
	//HLS remote record
	//public.GET("/stream/:uuid/channel/:channel/hls/rr/:s/:e/index.m3u8", HTTPAPIServerStreamRRM3U8)
	//public.GET("/stream/:uuid/channel/:channel/hls/rr/:s/:e/:seq/file.ts", HTTPAPIServerStreamRRTS)
	//HLS LL
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/index.m3u8", HTTPAPIServerStreamHLSLLM3U8)
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/init.mp4", HTTPAPIServerStreamHLSLLInit)
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/segment/:segment/:any", HTTPAPIServerStreamHLSLLM4Segment)
	public.GET("/stream/:uuid/channel/:channel/hlsll/live/fragment/:segment/:fragment/:any", HTTPAPIServerStreamHLSLLM4Fragment)
	//MSE
	public.GET("/stream/:uuid/channel/:channel/mse", HTTPAPIServerStreamMSE)
	public.POST("/stream/:uuid/channel/:channel/webrtc", HTTPAPIServerStreamWebRTC)
	//Save fragment to mp4
	public.GET("/stream/:uuid/channel/:channel/save/mp4/fragment/:duration", HTTPAPIServerStreamSaveToMP4)
	/*
		HTTPS Mode Cert
		# Key considerations for algorithm "RSA" ≥ 2048-bit
		openssl genrsa -out server.key 2048

		# Key considerations for algorithm "ECDSA" ≥ secp384r1
		# List ECDSA the supported curves (openssl ecparam -list_curves)
		#openssl ecparam -genkey -name secp384r1 -out server.key
		#Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)

		openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
	*/
	if Storage.ServerHTTPS() {
		if Storage.ServerHTTPSAutoTLSEnable() {
			go func() {
				err := autotls.Run(public, Storage.ServerHTTPSAutoTLSName()+Storage.ServerHTTPSPort())
				if err != nil {
					log.Println("Start HTTPS Server Error", err)
				}
			}()
		} else {
			go func() {
				err := public.RunTLS(Storage.ServerHTTPSPort(), Storage.ServerHTTPSCert(), Storage.ServerHTTPSKey())
				if err != nil {
					log.WithFields(logrus.Fields{
						"module": "http_router",
						"func":   "HTTPSAPIServer",
						"call":   "ServerHTTPSPort",
					}).Fatalln(err.Error())
					os.Exit(1)
				}
			}()
		}
	}
	err := public.Run(Storage.ServerHTTPPort())
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "http_router",
			"func":   "HTTPAPIServer",
			"call":   "ServerHTTPPort",
		}).Fatalln(err.Error())
		os.Exit(1)
	}

}

// HTTPAPIServerIndex index file
func HTTPAPIServerIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "index",
	})

}

func HTTPAPIServerDocumentation(c *gin.Context) {
	c.HTML(http.StatusOK, "documentation.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "documentation",
	})
}

func HTTPAPIStreamList(c *gin.Context) {
	c.HTML(http.StatusOK, "stream_list.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "stream_list",
	})
}

func HTTPAPIPlayHls(c *gin.Context) {
	c.HTML(http.StatusOK, "play_hls.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "play_hls",
		"uuid":    c.Param("uuid"),
		"channel": c.Param("channel"),
	})
}
func HTTPAPIPlayMse(c *gin.Context) {
	c.HTML(http.StatusOK, "play_mse.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "play_mse",
		"uuid":    c.Param("uuid"),
		"channel": c.Param("channel"),
	})
}
func HTTPAPIPlayWebrtc(c *gin.Context) {
	c.HTML(http.StatusOK, "play_webrtc.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "play_webrtc",
		"uuid":    c.Param("uuid"),
		"channel": c.Param("channel"),
	})
}
func HTTPAPIAddStream(c *gin.Context) {
	c.HTML(http.StatusOK, "add_stream.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "add_stream",
	})
}
func HTTPAPIEditStream(c *gin.Context) {
	c.HTML(http.StatusOK, "edit_stream.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "edit_stream",
		"uuid":    c.Param("uuid"),
	})
}

func HTTPAPIMultiview(c *gin.Context) {
	c.HTML(http.StatusOK, "multiview.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "multiview",
	})
}

func HTTPAPIPlayAll(c *gin.Context) {
	c.HTML(http.StatusOK, "play_all.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"page":    "play_all",
		"uuid":    c.Param("uuid"),
		"channel": c.Param("channel"),
	})
}

type MultiViewOptions struct {
	Grid   int                             `json:"grid"`
	Player map[string]MultiViewOptionsGrid `json:"player"`
}
type MultiViewOptionsGrid struct {
	UUID       string `json:"uuid"`
	Channel    int    `json:"channel"`
	PlayerType string `json:"playerType"`
}

func HTTPAPIFullScreenMultiView(c *gin.Context) {
	var createParams MultiViewOptions
	err := c.ShouldBindJSON(&createParams)
	if err != nil {
		log.WithFields(logrus.Fields{
			"module": "http_page",
			"func":   "HTTPAPIFullScreenMultiView",
			"call":   "BindJSON",
		}).Errorln(err.Error())
	}
	log.WithFields(logrus.Fields{
		"module": "http_page",
		"func":   "HTTPAPIFullScreenMultiView",
		"call":   "Options",
	}).Debugln(createParams)
	c.HTML(http.StatusOK, "fullscreenmulti.tmpl", gin.H{
		"port":    Storage.ServerHTTPPort(),
		"streams": Storage.Streams,
		"version": time.Now().String(),
		"options": createParams,
		"page":    "fullscreenmulti",
		"query":   c.Request.URL.Query(),
	})
}

// CrossOrigin Access-Control-Allow-Origin any methods
func CrossOrigin() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func aiEventHandler(c *gin.Context) {
	events := []Event{
		{
			ID:          "1",
			Thumbnail:   "/static/img/event1.jpg",
			Timestamp:   "2024-11-12 15:10:38",
			Description: "Car detected",
			Location:    "PHUCHAN_TEST2",
			Result:      "PID 7282 | 7420837",
			Object:      "Phát hiện người đi xe",
			Date:        "2024-11-12 15:10:30",
			Camera:      "PHUCHAN_TEST3",
			FullImage:   "/static/img/event1_full.jpg",
		},
		{
			ID:          "2",
			Thumbnail:   "/static/img/event2.jpg",
			Timestamp:   "2024-11-12 15:10:37",
			Description: "Motorbike detected",
			Location:    "PHUCHAN_TEST3",
			Result:      "PID 7281 | 7420472",
			Object:      "Phát hiện người đi xe máy",
			Date:        "2024-11-12 15:10:14",
			Camera:      "PHUCHAN_TEST3",
			FullImage:   "/static/img/event2_full.jpg",
		},
	}

	c.HTML(http.StatusOK, "ai_event_list.tmpl", gin.H{
		"events": events,
	})
}

type Event struct {
	ID          string
	Thumbnail   string
	Timestamp   string
	Description string
	Location    string
	Result      string
	Object      string
	Date        string
	Camera      string
	FullImage   string
}

func playbackHandler(c *gin.Context) {
	eventID := c.Param("id")

	events := map[string]Event{
		"1": {
			ID:          "1",
			Thumbnail:   "/static/img/event1.jpg",
			Timestamp:   "2024-11-12 15:10:38",
			Description: "Car detected",
			Location:    "PHUCHAN_TEST2",
			Result:      "PID 7282 | 7420837",
			Object:      "Phát hiện người đi xe",
			Date:        "2024-11-12 15:10:30",
			Camera:      "PHUCHAN_TEST3",
			FullImage:   "/static/img/event1_full.jpg",
		},
		"2": {
			ID:          "2",
			Thumbnail:   "/static/img/event2.jpg",
			Timestamp:   "2024-11-12 15:10:37",
			Description: "Motorbike detected",
			Location:    "PHUCHAN_TEST3",
			Result:      "PID 7281 | 7420472",
			Object:      "Phát hiện người đi xe máy",
			Date:        "2024-11-12 15:10:14",
			Camera:      "PHUCHAN_TEST3",
			FullImage:   "/static/img/event2_full.jpg",
		},
	}

	event, exists := events[eventID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.HTML(http.StatusOK, "playback.tmpl", gin.H{
		"event": event,
	})
}
