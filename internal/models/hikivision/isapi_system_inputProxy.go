package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type SourceInputPortDescriptor struct {
	AdminProtocol        string `xml:"adminProtocol" json:"adminProtocol"`                         // required
	AddressingFormatType string `xml:"addressingFormatType" json:"addressingFormatType"`           // required
	HostName             string `xml:"hostName,omitempty" json:"hostName,omitempty"`               // optional
	IPAddress            string `xml:"ipAddress,omitempty" json:"ipAddress,omitempty"`             // optional
	IPv6Address          string `xml:"ipv6Address,omitempty" json:"ipv6Address,omitempty"`         // optional
	ManagePortNo         int    `xml:"managePortNo" json:"managePortNo"`                           // required
	SrcInputPort         string `xml:"srcInputPort" json:"srcInputPort"`                           // required
	UserName             string `xml:"userName" json:"userName"`                                   // required
	Password             string `xml:"password" json:"password"`                                   // required
	StreamType           string `xml:"streamType,omitempty" json:"streamType,omitempty"`           // optional
	DeviceID             string `xml:"deviceID,omitempty" json:"deviceID,omitempty"`               // optional
	DeviceTypeName       string `xml:"deviceTypeName,omitempty" json:"deviceTypeName,omitempty"`   // optional & read-only
	SerialNumber         string `xml:"serialNumber,omitempty" json:"serialNumber,omitempty"`       // optional & read-only
	FirmwareVersion      string `xml:"firmwareVersion,omitempty" json:"firmwareVersion,omitempty"` // optional & read-only
	FirmwareCode         string `xml:"firmwareCode,omitempty" json:"firmwareCode,omitempty"`       // optional & read-only
}

type NVRInfo struct {
	IPAddressNVR string `xml:"ipAddressNVR,omitempty" json:"ipAddressNVR,omitempty"` // optional
	PortNVR      int    `xml:"portNVR,omitempty" json:"portNVR,omitempty"`           // optional
	IPCChannelNo int    `xml:"ipcChannelNo,omitempty" json:"ipcChannelNo,omitempty"` // optional
}

type InputProxyChannel struct {
	XMLName         xml.Name                  `xml:"InputProxyChannel,omitempty" json:"-"` // Skip JSON for XMLName
	XMLVersion      string                    `xml:"version,attr" json:"version"`          // version attribute
	XMLNamespace    string                    `xml:"xmlns,attr" json:"xmlns"`              // xmlns attribute
	ID              string                    `xml:"id" json:"id"`                         // required
	Name            string                    `xml:"name,omitempty" json:"name,omitempty"` // optional
	SourceInputPort SourceInputPortDescriptor `xml:"sourceInputPortDescriptor,omitempty" json:"sourceInputPortDescriptor,omitempty"`
	EnableAnr       *bool                     `xml:"enableAnr,omitempty" json:"enableAnr,omitempty"` // optional
	NVRInfo         NVRInfo                   `xml:"NVRInfo,omitempty" json:"NVRInfo,omitempty"`
}

func (c *Client) GetInputProxyChannel() (resp *InputProxyChannel, err error) {
	path := "/ISAPI/ContentMgmt/InputProxy/channels"
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

func (c *Client) PostInputProxyChannel(data *InputProxyChannel) (resp *ResponseStatus, err error) {
	path := "/ISAPI/ContentMgmt/InputProxy/channels"
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}
	fmt.Println("post XML ", c)
	body, err := c.PostXML(u, data)
	if err != nil {
		fmt.Println("======> POST XML Unsuccessful: ", err)
		return nil, err
	}
	err = xml.Unmarshal(body, &resp)
	if err != nil {
		fmt.Println("======> XML Unmarshal Unsuccessful: ", err)
		return nil, err
	}
	fmt.Println("=======> POST data success: ", resp)
	return resp, nil
}

func (c *Client) DeleteInputProxyChannel(indexCameraToNVR string) (resp *ResponseStatus, err error) {
	path := "/ISAPI/ContentMgmt/InputProxy/channels/" + indexCameraToNVR
	u, err := url.Parse(c.BaseURL + path)
	if err != nil {
		return nil, err
	}
	body, err := c.DeleteXML(u, nil)
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
