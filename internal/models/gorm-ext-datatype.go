package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// Tabler interface for set custom table name
type Tabler interface {
	TableName() string
}

/********* Keyvalue ***********/
type KeyValue struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Channel string `json:"channel,omitempty"`
	Code    string `json:"code,omitempty"`
}

func (sla *KeyValue) Scan(src interface{}) error {
	// For gorm:"type:jsonb"
	return json.Unmarshal(src.([]byte), &sla)
	// For gorm:"type:text"
	//return json.Unmarshal([]byte(src.(string)), &sla)
}

func (sla KeyValue) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* KeyValueArray ***********/
type KeyValueArray []KeyValue

func (sla *KeyValueArray) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, &sla)
	case string:
		return json.Unmarshal([]byte(src), &sla)
	default:
		return errors.New("unsupported type")
	}
}

func (sla KeyValueArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* CameraKeyValue List (Todo: fix front end for use KeyValue) ***********/
type CameraKeyValue struct {
	ID         string `json:"id"`
	DeviceName string `json:"deviceName"`
}

func (sla *CameraKeyValue) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla CameraKeyValue) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type CameraKeyValueArray []CameraKeyValue

func (sla *CameraKeyValueArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla CameraKeyValueArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* VideoStream List ***********/
type VideoStream struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	Type      string `json:"type,omitempty"`
	URL       string `json:"url,omitempty"`
	IsProxied bool   `json:"isProxied,omitempty"`
	IsDefault bool   `json:"isDefault,omitempty"`
	Channel   string `json:"channel"`
	Codec     string `json:"codec"`
}

func (sla *VideoStream) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla VideoStream) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type VideoStreamArray []VideoStream

func (sla *VideoStreamArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla VideoStreamArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* TimeBlock ***********/
type TimeBlock struct {
	Begin int64 `json:"begin"`
	End   int64 `json:"end"`
}

func (sla *TimeBlock) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla TimeBlock) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* TimeBlockArray ***********/
type TimeBlockArray []TimeBlock

func (sla *TimeBlockArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla TimeBlockArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* CameraKeyValue List (Todo: fix front end for use KeyValue) ***********/
type ListImgKeyValue struct {
	ID            string `json:"id"`
	StorageBucket string `json:"storageBucket,omitempty"`
	URLImage      string `json:"urlImage"`
	Type          string `json:"type"`
	ImageName     string `json:"name"`
}

func (sla *ListImgKeyValue) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ListImgKeyValue) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type ListImgKeyValueArray []ListImgKeyValue

func (sla *ListImgKeyValueArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ListImgKeyValueArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Network List (Todo: fix front end for use KeyValue) ***********/
func (sla *NetworkServer) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla NetworkServer) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type NetworkValueArray []NetworkServer

func (sla *NetworkValueArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla NetworkValueArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* NTPServer List (Todo: fix front end for use KeyValue) ***********/
func (sla *NTPServer) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla NTPServer) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* DDNSServer List (Todo: fix front end for use KeyValue) ***********/
func (sla *DDNSServer) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla DDNSServer) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Network List (Todo: fix front end for use KeyValue) ***********/
func (sla *AdminAccessProtocolServerList) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla AdminAccessProtocolServerList) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Time List (Todo: fix front end for use KeyValue) ***********/
func (sla *Time) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla Time) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* VideoOverlayServer List (Todo: fix front end for use KeyValue) ***********/
func (sla *VideoOverlayServer) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla VideoOverlayServer) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* StreamingChannel List (Todo: fix front end for use KeyValue) ***********/
func (sla *StreamingChannel) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla StreamingChannel) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* TimeBlock ***********/
type TimeSchedule struct {
	Begin int64 `json:"begin"`
	End   int64 `json:"end"`
}

func (sla *TimeSchedule) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla TimeSchedule) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* TimeBlockArray ***********/
type TimeScheduleArray struct {
	Day          string
	TimeSchedule []TimeSchedule
}

func (sla *TimeScheduleArray) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla TimeScheduleArray) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* UserCameraLocation ***********/
type UserCameraLocation struct {
	Location   int    `json:"location"`
	ID         string `json:"id"`
	DeviceName string `json:"deviceName"`
}

