package models

import (
	"time"

	uuid "github.com/google/uuid"
)

// Entity Model for NVR device
type NVRImport struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string    `json:"name"`
	Description string    `json:"description"`

	// General info
	Type                    KeyValue `json:"type" gorm:"embedded;type:jsonb"`
	Protocol                string   `json:"protocol" gorm:"column:protocol"`
	Model                   string   `json:"model" gorm:"column:model"`
	SerialNumber            string   `json:"serial" gorm:"column:serial_number"`
	IPAddress               string   `json:"ipAddress" gorm:"column:ip_address"`
	MACAddress              string   `json:"macAddress" gorm:"column:mac_address"`
	HttpPort                string   `json:"httpPort" gorm:"column:http_port"`
	OnvifPort               string   `json:"onvifPort" gorm:"column:onvif_port"`
	RtspPort                string   `json:"rtspPort,omitempty"`
	ManagementPort          string   `json:"managementPort" gorm:"column:management_port"`
	Username                string   `json:"username" gorm:"column:username"`
	Password                string   `json:"password" gorm:"column:password"`
	FirmwareVersion         string   `json:"firmwareVersion"`
	UseTLS                  bool     `json:"useTLS"`
	IsOfflineSetting        bool     `json:"isOfflineSetting"`
	IFrameURL               string   `json:"iframeURL"`
	LastImageURL            string   `json:"lastImageURL" gorm:"column:last_image_url"`
	FaceRecognition         bool     `json:"faceRecognition" gorm:"column:face_recognition;default:false"`
	LicensePlateRecognition bool     `json:"licensePlateRecognition" gorm:"column:license_plate_recognition;default:false"`
	InsertStatus            bool     `json:"insertstatus"`
	ExportStatus            bool     `json:"exportstatus"`

	// NVR group / location info
	GroupIDs   KeyValueArray `json:"groupIDs" gorm:"column:group_id;embedded;type:jsonb"`
	SiteInfo   KeyValue      `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Location   string        `json:"location" gorm:"column:location"`
	Lat        string        `json:"lat" gorm:"column:lat"`
	Long       string        `json:"long" gorm:"column:long"`
	Coordinate string        `json:"coordinate" gorm:"column:coordinate"`
	Position   string        `json:"position" gorm:"column:position"`

	// Added cameras
	Cameras KeyValueArray `json:"cameras" gorm:"embedded;type:jsonb"`
	Box     KeyValue      `json:"box" gorm:"column:box;embedded;type:jsonb"`

	// NVR config struct (foreign key)
	ConfigID uuid.UUID `json:"configID" gorm:"column:config_id;type:uuid;default:uuid_generate_v4()"`

	// Device's metadata, there are should be key/value for management I/O, relay, rs232, rs485...
	Metadata KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`

	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTONVRImport struct {
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
	RtspPort         string   `json:"rtspPort,omitempty"`
	ManagementPort   string   `json:"managementPort,omitempty"`
	Username         string   `json:"username,omitempty"`
	Password         string   `json:"password,omitempty"`
	UseTLS           bool     `json:"useTLS,omitempty"`
	IsOfflineSetting bool     `json:"isOfflineSetting"`
	IFrameURL        string   `json:"iframeURL,omitempty"`
	Lat              string   `json:"lat"`
	Long             string   `json:"long"`
	InsertStatus     bool     `json:"insertstatus"`
	ExportStatus     bool     `json:"exportstatus"`

	// Stream/recording info
	Streams VideoStreamArray `json:"streams,omitempty"`
	NVR     KeyValue         `json:"nvr"`
	Box     KeyValue         `json:"box"`
	// NVR group / location info
	GroupIDs                KeyValueArray `json:"groupIDs" gorm:"column:group_id;embedded;type:jsonb"`
	SiteInfo                KeyValue      `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Location                string        `json:"location" gorm:"column:location"`
	Coordinate              string        `json:"coordinate" gorm:"column:coordinate"`
	Position                string        `json:"position" gorm:"column:position"`
	FaceRecognition         bool          `json:"faceRecognition"`
	LicensePlateRecognition bool          `json:"licensePlateRecognition"`
	// NVR config struct (foreign key)
	ConfigID uuid.UUID `json:"configID"`
}

type DTONVRImport_Read_BasicInfo struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name,omitempty"`

	// General
	Type             KeyValue `json:"type,omitempty"`
	Protocol         KeyValue `json:"protocol,omitempty"`
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
	InsertStatus     bool     `json:"status"`
	ExportStatus     bool     `json:"exportstatus"`

	// Stream/recording info
	Streams                 VideoStreamArray `json:"streams,omitempty"`
	NVR                     KeyValue         `json:"nvr"`
	FaceRecognition         bool             `json:"faceRecognition"`
	LicensePlateRecognition bool             `json:"licensePlateRecognition"`

	// NVR group / location info
	GroupIDs   KeyValueArray `json:"groupIDs" gorm:"column:group_id;embedded;type:jsonb"`
	SiteInfo   KeyValue      `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Location   string        `json:"location" gorm:"column:location"`
	Coordinate string        `json:"coordinate" gorm:"column:coordinate"`
	Position   string        `json:"position" gorm:"column:position"`
}

type DTONVRImportDeviceConfig struct {
	ID      uuid.UUID     `json:"id"`
	Configs KeyValueArray `json:"configs,omitempty"`
}

type DTONVRImportDeviceConfigBatch struct {
	NVRIDs  []uuid.UUID   `json:"nvrIDs"`
	Configs KeyValueArray `json:"configs,omitempty"`
}
