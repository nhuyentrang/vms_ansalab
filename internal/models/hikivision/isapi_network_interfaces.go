package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type NetworkInterface struct {
	Version      string   `xml:"version,attr"`
	XMLName      xml.Name `xml:"NetworkInterface,omitempty"`
	XMLNamespace string   `xml:"xmlns,attr"`
	ID           int      `xml:"id"`
	IPAddress    struct {
		IPVersion      string `xml:"ipVersion"`
		AddressingType string `xml:"addressingType"`
		IPAddress      string `xml:"ipAddress"`
		SubnetMask     string `xml:"subnetMask"`
		IPv6Address    string `xml:"ipv6Address"`
		BitMask        string `xml:"bitMask"`
		DefaultGateway struct {
			IPAddress   string `xml:"ipAddress"`
			IPv6Address string `xml:"ipv6Address"`
		} `xml:"DefaultGateway"`
		PrimaryDNS struct {
			IPAddress string `xml:"ipAddress"`
		} `xml:"PrimaryDNS"`
		SecondaryDNS struct {
			IPAddress string `xml:"ipAddress"`
		} `xml:"SecondaryDNS"`
		DNSEnable bool `xml:"DNSEnable"`
	} `xml:"IPAddress"`
	Discovery struct {
		UPnP struct {
			Enabled bool `xml:"enabled"`
		} `xml:"UPnP"`
		Zeroconf struct {
			Enabled bool `xml:"enabled"`
		} `xml:"Zeroconf"`
	} `xml:"Discovery"`
	Link struct {
		MACAddress      string `xml:"MACAddress"`
		AutoNegotiation bool   `xml:"autoNegotiation"`
		Speed           int    `xml:"speed"`
		Duplex          string `xml:"duplex"`
		MTU             struct {
			Min   int `xml:"min,attr"`
			Max   int `xml:"max,attr"`
			Value int `xml:",chardata"`
		} `xml:"MTU"`
	} `xml:"Link"`
}

func (c *Client) GetNetworkInterface() (resp *NetworkInterface, err error) {
	path := "/ISAPI/System/Network/interfaces/1"
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

func (c *Client) PutNetworkInterface(data *NetworkInterface) (resp *ResponseStatus, err error) {
	path := "/ISAPI/System/Network/interfaces/1"
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