func (sla *UserCameraLocation) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla UserCameraLocation) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* UserCamera ***********/
type UserCameraView struct {
	DisplayName string               `json:"displayName"`
	Locations   []UserCameraLocation `json:"locations"`
}

func (sla *UserCameraView) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla UserCameraView) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type ListUserCameraView []UserCameraView

func (sla *ListUserCameraView) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ListUserCameraView) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Point ***********/
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (sla *Point) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla Point) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Virtual Fence ***********/
type VFence struct {
	DX     Point `json:"dx"`
	DY     Point `json:"dy"`
	Vector int   `json:"vector"`
}

func (sla *VFence) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla VFence) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type ListVFence []VFence

func (sla *ListVFence) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla ListVFence) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Virtual Fence ***********/
type VZoneCoordinate struct {
	Point_A Point `json:"point_a"`
	Point_B Point `json:"point_b"`
	Point_C Point `json:"point_c"`
	Point_D Point `json:"point_d"`
}

func (sla *VZoneCoordinate) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla VZoneCoordinate) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Virtual Fence ***********/
// type CameraSchedule struct {
// 	MondaySchedule    []TimeTable `json:"mondayschedule"`
// 	TuesdaySchedule   []TimeTable `json:"tuesdayschedule"`
// 	WednesdaySchedule []TimeTable `json:"wednesdayschedule"`
// 	ThursdaySchedule  []TimeTable `json:"thursdayschedule"`
// 	FridaySchedule    []TimeTable `json:"fridayschedule"`
// }
type TimePeriod struct {
	StartTimeInMinute int `json:"startTimeInMinute"`
	EndTimeInMinute   int `json:"endTimeInMinute"`
}

type DaySchedule struct {
	DayOfWeek   string       `json:"dayOfWeek"`
	TimePeriods []TimePeriod `json:"timePeriods"`
}

type CalendarDays struct {
	DayOfWeek   string       `json:"dayOfWeek"`
	TimePeriods []TimePeriod `json:"timePeriods"`
}

// type CalenderDays struct {
// 	DayOfWeek   string      `json:"dayOfWeek"`
// 	TimePeriods []TimeTable `json:"timePeriods"`
// }

func (sla *CalendarDays) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), sla)
}

func (sla CalendarDays) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Virtual Fence ***********/
// type TimeTable struct {
// 	StartTimeMinute int `json:"startTimeMinute"`
// 	EndTimeMinute   int `json:"endTimeMinute"`
// }

func (sla *TimePeriod) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), sla)
}

func (sla TimePeriod) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Camera AI Property ***********/
type CameraAIProperty struct {
	// PropertyID    uuid.UUID             `json:"propertyid,omitempty"`
	// AIType        string                `json:"aitype"`
	// Sensitivity   int                   `json:"sensitivity"`
	// Target        string                `json:"target"`
	CameraAIZone  CameraVirtualProperty `json:"cameraaizone"`
	CameraModelAI CameraModelAI         `json:"cameraModelAI"`
	CalendarDays  []CalendarDays        `json:"calendarDays"`
	IsActive      bool                  `json:"isactive"`
}

func (sla *CameraAIProperty) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla CameraAIProperty) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type CameraAIPropertyList []CameraAIProperty

func (sla *CameraAIPropertyList) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla CameraAIPropertyList) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Camera Virtual Property ***********/
type CameraVirtualProperty struct {
	Vfences ListVFence      `json:"vfences"`
	Vzone   VZoneCoordinate `json:"vzone"`
}

func (sla *CameraVirtualProperty) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla CameraVirtualProperty) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Camera model ai ***********/
func (sla *CameraModelAI) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla CameraModelAI) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type CameraModelAIList []CameraModelAI

func (sla *CameraModelAIList) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla CameraModelAIList) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

/********* Camera model ai ***********/

type Characteristic struct {
	CharacteristicType  string `json:"characteristicType" gorm:"column:characteristic_type"`
	CharacteristicName  string `json:"characteristicName" gorm:"column:characteristic_name"`
	CharacteristicValue string `json:"characteristicValue" gorm:"column:characteristic_value"`
}

func (sla *Characteristic) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla Characteristic) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type CharacteristicList []Characteristic

