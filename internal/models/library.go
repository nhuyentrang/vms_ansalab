package models

import (
	"time"

	"github.com/google/uuid"
)

type Library struct {
	ID                  uuid.UUID `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Atts                time.Time `json:"atts"`
	CompanyID           *string   `json:"company_id"`
	GetOriginalFilename *string   `json:"get_original_filename"`
	Hash                *string   `json:"hash"`
	IsPublic            bool      `json:"is_public"`
	Name                string    `json:"name"`
	Size                int       `json:"size"`
	Type                string    `json:"type"`
	UserID              string    `json:"user_id"`
	Username            string    `json:"username"`
}

type DTO_Library struct {
	ID                  uuid.UUID `json:"id"`
	CreatedAt           time.Time `json:"created_at"`
	CreatedBy           string    `json:"created_by"`
	UpdatedAt           time.Time `json:"updated_at"`
	UpdatedBy           string    `json:"updated_by"`
	Atts                time.Time `json:"atts"`
	CompanyID           *string   `json:"company_id"`
	GetOriginalFilename *string   `json:"get_original_filename"`
	Hash                *string   `json:"hash"`
	IsPublic            bool      `json:"is_public"`
	Name                string    `json:"name"`
	Size                int       `json:"size"`
	Type                string    `json:"type"`
	UserID              string    `json:"user_id"`
	Username            string    `json:"username"`
}
