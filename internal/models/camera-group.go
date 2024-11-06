package models

import (
	"time"

	uuid "github.com/google/uuid"
	//"github.com/lib/pq"
)

// Entity Model for device
type CameraGroup struct {
	ID          uuid.UUID           `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Cameras     CameraKeyValueArray `json:"cameras" gorm:"embedded;type:jsonb"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_CameraGroup struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Cameras     CameraKeyValueArray `json:"cameras,omitempty"`
}

type DTO_CameraGroup_Create struct {
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Cameras     CameraKeyValueArray `json:"cameras,omitempty"`
}

type DTO_CameraGroup_Read_BasicInfo struct {
	ID          uuid.UUID           `json:"id"`
	Name        string              `json:"name,omitempty"`
	Description string              `json:"description,omitempty"`
	Cameras     CameraKeyValueArray `json:"cameras,omitempty"`
}
