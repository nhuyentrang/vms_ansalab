package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type Event struct {
	ID                 uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	ImageURL           string    `json:"imageUrl"`
	TimeOccur          int       `json:"timeOccur"`
	NotificationStatus int       `json:"notificationStatus"`
	EventType          string    `json:"eventType"`
	EventName          string    `json:"eventName"`
	DeviceId           string    `json:"deviceId"`
	DeviceName         string    `json:"deviceName"`
	AreaId             string    `json:"areaId"`
	AreaName           string    `json:"areaName"`
	Status             string    `json:"status"`
	RawMsg             string    `json:"rawMsgs"`
	Description        string    `json:"description"`

	CamIP     string    `json:"camIP"`
	CamName   string    `json:"camName"`
	CameraId  string    `json:"cameraId"`
	Image     string    `json:"image"`
	ImageID   uuid.UUID `json:"imageID"`
	Timestamp string    `json:"timestamp"`

	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_Event struct {
	ID                 uuid.UUID `json:"id"`
	ImageURL           string    `json:"imageUrl"`
	TimeOccur          int       `json:"timeOccur"`
	NotificationStatus int       `json:"notificationStatus"`
	EventType          string    `json:"eventType"`
	EventName          string    `json:"eventName"`
	DeviceId           string    `json:"deviceId"`
	DeviceName         string    `json:"deviceName"`
	AreaId             string    `json:"areaId"`
	AreaName           string    `json:"areaName"`
	Status             string    `json:"status"`
	RawMsg             string    `json:"rawMsgs"`
}
type DTO_Event_Created struct {
	CamIP     string    `json:"camIP"`
	CamName   string    `json:"camName"`
	CameraId  string    `json:"cameraId"`
	ImageID   uuid.UUID `json:"imageID"`
	Image     string    `json:"image"`
	Timestamp string    `json:"timestamp"`
	RawMsg    string    `json:"rawMsgs"`
}

type DTO_Event_Read_BasicInfo struct {
	ID                 uuid.UUID `json:"id"`
	ImageURL           string    `json:"imageUrl"`
	ImageID            uuid.UUID `json:"imageID"`
	TimeOccur          int       `json:"timeOccur"`
	NotificationStatus int       `json:"notificationStatus"`
	EventType          string    `json:"eventType"`
	EventName          string    `json:"eventName"`
	DeviceId           string    `json:"deviceId"`
	DeviceName         string    `json:"deviceName"`
	AreaId             string    `json:"areaId"`
	AreaName           string    `json:"areaName"`
	Status             string    `json:"status"`
	RawMsg             string    `json:"rawMsgs"`
}

type DTO_Event_UpdateDescription struct {
	Description string `json:"description"`
}

type DTO_Event_UpdateStatus struct {
	Status string `json:"status"`
}

type DTO_Event_Update struct {
	Status string `json:"status"`
}

type DTO_Event_ImageURL struct {
	ImageURL string `json:"imageUrl"`
}
