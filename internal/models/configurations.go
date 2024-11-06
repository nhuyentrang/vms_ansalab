package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type ImageConfig struct {
	ID           uuid.UUID     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Metadata     KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`
	DeleteMark   bool          `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	RecentUpdate bool          `json:"recentUpdate" gorm:"column:recent_update;default:true"`
	CreatedAt    time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt    time.Time     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt    time.Time     `json:"deletedAt" gorm:"column:deleted_at"`

	CameraID    uuid.UUID `json:"cameraID" gorm:"column:camera_id"`
	DisableName bool      `json:"disableName" gorm:"column:disable_name;default:true"`
	DisableDate bool      `json:"disableDate" gorm:"column:disable_date;default:true"`
	DisableWeek bool      `json:"disableWeek" gorm:"column:disable_week;default:true"`
	DateFormat  string    `json:"dateFormat" gorm:"column:date_format"`
	TimeFormat  string    `json:"timeFormat" gorm:"column:time_format"`
	NameX       string    `json:"nameX" gorm:"column:name_x"`
	NameY       string    `json:"nameY" gorm:"column:name_y"`
	DateX       string    `json:"dateX" gorm:"column:date_x"`
	DateY       string    `json:"dateY" gorm:"column:date_y"`
	WeekX       string    `json:"weekX" gorm:"column:week_x"`
	WeekY       string    `json:"weekY" gorm:"column:week_y"`
}

// type ImageConfig struct {
// 	ID           uuid.UUID     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
// 	Metadata     KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`
// 	DeleteMark   bool          `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
// 	RecentUpdate bool          `json:"recentUpdate" gorm:"column:recent_update;default:true"`
// 	CreatedAt    time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
// 	UpdatedAt    time.Time     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
// 	DeletedAt    time.Time     `json:"deletedAt" gorm:"column:deleted_at"`

// 	CameraID    uuid.UUID    `json:"cameraID" gorm:"column:camera_id"`
// 	ImageConfig *ImageConfig `json:"imageConfig,omitempty"`
// }

// type ImageChannels struct {
// 	Version            string             `xml:"version,attr"`
// 	XMLName            xml.Name           `xml:"ImageChannel,omitempty"`
// 	XMLNamespace       string             `xml:"xmlns,attr"`
// 	ID                 int                `xml:"id"`
// 	Enabled            bool               `xml:"enabled"`
// 	VideoInputID       int                `xml:"videoInputID"`
// 	ImageFlip          ImageFlip          `xml:"ImageFlip"`
// 	WDR                WDR                `xml:"WDR"`
// 	BLC                BLC                `xml:"BLC"`
// 	IrcutFilter        IrcutFilter        `xml:"IrcutFilter"`
// 	WhiteBalance       WhiteBalance       `xml:"WhiteBalance"`
// 	Exposure           Exposure           `xml:"Exposure"`
// 	Sharpness          Sharpness          `xml:"Sharpness"`
// 	Shutter            Shutter            `xml:"Shutter"`
// 	PowerLineFrequency PowerLineFrequency `xml:"powerLineFrequency"`
// 	Color              Color              `xml:"Color"`
// 	NoiseReduce        NoiseReduce        `xml:"NoiseReduce"`
// 	HLC                HLC                `xml:"HLC"`
// 	SupplementLight    SupplementLight    `xml:"SupplementLight"`
// }
// type ImageFlip struct {
// 	Enabled        bool   `xml:"enabled"`
// 	ImageFlipStyle string `xml:"ImageFlipStyle"`
// }

// type WDR struct {
// 	Mode     string `xml:"mode"`
// 	WDRLevel int    `xml:"WDRLevel"`
// }

// type BLC struct {
// 	Enabled bool   `xml:"enabled"`
// 	BLCMode string `xml:"BLCMode"`
// }

// type IrcutFilter struct {
// 	IrcutFilterType       string   `xml:"IrcutFilterType"`
// 	NightToDayFilterLevel int      `xml:"nightToDayFilterLevel"`
// 	NightToDayFilterTime  int      `xml:"nightToDayFilterTime"`
// 	Schedule              Schedule `xml:"Schedule"`
// }

// type Schedule struct {
// 	ScheduleType string `xml:"scheduleType"`
// 	BeginTime    string `xml:"TimeRange>beginTime"`
// 	EndTime      string `xml:"TimeRange>endTime"`
// }

