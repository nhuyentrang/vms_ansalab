package models

import (
	"time"

	uuid "github.com/google/uuid"
	//"github.com/lib/pq"
)

// Entity Model for camera device
type Camera struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string    `json:"name"`
	Description string    `json:"description"`

	// General info
	Status                  KeyValue `json:"status" gorm:"embedded;type:jsonb"`
	Type                    KeyValue `json:"type" gorm:"embedded;type:jsonb"`
	Protocol                string   `json:"protocol" gorm:"column:protocol"`
	Model                   string   `json:"model" gorm:"column:model"`
	SerialNumber            string   `json:"serial" gorm:"column:serial_number"`
	IPAddress               string   `json:"ipAddress" gorm:"column:ip_address"`
	MACAddress              string   `json:"macAddress" gorm:"column:mac_address"`
	HttpPort                string   `json:"httpPort" gorm:"column:http_port"`
	RTSPPort                string   `json:"rtspPort" gorm:"column:rtsp_port"`
	OnvifPort               string   `json:"onvifPort" gorm:"column:onvif_port"`
	ManagementPort          string   `json:"managementPort" gorm:"column:management_port"`
	Username                string   `json:"username" gorm:"column:username"`
	Password                string   `json:"password" gorm:"column:password"`
	FirmwareVersion         string   `json:"firmwareVersion"`
	UseTLS                  bool     `json:"useTLS"`
	IsOfflineSetting        *bool    `json:"isOfflineSetting" gorm:"type:boolean;column:is_offline_setting"`
	IFrameURL               string   `json:"iframeURL"`
	LastImageURL            string   `json:"lastImageURL" gorm:"column:last_image_url"`
	FaceRecognition         bool     `json:"faceRecognition" gorm:"column:face_recognition;default:false"`
	LicensePlateRecognition bool     `json:"licensePlateRecognition" gorm:"column:license_plate_recognition;default:false"`

	// Camera group / location info
	GroupIDs   KeyValueArray `json:"groupIDs" gorm:"column:group_id;embedded;type:jsonb"`
	SiteInfo   KeyValue      `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Location   string        `json:"location" gorm:"column:location"`
	Lat        string        `json:"lat" gorm:"column:lat"`
	Long       string        `json:"long" gorm:"column:long"`
	Coordinate string        `json:"coordinate" gorm:"column:coordinate"`
	Position   string        `json:"position" gorm:"column:position"`

	// Video channels / nvr info
	Streams    VideoStreamArray `json:"streams" gorm:"embedded;type:jsonb"`
	NVR        KeyValue         `json:"nvr" gorm:"column:nvr;embedded;type:jsonb"`
	Box        KeyValue         `json:"box" gorm:"column:box;embedded;type:jsonb"`
	IndexNVR   string           `json:"indexNVR" gorm:"column:index_nvr"`
	NameDevice string           `json:"nameDevice,omitempty"`
	// Todo: site thì đơn vị nhỏ nhất nên là một unit, khi tạo cam, chọn site, sẽ load ra map
	// hoặc bản vẽ tuỳ theo cấu hình site,  click vị trí trên map/bản vẽ sẽ hiển thị
	// decimal degrees hoặc toạ độ x/y vào 1 ô textbox
	// ví dụ: 41.40338, 2.17403 hoặc 123.23, 132.2

	// Camera config struct (foreign key)
	ConfigID uuid.UUID `json:"configID" gorm:"column:config_id;type:uuid;default:uuid_generate_v4()"`

	// Device's metadata, there are should be key/value for management I/O, relay, rs232, rs485...
	Metadata KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`

	// Camera has one CameraAIEventProperty, camera_id is the foreign key
	//CameraAIEventProperty []*CameraAIEventProperty `gorm:"foreignKey:camera_id"`

	// AIEngine back-referencer table "cam_ai_engine_planning"
	AIEngines []*AIEngine `gorm:"many2many:cam_ai_engine_planning;"`

	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	LastPing  time.Time `json:"lastPing"`
}

