package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type CameraAIEventProperty struct {
	ID                   uuid.UUID            `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CameraID             uuid.UUID            `json:"camera_id" gorm:"column:camera_id;type:uuid;default:uuid_generate_v4()"`
	CameraAIPropertyList CameraAIPropertyList `json:"cameraaiproperty" gorm:"embedded;type:jsonb;column:cameraaiproperty"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

// Todo, fix FE and this model for the correct json name camera_id and cameraaiproperty
type DTO_Camera_AI_Property_BasicInfo struct {
	ID                   uuid.UUID            `json:"id"`
	CameraID             uuid.UUID            `json:"camera_id"`
	CameraAIPropertyList CameraAIPropertyList `json:"cameraaiproperty"`
}
