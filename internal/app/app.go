/*
Package app is the primary runtime service.
*/
package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/madflojo/hord"
	"github.com/madflojo/hord/drivers/cassandra"
	"github.com/madflojo/hord/drivers/redis"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"net/http"
	"net/http/pprof"

	"vms/appconfig"
	"vms/internal/controllers"
	"vms/internal/models"
	services "vms/internal/services/ws"
	"vms/wsnode"

	"vms/wssignaling"

	"vms/comongo/auth"
	"vms/comongo/kafkaclient"
	"vms/comongo/minioclient"
	"vms/comongo/reposity"
	"vms/comongo/telemetry"
)

// Common errors returned by this app.
var (
	ErrShutdown = fmt.Errorf("application shutdown gracefully")
)

// srv is the global reference for the HTTP Server.
var srv *http.Server

// kv is the global reference for the K/V Store.
var kv hord.Database

// runCtx is a global context used to control shutdown of the application.
var runCtx context.Context

// runCancel is a global context cancelFunc used to trigger the shutdown of applications.
var runCancel context.CancelFunc

// cfg is used across the app package to contain configuration.
var cfg *viper.Viper

// log is used across the app package for logging.
var log *logrus.Logger

// scheduler is a internal task scheduler for recurring tasks
//var scheduler *tasks.Scheduler

// stats is used across the app package to manage and access system metrics.
var stats = telemetry.New()

