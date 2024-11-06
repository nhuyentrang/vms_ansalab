package models

// DTO Model for Sensor Scan
type DTO_WebsocketMessage struct {
	DataType string `json:"dataType" example:"cabin-list;cabin-status;report;devicelog"`
	Data     string `json:"data" example:"string of json object"`
}
