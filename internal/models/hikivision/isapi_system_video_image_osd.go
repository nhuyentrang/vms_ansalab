package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type VideoOverlay struct {
	Version              string               `xml:"version,attr"`
	XMLName              xml.Name             `xml:"VideoOverlay,omitempty"`
	XMLNamespace         string               `xml:"xmlns,attr"`
	NormalizedScreenSize NormalizedScreenSize `xml:"normalizedScreenSize"`
	Attribute            Attribute            `xml:"attribute"`
	TextOverlayList      TextOverlayList      `xml:"TextOverlayList"`
	DateTimeOverlay      DateTimeOverlay      `xml:"DateTimeOverlay"`
	ChannelNameOverlay   ChannelNameOverlay   `xml:"channelNameOverlay"`
}

type NormalizedScreenSize struct {
	NormalizedScreenWidth  int `xml:"normalizedScreenWidth"`
	NormalizedScreenHeight int `xml:"normalizedScreenHeight"`
}

type Attribute struct {
	Transparent bool `xml:"transparent"`
	Flashing    bool `xml:"flashing"`
}

type TextOverlayList struct {
	Size int `xml:"size,attr"`
}

type DateTimeOverlay struct {
	Enabled     bool   `xml:"enabled"`
	PositionX   int    `xml:"positionX"`
	PositionY   int    `xml:"positionY"`
	DateStyle   string `xml:"dateStyle"`
	TimeStyle   string `xml:"timeStyle"`
	DisplayWeek bool   `xml:"displayWeek"`
}

type ChannelNameOverlay struct {
	Enabled   bool `xml:"enabled"`
	PositionX int  `xml:"positionX"`
	PositionY int  `xml:"positionY"`
}

func (c *Client) GetVideoOverlay(channel string) (resp *VideoOverlay, err error) {
	path := "/ISAPI/ContentMgmt/InputProxy/channels/" + channel + "/video/overlays"
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

func (c *Client) PutVideoOverlay(data *VideoOverlay) (resp *ResponseStatus, err error) {
	path := "/ISAPI/ContentMgmt/InputProxy/channels/1/video/overlays"
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