// type WhiteBalance struct {
// 	WhiteBalanceStyle string `xml:"WhiteBalanceStyle"`
// 	WhiteBalanceRed   int    `xml:"WhiteBalanceRed"`
// 	WhiteBalanceBlue  int    `xml:"WhiteBalanceBlue"`
// }

// type Exposure struct {
// 	ExposureType       string             `xml:"ExposureType"`
// 	OverexposeSuppress OverexposeSuppress `xml:"OverexposeSuppress"`
// }

// type OverexposeSuppress struct {
// 	Enabled bool `xml:"enabled"`
// }

// type Sharpness struct {
// 	SharpnessLevel int `xml:"SharpnessLevel"`
// }

// type Shutter struct {
// 	ShutterLevel string `xml:"ShutterLevel"`
// }

// type PowerLineFrequency struct {
// 	PowerLineFrequencyMode string `xml:"powerLineFrequencyMode"`
// }

// type Color struct {
// 	BrightnessLevel int `xml:"brightnessLevel"`
// 	ContrastLevel   int `xml:"contrastLevel"`
// 	SaturationLevel int `xml:"saturationLevel"`
// }

// type NoiseReduce struct {
// 	Mode         string       `xml:"mode"`
// 	GeneralMode  GeneralMode  `xml:"GeneralMode"`
// 	AdvancedMode AdvancedMode `xml:"AdvancedMode"`
// }

// type GeneralMode struct {
// 	GeneralLevel int `xml:"generalLevel"`
// }

// type AdvancedMode struct {
// 	FrameNoiseReduceLevel      int `xml:"FrameNoiseReduceLevel"`
// 	InterFrameNoiseReduceLevel int `xml:"InterFrameNoiseReduceLevel"`
// }

// type HLC struct {
// 	Enabled  bool `xml:"enabled"`
// 	HLCLevel int  `xml:"HLCLevel"`
// }

// type SupplementLight struct {
// 	SupplementLightMode             string `xml:"supplementLightMode"`
// 	MixedLightBrightnessRegulatMode string `xml:"mixedLightBrightnessRegulatMode"`
// 	IrLightBrightness               int    `xml:"irLightBrightness"`
// 	IsAutoModeBrightnessCfg         bool   `xml:"isAutoModeBrightnessCfg"`
// }

type VideoConfig struct {
	ID           uuid.UUID     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Metadata     KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`
	DeleteMark   bool          `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	RecentUpdate bool          `json:"recentUpdate" gorm:"column:recent_update;default:true"`
	CreatedAt    time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt    time.Time     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt    time.Time     `json:"deletedAt" gorm:"column:deleted_at"`

	CameraID        uuid.UUID          `json:"cameraID" gorm:"column:camera_id"`
	VideoConfigInfo VideoConfigInfoArr `json:"videoConfig" gorm:"column:video_config;embedded;type:jsonb"`
}

