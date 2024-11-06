package models

import (
	"time"

	uuid "github.com/google/uuid"
	//"github.com/lib/pq"
)

// Entity Model for device
type RecordingSchedule struct {
	ID          uuid.UUID      `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Sunday      TimeBlockArray `json:"sunday" gorm:"embedded;type:jsonb"`
	Monday      TimeBlockArray `json:"monday" gorm:"embedded;type:jsonb"`
	Tuesday     TimeBlockArray `json:"tuesday" gorm:"embedded;type:jsonb"`
	Wednesday   TimeBlockArray `json:"wednesday" gorm:"embedded;type:jsonb"`
	Thursday    TimeBlockArray `json:"thursday" gorm:"embedded;type:jsonb"`
	Saturday    TimeBlockArray `json:"saturday" gorm:"embedded;type:jsonb"`
	Friday      TimeBlockArray `json:"friday" gorm:"embedded;type:jsonb"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_RecordingSchedule struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Sunday      TimeBlockArray `json:"sunday"`
	Monday      TimeBlockArray `json:"monday"`
	Tuesday     TimeBlockArray `json:"tuesday"`
	Wednesday   TimeBlockArray `json:"wednesday"`
	Thursday    TimeBlockArray `json:"thursday"`
	Saturday    TimeBlockArray `json:"saturday"`
	Friday      TimeBlockArray `json:"friday"`
}

type DTO_RecordingSchedule_Create struct {
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Sunday      TimeBlockArray `json:"sunday"`
	Monday      TimeBlockArray `json:"monday"`
	Tuesday     TimeBlockArray `json:"tuesday"`
	Wednesday   TimeBlockArray `json:"wednesday"`
	Thursday    TimeBlockArray `json:"thursday"`
	Saturday    TimeBlockArray `json:"saturday"`
	Friday      TimeBlockArray `json:"friday"`
}

type DTO_RecordingSchedule_Read_BasicInfo struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Sunday      TimeBlockArray `json:"sunday"`
	Monday      TimeBlockArray `json:"monday"`
	Tuesday     TimeBlockArray `json:"tuesday"`
	Wednesday   TimeBlockArray `json:"wednesday"`
	Thursday    TimeBlockArray `json:"thursday"`
	Saturday    TimeBlockArray `json:"saturday"`
	Friday      TimeBlockArray `json:"friday"`
}
