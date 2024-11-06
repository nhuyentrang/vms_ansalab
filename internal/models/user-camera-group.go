package models

import (
	"time"

	uuid "github.com/google/uuid"
)

// Entity Model for camera device
type UserCameraGroup struct {
	ID      uuid.UUID          `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	UserID  string             `json:"userid"`
	Cameras ListUserCameraView `json:"cameras" gorm:"embedded;type:jsonb"`

	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_User_Camera_Group_BasicInfo struct {
	ID      uuid.UUID          `json:"id"`
	UserID  string             `json:"userid"`
	Cameras ListUserCameraView `json:"cameras,omitempty"`
}
