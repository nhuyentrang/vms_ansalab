package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type CameraModelAI struct {
	ID             uuid.UUID          `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Type           string             `json:"type" gorm:"column:type"`
	ModelName      string             `json:"modelName" gorm:"column:model_name"`
	Characteristic CharacteristicList `json:"characteristic" gorm:"column:characteristic;embedded;type:jsonb"`
	CreatedAt      time.Time          `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt      time.Time          `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt      time.Time          `json:"deletedAt" gorm:"column:deleted_at"`

	Count int64 `json:"countCamera" gorm:"column:model_Count"`
}

type DTO_CameraModelAI struct {
	ID             uuid.UUID          `json:"id,omitempty"`
	Type           string             `json:"type,omitempty"`
	ModelName      string             `json:"modelName,omitempty"`
	Characteristic CharacteristicList `json:"characteristic"`
	CreatedAt      time.Time          `json:"createdAt,omitempty"`
	UpdatedAt      time.Time          `json:"updatedAt,omitempty"`
	DeletedAt      time.Time          `json:"deletedAt,omitempty"`
	Count          int64              `json:"countCamera" gorm:"column:model_Count"`
}
