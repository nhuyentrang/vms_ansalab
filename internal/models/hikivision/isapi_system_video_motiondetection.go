package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type MotionDetection struct {
	Version          string   `xml:"version,attr"`
	XMLName          xml.Name `xml:"MotionDetection,omitempty"`
	XMLNamespace     string   `xml:"xmlns,attr"`
	Enabled          bool     `xml:"enabled"`
	EnableHighlight  bool     `xml:"enableHighlight"`
	SamplingInterval int      `xml:"samplingInterval"`
	StartTriggerTime int      `xml:"startTriggerTime"`
	EndTriggerTime   int      `xml:"endTriggerTime"`
	RegionType       string   `xml:"regionType"`
	Grid             struct {
		RowGranularity    int `xml:"rowGranularity"`
		ColumnGranularity int `xml:"columnGranularity"`
	} `xml:"Grid"`
	MotionDetectionLayout struct {
		XMLName          xml.Name `xml:"MotionDetectionLayout"`
		SensitivityLevel int      `xml:"sensitivityLevel"`
		Layout           struct {
			GridMap string `xml:"gridMap"`
		} `xml:"layout"`
	} `xml:"MotionDetectionLayout"`
}

// Đổi cấu hình motiondetection
type SetMotionDetection struct {
	Sensitivity int    `json:"sensitivity,omitempty"`
	Object      string `json:"object,omitempty"`
	GridMap     string `json:"gridMap,omitempty"`
}

func (c *Client) GetMotionDetection(indexNVR string) (resp *MotionDetection, err error) {
	path := fmt.Sprintf("/ISAPI/System/Video/inputs/channels/%s/motionDetection", indexNVR)
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		fmt.Println("======> 1 Put XML Unsuccessful: ", err)
		return nil, err
	}
	fmt.Println("Data:", u)
	body, err := c.Get(u)
	if err != nil {
		fmt.Println("======> 2 Put XML Unsuccessful: ", err)
		return nil, err
	}
	err = xml.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("======> 3 Put XML Unsuccessful: ", err)
		return nil, err
	}
	return resp, nil
}

func (c *Client) PutMotionDetection(indexNVR string, data *MotionDetection) (resp *ResponseStatus, err error) {
	path := fmt.Sprintf("/ISAPI/System/Video/inputs/channels/%s/motionDetection", indexNVR)
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
