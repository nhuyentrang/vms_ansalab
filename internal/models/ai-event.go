package models

import (
	"time"

	uuid "github.com/google/uuid"
	//"github.com/lib/pq"
)

// Entity Model for device
type AIEvent struct {
	ID uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	//Name         string        `json:"name"`
	Type         KeyValue      `json:"type" gorm:"embedded;type:jsonb"`
	SourceDevice KeyValue      `json:"sourceDevice" gorm:"column:source_device;embedded;type:jsonb"`
	SiteInfo     KeyValue      `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Status       KeyValue      `json:"status" gorm:"embedded;type:jsonb"`
	Level        KeyValue      `json:"level" gorm:"embedded;type:jsonb"`
	Atts         int64         `json:"atts"`
	Metadata     KeyValueArray `json:"metadata" gorm:"embedded;type:jsonb"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_AIEvent struct {
	ID uuid.UUID `json:"id"`
	//Name         string        `json:"name,omitempty"`
	Type         KeyValue      `json:"type,omitempty"`
	SourceDevice KeyValue      `json:"sourceDevice,omitempty"`
	DeviceName   string        `json:"deviceName,omitempty"`
	SiteInfo     KeyValue      `json:"siteInfo,omitempty"`
	Status       KeyValue      `json:"status,omitempty"`
	Level        KeyValue      `json:"level,omitempty"`
	Atts         int64         `json:"atts,omitempty"`
	Metadata     KeyValueArray `json:"metadata,omitempty"`
}

type DTO_AIEvent_Create struct {
	//Name         string        `json:"name,omitempty"`
	Type         KeyValue      `json:"type,omitempty"`
	SourceDevice KeyValue      `json:"sourceDevice,omitempty"`
	SiteInfo     KeyValue      `json:"siteInfo,omitempty"`
	Level        KeyValue      `json:"level,omitempty"`
	Atts         int64         `json:"atts,omitempty"`
	Metadata     KeyValueArray `json:"metadata,omitempty"`
}

type DTO_AIEvent_Read_BasicInfo struct {
	ID uuid.UUID `json:"id"`
	//Name         string        `json:"name,omitempty"`
	Type         KeyValue      `json:"type,omitempty"`
	SourceDevice KeyValue      `json:"sourceDevice,omitempty"`
	SiteInfo     KeyValue      `json:"siteInfo,omitempty"`
	Status       KeyValue      `json:"status,omitempty"`
	Level        KeyValue      `json:"level,omitempty"`
	Atts         int64         `json:"atts,omitempty"`
	Metadata     KeyValueArray `json:"metadata,omitempty"`
}

// Json for response of cctv device types
type JsonAIDeviceTypeRsp struct {
	Code    int        `json:"code"`
	Count   int        `json:"count"`
	Data    []KeyValue `json:"data"`
	Message string     `json:"message,omitempty"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
}

// Json for response of cctv event types
type JsonAIEventTypeRsp struct {
	Code    int        `json:"code"`
	Count   int        `json:"count"`
	Data    []KeyValue `json:"data"`
	Message string     `json:"message,omitempty"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
}

// Json for response of cctv event status type
type JsonAIEventStatusTypeRsp struct {
	Code    int        `json:"code"`
	Count   int        `json:"count"`
	Data    []KeyValue `json:"data"`
	Message string     `json:"message,omitempty"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
}

// Json for response of cctv event levels
type JsonAIEventLevelTypeRsp struct {
	Code    int        `json:"code"`
	Count   int        `json:"count"`
	Data    []KeyValue `json:"data"`
	Message string     `json:"message,omitempty"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
}
