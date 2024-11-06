package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type DDNS struct {
	Version       string   `xml:"version,attr"`
	XMLName       xml.Name `xml:"DDNS,omitempty"`
	XMLNamespace  string   `xml:"xmlns,attr"`
	ID            int      `xml:"id"`
	Enabled       bool     `xml:"enabled"`
	Provider      string   `xml:"provider"`
	ServerAddress struct {
		AddressingFormatType string `xml:"addressingFormatType"`
		HostName             string `xml:"hostName"`
	} `xml:"serverAddress"`
	PortNo           int    `xml:"portNo"`
	DeviceDomainName string `xml:"deviceDomainName"`
	UserName         string `xml:"userName"`
	CountryID        int    `xml:"countryID"`
	Status           string `xml:"status"`
}

func (c *Client) GetDDNS() (resp *DDNS, err error) {
	path := "/ISAPI/System/Network/DDNS/1"
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

func (c *Client) PutDDNS(data *DDNS) (resp *ResponseStatus, err error) {
	path := "/ISAPI/System/Network/DDNS/1"
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
