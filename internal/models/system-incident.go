package models

import (
	"time"

	uuid "github.com/google/uuid"
)

// Entity Model for camera device
type SystemIncident struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	DeviceType string    `json:"deviceType"`
	DeviceID   string    `json:"deviceID"`
	DeviceName string    `json:"deviceName"`
	EventType  string    `json:"eventType"`
	EventName  string    `json:"eventName"`
	Severity   string    `json:"severity"`
	Status     string    `json:"status"`
	Location   string    `json:"location"`
	Source     string    `json:"source"`
	Type       string    `json:"type"`

	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_System_Incident_BasicInfo struct {
	ID         uuid.UUID `json:"id"`
	DeviceType string    `json:"deviceType"`
	DeviceID   string    `json:"deviceID"`
	EventType  string    `json:"eventType"`
	EventName  string    `json:"eventName,omitempty"`
	Severity   string    `json:"severity,omitempty"`
	Status     string    `json:"status,omitempty"`
	Location   string    `json:"location,omitempty"`
	Source     string    `json:"source,omitempty"`
	Type       string    `json:"type,omitempty"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
}
