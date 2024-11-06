package models

/*
// Table for reference between camera and AI engine (AI planning)
type CamAIEnginePlanning struct {
	//ID         uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CameraID   uuid.UUID `json:"cameraID" gorm:"primary_key;type:uuid;default:uuid_generate_v4();column:camera_id"`
	AIEngineID uuid.UUID `json:"aiEngineID" gorm:"primary_key;type:uuid;default:uuid_generate_v4();column:ai_engine_id"`
	StreamURL  string    `json:"streamURL" gorm:"column:stream_url"`

	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}
*/
