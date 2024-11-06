package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type StreamingChannelList struct {
	Version          string                  `xml:"version,attr"`
	XMLName          xml.Name                `xml:"StreamingChannelList,omitempty"`
	XMLNamespace     string                  `xml:"xmlns,attr"`
	StreamingChannel []StreamingChannelVideo `xml:"StreamingChannel"`
}
type DTO_StreamingChannelList struct {
	Version          string                  `xml:"version,attr"`
	XMLName          xml.Name                `xml:"StreamingChannelList,omitempty"`
	XMLNamespace     string                  `xml:"xmlns,attr"`
	StreamingChannel []StreamingChannelVideo `xml:"StreamingChannel"`
}

type StreamingChannelVideo struct {
	// Version      string    `xml:"version,attr"`
	// XMLName      xml.Name  `xml:"StreamingChannel,omitempty"`
	// XMLNamespace string    `xml:"xmlns,attr"`
	ID          int       `xml:"id"`
	ChannelName string    `xml:"channelName"`
	Enabled     bool      `xml:"enabled"`
	Transport   Transport `xml:"Transport"`
	Video       Video     `xml:"Video"`
}

type Transport struct {
	MaxPacketSize       int                 `xml:"maxPacketSize"`
	ControlProtocolList ControlProtocolList `xml:"ControlProtocolList"`
	Unicast             Unicast             `xml:"Unicast"`
	Multicast           Multicast           `xml:"Multicast"`
	Security            Security            `xml:"Security"`
}

type ControlProtocolList struct {
	ControlProtocols []ControlProtocol `xml:"ControlProtocol"`
}

type ControlProtocol struct {
	StreamingTransport string `xml:"streamingTransport"`
}

type Unicast struct {
	Enabled          bool   `xml:"enabled"`
	RTPTransportType string `xml:"rtpTransportType"`
}

type Multicast struct {
	Enabled         bool   `xml:"enabled"`
	DestIPAddress   string `xml:"destIPAddress"`
	VideoDestPortNo int    `xml:"videoDestPortNo"`
	AudioDestPortNo int    `xml:"audioDestPortNo"`
}

type Security struct {
	Enabled           bool              `xml:"enabled"`
	CertificateType   string            `xml:"certificateType"`
	SecurityAlgorithm SecurityAlgorithm `xml:"SecurityAlgorithm"`
}

type SecurityAlgorithm struct {
	AlgorithmType string `xml:"algorithmType"`
}

type Video struct {
	Enabled                 bool       `xml:"enabled"`
	VideoInputChannelID     int        `xml:"videoInputChannelID"`
	VideoCodecType          string     `xml:"videoCodecType"`
	VideoScanType           string     `xml:"videoScanType"`
	VideoResolutionWidth    int        `xml:"videoResolutionWidth"`
	VideoResolutionHeight   int        `xml:"videoResolutionHeight"`
	VideoQualityControlType string     `xml:"videoQualityControlType"`
	ConstantBitRate         int        `xml:"constantBitRate"`
	FixedQuality            int        `xml:"fixedQuality"`
	VbrUpperCap             int        `xml:"vbrUpperCap"`
	VbrLowerCap             int        `xml:"vbrLowerCap"`
	MaxFrameRate            int        `xml:"maxFrameRate"`
	KeyFrameInterval        int        `xml:"keyFrameInterval"`
	SnapShotImageType       string     `xml:"snapShotImageType"`
	GovLength               int        `xml:"GovLength"`
	PacketTypes             []string   `xml:"PacketType"`
	Smoothing               int        `xml:"smoothing"`
	H265Profile             string     `xml:"H265Profile"`
	SmartCodec              SmartCodec `xml:"SmartCodec"`
}
type SmartCodec struct {
	Enabled bool `xml:"enabled"`
}

func (c *Client) GetVideoConfig(channel string) (resp *StreamingChannelList, err error) {
	path := "/ISAPI/Streaming/channels"
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

func (c *Client) PutVideoConfig(data *StreamingChannelList) (resp *ResponseStatus, err error) {
	path := "/ISAPI/Streaming/channels"
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
