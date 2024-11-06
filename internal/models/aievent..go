package models

import (
	"time"

	uuid "github.com/google/uuid"
)

type DTO_AI_Event struct {
	ID              uuid.UUID `json:"id,omitempty"`
	MessageID       string    `json:"messageID,omitempty"`
	MsVersion       string    `json:"msVersion,omitempty"`
	SensorID        string    `json:"sensorID,omitempty"`
	Description     string    `json:"description,omitempty"`
	Timestamp       int64     `json:"timestamp,omitempty"`
	TimeStart       int64     `json:"timeStart,omitempty"`
	TimeEnd         int64     `json:"timeEnd,omitempty"`
	Image           string    `json:"image,omitempty"`
	ImageResult     string    `json:"imageResult,omitempty"`
	ImageObject     string    `json:"imageObject,omitempty"`
	Video           string    `json:"video,omitempty"`
	StorageBucket   string    `json:"storageBucket,omitempty"`
	EventType       string    `json:"eventType,omitempty"`
	CamIP           string    `json:"camIP,omitempty"`
	CamName         string    `json:"camName,omitempty"`
	CameraId        string    `json:"cameraId,omitempty"`
	MemberID        string    `json:"memberID,omitempty"`
	EventTypeString string    `json:"eventTypeString,omitempty"`
	Location        string    `json:"location,omitempty"`
	Status          string    `json:"status,omitempty"`
	CabinID         uuid.UUID `json:"cabinID,omitempty"`
	CabinName       string    `json:"cabinName,omitempty"`
	ConverTimestamp time.Time `json:"converTimestamp,omitempty"`
	TypeOfAIEvent   string    `json:"TypeOfAIEvent,omitempty"`
	Longtitude      float64   `json:"LongtitudeOfCam"`
	Latitude        float64   `json:"LatitudeOfCam"`
}

type DTO_AIEvent_Chart struct {
	DataAIEvent []DTO_AIEvent_Count `json:"dataAIEvent,omitempty"`
}

type DTO_AIEvent_Count struct {
	Date      time.Time `json:"date"`
	SABOTAGE  int       `json:"sabotage"`
	DANGEROUS int       `json:"dangerous"`
	ABNORMAL  int       `json:"abnormal"`
	Location  int       `json:"location"`
}

type DTO_AIEvent_Image struct {
	Time        time.Time               `json:"time,omitempty"`
	DataAIEvent []DTO_AIEvent_ImageItem `json:"dataAIEvent,omitempty"`
}

type DTO_AIEvent_ImageItem struct {
	Image         string `json:"image,omitempty"`
	ImageResult   string `json:"imageResult,omitempty"`
	ImageObject   string `json:"imageObject,omitempty"`
	Video         string `json:"video,omitempty"`
	StorageBucket string `json:"storageBucket,omitempty"`
}

type DTO_CameraStatus struct {
	CabinID             uuid.UUID `json:"cabinID,omitempty"`
	CameraID            uuid.UUID `json:"cameraID,omitempty"`
	CameraStatusConnect string    `json:"cameraStatusConnect,omitempty"`
	CameraName          string    `json:"cameraName,omitempty"`
	EpochTime           int       `json:"epochTime,omitempty"`
}

type Top5Event struct {
	Name  string `json:"name,omitempty"`
	Count string `json:"count,omitempty"`
}

// Entity Model
type AIWaring struct {
	ID            uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	MessageID     string    `json:"messageID,omitempty"`
	MsVersion     string    `json:"msVersion,omitempty"`
	SensorID      string    `json:"sensorID,omitempty"`
	Description   string    `json:"description,omitempty"`
	Timestamp     int64     `json:"timestamp,omitempty"`
	TimeStart     int64     `json:"timeStart,omitempty"`
	TimeEnd       int64     `json:"timeEnd,omitempty"`
	Image         string    `json:"image,omitempty"`
	ImageResult   string    `json:"imageResult,omitempty"`
	ImageObject   string    `json:"imageObject,omitempty"`
	Video         string    `json:"video,omitempty"`
	StorageBucket string    `json:"storageBucket,omitempty"`
	EventType     string    `json:"eventType,omitempty"`
	CamIP         string    `json:"camIP,omitempty"`
	CamName       string    `json:"camName,omitempty"`
	CameraId      string    `json:"cameraId,omitempty"`
	MemberID      string    `json:"memberID,omitempty"`

	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeleteMark bool      `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	DeletedAt  time.Time `json:"deletedAt" gorm:"column:deleted_at"`

	EventTypeString string    `json:"eventTypeString,omitempty"`
	Location        string    `json:"location,omitempty"`
	Status          string    `json:"status,omitempty"`
	CabinName       string    `json:"cabinName,omitempty"`
	ConverTimestamp time.Time `json:"converTimestamp,omitempty"`
	CabinID         uuid.UUID `json:"cabinID,omitempty" gorm:"column:cabin_id"`
	TypeOfAIEvent   string    `json:"typeOfAIEvent,omitempty"`
	Longtitude      float64   `json:"longtitudeOfCam"`
	Latitude        float64   `json:"latitudeOfCam"`

	Result string `json:"result"`
}

type DTO_AIWaring struct {
	ID            uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	MessageID     string    `json:"messageID,omitempty"`
	MsVersion     string    `json:"msVersion,omitempty"`
	SensorID      string    `json:"sensorID,omitempty"`
	Description   string    `json:"description,omitempty"`
	Timestamp     int64     `json:"timestamp,omitempty"`
	TimeStart     int64     `json:"timeStart,omitempty"`
	TimeEnd       int64     `json:"timeEnd,omitempty"`
	Image         string    `json:"image,omitempty"`
	ImageResult   string    `json:"imageResult,omitempty"`
	ImageObject   string    `json:"imageObject,omitempty"`
	Video         string    `json:"video,omitempty"`
	StorageBucket string    `json:"storageBucket,omitempty"`
	EventType     string    `json:"eventType,omitempty"`
	CamIP         string    `json:"camIP,omitempty"`
	CamName       string    `json:"camName,omitempty"`
	CameraId      string    `json:"cameraId,omitempty"`
	MemberID      string    `json:"memberID,omitempty"`

	CreatedAt       time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt       time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeleteMark      bool      `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	DeletedAt       time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	Longtitude      float64   `json:"longtitudeOfCam"`
	Latitude        float64   `json:"latitudeOfCam"`
	EventTypeString string    `json:"eventTypeString,omitempty"`
	Location        string    `json:"location,omitempty"`
	Status          string    `json:"status,omitempty"`
	CabinName       string    `json:"cabinName,omitempty"`
	ConverTimestamp time.Time `json:"converTimestamp,omitempty"`
	CabinID         uuid.UUID `json:"cabinID,omitempty" gorm:"column:cabin_id"`
	Result          string    `json:"result,omitempty" gorm:"column:result"`
}

