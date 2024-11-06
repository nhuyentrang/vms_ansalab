package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type Track struct {
	Version          string         `xml:"version,attr"`
	XMLName          xml.Name       `xml:"Track,omitempty"`
	XMLNamespace     string         `xml:"xmlns,attr"`
	ID               int            `xml:"id"`
	Channel          int            `xml:"Channel"`
	Enable           Enable         `xml:"Enable"`
	Description      string         `xml:"Description"`
	TrackGUID        string         `xml:"TrackGUID"`
	DefaultRecording DefaultRecMode `xml:"DefaultRecordingMode"`
	LoopEnable       bool           `xml:"LoopEnable"`
	SrcDescriptor    SrcDescriptor  `xml:"SrcDescriptor"`
	TrackSchedule    TrackSchedule  `xml:"TrackSchedule"`
}

type Enable struct {
	Opt  string `xml:"opt,attr"`
	Text string `xml:",chardata"`
}

type DefaultRecMode struct {
	Opt  string `xml:"opt,attr"`
	Text string `xml:",chardata"`
}

type SrcDescriptor struct {
	SrcGUID    string `xml:"SrcGUID"`
	SrcChannel int    `xml:"SrcChannel"`
	StreamHint string `xml:"StreamHint"`
	SrcDriver  string `xml:"SrcDriver"`
	SrcType    string `xml:"SrcType"`
	SrcURL     string `xml:"SrcUrl"`
	SrcTypes   string `xml:"SrcTypes"`
	SrcLogin   string `xml:"SrcLogin"`
}

type TrackSchedule struct {
	ScheduleBlockList ScheduleBlockList `xml:"ScheduleBlockList"`
}

type ScheduleBlockList struct {
	ScheduleBlock ScheduleBlock `xml:"ScheduleBlock"`
}

type ScheduleBlock struct {
	ScheduleBlockGUID string           `xml:"ScheduleBlockGUID"`
	ScheduleBlockType string           `xml:"ScheduleBlockType"`
	ScheduleActions   []ScheduleAction `xml:"ScheduleAction"`
}
type ScheduleAction struct {
	ID                      int             `xml:"id"`
	ScheduleActionStartTime ScheduleStogare `xml:"ScheduleActionStartTime"`
	ScheduleActionEndTime   ScheduleStogare `xml:"ScheduleActionEndTime"`
	ScheduleDSTEnable       bool            `xml:"ScheduleDSTEnable"`
	Description             string          `xml:"Description"`
	Actions                 Actions         `xml:"Actions"`
}

type ScheduleStogare struct {
	DayOfWeek string `xml:"DayOfWeek"`
	TimeOfDay string `xml:"TimeOfDay"`
}

type Actions struct {
	Record              bool   `xml:"Record"`
	ActionRecordingMode string `xml:"ActionRecordingMode"`
}

func (c *Client) GetStogare(track string) (resp *Track, err error) {
	path := "/ISAPI/ContentMgmt/record/tracks/" + track + "/capabilities"
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}
	body, err := c.Get(u)
	if err != nil {
		return nil, err
	}
	err = xml.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) SetStogare(data *Track) (resp *ResponseStatus, err error) {
	path := "/ISAPI/ContentMgmt/record/tracks/103/capabilities"
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}
	body, err := c.PutXML(u, data)
	if err != nil {
		fmt.Println("======> Put XML Unsuccessful: ", err)
		return nil, err
	}
	err = xml.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("======> XML Unmarshal Unsuccessful: ", err)
		return nil, err
	}
	fmt.Println("=======> Put data success: ", resp)
	return resp, nil
}
