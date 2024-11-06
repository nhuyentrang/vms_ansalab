package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type BlackList struct {
	ID uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	// ImageUrls    string               `json:"imageUrls"`
	// Images       string               `json:"images"`
	MainImageID    string               `json:"mainImageId"`
	MainImageURL   string               `json:"mainImageUrl"`
	Name           string               `json:"name"`
	Note           string               `json:"note"`
	Imgs           ListImgKeyValueArray `json:"imgs" gorm:"embedded;type:jsonb"`
	LastAppearance time.Time            `json:"lastAppearance" gorm:"column:last_appearance;autoCreateTime:true"`
	Status         bool                 `json:"status" gorm:"column:status;default:false"`
	Type           KeyValue             `json:"type" gorm:"embedded;type:jsonb"`
	DeleteMark     bool                 `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	CreatedAt      time.Time            `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt      time.Time            `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt      time.Time            `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_BlackList struct {
	ID uuid.UUID `json:"id"`
	// ImageUrls    string               `json:"imageUrls"`
	// Images       string               `json:"images"`
	MainImageID    string               `json:"mainImageId"`
	MainImageURL   string               `json:"mainImageUrl"`
	Name           string               `json:"name"`
	Status         bool                 `json:"status"`
	LastAppearance time.Time            `json:"lastAppearance"`
	Note           string               `json:"note"`
	Imgs           ListImgKeyValueArray `json:"imgs,omitempty"`
	DeleteMark     bool                 `json:"deletedMark"`
	Type           KeyValue             `json:"type,omitempty"`
}

type DTO_BlackList_Created struct {
	// ImageUrls    string               `json:"imageUrls"`
	// Images       string               `json:"images"`
	ID           uuid.UUID            `json:"id"`
	MainImageID  string               `json:"mainImageId"`
	MainImageURL string               `json:"mainImageUrl"`
	Name         string               `json:"name"`
	Note         string               `json:"note"`
	Imgs         ListImgKeyValueArray `json:"imgs,omitempty"`
	Status       bool                 `json:"status"`
	Type         KeyValue             `json:"type,omitempty"`
}

type DTO_BlackList_Edit struct {
	Status       bool                 `json:"status" gorm:"column:status;default:false"`
	Name         string               `json:"name"`
	MainImageID  string               `json:"mainImageId"`
	MainImageURL string               `json:"mainImageUrl"`
	Note         string               `json:"note"`
	Imgs         ListImgKeyValueArray `json:"imgs" gorm:"embedded;type:jsonb"`
	Type         KeyValue             `json:"type,omitempty"`
}

type DTO_BlackList_Ids struct {
	IDs []string `json:"ids"`
}

// Face Regesiter
type FaceRegData struct {
	Member        int            `json:"member"`
	Blacklist     int            `json:"blacklist"`
	ImageInfo     []ImageDetails `json:"imageInfo"`
	MemberID      string         `json:"userId"`
	Threshold     int            `json:"threshold"`
	Topk          int            `json:"topk"`
	IsForceUpdate bool           `json:"isForceUpdate"`
	TenantID      string         `json:"tenantId,omitempty"`
}

type ImageDetails struct {
	URL string `json:"url"`
	ID  string `json:"id"`
	// Type string `json:"type"`
}

type DeleteFaceVectorRequest struct {
	Member    int    `json:"member"`
	Blacklist int    `json:"blacklist"`
	MemberID  string `json:"memberId"`
	TenantID  string `json:"tenantId,omitempty"`
}

type DeleteFaceVectorResponse struct {
	Code     int    `json:"code"`
	MemberID string `json:"member_id"`
}

type DataRespFaceReg struct {
	Code              int              `json:"code"`
	Status            string           `json:"status"`
	Detail            string           `json:"detail"`
	MemberID          string           `json:"member_id"`
	StandardDeviation float64          `json:"standard_deviation"`
	TopkBlacklists    []TopkBlacklists `json:"topk_blacklists"`
}

type TopkBlacklists struct {
	Dataistance float64 `json:"distance"`
	MemberID    string  `json:"member_id"`
}

type TopkMember struct {
	Distance float64 `json:"distance"`
	MemberID string  `json:"member_id"`
}

type UserDetail struct {
	UserID      string `json:"userId"`
	Username    string `json:"username"`
	FullName    string `json:"fullName"`
	Gender      string `json:"gender"`
	PhoneNumber string `json:"phoneNumber"`
	Email       string `json:"email"`
}

type MemberFaceImageData struct {
	ID           string `json:"id"`
	FaceFeature  string `json:"faceFeature"`
	FaceType     string `json:"faceType"`
	FeatureType  string `json:"featureType"`
	ImageBase64  string `json:"imageBase64"`
	ImageFileID  string `json:"imageFileId"`
	ImageFileURL string `json:"imageFileUrl"`
	UserID       string `json:"userId"`
}

// Prepare the response
type VectorSearchMemberData struct {
	Distance   float64               `json:"distance"`
	MemberID   string                `json:"member_id"`
	Data       []MemberFaceImageData `json:"data"`
	UserDetail UserDetail            `json:"user_detail"`
}
type APISearchBlacklistResponse struct {
	TopkMembers    []VectorSearchMemberData `json:"topk_members"`
	TopkBlacklists []VectorSearchMemberData `json:"topk_blacklists"`
}
