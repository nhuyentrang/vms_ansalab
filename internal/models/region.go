package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type Region struct {
	ID uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`

	Type  string `json:"type,omitempty"`
	Code  string `json:"code,omitempty"`
	Name  string `json:"name,omitempty"`
	Level string `json:"level,omitempty"`

	ParentID   string `json:"parentID,omitempty"`
	ParentName string `json:"parentName,omitempty"`
	ParentCode string `json:"parentCode,omitempty"`
	Path       string `json:"path,omitempty"`
	IsParent   bool   `json:"isParent"`
	Latitude   string `json:"latitude,omitempty"`
	Longitude  string `json:"longitude,omitempty"`
	Detail     string `json:"detail,omitempty"`

	IndentifedNumber string `json:"identifiedNumber"`
	Deleted          bool   `json:"deleted,omitempty" default:"false"`

	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt,omitempty"`
}

type DTORegion struct {
	ID uuid.UUID `json:"id"`

	Type  string `json:"type,omitempty"`
	Code  string `json:"code,omitempty"`
	Name  string `json:"name,omitempty"`
	Level string `json:"level,omitempty"`

	ParentID   string `json:"parentID,omitempty"`
	ParentName string `json:"parentName,omitempty"`
	ParentCode string `json:"parentCode,omitempty"`
	Path       string `json:"path,omitempty"`
	IsParent   bool   `json:"isParent"`

	Latitude  string `json:"latitude,omitempty"`
	Longitude string `json:"longitude,omitempty"`
	Detail    string `json:"detail,omitempty"`
}
