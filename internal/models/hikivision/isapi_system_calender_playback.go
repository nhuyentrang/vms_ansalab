package hikivision

import (
	"encoding/xml"
)

// TrackDailyDistribution đại diện cho phần tử gốc
type TrackDailyDistribution struct {
	Version      string   `xml:"version,attr"`
	XMLName      xml.Name `xml:"trackDailyDistribution,omitempty"`
	XMLNamespace string   `xml:"xmlns,attr"`
	DayList      DayList  `xml:"dayList"`
}

// DayList đại diện cho danh sách các ngày
type DayList struct {
	Days []Day `xml:"day"`
}

// Day đại diện cho mỗi ngày trong danh sách
type Day struct {
	ID         int    `xml:"id"`
	DayOfMonth int    `xml:"dayOfMonth"`
	Record     bool   `xml:"record"`
	RecordType string `xml:"recordType,omitempty"`
}

type TrackDailyParam struct {
	Year        int `xml:"year,omitempty"`
	MonthOfYear int `xml:"monthOfYear,omitempty"`
}
