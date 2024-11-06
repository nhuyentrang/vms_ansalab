package hikivision

import (
	"encoding/xml"
	"fmt"
	"net/url"
)

type NTPServer struct {
	XMLName              xml.Name `xml:"NTPServer,omitempty"`
	XMLVersion           string   `xml:"version,attr"`
	XMLNamespace         string   `xml:"xmlns,attr"`
	ID                   int      `xml:"id"`
	AddressingFormatType string   `xml:"addressingFormatType"`
	HostName             string   `xml:"hostName"`
	PortNo               int      `xml:"portNo"`
	SynchronizeInterval  int      `xml:"synchronizeInterval"`
}

func (c *Client) GetNTP() (resp *NTPServer, err error) {
	path := "/ISAPI/System/time/ntpServers/1"
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

func (c *Client) PutNTP(data *NTPServer) (resp *ResponseStatus, err error) {
	path := "/ISAPI/System/time/ntpServers/1"
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
