package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type AdminAccessProtocolList struct {
	Version              string                `xml:"version,attr"`
	XMLName              xml.Name              `xml:"AdminAccessProtocolList,omitempty"`
	XMLNamespace         string                `xml:"xmlns,attr"`
	AdminAccessProtocols []AdminAccessProtocol `xml:"AdminAccessProtocol"`
}

type AdminAccessProtocol struct {
	Version  string `xml:"version,attr"`
	ID       int    `xml:"id"`
	Enabled  bool   `xml:"enabled"`
	Protocol string `xml:"protocol"`
	PortNo   int    `xml:"portNo"`
}

func (c *Client) GetAdminAccessProtocolList() (resp *AdminAccessProtocolList, err error) {
	path := "/ISAPI/Security/adminAccesses"
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

func (c *Client) PutAdminAccessProtocolList(data *AdminAccessProtocolList) (resp *ResponseStatus, err error) {
	path := "/ISAPI/Security/adminAccesses"
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