// Run starts the primary application. It handles starting background services,
// populating package globals & structures, and clean up tasks.
func Run(c *viper.Viper) error {
	var err error

	// Create App Context
	runCtx, runCancel = context.WithCancel(context.Background())

	// Apply config provided by config package global
	cfg = c

	// Initiate a new logger
	log = logrus.New()
	if cfg.GetBool("debug") {
		log.Level = logrus.DebugLevel
		log.Debug("Enabling Debug Logging")
	}
	if cfg.GetBool("trace") {
		log.Level = logrus.TraceLevel
		log.Debug("Enabling Trace Logging")
	}
	if cfg.GetBool("disable_logging") {
		log.Level = logrus.FatalLevel
	}

	// Setup the KV Connection
	if cfg.GetBool("enable_kvstore") {
		log.Infof("Connecting to KV Store")
		switch cfg.GetString("kvstore_type") {
		case "redis":
			kv, err = redis.Dial(redis.Config{
				Server:   cfg.GetString("redis_server"),
				Password: cfg.GetString("redis_password"),
				SentinelConfig: redis.SentinelConfig{
					Servers: cfg.GetStringSlice("redis_sentinel_servers"),
					Master:  cfg.GetString("redis_sentinel_master"),
				},
				ConnectTimeout: time.Duration(cfg.GetInt("redis_connect_timeout")) * time.Second,
				Database:       cfg.GetInt("redis_database"),
				SkipTLSVerify:  cfg.GetBool("redis_hostname_verify"),
				KeepAlive:      time.Duration(cfg.GetInt("redis_keepalive")) * time.Second,
				MaxActive:      cfg.GetInt("redis_max_active"),
				ReadTimeout:    time.Duration(cfg.GetInt("redis_read_timeout")) * time.Second,
				WriteTimeout:   time.Duration(cfg.GetInt("redis_write_timeout")) * time.Second,
			})
			if err != nil {
				return fmt.Errorf("could not establish kvstore connection - %s", err)
			}
		case "cassandra":
			kv, err = cassandra.Dial(cassandra.Config{
				Hosts:                      cfg.GetStringSlice("cassandra_hosts"),
				Port:                       cfg.GetInt("cassandra_port"),
				Keyspace:                   cfg.GetString("cassandra_keyspace"),
				Consistency:                cfg.GetString("cassandra_consistency"),
				ReplicationStrategy:        cfg.GetString("cassandra_repl_strategy"),
				Replicas:                   cfg.GetInt("cassandra_replicas"),
				User:                       cfg.GetString("cassandra_user"),
				Password:                   cfg.GetString("cassandra_password"),
				EnableHostnameVerification: cfg.GetBool("cassandra_hostname_verify"),
			})
			if err != nil {
				return fmt.Errorf("could not establish kvstore connection - %s", err)
			}
		default:
			return fmt.Errorf("unknown kvstore specified - %s", cfg.GetString("kvstore_type"))
		}

		// Clean up KV Store connections on shutdown
		defer kv.Close()

		// Initialize the KV
		err = kv.Setup()
		if err != nil {
			return fmt.Errorf("could not setup kvstore - %s", err)
		}
	}

	if kv == nil {
		log.Infof("KV Store not configured, skipping")
	}

	// Setup the SQL Connection
	if appconfig.Enable_sql {
		log.Infof("Connecting to DB")
		err := reposity.Connect(
			appconfig.SQL_Host,
			appconfig.SQL_Port,
			appconfig.SQL_dbname,
			appconfig.SQL_sslmode,
			appconfig.SQL_Username,
			appconfig.SQL_Password,
			appconfig.SQL_Schema)
		if err != nil {
			return fmt.Errorf("could not establish database connection - %s", err)
		}
	}
	if !reposity.Connected {
		log.Infof("SQL DB not configured, skipping")
	} else {
		err = reposity.Ping()
		if err != nil {
			log.Errorf(err.Error())
		} else {
			log.Info("Successfully connected to database")
			// Migrate tables
			err = reposity.Migrate(
				&models.ImageConfig{},
				&models.VideoConfig{},
				&models.AudioConfig{},
				&models.NetworkConfig{},
				&models.StorageConfig{},
				&models.RecordingSchedule{},
				&models.StreamingConfig{},
				&models.AIConfig{},
				&models.PTZConfig{},
				&models.CameraConfig{},
				&models.Camera{},
				&models.CameraGroup{},
				&models.NVRConfig{},
				&models.NVR{},
				&models.CCTVEvent{},
				&models.AIEvent{},
				&models.Device{},
				&models.DeviceProduct{},
				&models.SystemIncident{},
				&models.MediaLibrary{},
				&models.User{},
				&models.UserHome{},
				&models.Login{},
				&models.Event{},
				&models.BlackList{},
				&models.NetworkInterfaces{},
				&models.EventConfig{},
				&models.UserCameraGroup{},
				&models.CameraAIEventProperty{},
				&models.LicensePlates{},
				&models.DTO_AI_Event{},
				&models.AIWaring{},
				&models.CameraModelAI{},
				&models.PcInfo{},
				&models.NVRImport{},
				&models.AIEngine{},
			)
			if err != nil {
				panic("Failed to AutoMigrate table! err: " + err.Error())
			}
		}
	}

	defer reposity.Close()

	// Run ws node
	wsnode.Run()
	// Register processing functions coresponding to channel prefix for wsnode
	wsnode.RegisterProcessingFunction("aieagent", services.WSHanleOnSubscribeAIEAgent, services.WSHandleOnDisconnectAIEAgent, services.WSHandleOnPublishAIEAgent)
	wsnode.RegisterProcessingFunction("cam", services.WSHanleOnSubscribeCam, services.WSHandleOnDisconnectCam, services.WSHandleOnPublishCam)
	wsnode.RegisterProcessingFunction("NVR", services.WSHanleOnSubscribeNVR, services.WSHandleOnDisconnectNVR, services.WSHandleOnPublishNVR)

	// Setup kafka client
	// Init kafka client
	err = kafkaclient.Init(appconfig.Kafka_Bootstrap_Ivi, true, appconfig.Kafka_Topic_Producer_Device_Command_Ivi, true,
		appconfig.Kafka_Topic_Consumer_Group)
	if err != nil {
		panic("Failed to setup kafka client, err: " + err.Error())
	}

	// Init minio client

	err = minioclient.Connect(
		appconfig.Minio_Endpoint,
		appconfig.Minio_Accesskey_ID,
		appconfig.Minio_Secret_Access_Key,
		appconfig.Minio_UseSSL,
	)
	if err != nil {
		panic("Failed to connect to minio, err: " + err.Error())
	}
	// Test minio bucket
	err = minioclient.TestBucketAccess("vms-dev")
	if err != nil {
		panic("Failed to test minio bucket access, err: " + err.Error())
	}
	// Create consumer for topic device event and ai vms event
	deviceEventConsumerName, err := kafkaclient.CreateNewConsumer(
		appconfig.Kafka_Topic_Consumer_Device_Feedback_Ivi)
	if err != nil {
		panic("Failed to create kafka consumer for device event, err: " + err.Error())
	}
	// Run device event message processing
	controllers.RunDeviceEventProcessing(deviceEventConsumerName, appconfig.Kafka_Topic_Producer_Device_Command_Ivi)

	// Create consumer for topic device event dmp
	dmpEventConsumerName, err := kafkaclient.CreateNewConsumer(appconfig.Kafka_Topic_Consumer_Device_Info_Ivi)
	if err != nil {
		panic("Failed to create kafka consumer for device dmp event, err: " + err.Error())
	}
	// Run device event message processing
	controllers.RunDeviceDMPEventProcessing(dmpEventConsumerName, appconfig.Kafka_Topic_Producer_Device_Command_Ivi_Face_Register)

	// Create consumer for topic ai core event
	aicoreEventConsumerName, err := kafkaclient.CreateNewConsumer(appconfig.Kafka_Topic_Consumer_AI_VMS_Event)
	if err != nil {
		panic("Failed to create kafka consumer for ai-core event, err: " + err.Error())
	}
	// Run event message monitor
	controllers.RunAICoreEventProcessing(aicoreEventConsumerName)

	// Create new websocket hub for webRTC signaling and run it
	wssignaling.Start()
	wssignaling.StartTopic()

	/*
		// Create ivis sync and run it
		ivs := ivissync.NewIVISSync(log)
		ivs.SetAPILoginURL(cfg.GetString("client_iam_login_uri"))
		ivs.SetAPIGetAllCabinURL(cfg.GetString("client_cabin_sync_uri"))
		ivs.SetLoginInfo(cfg.GetString("client_iam_clientid"), cfg.GetString("client_iam_clientsecret"))
		ivissync.Start()
		ivissync.Login()
	*/

	// Setup router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowCredentials:       true,
		AllowWildcard:          true,
		AllowBrowserExtensions: true,
		AllowOrigins:           []string{"*"},
		//AllowOrigins:     []string{"https://foo.com"},
		AllowMethods: []string{"POST", "PUT", "PATCH", "DELETE", "GET", "OPTIONS", "UPDATE"},
		AllowHeaders: []string{
			"Content-Type, content-length, accept-encoding, X-CSRF-Token, " +
				"access-control-allow-origin, Authorization, X-Max, access-control-allow-headers, " +
				"accept, origin, Cache-Control, X-Requested-With, X-Request-Source"},
		/*
			AllowOriginFunc: func(origin string) bool {
				return origin == "https://bar.com"
			},
		*/
		MaxAge: 12 * time.Hour,
	}))

	// RestAPI service's router
	apiV0 := router.Group(appconfig.Service_path)
	{
		var tokenRequired bool = true
		if strings.Contains(appconfig.Base_URL, "127.0.0.1") {
			tokenRequired = false
		}

		// swagger
		apiV0.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		// health check
		apiV0.GET("/health", controllers.Health)
		apiV0.GET("/ready", controllers.Ready)
		apiV0.GET("/metrics", gin.WrapH(promhttp.Handler()))

		// pprof handlers
		apiV0.GET("/debug/pprof/", gin.WrapH(http.HandlerFunc(pprof.Index)))
		apiV0.GET("/debug/pprof/cmdline", gin.WrapH(http.HandlerFunc(pprof.Cmdline)))
		apiV0.GET("/debug/pprof/profile", gin.WrapH(http.HandlerFunc(pprof.Profile)))
		apiV0.GET("/debug/pprof/symbol", gin.WrapH(http.HandlerFunc(pprof.Symbol)))
		apiV0.GET("/debug/pprof/trace", gin.WrapH(http.HandlerFunc(pprof.Trace)))
		apiV0.GET("/debug/pprof/allocs", gin.WrapH(pprof.Handler("allocs")))
		apiV0.GET("/debug/pprof/mutex", gin.WrapH(pprof.Handler("mutex")))
		apiV0.GET("/debug/pprof/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		apiV0.GET("/debug/pprof/heap", gin.WrapH(pprof.Handler("heap")))
		apiV0.GET("/debug/pprof/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
		apiV0.GET("/debug/pprof/block", gin.WrapH(pprof.Handler("block")))

		// cameras
		apiV0.POST("/cameras", handleWrapper(controllers.CreateCamera, tokenRequired))
		apiV0.GET("/cameras", handleWrapper(controllers.GetCameras, tokenRequired))
		apiV0.GET("/cameras/:id", handleWrapper(controllers.ReadCamera, tokenRequired))
		apiV0.PUT("/cameras/:id", handleWrapper(controllers.UpdateCamera, tokenRequired))
		apiV0.DELETE("/cameras/:id", handleWrapper(controllers.DeleteCamera, tokenRequired))
		apiV0.GET("/cameras/serialnumber/:serial", handleWrapper(controllers.ReadSerialCamera, tokenRequired))
		apiV0.POST("/cameras/import", handleWrapper(controllers.ImportCamerasFromCSV, tokenRequired))
		apiV0.GET("/cameras/sample-data", handleWrapper(controllers.GenerateSampleCameraData, tokenRequired))
		apiV0.GET("/cameras/imported/download", handleWrapper(controllers.DownloadImportedCameras, tokenRequired))
		apiV0.PUT("/cameras/config/changepasswordseries", handleWrapper(controllers.ChangePasswordCameraSeries, tokenRequired))

		apiV0.GET("/cameras/options/combox", handleWrapper(controllers.GetCameraCombox, tokenRequired))
		apiV0.GET("/cameras/options/config-types", handleWrapper(controllers.GetCameraConfigTypeOptions, tokenRequired))
		apiV0.GET("/cameras/options/protocol-types", handleWrapper(controllers.GetCameraProtocolTypes, tokenRequired))
		apiV0.GET("/cameras/options/stream-types", handleWrapper(controllers.GetCameraStreamTypes, tokenRequired))
		apiV0.GET("/cameras/options/types", handleWrapper(controllers.GetCameraTypes, tokenRequired))
		apiV0.POST("/cameras/set-config/all", handleWrapper(controllers.SetCameraConfigAll, tokenRequired))
		apiV0.POST("/cameras/set-config/batch", handleWrapper(controllers.SetCameraConfigBatch, tokenRequired))
		apiV0.POST("/cameras/set-config/:id", handleWrapper(controllers.SetCameraConfig, tokenRequired))
		apiV0.POST("/cameras/get-config/all", handleWrapper(controllers.GetCameraConfigAll, tokenRequired))
		apiV0.POST("/cameras/get-config/batch", handleWrapper(controllers.GetCameraConfigBatch, tokenRequired))
		apiV0.GET("/cameras/get-config/:id", handleWrapper(controllers.GetCameraConfig, tokenRequired))
		apiV0.PUT("/cameras/:id/device", handleWrapper(controllers.ChangeCameraDevice, tokenRequired))

		apiV0.GET("/cameras/config/networkconfig/:id", handleWrapper(controllers.GetNetworkConfigCamera, tokenRequired))
		apiV0.PUT("/cameras/config/networkconfig/:id", handleWrapper(controllers.UpdateNetworkConfigCamera, tokenRequired))
		apiV0.PUT("/cameras/config/networkconfig/update", handleWrapper(controllers.UpdateNetworkConfigCameras, tokenRequired))
		apiV0.GET("/cameras/config/imageconfig/:idCamera", handleWrapper(controllers.GetImageConfigCamera, tokenRequired))
		apiV0.PUT("/cameras/config/imageconfig/:id", handleWrapper(controllers.UpdateImageConfigCamera, tokenRequired))
		apiV0.PUT("/cameras/config/imageconfigs", handleWrapper(controllers.UpdateImageConfigCameras, tokenRequired))
		apiV0.PUT("/cameras/config/changepassword", handleWrapper(controllers.ChangePasswordCamera, tokenRequired))
		apiV0.GET("/cameras/config/videoconfig/:id", handleWrapper(controllers.GetVideoConfigCamera, tokenRequired))
		apiV0.PUT("/cameras/config/videoconfig/:id", handleWrapper(controllers.UpdateVideoConfigCamera, tokenRequired))
		apiV0.PUT("/cameras/config/videoconfigs", handleWrapper(controllers.UpdateVideoConfigCameras, tokenRequired))

		apiV0.GET("/cameras/user/:id", handleWrapper(controllers.ReadUserCameraGroup, tokenRequired))
		apiV0.GET("/cameras/user", handleWrapper(controllers.GetUserCameras, tokenRequired))
		apiV0.POST("/cameras/user", handleWrapper(controllers.CreateUserCameraGroup, tokenRequired))
		apiV0.PUT("/cameras/user/:id", handleWrapper(controllers.UpdateUserCameraGroup, tokenRequired))
		apiV0.GET("/cameras/snapshot/:id", handleWrapper(controllers.GetCameraSnapshot, tokenRequired))

		apiV0.GET("/cameras/user/ai-properties", handleWrapper(controllers.GetUserCameraAIProperty, tokenRequired))
		apiV0.GET("/cameras/user/ai-properties/:id", handleWrapper(controllers.ReadUserCameraAIProperty, tokenRequired))
		apiV0.GET("/cameras/user/ai-properties/camera/:camera_id", handleWrapper(controllers.ReadUserCameraAIPropertybyCameraID, tokenRequired))
		apiV0.POST("/cameras/user/ai-properties", handleWrapper(controllers.CreateUserCameraAIProperty, tokenRequired))
		apiV0.PUT("/cameras/user/ai-properties/:id", handleWrapper(controllers.UpdateUserCameraAIProperty, tokenRequired))
		apiV0.PUT("/cameras/user/ai-properties/cameras", handleWrapper(controllers.UpdateUserCamerasAIProperties, tokenRequired))
		apiV0.GET("/device/nvr/synchronize", handleWrapper(controllers.SynchronizeNVRDataConfigForNVR, tokenRequired))
		apiV0.GET("/device/camera/synchronize", handleWrapper(controllers.SynchronizeCameraDataConfigForCamera, tokenRequired))

		// NVRs
		apiV0.POST("/nvrs", handleWrapper(controllers.CreateNVR, tokenRequired))
		apiV0.GET("/nvrs/:id", handleWrapper(controllers.ReadNVR, tokenRequired))
		apiV0.PUT("/nvrs/:id", handleWrapper(controllers.UpdateNVR, tokenRequired))
		apiV0.DELETE("/nvrs/:id", handleWrapper(controllers.DeleteNVR, tokenRequired))
		apiV0.GET("/nvrs", handleWrapper(controllers.GetNVRs, tokenRequired))
		apiV0.POST("/nvrs/set-config/all", handleWrapper(controllers.SetNVRConfigAll, tokenRequired))
		apiV0.POST("/nvrs/set-config/batch", handleWrapper(controllers.SetNVRConfigBatch, tokenRequired))
		apiV0.POST("/nvrs/set-config/:id", handleWrapper(controllers.SetNVRConfig, tokenRequired))
		apiV0.POST("/nvrs/get-config/all", handleWrapper(controllers.GetNVRConfigAll, tokenRequired))
		apiV0.POST("/nvrs/get-config/batch", handleWrapper(controllers.GetNVRConfigBatch, tokenRequired))
		apiV0.GET("/nvrs/get-config/:id", handleWrapper(controllers.GetNVRConfig, tokenRequired))
		apiV0.GET("/nvrs/options/config-types", handleWrapper(controllers.GetNVRConfigTypeOptions, tokenRequired))
		apiV0.GET("/nvrs/options/protocol-types", handleWrapper(controllers.GetNVRProtocolTypes, tokenRequired))
		apiV0.GET("/nvrs/options/stream-types", handleWrapper(controllers.GetNVRStreamTypes, tokenRequired))
		apiV0.GET("/nvrs/options/types", handleWrapper(controllers.GetNVRTypes, tokenRequired))
		apiV0.POST("/nvrs/import", handleWrapper(controllers.ImportNVRsFromCSV, tokenRequired))
		apiV0.GET("/nvrs/sample-data", handleWrapper(controllers.GenerateSampleNVRData, tokenRequired))
		apiV0.GET("/nvrs/imported/download", handleWrapper(controllers.DownloadImportedNVRs, tokenRequired))

		apiV0.GET("/nvrs/config/networkconfig/:idNVR", handleWrapper(controllers.GetNetworkConfig, tokenRequired))
		apiV0.PUT("/nvrs/config/networkconfig/:idNetWorkConfig", handleWrapper(controllers.UpdateNetworkConfig, tokenRequired))
		apiV0.GET("/nvrs/config/imageconfig/:idNVR", handleWrapper(controllers.GetImageConfigNVR, tokenRequired))
		apiV0.PUT("/nvrs/config/imageconfigNVR/:idNVR", handleWrapper(controllers.UpdateImageConfigNVR, tokenRequired))
		apiV0.PUT("/nvrs/config/imageconfig/:idimageconfig", handleWrapper(controllers.UpdateImageConfig, tokenRequired))
		apiV0.PUT("/nvrs/config/changepassword", handleWrapper(controllers.ChangePassword, tokenRequired))
		apiV0.PUT("/nvrs/camera/:id", handleWrapper(controllers.AddCameraToNVR, tokenRequired))
		apiV0.GET("/nvrs/:id/cameras", handleWrapper(controllers.GetAttachedCamera, tokenRequired))
		apiV0.PUT("/nvrs/config/videoconfig/:id", handleWrapper(controllers.UpdateVideoConfigNVR, tokenRequired))

		apiV0.GET("/nvrs/config/videoconfigNVR/:id", handleWrapper(controllers.GetVideoConfigNVR, tokenRequired))
		//Camera-Configs

		//Scan NVRs

		// camera-groups
		apiV0.POST("/camera-groups", handleWrapper(controllers.CreateCameraGroup, tokenRequired))
		apiV0.GET("/camera-groups/:id", handleWrapper(controllers.ReadCameraGroup, tokenRequired))
		apiV0.PUT("/camera-groups/:id", handleWrapper(controllers.UpdateCameraGroup, tokenRequired))
		apiV0.DELETE("/camera-groups/:id", handleWrapper(controllers.DeleteCameraGroup, tokenRequired))
		apiV0.GET("/camera-groups", handleWrapper(controllers.GetCameraGroups, tokenRequired))
		apiV0.GET("/camera-groups/options/sync-cabin", handleWrapper(controllers.GetCameraTypes, tokenRequired))
		apiV0.GET("/camera-groups/cameras/:id", handleWrapper(controllers.GetCameraFromGroups, tokenRequired))

		// media-libraries
		apiV0.POST("/media-libraries", handleWrapper(controllers.CreateMediaLibrary, tokenRequired))
		apiV0.GET("/media-libraries/:id", handleWrapper(controllers.ReadMediaLibrary, tokenRequired))
		apiV0.PUT("/media-libraries/:id", handleWrapper(controllers.UpdateMediaLibrary, tokenRequired))
		apiV0.DELETE("/media-libraries/:id", handleWrapper(controllers.DeleteMediaLibrary, tokenRequired))
		apiV0.GET("/media-libraries", handleWrapper(controllers.GetMediaLibraries, tokenRequired))
		apiV0.GET("/media-libraries/options/file-types", handleWrapper(controllers.GetMediaFileTypes, tokenRequired))

		// camera-groups
		apiV0.POST("/recording-schedules", handleWrapper(controllers.CreateRecordingSchedule, tokenRequired))
		apiV0.GET("/recording-schedules/:id", handleWrapper(controllers.ReadRecordingSchedule, tokenRequired))
		apiV0.PUT("/recording-schedules/:id", handleWrapper(controllers.UpdateRecordingSchedule, tokenRequired))
		apiV0.DELETE("/recording-schedules/:id", handleWrapper(controllers.DeleteRecordingSchedule, tokenRequired))
		apiV0.GET("/recording-schedules", handleWrapper(controllers.GetRecordingSchedules, tokenRequired))

		// cctv-events
		apiV0.POST("/cctv-events", handleWrapper(controllers.CreateCCTVEvent, tokenRequired))
		apiV0.GET("/cctv-events/:id", handleWrapper(controllers.ReadCCTVEvent, tokenRequired))
		apiV0.PUT("/cctv-events/:id", handleWrapper(controllers.UpdateCCTVEvent, tokenRequired))
		apiV0.DELETE("/cctv-events/:id", handleWrapper(controllers.DeleteCCTVEvent, tokenRequired))
		apiV0.GET("/cctv-events", handleWrapper(controllers.GetCCTVEvents, tokenRequired))
		apiV0.GET("/cctv-events/options/device-types", handleWrapper(controllers.GetCCTVDeviceTypes, tokenRequired))
		apiV0.GET("/cctv-events/options/types", handleWrapper(controllers.GetCCTVEventTypes, tokenRequired))
		apiV0.GET("/cctv-events/options/status-types", handleWrapper(controllers.GetCCTVEventStatusTypes, tokenRequired))
		apiV0.GET("/cctv-events/options/level-types", handleWrapper(controllers.GetCCTVEventLevelTypes, tokenRequired))

		// ai-events
		apiV0.POST("/ai-events", handleWrapper(controllers.CreateAIEvent, tokenRequired))
		apiV0.GET("/ai-events/:id", handleWrapper(controllers.ReadAIEvent, tokenRequired))
		apiV0.PUT("/ai-events/:id", handleWrapper(controllers.UpdateAIEvent, tokenRequired))
		apiV0.DELETE("/ai-events/:id", handleWrapper(controllers.DeleteAIWarning, tokenRequired))       // Todo: need to refactor ai-events vs ai-waring ...
		apiV0.DELETE("/aievents/delete/:id", handleWrapper(controllers.DeleteAIWarning, tokenRequired)) // Todo: need to remove this from FE
		apiV0.GET("/ai-events", handleWrapper(controllers.GetAIEvents, tokenRequired))
		apiV0.GET("/ai-events/options/device-types", handleWrapper(controllers.GetAIDeviceTypes, tokenRequired))
		apiV0.GET("/ai-events/options/types", handleWrapper(controllers.GetAIEventTypes, tokenRequired))
		apiV0.GET("/ai-events/options/status-types", handleWrapper(controllers.GetAIEventStatusTypes, tokenRequired))
		apiV0.GET("/ai-events/options/level-types", handleWrapper(controllers.GetAIEventLevelTypes, tokenRequired))
		apiV0.POST("/aievents/camerastatus", handleWrapper(controllers.CameraStatus, tokenRequired))
		apiV0.GET("/aievents", handleWrapper(controllers.GetAIEvent, tokenRequired))
		apiV0.GET("/aievents/genimage", handleWrapper(controllers.GetAIEventGenImage, tokenRequired))
		apiV0.GET("/aievents/searching", handleWrapper(controllers.SearchAIEvent, tokenRequired))
		apiV0.GET("/aievents/routine", handleWrapper(controllers.AIEventRoutineHandler, tokenRequired))
		apiV0.POST("/aievents", handleWrapper(controllers.CreateAIWarning, tokenRequired))
		apiV0.GET("/aievents/:id", handleWrapper(controllers.ReadAIWarning, tokenRequired))
		apiV0.PUT("/aievents/:id", handleWrapper(controllers.UpdateAIWarning, tokenRequired))
		apiV0.DELETE("/aievents/:id", handleWrapper(controllers.DeleteAIWarning, tokenRequired))

		// WebRTC websocket signaling
		apiV0.GET("/ws/signaling/:id", controllers.ServeWS)
		apiV0.GET("/websocket/topic/:topic/:client-uuid", controllers.WSServe)
		apiV0.GET("/ws/signaling/video/:id", controllers.ServeWSS)

		// WS connection for ai-engine agent
		apiV0.GET("/connection/websocket", gin.WrapH(wsnode.WebsocketHandler))

		//http.Handle("/connection/websocket", websocketHandler)

		//data mockup
		apiV0.GET("/datamockup/tcpip", handleWrapper(controllers.GetTCPIP, false))
		apiV0.GET("/datamockup/ddns", handleWrapper(controllers.DDNS, false))
		apiV0.GET("/datamockup/port", handleWrapper(controllers.Port, false))
		apiV0.GET("/datamockup/ntp", handleWrapper(controllers.NTP, false))
		apiV0.GET("/datamockup/osd", handleWrapper(controllers.OSD, false))

		// events
		apiV0.POST("/events", handleWrapper(controllers.CreateEvent, tokenRequired))
		apiV0.GET("/events/:id", handleWrapper(controllers.ReadEvent, tokenRequired))
		apiV0.GET("/events", handleWrapper(controllers.GetEvents, tokenRequired))
		apiV0.GET("/events/list-event-types", handleWrapper(controllers.GetListTypeEvent, tokenRequired))
		apiV0.PUT("/events/:id/description", handleWrapper(controllers.UpdateDescription, tokenRequired))
		apiV0.PUT("/events/:id/update-event-status", handleWrapper(controllers.UpdateStatus, tokenRequired))
		apiV0.GET("/events/filter", handleWrapper(controllers.GetEventFilter, tokenRequired))
		apiV0.PUT("/events/:id/update-image", handleWrapper(controllers.UpdateImageURL, tokenRequired))

		//blacklist
		apiV0.POST("/blacklist", handleWrapper(controllers.CreateBlacklist, tokenRequired))
		apiV0.PUT("/blacklist/:id", handleWrapper(controllers.UpdateName, tokenRequired))
		apiV0.GET("/blacklist/:id", handleWrapper(controllers.ReadBlackList, tokenRequired))
		apiV0.DELETE("/blacklist/:id", handleWrapper(controllers.DeleteBlackList, tokenRequired))
		apiV0.DELETE("/blacklist", handleWrapper(controllers.DeleteBlackLists, tokenRequired))
		apiV0.POST("/blacklist/search", handleWrapper(controllers.SearchBlacklist, tokenRequired))
		apiV0.GET("/blacklist", handleWrapper(controllers.GetBlackList, tokenRequired))
		apiV0.GET("/blacklist/images/:member_id", handleWrapper(controllers.GetBlacklistImage, tokenRequired))
		apiV0.POST("/images/upload", handleWrapper(controllers.UploadImage, tokenRequired))
		//device
		apiV0.GET("/device/hik", handleWrapper(controllers.ReadNetworkInterfaces, tokenRequired))
		apiV0.GET("/device/:deviceID/cameras", handleWrapper(controllers.GetCamerasByDeviceID, tokenRequired))
		apiV0.PUT("/device/hik/:id/:command", handleWrapper(controllers.UpdateNetworkInterfaces, tokenRequired))
		apiV0.POST("/device/onvif/scandevice", handleWrapper(controllers.DeviceScanOnvif, tokenRequired))
		apiV0.GET("/device/active", handleWrapper(controllers.GetActiveDevices, tokenRequired))
		apiV0.POST("/device/scandevice", handleWrapper(controllers.DeviceScan, tokenRequired))
		apiV0.POST("/device/onvif/scandevicestaticip", handleWrapper(controllers.DeviceScanStaticIP, tokenRequired))
		apiV0.POST("/device/onvif/devicescanlistip", handleWrapper(controllers.DeviceScanIPList, tokenRequired))
		apiV0.POST("/device/onvif/devicescanhikivision", handleWrapper(controllers.DeviceScanHikivision, tokenRequired))
		apiV0.GET("/device/nvrs/synchronize", handleWrapper(controllers.SynchronizeNVRDataConfig, tokenRequired))
		apiV0.GET("/device/cameras/synchronize", handleWrapper(controllers.SynchronizeCameraDataConfig, tokenRequired))

		//Incidents
		apiV0.POST("/incidents/logs", handleWrapper(controllers.CreateSystemIncidentLog, tokenRequired))
		apiV0.GET("/incidents/logs/:id", handleWrapper(controllers.ReadSystemIncidentLog, tokenRequired))
		apiV0.PUT("/incidents/logs/:id", handleWrapper(controllers.UpdateSystemIncidentLog, tokenRequired))
		apiV0.GET("/incidents/logs", handleWrapper(controllers.GetSystemIncidentLogs, tokenRequired))
		apiV0.PUT("/incidents/logs/status/:id", handleWrapper(controllers.UpdateSystemIncidentLogStatus, tokenRequired))
		apiV0.DELETE("/incidents/:id", handleWrapper(controllers.DeleteSystemIncident, tokenRequired))

		//device box

		apiV0.GET("/device", handleWrapper(controllers.GetDevices, tokenRequired))
		apiV0.PUT("/device/:id", handleWrapper(controllers.UpdateDevice, tokenRequired))
		apiV0.POST("/device", handleWrapper(controllers.CreateDevice, tokenRequired))
		apiV0.GET("/device/search/:id", handleWrapper(controllers.ReadDevice, tokenRequired))
		apiV0.DELETE("/device/:id", handleWrapper(controllers.DeleteDevice, tokenRequired))
		//eventconfig
		apiV0.POST("/event-config", handleWrapper(controllers.CreateEventconfig, tokenRequired))
		apiV0.GET("/event-config/:id", handleWrapper(controllers.ReadEventconfig, tokenRequired))
		apiV0.PUT("/event-config/:id", handleWrapper(controllers.UpdateEventconfig, tokenRequired))
		apiV0.DELETE("/event-config/:id", handleWrapper(controllers.DeleteEventconfig, tokenRequired))
		apiV0.POST("/event-config/scan", handleWrapper(controllers.PostScanEventModel, tokenRequired))

		//videoconfigNVR
		apiV0.PUT("/nvrs/config/videoconfig", handleWrapper(controllers.UpdateVideoConfig, tokenRequired))
		apiV0.GET("/nvrs/config/videoconfig/:id", handleWrapper(controllers.GetVideoConfig, tokenRequired))

		//LicensePlates
		apiV0.POST("/licenseplates", handleWrapper(controllers.CreateLicensePlates, tokenRequired))
		apiV0.PUT("/licenseplates/:id", handleWrapper(controllers.UpdateLicensePlates, tokenRequired))
		apiV0.GET("/licenseplates/:id", handleWrapper(controllers.ReadLicensePlates, tokenRequired))
		apiV0.GET("/licenseplates", handleWrapper(controllers.GetLicensePlates, tokenRequired))
		apiV0.DELETE("/licenseplates/:id", handleWrapper(controllers.DeleteLicensePlates, tokenRequired))
		apiV0.DELETE("/licenseplates", handleWrapper(controllers.DeleteLicensePlatess, tokenRequired))

		//Camera model AI
		apiV0.POST("/camera-model-ai", handleWrapper(controllers.CreateCameraModelAI, tokenRequired))
		apiV0.PUT("/camera-model-ai/:id", handleWrapper(controllers.UpdateCameraModelAI, tokenRequired))
		apiV0.GET("/camera-model-ai/:id", handleWrapper(controllers.ReadCameraModelAI, tokenRequired))
		apiV0.GET("/camera-model-ai", handleWrapper(controllers.GetCameraModelAIs, tokenRequired))
		apiV0.DELETE("/camera-model-ai/:id", handleWrapper(controllers.DeleteCameraModelAI, tokenRequired))
		apiV0.GET("/camera-model-ai/count-camera-of-model-ai", handleWrapper(controllers.GetCountOfModelAI, tokenRequired))

		//Camera playback
		apiV0.GET("/playback/camera/:id", handleWrapper(controllers.GetStoragePlayback, tokenRequired))
		apiV0.POST("/playback/camera/:id", handleWrapper(controllers.GetCalenderPlayback, tokenRequired))

		apiV0.POST("/pc-info", handleWrapper(controllers.CreatePCInfo, tokenRequired))

		//Download video
		apiV0.GET("/downloadVideo", handleWrapper(controllers.DownloadVideo, tokenRequired))
		apiV0.GET("/video/download", handleWrapper(controllers.DownLoadVideo, tokenRequired))
		apiV0.GET("/video/downloadchoice", handleWrapper(controllers.ExtractClipWithChoice, tokenRequired))

	}

	// Setup the HTTP Server
	srv = &http.Server{
		Addr:    appconfig.Listen_addr,
		Handler: router,
	}

	// Kick off Graceful Shutdown Go Routine
	go func() {
		// Make the Trap
		trap := make(chan os.Signal, 1)
		signal.Notify(trap, syscall.SIGTERM)

		// Wait for a signal then action
		s := <-trap
		log.Infof("Received shutdown signal %s", s)

		defer Stop()
	}()

	// Start HTTP Listener
	log.Infof("Starting HTTP Listener on %s, service path is %s", appconfig.Listen_addr, appconfig.Service_path)
	err = srv.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			// Wait until all outstanding requests are done
			<-runCtx.Done()
			return ErrShutdown
		}
		return fmt.Errorf("unable to start HTTP Server - %s", err)
	}

	// Stop websoket node
	wsnode.WaitExitSignal()

	return nil
}

