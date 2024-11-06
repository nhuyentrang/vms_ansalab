package models

import (
	"time"

	uuid "github.com/google/uuid"
	//"github.com/lib/pq"
)

// Entity Model for camera device
type NVR struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string    `json:"name"`
	Description string    `json:"description"`

	// General info
	Status           KeyValue `json:"status" gorm:"embedded;type:jsonb"`
	Type             KeyValue `json:"type" gorm:"embedded;type:jsonb"`
	Protocol         string   `json:"protocol" gorm:"column:protocol"`
	Model            string   `json:"model" gorm:"column:model"`
	IPAddress        string   `json:"ipAddress" gorm:"column:ip_address"`
	MACAddress       string   `json:"macAddress" gorm:"column:mac_address"`
	HttpPort         string   `json:"httpPort" gorm:"column:http_port"`
	RtspPort         string   `json:"rtspPort" gorm:"column:rtsp_port"`
	OnvifPort        string   `json:"onvifPort" gorm:"column:onvif_port"`
	ManagementPort   string   `json:"managementPort" gorm:"column:management_port"`
	Username         string   `json:"username" gorm:"column:username"`
	Password         string   `json:"password" gorm:"column:password"`
	FirmwareVersion  string   `json:"firmwareVersion"`
	UseTLS           bool     `json:"useTLS"`
	IsOfflineSetting *bool    `json:"isOfflineSetting"`
	IFrameURL        string   `json:"iframeURL"`

	LastImageURL string `json:"lastImageURL" gorm:"column:last_image_url"`

	// NVR group / location info
	SiteInfo   KeyValue `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Location   string   `json:"location" gorm:"column:location"`
	Coordinate string   `json:"coordinate" gorm:"column:coordinate"`
	Position   string   `json:"position" gorm:"column:position"`

	// Added cameras
	Cameras *KeyValueArray `json:"cameras" gorm:"embedded;type:jsonb"`
	Box     KeyValue       `json:"box" gorm:"column:box;embedded;type:jsonb"`

	// NVR config struct (foreign key)
	ConfigID uuid.UUID `json:"configID" gorm:"column:config_id;type:uuid;default:uuid_generate_v4()"`

	// Device's metadata, there are should be key/value for management I/O, relay, rs232, rs485...
	Metadata KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`

	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	LastPing  time.Time `json:"lastPing"`
}

type NVRConfig struct {
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
}

type DTO_NVR struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name,omitempty"`

	// General
	Type             KeyValue `json:"type,omitempty"`
	Protocol         string   `json:"	,omitempty"`
	Model            string   `json:"model,omitempty"`
	FirmwareVersion  string   `json:"firmwareVersion,omitempty"`
	IPAddress        string   `json:"ipAddress,omitempty"`
	MACAddress       string   `json:"macAddress,omitempty"`
	HttpPort         string   `json:"httpPort,omitempty"`
	OnvifPort        string   `json:"onvifPort,omitempty"`
	RtspPort         string   `json:"rtspPort,omitempty"`
	ManagementPort   string   `json:"managementPort,omitempty"`
	Username         string   `json:"username,omitempty"`
	Password         string   `json:"password,omitempty"`
	UseTLS           bool     `json:"useTLS,omitempty"`
	IsOfflineSetting *bool    `json:"isOfflineSetting"`
	IFrameURL        string   `json:"iframeURL,omitempty"`

	Status KeyValue `json:"status,omitempty"`

	// Added cameras
	Cameras     *KeyValueArray `json:"cameras,omitempty" gorm:"embedded;type:jsonb"`
	NumberOfCam int            `json:"numberOfCam,omitempty"`
	Box         KeyValue       `json:"box" gorm:"embedded;type:jsonb"`

	// NVR group / location info
	SiteInfo   KeyValue `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Location   string   `json:"location" gorm:"column:location"`
	Coordinate string   `json:"coordinate" gorm:"column:coordinate"`
	Position   string   `json:"position" gorm:"column:position"`

	// NVR config struct (foreign key)
	ConfigID uuid.UUID `json:"configID" gorm:"column:config_id"`
	LastPing time.Time `json:"lastPing"`
}

type DTO_NVR_Read_BasicInfo struct {
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

	Status KeyValue `json:"status,omitempty"`

	// Stream/recording info
	Cameras     KeyValueArray `json:"cameras,omitempty" gorm:"embedded;type:jsonb"`
	NumberOfCam int           `json:"numberOfCam,omitempty"`

	// NVR group / location info
	GroupIDs   KeyValueArray `json:"groupIDs" gorm:"column:group_id;embedded;type:jsonb"`
	SiteInfo   KeyValue      `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Location   string        `json:"location" gorm:"column:location"`
	Coordinate string        `json:"coordinate" gorm:"column:coordinate"`
	Position   string        `json:"position" gorm:"column:position"`
	LastPing   time.Time     `json:"lastPing"`
}

// Đổi mật khẩu
type DTO_ChangePassword struct {
	ID          string `json:"id,omitempty"`
	UserName    string `json:"username,omitempty"`
	PasswordOld string `json:"passwordOld,omitempty"`
	PasswordNew string `json:"passwordNew,omitempty"`
	IPAddress   string `json:"ipAddress,omitempty"`
	HttpPort    string `json:"httpPort,omitempty"`
	Status      bool   `json:"status,omitempty"`
}

type DTO_NVR_DeviceConfig struct {
	ID      uuid.UUID     `json:"id"`
	Configs KeyValueArray `json:"configs,omitempty"`
}

type DTO_NVR_DeviceConfigBatch struct {
	NVRIDs  []uuid.UUID   `json:"nvrIDs"`
	Configs KeyValueArray `json:"configs,omitempty"`
}

type DTO_NVRConfig struct {
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
