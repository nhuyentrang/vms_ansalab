package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type ImageChannel struct {
	Version            string             `xml:"version,attr"`
	XMLName            xml.Name           `xml:"ImageChannel,omitempty"`
	XMLNamespace       string             `xml:"xmlns,attr"`
	ID                 int                `xml:"id"`
	Enabled            bool               `xml:"enabled"`
	VideoInputID       int                `xml:"videoInputID"`
	ImageFlip          ImageFlip          `xml:"ImageFlip"`
	WDR                WDR                `xml:"WDR"`
	BLC                BLC                `xml:"BLC"`
	IrcutFilter        IrcutFilter        `xml:"IrcutFilter"`
	WhiteBalance       WhiteBalance       `xml:"WhiteBalance"`
	Exposure           Exposure           `xml:"Exposure"`
	Sharpness          Sharpness          `xml:"Sharpness"`
	Shutter            Shutter            `xml:"Shutter"`
	PowerLineFrequency PowerLineFrequency `xml:"powerLineFrequency"`
	Color              Color              `xml:"Color"`
	NoiseReduce        NoiseReduce        `xml:"NoiseReduce"`
	HLC                HLC                `xml:"HLC"`
	SupplementLight    SupplementLight    `xml:"SupplementLight"`
}

type ISPMode struct {
	Version      string          `xml:"version,attr"`
	XMLName      xml.Name        `xml:"ISPMode,omitempty"`
	XMLNamespace string          `xml:"xmlns,attr"`
	Mode         string          `xml:"mode"`
	Schedule     ScheduleISPMode `xml:"Schedule"`
}

type ImageFlip struct {
	Enabled        bool   `xml:"enabled"`
	ImageFlipStyle string `xml:"ImageFlipStyle"`
}

type WDR struct {
	Mode     string `xml:"mode"`
	WDRLevel int    `xml:"WDRLevel"`
}

type BLC struct {
	Enabled bool   `xml:"enabled"`
	BLCMode string `xml:"BLCMode"`
}

type IrcutFilter struct {
	IrcutFilterType       string   `xml:"IrcutFilterType"`
	NightToDayFilterLevel int      `xml:"nightToDayFilterLevel"`
	NightToDayFilterTime  int      `xml:"nightToDayFilterTime"`
	Schedule              Schedule `xml:"Schedule"`
}

type Schedule struct {
	ScheduleType string `xml:"scheduleType"`
	BeginTime    string `xml:"TimeRange>beginTime"`
	EndTime      string `xml:"TimeRange>endTime"`
}

type WhiteBalance struct {
	WhiteBalanceStyle string `xml:"WhiteBalanceStyle"`
	WhiteBalanceRed   int    `xml:"WhiteBalanceRed"`
	WhiteBalanceBlue  int    `xml:"WhiteBalanceBlue"`
}

type Exposure struct {
	ExposureType       string             `xml:"ExposureType"`
	OverexposeSuppress OverexposeSuppress `xml:"OverexposeSuppress"`
}

type OverexposeSuppress struct {
	Enabled bool `xml:"enabled"`
}

type Sharpness struct {
	SharpnessLevel int `xml:"SharpnessLevel"`
}

type Shutter struct {
	ShutterLevel string `xml:"ShutterLevel"`
}

type PowerLineFrequency struct {
	PowerLineFrequencyMode string `xml:"powerLineFrequencyMode"`
}

type Color struct {
	BrightnessLevel int `xml:"brightnessLevel"`
	ContrastLevel   int `xml:"contrastLevel"`
	SaturationLevel int `xml:"saturationLevel"`
}

type NoiseReduce struct {
	Mode         string       `xml:"mode"`
	GeneralMode  GeneralMode  `xml:"GeneralMode"`
	AdvancedMode AdvancedMode `xml:"AdvancedMode"`
}

type GeneralMode struct {
	GeneralLevel int `xml:"generalLevel"`
}

type AdvancedMode struct {
	FrameNoiseReduceLevel      int `xml:"FrameNoiseReduceLevel"`
	InterFrameNoiseReduceLevel int `xml:"InterFrameNoiseReduceLevel"`
}

type HLC struct {
	Enabled  bool `xml:"enabled"`
	HLCLevel int  `xml:"HLCLevel"`
}

type SupplementLight struct {
	SupplementLightMode             string `xml:"supplementLightMode"`
	MixedLightBrightnessRegulatMode string `xml:"mixedLightBrightnessRegulatMode"`
	IrLightBrightness               int    `xml:"irLightBrightness"`
	IsAutoModeBrightnessCfg         bool   `xml:"isAutoModeBrightnessCfg"`
}

type ScheduleISPMode struct {
	ScheduleType string    `xml:"scheduleType"`
	TimeRange    TimeRange `xml:"TimeRange"`
}

type TimeRange struct {
	BeginTime string `xml:"beginTime"`
	EndTime   string `xml:"endTime"`
}

func (c *Client) GetImageConfig(channel string) (resp *ImageChannel, err error) {
	path := "/ISAPI/Image/channels/1/capabilities"
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

func (c *Client) SetImageConfig(data *ImageChannel) (resp *ResponseStatus, err error) {
	path := "/ISAPI/Image/channels/1/capabilities"
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

func (c *Client) GetIsPmodel(channel string) (resp *ISPMode, err error) {
	path := "/ISAPI/Image/channels/" + channel + "/ISPMode"
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

func (c *Client) SetIsPmodel(data *ISPMode) (resp *ResponseStatus, err error) {
	path := "/ISAPI/Image/channels/1/ISPMode"
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