type AudioConfig struct {
	ID       uuid.UUID     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Metadata KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`
}
type NetworkConfig struct {
	ID uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	// Metadata     KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`
	DeleteMark   bool      `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	RecentUpdate bool      `json:"recentUpdate" gorm:"column:recent_update;default:true"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt    time.Time `json:"deletedAt" gorm:"column:deleted_at;autoCreateTime:true"`
	// Config TCP/IP
	NicType            KeyValue `json:"nicType" gorm:"column:nic_type;embedded;type:jsonb"`
	DHCP               bool     `json:"dhcp" gorm:"column:dhcp"`
	IPv4SubnetMask     string   `json:"ipv4SubnetMask" gorm:"column:ipv4_subnet_mask"`
	IPv4DefaultGateway string   `json:"ipv4DefaultGateway" gorm:"column:ipv4_default_gateway"`
	SubnetPrefixLength string   `json:"subnetPrefixLength" gorm:"column:subnet_prefix_length"`
	IPv6DefaultGateway string   `json:"ipv6DefaultGateway" gorm:"column:ipv6_default_gateway"`
	MacAddress         string   `json:"macAddress" gorm:"column:mac_address"`
	MTU                int      `json:"mtu" gorm:"column:mtu"`
	AutoDNS            bool     `json:"autoDNS" gorm:"column:auto_dns"`
	MulticastAddress   string   `json:"multicastAddress" gorm:"column:multicast_address"`
	PrefDNS            string   `json:"prefDNS" gorm:"column:pref_dns"`
	AlterDNS           string   `json:"alterDNS" gorm:"column:alter_dns"`
	// Config DDNS
	DDNS              bool     `json:"ddns" gorm:"column:ddns"`
	DDNSType          KeyValue `json:"ddnsType" gorm:"column:ddns_type;embedded;type:jsonb"`
	ServerAddressDDNS string   `json:"serverAddress" gorm:"column:server_address"`
	Domain            string   `json:"domain" gorm:"column:domain"`
	Port              string   `json:"port" gorm:"column:port"`
	UserName          string   `json:"userName" gorm:"column:user_name"`
	Password          string   `json:"password" gorm:"column:password"`
	// Config Port
	HTTP   int `json:"http" gorm:"column:http"`
	HTTPS  int `json:"https" gorm:"column:https"`
	RTSP   int `json:"rtsp" gorm:"column:rtsp"`
	Server int `json:"server" gorm:"column:server"`
	// Config NTP
	NTP              bool      `json:"ntp" gorm:"column:ntp"`
	ServerAddressNTP string    `json:"serverAddressNTP" gorm:"column:server_address_ntp"`
	NTPPort          string    `json:"ntpPort" gorm:"column:ntp_port"`
	Interval         string    `json:"interval" gorm:"column:interval"`
	ManualTime       bool      `json:"manualTime" gorm:"column:manual_time"`
	SetTime          time.Time `json:"setTime" gorm:"column:set_time"`

	//Token
	TokenSetNetworkInterfaces string `json:"tokenSetNetworkInterfaces" gorm:"column:token_setnetwork_interfaces"`
}

type StorageConfig struct {
	ID           uuid.UUID     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Metadata     KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`
	DeleteMark   bool          `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	RecentUpdate bool          `json:"recentUpdate" gorm:"column:recent_update;default:true"`
	CreatedAt    time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt    time.Time     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt    time.Time     `json:"deletedAt" gorm:"column:deleted_at"`

	IsConfig       bool              `json:"isConfig" gorm:"column:is_config"`
	ConfigName     bool              `json:"configName" gorm:"column:config_name"`
	ConfigDescribe bool              `json:"configDescribe" gorm:"column:config_describe"`
	MaxMemory      bool              `json:"maxMemory" gorm:"column:max_memory"`
	ExpectedMemory bool              `json:"expectedMemory" gorm:"column:expected_memory"`
	Schedule       TimeScheduleArray `json:"schedule" gorm:"embedded;type:jsonb"`
}

type StreamingConfig struct {
	// Config for signaling, p2p, turn/stun
	ID       uuid.UUID     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Metadata KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`
}

type AIConfig struct {
	// Config for signaling, p2p, turn/stun
	ID       uuid.UUID     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Metadata KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`
}

type PTZConfig struct {
	// Config for ptz
	ID       uuid.UUID     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Metadata KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`
}

type EventConfig struct {
	// Config for event config
	ID          uuid.UUID     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CreatedAt   time.Time     `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt   time.Time     `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeleteMark  bool          `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	DeletedAt   time.Time     `json:"deletedAt" gorm:"column:deleted_at"`
	Metadata    KeyValueArray `json:"metadata" gorm:"column:metadata;embedded;type:jsonb"`
	CameraID    uuid.UUID     `json:"cameraID" gorm:"type:uuid;default:uuid_generate_v4()"`
	Type        string        `json:"type" gorm:"column:type"`
	Sensitivity int           `json:"sensitivity" gorm:"column:sensitivity"`
	Object      string        `json:"object" gorm:"column:object"`
	GridMap     string        `json:"gridMap" gorm:"column:grid_map"`
}