// Stop is used to gracefully shutdown the server.
func Stop() {
	err := srv.Shutdown(context.Background())
	if err != nil {
		log.Errorf("Unexpected error while shutting down HTTP server - %s", err)
	}
	defer runCancel()
}

// isPProf is a regex that validates if the given path is used for PProf
var isPProf = regexp.MustCompile(`.*debug\/pprof.*`)

// middleware is used to intercept incoming HTTP calls and apply general functions upon
// them. e.g. Metrics, Logging...
func handleWrapper(n gin.HandlerFunc, tokenRequired bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()

		// Set the Tarmac server response header
		//w.Header().Set("Server", "tarmac")

		// Log the basics
		log.WithFields(logrus.Fields{
			"method":         c.Request.Method,
			"remote-addr":    c.Request.RemoteAddr,
			"http-protocol":  c.Request.Proto,
			"headers":        c.Request.Header,
			"content-length": c.Request.ContentLength,
		}).Debugf("HTTP Request to %s", c.Request.URL)

		// Verify if PProf
		if isPProf.MatchString(c.Request.URL.Path) && !cfg.GetBool("enable_pprof") {
			log.WithFields(logrus.Fields{
				"method":         c.Request.Method,
				"remote-addr":    c.Request.RemoteAddr,
				"http-protocol":  c.Request.Proto,
				"headers":        c.Request.Header,
				"content-length": c.Request.ContentLength,
			}).Debugf("Request to PProf Address failed, PProf disabled")
			c.AbortWithStatus(http.StatusForbidden)

			stats.Srv.WithLabelValues(c.Request.URL.Path).Observe(time.Since(now).Seconds())
			return
		}

		if tokenRequired {
			if !auth.ValidateToken(c, appconfig.IWK_set_uri) {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		// Call actual handler
		n(c)
		stats.Srv.WithLabelValues(c.Request.URL.Path).Observe(time.Since(now).Seconds())
	}
}
