package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	uuid "github.com/google/uuid"
)

// Model for device

/*************************************************************** Devices ***************************************************************/
type Device struct {
	ID              uuid.UUID   `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	NameDevice      string      `json:"nameDevice"`
	UserID          int64       `json:"user_id" gorm:"column:user_id"`
	Serial          string      `json:"serial"`
	DeviceType      string      `json:"deviceType"`
	ModelID         string      `json:"modelId"`
	AppVersion      string      `json:"appVersion"`
	AreaId          string      `json:"areaId"`
	AreaName        string      `json:"areaName"`
	HttpPort        string      `json:"httpPort" gorm:"column:http_port"`
	OnvifPort       string      `json:"onvifPort" gorm:"column:onvif_port"`
	DeviceCode      string      `json:"deviceCode"`
	Status          string      `json:"status"`
	Location        string      `json:"location"`
	IPAddress       string      `json:"ipAddress"`
	HardwareVersion string      `json:"hardwareVersion"`
	MacAddress      string      `json:"macAddress"`
	MqttAccount     MqttAccount `json:"mqttAccount" gorm:"embedded;type:jsonb"`
	LastPing        time.Time   `json:"lastPing"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_Device struct {
	ID              uuid.UUID   `json:"id"`
	NameDevice      string      `json:"nameDevice"`
	UserID          int64       `json:"user_id"`
	Serial          string      `json:"serial"`
	Token           string      `json:"token"`
	DeviceType      string      `json:"deviceType"`
	ModelID         string      `json:"modelId"`
	DeviceCode      string      `json:"deviceCode"`
	Status          string      `json:"status"`
	AppVersion      string      `json:"appVersion"`
	HardwareVersion string      `json:"hardwareVersion"`
	AreaId          string      `json:"areaId"`
	AreaName        string      `json:"areaName"`
	HttpPort        string      `json:"httpPort" gorm:"column:http_port"`
	OnvifPort       string      `json:"onvifPort" gorm:"column:onvif_port"`
	Location        string      `json:"location"`
	IPAddress       string      `json:"ipAddress"`
	MacAddress      string      `json:"macAddress"`
	Metadata        MqttAccount `json:"metadata"`
	Telemetry       MqttAccount `json:"telemetry"`
	MqttAccount     MqttAccount `json:"mqttAccount"`
	LastPing        time.Time   `json:"lastPing"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
}

func (sla *MqttAccount) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla MqttAccount) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type DeviceProduct struct {
	ID         int64  `json:"id"`
	Serial     string `json:"serial"`
	Token      string `json:"token"`
	Type       string `json:"type"`
	IsBlocking int64  `json:"is_blocking"`
	CreatedAt  string `json:"created_at" gorm:"autoCreateTime:true"`
	UpdatedAt  string `json:"updated_at" gorm:"autoUpdateTime:true"`
}

type DTO_DeviceProduct struct {
	ID     int64  `json:"id" validate:"required"`
	Serial string `json:"serial" validate:"required"`
	Token  string `json:"token"`
}

type Command struct {
	CommandID       string `json:"commandId,omitempty"`
	Cmd             string `json:"cmd,omitempty"`
	Status          string `json:"status,omitempty"`
	EventTime       string `json:"eventTime,omitempty"`
	SoftwareVersion string `json:"softwareVersion,omitempty"`
	LinkMinio       string `json:"linkMinio,omitempty"`
	Token           string `json:"token,omitempty"`
}

// Onvif
type ScanDevice struct {
	URL             string `json:"url,omitempty"`
	IDBox           string `json:"idBox,omitempty"`
	Host            string `json:"host,omitempty"`
	UserName        string `json:"userName,omitempty"`
	PassWord        string `json:"passWord,omitempty"`
	StartURL        string `json:"startURL,omitempty"`
	EndURL          string `json:"endURL,omitempty"`
	StartHost       int    `json:"startHost,omitempty"`
	EndHost         int    `json:"endHost,omitempty"`
	Name            string `json:"name"`
	Type            string `json:"type"`
	Location        string `json:"location"`
	Status          string `json:"status"`
	MacAddress      string `json:"macAddress"`
	FirmwareVersion string `xml:"firmwareVersion,omitempty" json:"firmwareVersion,omitempty"`
	Protocol        string `json:"protocol,omitempty"`
}
