package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type NetworkInterfaceList struct {
	Version           string            `xml:"version,attr"`
	XMLName           xml.Name          `xml:"NetworkInterfaceList,omitempty"`
	XMLNamespace      string            `xml:"xmlns,attr"`
	NetworkInterfaces NetworkInterfaces `xml:"NetworkInterface"`
}

type NetworkInterfaces struct {
	ID        int       `xml:"id"`
	IPAddress IPAddress `xml:"IPAddress"`
	Discovery Discovery `xml:"Discovery"`
	Link      Link      `xml:"Link"`
}

type IPAddress struct {
	Version        string `xml:"version,attr"`
	XMLNamespace   string `xml:"xmlns,attr"`
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
	DNSEnable string `xml:"DNSEnable"`
}

type Discovery struct {
	Version      string `xml:"version,attr"`
	XMLNamespace string `xml:"xmlns,attr"`
	UPnP         struct {
		Enabled string `xml:"enabled"`
	} `xml:"UPnP"`
	Zeroconf struct {
		Enabled string `xml:"enabled"`
	} `xml:"Zeroconf"`
}

type Link struct {
	Version         string `xml:"version,attr"`
	XMLNamespace    string `xml:"xmlns,attr"`
	MACAddress      string `xml:"MACAddress"`
	AutoNegotiation string `xml:"autoNegotiation"`
	Speed           int    `xml:"speed"`
	Duplex          string `xml:"duplex"`
	MTU             int    `xml:"MTU"`
}

func (c *Client) GetNetworkInterfaceList() (resp *NetworkInterfaceList, err error) {
	path := "/ISAPI/System/Network/interfaces"
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
	fmt.Println("======> ", resp)
	return resp, nil
}
