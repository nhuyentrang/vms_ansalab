package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type LicensePlates struct {
	ID uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	// ImageUrls    string               `json:"imageUrls"`
	// Images       string               `json:"images"`
	MainImageID    string               `json:"mainImageId"`
	VehicleType    string               `json:"vehicleType"`
	MainImageURL   string               `json:"mainImageUrl"`
	Name           string               `json:"name"`
	Note           string               `json:"note"`
	Imgs           ListImgKeyValueArray `json:"imgs" gorm:"embedded;type:jsonb"`
	LastAppearance time.Time            `json:"lastAppearance" gorm:"column:last_appearance"`
	Status         bool                 `json:"status" gorm:"column:status;default:false"`
	DeleteMark     bool                 `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	CreatedAt      time.Time            `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt      time.Time            `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt      time.Time            `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_LicensePlates struct {
	ID uuid.UUID `json:"id"`
	// ImageUrls    string               `json:"imageUrls"`
	// Images       string               `json:"images"`
	MainImageID    string               `json:"mainImageId"`
	MainImageURL   string               `json:"mainImageUrl"`
	VehicleType    string               `json:"vehicleType"`
	Name           string               `json:"name"`
	Status         bool                 `json:"status"`
	Note           string               `json:"note"`
	LastAppearance time.Time            `json:"lastAppearance"`
	Imgs           ListImgKeyValueArray `json:"imgs,omitempty"`
	DeleteMark     bool                 `json:"deletedMark"`
}

type DTO_LicensePlates_Ids struct {
	IDs []string `json:"ids"`
}
