package models

import (
	"time"

	uuid "github.com/google/uuid"
)

// TableName overrides the table name
func (Cabin) TableName() string {
	return "ivis_vms.cabin_tbl"
}

// Entity Model for cabin
type Cabin struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	// General info
	CabinCode       string `json:"cabinCode" gorm:"column:cabin_code"`
	SerialNumber    string `json:"serialNumber" gorm:"column:serial_number"`
	MACAddress      string `json:"macAddress" gorm:"column:mac_address"`
	IPAddress       string `json:"ipAddress" gorm:"column:ip_address"`
	SensorCount     int    `json:"sensorCount" gorm:"column:sensor_count"`
	CameraCount     int    `json:"cameraCount" gorm:"column:camera_count"`
	Source          string `json:"source" gorm:"column:source"`
	Connect         int    `json:"connect" gorm:"column:connect"`
	PeripheralID    string `json:"peripheralID" gorm:"column:peripheral_id"`
	Username        string `json:"username" gorm:"column:username"`
	Password        string `json:"password" gorm:"column:password"`
	FirmwareVersion string `json:"firmwareVersion" gorm:"column:firmware_version"`

	// Location info
	SiteInfo       KeyValue `json:"siteInfo" gorm:"column:site_info;embedded;type:jsonb"`
	Location       string   `json:"location" gorm:"column:location"`
	LocationDetail string   `json:"locationDetail" gorm:"column:location_detail"`
	Coordinate     string   `json:"coordinate" gorm:"column:coordinate"`
	Position       string   `json:"position" gorm:"column:position"`
	Longitude      string   `json:"longitude" gorm:"column:longitude"`
	Latitude       string   `json:"latitude" gorm:"column:latitude"`
	Area           string   `json:"area" gorm:"column:area"`

	// For compatible with controlbox
	AlertStatus                string `json:"alertStatus" gorm:"column:alert_status"`
	StatusAI                   string `json:"statusAI" gorm:"column:statusai"`
	VibrationSensorLevel       string `json:"vibrationSensorLevel" gorm:"column:vibration_sensor_level"`
	ElectricLeakageSensorLevel string `json:"electronicLeakSensorLevel" gorm:"column:electronic_leak_sensor_level"`
	MovementSensorLevel        string `json:"movementSensorLevel" gorm:"column:movement_sensor_level"`
	DoorSensorLevel            string `json:"doorSensorLevel" gorm:"column:door_sensor_level"`
	TemperatureSensorLevel     string `json:"temperatureSensorLevel" gorm:"column:temperature_sensor_level"`
	SmokeSensorLevel           string `json:"smokeSensorLevel" gorm:"column:smoke_sensor_level"`
	FireSensorLevel            string `json:"fireSensorLevel" gorm:"column:fire_sensor_level"`
	AlarmSensorLevel           string `json:"alarmSensorLevel" gorm:"column:alarm_sensor_level"`
	GpsSensorLevel             string `json:"gpsSensorLevel" gorm:"column:gps_sensor_level"`

	// Action
	Action                      string `json:"action"`
	LampAction                  string `json:"lampAction" gorm:"column:lamp_action"`
	VibrationSensorAction       string `json:"vibrationSensorAction" gorm:"column:vibration_sensor_action"`
	ElectricLeakageSensorAction string `json:"electronicLeakSensorAction" gorm:"column:electronic_leak_sensor_action"`
	MovementSensorAction        string `json:"movementSensorAction" gorm:"column:movement_sensor_action"`
	DoorSensorAction            string `json:"doorSensorAction" gorm:"column:door_sensor_action"`
	TemperatureSensorAction     string `json:"temperatureSensorAction" gorm:"column:temperature_sensor_action"`
	SmokeSensorAction           string `json:"smokeSensorAction" gorm:"column:smoke_sensor_action"`
	FireSensorAction            string `json:"fireSensorAction" gorm:"column:fire_sensor_action"`
	AlarmSensorAction           string `json:"alarmSensorAction" gorm:"column:alarm_sensor_action"`
	GpsSensorAction             string `json:"gpsSensorAction" gorm:"column:gps_sensor_action"`
	AIEngine                    string `json:"aiEngine" gorm:"column:ai_engine"`
	LockDoor                    string `json:"lockDoor" gorm:"column:lock_door"`
	AccessControl               string `json:"accessControl" gorm:"column:access_control"`
	DeleteMark                  bool   `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	DeviceUpdate                bool   `json:"deviceUpdate" gorm:"column:device_update;default:false"`

	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_Cabin struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	Area           string    `json:"area,omitempty"`
	CabinCode      string    `json:"cabinCode,omitempty"`
	CameraCount    int       `json:"cameraCount,omitempty"`
	Connect        int       `json:"connect,omitempty"`
	Description    string    `json:"description,omitempty"`
	Latitude       string    `json:"latitude,omitempty"`
	Location       string    `json:"location,omitempty"`
	LocationDetail string    `json:"locationDetail,omitempty"`
	Longitude      string    `json:"longitude,omitempty"`
	MacAddress     string    `json:"macAddress,omitempty"`
	SensorCount    int       `json:"sensorCount,omitempty"`
	SerialNumber   string    `json:"serialNumber,omitempty"`
	Source         string    `json:"source,omitempty"`
	// For compatible with controlbox / Status
	AlertStatus string `json:"alertStatus,omitempty"`
	StatusAI    string `json:"statusAI,omitempty"`
	// For compatible with controlbox / Sensivity level
	VibrationSensorLevel       string `json:"vibrationSensorLevel,omitempty"`
	ElectricLeakageSensorLevel string `json:"electronicLeakSensorLevel,omitempty"`
	MovementSensorLevel        string `json:"movementSensorLevel,omitempty"`
	DoorSensorLevel            string `json:"doorSensorLevel,omitempty"`
	TemperatureSensorLevel     string `json:"temperatureSensorLevel,omitempty"`
	SmokeSensorLevel           string `json:"smokeSensorLevel,omitempty"`
	FireSensorLevel            string `json:"fireSensorLevel,omitempty"`
	AlarmSensorLevel           string `json:"alarmSensorLevel,omitempty"`
	GpsSensorLevel             string `json:"gpsSensorLevel,omitempty"`

	// For compatible with controlbox / Action
	Action                      string `json:"action,omitempty"`
	LampAction                  string `json:"lampAction,omitempty"`
	VibrationSensorAction       string `json:"vibrationSensorAction,omitempty"`
	ElectricLeakageSensorAction string `json:"electronicLeakSensorAction,omitempty"`
	MovementSensorAction        string `json:"movementSensorAction,omitempty"`
	DoorSensorAction            string `json:"doorSensorAction,omitempty"`
	TemperatureSensorAction     string `json:"temperatureSensorAction,omitempty"`
	SmokeSensorAction           string `json:"smokeSensorAction,omitempty"`
	FireSensorAction            string `json:"fireSensorAction,omitempty"`
	AlarmSensorAction           string `json:"alarmSensorAction,omitempty"`
	GpsSensorAction             string `json:"gpsSensorAction,omitempty"`
	AIEngine                    string `json:"aiEngine,omitempty"`
	LockDoor                    string `json:"lockDoor,omitempty"`
	AccessControl               string `json:"accessControl,omitempty"`
	DeviceUpdate                bool   `json:"deviceUpdate"`
}

type DTO_Cabin_Create struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	AccessControl  string    `json:"accessControl,omitempty"`
	AlertStatus    string    `json:"alertStatus,omitempty"`
	Area           string    `json:"area,omitempty"`
	CabinCode      string    `json:"cabinCode,omitempty"`
	CameraCount    int       `json:"cameraCount,omitempty"`
	Connect        int       `json:"connect,omitempty"`
	Description    string    `json:"description,omitempty"`
	Latitude       string    `json:"latitude,omitempty"`
	Location       string    `json:"location,omitempty"`
	LocationDetail string    `json:"locationDetail,omitempty"`
	LockDoor       string    `json:"lockDoor,omitempty"`
	Longitude      string    `json:"longitude,omitempty"`
	MacAddress     string    `json:"macAddress,omitempty"`
	SensorCount    int       `json:"sensorCount,omitempty"`
	SerialNumber   string    `json:"serialNumber,omitempty"`
	Source         string    `json:"source,omitempty"`
	StatusAI       string    `json:"statusAI,omitempty"`
}

type DTO_Cabin_Read_BasicInfo struct {
	ID             uuid.UUID `json:"id,omitempty"`
	Name           string    `json:"name,omitempty"`
	AccessControl  string    `json:"accessControl,omitempty"`
	AlertStatus    string    `json:"alertStatus,omitempty"`
	Area           string    `json:"area,omitempty"`
	CabinCode      string    `json:"cabinCode,omitempty"`
	CameraCount    int       `json:"cameraCount,omitempty"`
	Connect        int       `json:"connect,omitempty"`
	Description    string    `json:"description,omitempty"`
	Latitude       string    `json:"latitude,omitempty"`
	Location       string    `json:"location,omitempty"`
	LocationDetail string    `json:"locationDetail,omitempty"`
	LockDoor       string    `json:"lockDoor,omitempty"`
	Longitude      string    `json:"longitude,omitempty"`
	MacAddress     string    `json:"macAddress,omitempty"`
	ModelAIs       []string  `json:"modelAIs,omitempty"`
	Peripheral     []string  `json:"peripheral,omitempty"`
	SensorCount    int       `json:"sensorCount,omitempty"`
	SerialNumber   string    `json:"serialNumber,omitempty"`
	Source         string    `json:"source,omitempty"`
	StatusAI       string    `json:"statusAI,omitempty"`
	Type           string    `json:"type"`
}

type DTO_Cabin_Read_StatisticMapLabel struct {
	// Status
	NumberRunning      int64 `json:"numberRunning,omitempty"`
	NumberStopped      int64 `json:"numberStopped,omitempty"`
	NumberError        int64 `json:"numberError,omitempty"`
	NumberDisconnected int64 `json:"numberDisconnected,omitempty"`
}

type DTO_UpdateCabinTimeOut struct {
	Connect     int    `json:"connect,omitempty"`
	AlertStatus string `json:"alertStatus,omitempty"`
}

type DTO_CabinGPS struct {
	CabinID   uuid.UUID `json:"cabinID,omitempty"`
	Longitude string    `json:"longitude,omitempty"`
	Latitude  string    `json:"latitude,omitempty"`
	Name      string    `json:"name,omitempty"`
	Source    string    `json:"source,omitempty"`
}