type CameraConfig struct {
	// Cấu hình hình ảnh, video, mạng, lưu trữ, ptz, cctv/ai event...
	// Cấu hình p2p, liveview, playback, AI (model, server, enpoint...), ROI, event, log...
	// Có thể nhiều camera chung cấu hình (cấu hình hàng loạt...)
	ID uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	// Image config
	ImageConfigID uuid.UUID `json:"imageConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Video config
	VideoConfigID uuid.UUID `json:"videoConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Audio config
	AudioConfigID uuid.UUID `json:"audioConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Network config
	NetworkConfigID uuid.UUID `json:"networkConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Storage config
	StorageConfigID uuid.UUID `json:"storageConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Recording schedule
	RecordingScheduleID uuid.UUID `json:"recordingScheduleID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Streaming config
	StreamingConfigID uuid.UUID `json:"streamingConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// AI config
	AIConfigID uuid.UUID `json:"aiConfigID" gorm:"column:ai_config_id;type:uuid;default:uuid_generate_v4()"`
	// PTZ config
	PTZConfigID uuid.UUID `json:"ptzConfigID" gorm:"column:ptz_config_id;type:uuid;default:uuid_generate_v4()"`
}

type DTOCamera struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name,omitempty"`
	IndexNVR string    `json:"indexNVR"`

	// General
	Type             KeyValue `json:"type,omitempty"`
	Protocol         string   `json:"protocol,omitempty"`
	Model            string   `json:"model,omitempty"`
	SerialNumber     string   `json:"serial"`
	FirmwareVersion  string   `json:"firmwareVersion,omitempty"`
	IPAddress        string   `json:"ipAddress,omitempty"`
	MACAddress       string   `json:"macAddress,omitempty"`
	HttpPort         string   `json:"httpPort,omitempty"`
	OnvifPort        string   `json:"onvifPort,omitempty"`
	RTSPPort         string   `json:"rtspPort,omitempty"`
	ManagementPort   string   `json:"managementPort,omitempty"`
	Username         string   `json:"username,omitempty"`
	Password         string   `json:"password,omitempty"`
	UseTLS           bool     `json:"useTLS,omitempty"`
	IsOfflineSetting *bool    `json:"isOfflineSetting" gorm:"type:boolean;column:is_offline_setting"`
	IFrameURL        string   `json:"iframeURL,omitempty"`
	Lat              string   `json:"lat"`
	Long             string   `json:"long"`
	Status           KeyValue `json:"status,omitempty"`

	// Stream/recording info
	Streams    VideoStreamArray `json:"streams,omitempty" gorm:"embedded;type:jsonb"`
	NVR        KeyValue         `json:"nvr"`
	Box        KeyValue         `json:"box"`
	NameDevice string           `json:"nameDevice,omitempty"`
	// Camera group / location info
	GroupIDs                KeyValueArray `json:"groupIDs" gorm:"column:group_id;embedded;type:jsonb"`
	SiteInfo                KeyValue      `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Location                string        `json:"location" gorm:"column:location"`
	Coordinate              string        `json:"coordinate" gorm:"column:coordinate"`
	Position                string        `json:"position" gorm:"column:position"`
	FaceRecognition         bool          `json:"faceRecognition"`
	LicensePlateRecognition bool          `json:"licensePlateRecognition"`
	LastPing                time.Time     `json:"lastPing"`

	// Camera config struct (foreign key)
	ConfigID uuid.UUID `json:"configID"` //WHY IT CREATED RANDOMLY
}

type DTO_Camera_Serial struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name,omitempty"`
	SerialNumber string    `json:"serial"`
}

type DTO_Camera_Read_BasicInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name,omitempty"`

	// General
	Type             KeyValue `json:"type,omitempty"`
	Protocol         string   `json:"protocol,omitempty"`
	Model            string   `json:"model,omitempty"`
	FirmwareVersion  string   `json:"firmwareVersion,omitempty"`
	IPAddress        string   `json:"ipAddress,omitempty"`
	MACAddress       string   `json:"macAddress,omitempty"`
	HttpPort         string   `json:"httpPort,omitempty"`
	OnvifPort        string   `json:"onvifPort,omitempty"`
	ManagementPort   string   `json:"managementPort,omitempty"`
	Username         string   `json:"username,omitempty"`
	Password         string   `json:"password,omitempty"`
	UseTLS           bool     `json:"useTLS,omitempty"`
	IsOfflineSetting bool     `json:"isOfflineSetting"`
	IFrameURL        string   `json:"iframeURL,omitempty"`
	Lat              string   `json:"lat"`
	Long             string   `json:"long"`
	Status           KeyValue `json:"status,omitempty"`

	// Stream/recording info
	Streams                 VideoStreamArray `json:"streams,omitempty"`
	NVR                     KeyValue         `json:"nvr"`
	FaceRecognition         bool             `json:"faceRecognition"`
	LicensePlateRecognition bool             `json:"licensePlateRecognition"`
	NameDevice              string           `json:"nameDevice,omitempty"`

	// Camera group / location info
	GroupIDs   KeyValueArray `json:"groupIDs" gorm:"column:group_id;embedded;type:jsonb"`
	SiteInfo   KeyValue      `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Location   string        `json:"location" gorm:"column:location"`
	Coordinate string        `json:"coordinate" gorm:"column:coordinate"`
	Position   string        `json:"position" gorm:"column:position"`
}

type DTOCameraDeviceConfig struct {
	ID      uuid.UUID     `json:"id"`
	Configs KeyValueArray `json:"configs,omitempty"`
}

type DTOCameraDeviceConfigBatch struct {
	CameraIDs []uuid.UUID   `json:"cameraIDs"`
	Configs   KeyValueArray `json:"configs,omitempty"`
}

type DTO_CameraConfig struct {
	// Cấu hình hình ảnh, video, mạng, lưu trữ, ptz, cctv/ai event...
	// Cấu hình p2p, liveview, playback, AI (model, server, enpoint...), ROI, event, log...
	// Có thể nhiều camera chung cấu hình (cấu hình hàng loạt...)
	ID uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	// Image config
	ImageConfigID uuid.UUID `json:"imageConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Video config
	VideoConfigID uuid.UUID `json:"videoConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Audio config
	AudioConfigID uuid.UUID `json:"audioConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Network config
	NetworkConfigID uuid.UUID `json:"networkConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Storage config
	StorageConfigID uuid.UUID `json:"storageConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// Streaming config
	StreamingConfigID uuid.UUID `json:"streamingConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// AI config
	AIConfigID uuid.UUID `json:"aiConfigID" gorm:"column:ai_config_id;type:uuid;default:uuid_generate_v4()"`
	// Recording schedule
	RecordingScheduleID uuid.UUID `json:"recordingScheduleID" gorm:"type:uuid;default:uuid_generate_v4()"`
	// PTZ config
	PTZConfigID uuid.UUID `json:"ptzConfigID" gorm:"column:ptz_config_id;type:uuid;default:uuid_generate_v4()"`
}

type DTO_Camera_ChangePassword struct {
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	IPAddress string `json:"ipAddress,omitempty"`
	HttpPort  string `json:"httpPort,omitempty"`
}
