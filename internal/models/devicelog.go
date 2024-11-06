package models

import (
	"time"

	uuid "github.com/google/uuid"
)

// TableName overrides the table name
func (DeviceLog) TableName() string {
	return "ivis_vms.device_log"
}

// Entity Model for DeviceLog
type DeviceLog struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Atts       int64     `json:"atts"`
	Command    string    `json:"commnad"`
	DeviceCode string    `json:"deviceCode" gorm:"column:device_code"`
	DeviceType string    `json:"deviceType" gorm:"column:device_type"`
	Action     string    `json:"action"`
	Detail     string    `json:"detail"`
	Protocol   string    `json:"Protocol"`
	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

// DTO Model for DeviceLog Read
type DTO_DeviceLog struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Atts       int64     `json:"atts,omitempty"`
	Command    string    `json:"commnad,omitempty"`
	DeviceCode string    `json:"deviceCode,omitempty"`
	DeviceType string    `json:"deviceType,omitempty"`
	Action     string    `json:"action,omitempty"`
	Payload    string    `json:"payload,omitempty"`
	Protocol   string    `json:"Protocol,omitempty"`
	Detail     string    `json:"detail,omitempty"`
}

// DTO Model for DeviceLog
type DTO_DeviceLog_Read_BasicInfo struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Atts       int64     `json:"atts,omitempty"`
	Command    string    `json:"commnad,omitempty"`
	DeviceCode string    `json:"deviceCode,omitempty"`
	DeviceType string    `json:"deviceType,omitempty"`
	Action     string    `json:"action,omitempty"`
	Payload    string    `json:"payload,omitempty"`
	Protocol   string    `json:"Protocol,omitempty"`
	Detail     string    `json:"detail,omitempty"`
}
