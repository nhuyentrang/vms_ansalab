package models

import (
	"time"

	"github.com/google/uuid"
)

type PcInfo struct {
	ID   uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Time int64     `json:"time"`
	CPU  float64   `json:"cpu"`
	GPU  float64   `json:"gpu"`
	Ram  float64   `json:"Ram"`
	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
}

type DTO_PcInfo struct {
	ID   uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Time int64     `json:"time"`
	CPU  float64   `json:"cpu"`
	GPU  float64   `json:"gpu"`
	Ram  float64   `json:"Ram"`
	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
}