type DTO_NetworkConfig struct {
	ID uuid.UUID `json:"id,omitempty" `
	// Config TCP/IP
	NicType            KeyValue `json:"nicType,omitempty"`
	DHCP               bool     `json:"dhcp"`
	IPv4SubnetMask     string   `json:"ipv4SubnetMask,omitempty"`
	IPv4DefaultGateway string   `json:"ipv4DefaultGateway,omitempty"`
	SubnetPrefixLength string   `json:"subnetPrefixLength,omitempty"`
	IPv6DefaultGateway string   `json:"ipv6DefaultGateway,omitempty"`
	MacAddress         string   `json:"macAddress,omitempty"`
	MTU                int      `json:"mtu,omitempty"`
	MulticastAddress   string   `json:"multicastAddress,omitempty"`
	AutoDNS            bool     `json:"autoDNS"`
	PrefDNS            string   `json:"prefDNS,omitempty"`
	AlterDNS           string   `json:"alterDNS,omitempty"`
	// Config DDNS
	DDNS              bool     `json:"ddns"`
	DDNSType          KeyValue `json:"ddnsType,omitempty"`
	ServerAddressDDNS string   `json:"serverAddress,omitempty"`
	Domain            string   `json:"domain,omitempty"`
	Port              string   `json:"port,omitempty"`
	UserName          string   `json:"userName,omitempty"`
	Password          string   `json:"password,omitempty"`
	// Config Port
	HTTP   int `json:"http,omitempty"`
	HTTPS  int `json:"https,omitempty"`
	RTSP   int `json:"rtsp,omitempty"`
	ONVIF  int `json:"onvif,omitempty"`
	Server int `json:"server,omitempty"`
	// Config NTP
	NTP              bool      `json:"ntp"`
	ServerAddressNTP string    `json:"serverAddressNTP,omitempty"`
	NTPPort          string    `json:"ntpPort,omitempty"`
	Interval         string    `json:"interval,omitempty"`
	ManualTime       bool      `json:"manualTime"`
	SetTime          time.Time `json:"setTime,omitempty"`
	// Token
	TokenSetNetworkInterfaces string `json:"tokenSetNetworkInterfaces,omitempty"`
}

type DTO_ImageConfig struct {
	ID          uuid.UUID `json:"id,omitempty" `
	CameraID    uuid.UUID `json:"cameraID,omitempty"`
	DisableName bool      `json:"disableName"`
	DisableDate bool      `json:"disableDate"`
	DisableWeek bool      `json:"disableWeek"`
	DateFormat  string    `json:"dateFormat,omitempty"`
	TimeFormat  string    `json:"timeFormat,omitempty"`
	NameX       string    `json:"nameX,omitempty"`
	NameY       string    `json:"nameY,omitempty"`
	DateX       string    `json:"dateX,omitempty"`
	DateY       string    `json:"dateY,omitempty"`
	WeekX       string    `json:"weekX,omitempty"`
	WeekY       string    `json:"weekY,omitempty"`
}

type DTO_EventConfig struct {
	ID          uuid.UUID     `json:"id,omitempty"`
	Metadata    KeyValueArray `json:"metadata,omitempty"`
	CameraID    uuid.UUID     `json:"cameraID,omitempty"`
	Type        string        `json:"type,omitempty"`
	Sensitivity int           `json:"sensitivity,omitempty"`
	Object      string        `json:"object,omitempty"`
	GridMap     string        `json:"gridMap,omitempty"`
}

type DTO_VideoConfig struct {
	ID              uuid.UUID          `json:"id"`
	CameraID        uuid.UUID          `json:"cameraID"`
	VideoConfigInfo VideoConfigInfoArr `json:"videoConfig" gorm:"column:video_config;embedded;type:jsonb"`
}
type DTO_VideoConfigCamera struct {
	ID          string            `json:"id"`
	CameraID    string            `json:"cameraID"`
	VideoConfig []VideoConfigInfo `json:"videoConfig"`
}

type DTO_VideoConfigNVR struct {
	ID              uuid.UUID          `json:"id"`
	CameraID        uuid.UUID          `json:"cameraID"`
	VideoConfigInfo VideoConfigInfoArr `json:"videoConfigInfo" gorm:"column:video_config;embedded;type:jsonb"`
}

type DTO_CmdEvent struct {
	CommandID      string  `json:"commandId,omitempty"`
	ClientID       string  `json:"clientID,omitempty"`
	EventTimeEpoch float64 `json:"eventTimeEpoch,omitempty"`
	EventTime      string  `json:"eventTime,omitempty"`
	Cmd            string  `json:"cmd,omitempty"`
	Code           string  `json:"code,omitempty"`
	Type           string  `json:"type,omitempty"`
	X              int16   `json:"x,omitempty"`
	Y              int16   `json:"y,omitempty"`
	W              int16   `json:"w,omitempty"`
	H              int16   `json:"h,omitempty"`
	Sensitivisy    string  `json:"sensitivisy,omitempty"`
	GridMap        string  `json:"gridMap,omitempty"`
}