type DTO_AIWaring_Update struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Description string    `json:"description,omitempty"` //
	CamIP       string    `json:"camIP,omitempty"`       //
	CamName     string    `json:"camName,omitempty"`     //
	CameraId    string    `json:"cameraId,omitempty"`
	MemberID    string    `json:"memberID,omitempty"` //

	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeleteMark bool      `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	DeletedAt  time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	Location   string    `json:"location,omitempty"`                    //
	Result     string    `json:"result,omitempty" gorm:"column:result"` //
}

// DTO Model
type DTO_Report struct {
	ID         uuid.UUID `json:"id,omitempty"`
	CabinID    uuid.UUID `json:"cabinId,omitempty"`
	SensorID   uuid.UUID `json:"sensor_id,omitempty"`
	Deleted    bool      `json:"deleted,omitempty"`
	Status     string    `json:"status,omitempty"`
	AlertDate  time.Time `json:"alertDate,omitempty"`
	AlertType  string    `json:"alertType,omitempty"`
	CabinName  string    `json:"cabinName,omitempty"`
	DeviceName string    `json:"deviceName,omitempty"`
	Location   string    `json:"location,omitempty"`
}

// DTO Model
type DTO_SystemWaring struct {
	ID         uuid.UUID `json:"id,omitempty"`
	CabinID    uuid.UUID `json:"cabinId,omitempty"`
	SensorID   uuid.UUID `json:"sensor_id,omitempty"`
	Deleted    bool      `json:"deleted,omitempty"`
	Status     string    `json:"status,omitempty"`
	AlertDate  time.Time `json:"alertDate,omitempty"`
	AlertType  string    `json:"alertType,omitempty"`
	CabinName  string    `json:"cabinName,omitempty"`
	DeviceName string    `json:"deviceName,omitempty"`
	Location   string    `json:"location,omitempty"`
}

// Entity Model
type Report struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CabinID    uuid.UUID `json:"cabinID" gorm:"type:uuid;default:uuid_generate_v4();column:cabin_id"`
	SensorID   uuid.UUID `json:"sensorID" gorm:"type:uuid;default:uuid_generate_v4();column:sensor_id"`
	Deleted    bool      `json:"deleted" gorm:"column:deleted"`
	Status     string    `json:"status" gorm:"column:status"`
	AlertDate  time.Time `json:"alertDate" gorm:"column:alert_date;autoCreateTime:true"`
	AlertType  string    `json:"alertType" gorm:"column:alert_type"`
	CabinName  string    `json:"cabinName" gorm:"column:cabin_name"`
	DeviceName string    `json:"deviceName" gorm:"column:device_name"`
	Location   string    `json:"location" gorm:"column:location"`

	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

// Entity Model
type SystemWaring struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	CabinID    uuid.UUID `json:"cabinID" gorm:"type:uuid;default:uuid_generate_v4();column:cabin_id"`
	SensorID   uuid.UUID `json:"sensorID" gorm:"type:uuid;default:uuid_generate_v4();column:sensor_id"`
	Deleted    bool      `json:"deleted" gorm:"column:deleted"`
	Status     string    `json:"status" gorm:"column:status"`
	AlertDate  time.Time `json:"alertDate" gorm:"column:alert_date;autoCreateTime:true"`
	AlertType  string    `json:"alertType" gorm:"column:alert_type"`
	CabinName  string    `json:"cabinName" gorm:"column:cabin_name"`
	DeviceName string    `json:"deviceName" gorm:"column:device_name"`
	Location   string    `json:"location" gorm:"column:location"`

	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}
