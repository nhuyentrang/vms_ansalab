package models

import (
	"time"

	uuid "github.com/google/uuid"
	//"github.com/lib/pq"
)

// Entity Model for device
type MediaLibrary struct {
	ID       uuid.UUID     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name     string        `json:"name"`
	Type     KeyValue      `json:"type" gorm:"embedded;type:jsonb"`
	UserID   string        `json:"userID" gorm:"column:user_id"`
	Size     int64         `json:"size"`
	Atts     int64         `json:"atts"`
	Metadata KeyValueArray `json:"metadata" gorm:"embedded;type:jsonb"`
	URL      string        `json:"url"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_MediaLibrary struct {
	ID       uuid.UUID     `json:"id"`
	Name     string        `json:"name,omitempty"`
	Type     KeyValue      `json:"type,omitempty"`
	UserID   string        `json:"userID,omitempty"`
	Size     int64         `json:"size,omitempty"`
	Atts     int64         `json:"atts,omitempty"`
	Metadata KeyValueArray `json:"metadata,omitempty"`
	URL      string        `json:"url,omitempty"`
}

type DTO_MediaLibrary_Create struct {
	Name     string        `json:"name,omitempty"`
	Type     KeyValue      `json:"type,omitempty"`
	UserID   string        `json:"userID,omitempty"`
	Size     int64         `json:"size,omitempty"`
	Atts     int64         `json:"atts,omitempty"`
	Metadata KeyValueArray `json:"metadata,omitempty"`
	URL      string        `json:"url,omitempty"`
}

type DTO_MediaLibrary_Read_BasicInfo struct {
	ID       uuid.UUID     `json:"id"`
	Name     string        `json:"name,omitempty"`
	Type     KeyValue      `json:"type,omitempty"`
	UserID   string        `json:"userID,omitempty"`
	Size     int64         `json:"size,omitempty"`
	Atts     int64         `json:"atts,omitempty"`
	Metadata KeyValueArray `json:"metadata,omitempty"`
	URL      string        `json:"url,omitempty"`
}

// Json for response of media file types
type JsonMediaFileTypeRsp struct {
	Code    int        `json:"code"`
	Count   int        `json:"count"`
	Data    []KeyValue `json:"data"`
	Message string     `json:"message,omitempty"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
}