func (sla *CharacteristicList) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}

func (sla CharacteristicList) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type VideoConfigInfo struct {
	CameraName     string `json:"cameraName,omitempty"`
	CameraChannel  string `json:"dynVideoInputChannelID,omitempty"`
	DataStreamType string `json:"dataStreamType,omitempty"`
	Resolution     string `json:"resolution,omitempty"`
	BitrateType    string `json:"bitrateType,omitempty"`
	VideoQuality   string `json:"videoQuality,omitempty"`
	FrameRate      string `json:"frameRate,omitempty"`
	MaxBitrate     string `json:"maxBitrate,omitempty"`
	VideoEncoding  string `json:"videoEncoding,omitempty"`
	H265           string `json:"h265,omitempty"`
}

func (sla *VideoConfigInfo) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, &sla)
	case string:
		return json.Unmarshal([]byte(src), &sla)
	default:
		return errors.New("unsupported type")
	}
}

func (sla VideoConfigInfo) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type VideoConfigInfoArr []VideoConfigInfo

func (s *VideoConfigInfoArr) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, &s)
	case string:
		return json.Unmarshal([]byte(src), &s)
	default:
		return errors.New("unsupported type")
	}
}

func (s VideoConfigInfoArr) Value() (driver.Value, error) {
	val, err := json.Marshal(s)
	return string(val), err
}

/********* Camera AI Event Coordiante ***********/
// type EventZoneCoordinate struct {
// 	XCord  int `json:"xcord"`
// 	YCord  int `json:"ycord"
// }

/*
// JSONB type
type JSONBMap map[string]interface{}

// Value return json value, implement driver.Valuer interface
func (m JSONBMap) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *JSONBMap) Scan(val interface{}) error {
	if val == nil {
		*m = make(JSONBMap)
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
	}
	t := map[string]interface{}{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m JSONBMap) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	t := (map[string]interface{})(m)
	return json.Marshal(t)
}

// UnmarshalJSON to deserialize []byte
func (m *JSONBMap) UnmarshalJSON(b []byte) error {
	t := map[string]interface{}{}
	err := json.Unmarshal(b, &t)
	*m = JSONBMap(t)
	return err
}

// GormDataType gorm common data type
func (m JSONBMap) GormDataType() string {
	return "jsonmap"
}

// GormDBDataType gorm db data type
func (JSONBMap) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

func (jm JSONBMap) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}

type JSONBMapArray []map[string]interface{}

// Value return json value, implement driver.Valuer interface
func (m JSONBMapArray) Value() (driver.Value, error) {
	if m == nil {
		return nil, nil
	}
	ba, err := m.MarshalJSON()
	return string(ba), err
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (m *JSONBMapArray) Scan(val interface{}) error {
	if val == nil {
		*m = make(JSONBMapArray, 0)
		return nil
	}
	var ba []byte
	switch v := val.(type) {
	case []byte:
		ba = v
	case string:
		ba = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", val))
	}
	t := []map[string]interface{}{}
	err := json.Unmarshal(ba, &t)
	*m = t
	return err
}

// MarshalJSON to output non base64 encoded []byte
func (m JSONBMapArray) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	t := ([]map[string]interface{})(m)
	return json.Marshal(t)
}

// UnmarshalJSON to deserialize []byte
func (m *JSONBMapArray) UnmarshalJSON(b []byte) error {
	t := []map[string]interface{}{}
	err := json.Unmarshal(b, &t)
	*m = JSONBMapArray(t)
	return err
}

// GormDataType gorm common data type
func (m JSONBMapArray) GormDataType() string {
	return "jsonmap"
}

// GormDBDataType gorm db data type
func (JSONBMapArray) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	case "sqlserver":
		return "NVARCHAR(MAX)"
	}
	return ""
}

func (jm JSONBMapArray) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	data, _ := jm.MarshalJSON()
	switch db.Dialector.Name() {
	case "mysql":
		if v, ok := db.Dialector.(*mysql.Dialector); ok && !strings.Contains(v.ServerVersion, "MariaDB") {
			return gorm.Expr("CAST(? AS JSON)", string(data))
		}
	}
	return gorm.Expr("?", string(data))
}
*/
