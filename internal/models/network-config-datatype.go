package models

import (
	"database/sql/driver"
	"encoding/json"
)

/****************************************************** Basic Network Setting ********************************************************/
// BasicTcpIp
type BasicTCPIP struct {
	// Common
	NICType            KeyValue `default:"auto" json:"nicType,omitempty"`
	DHCP               KeyValue `json:"dhcp,omitempty"`
	IPVersion          string   `default:"dual" json:"ipVersion,omitempty"`
	PreferredDNSServer string   `json:"preferredDNSServer,omitempty"`
	AlternateDNSServer string   `json:"alternateDNSServer,omitempty"`
	MACAddress         string   `json:"macAddress,omitempty"`
	// IPv4
	IPv4AddressingType string `default:"dynamic" json:"ipv4AddressingType,omitempty"`
	IPv4Address        string `json:"ipv4Address,omitempty"`
	IPv4DefaultGateway string `json:"ipv4DefaultGateway,omitempty"`
	IPv4SubnetMask     string `json:"ipv4SubnetMask,omitempty"`
	// IPv6
	IPv6AddressingType string `json:"ipv6AddressingType,omitempty"`
	IPv6Address        string `json:"ipv6Address,omitempty"`
	IPv6DefaultGateway string `json:"ipv6DefaultGateway,omitempty"`
	IPv6SubnetMask     string `json:"ipv6SubnetMask,omitempty"`
	IPv6BitMask        string `json:"ipv6BitMask,omitempty"`
	// Discovery
	UPnP     KeyValue `json:"upnp,omitempty"`
	Zeroconf KeyValue `json:"zeroconf,omitempty"`
	// Link
	AutoNegotiation string `json:"autoNegotiation,omitempty"`
	Speed           string `json:"speed,omitempty"`
	Duplex          string `json:"duplex,omitempty"`
	MTU             string `json:"mtu,omitempty"`
	// Wireless
	WirelessEnabled         KeyValue `json:"wirelessEnabled,omitempty"`
	WirelessNetworkMode     KeyValue `json:"wirelessNetworkMode,omitempty"`
	WirelessChannel         KeyValue `json:"wirelessChannel,omitempty"`
	SSID                    string   `json:"ssid,omitempty"`
	WmmEnabled              KeyValue `json:"wmmEnabled,omitempty"`
	WirelessSecurityMode    KeyValue `json:"wirelessSecurityMode,omitempty"`
	WirelessAlgorithmType   KeyValue `json:"wirelessAlgorithmType,omitempty"`
	WirelessSharedKey       string   `json:"wirelessSharedKey,omitempty"`
	WpaKeyLength            int64    `json:"wpaKeyLength,omitempty"`
	WirelessSupport64bitKey KeyValue `json:"wirelessSupport64bitKey,omitempty"`
	// Multicast
	MulticastAddress   string   `json:"multicastAddress,omitempty"`
	MulticastDiscovery KeyValue `json:"multicastDiscovery,omitempty"`
}

func (sla *BasicTCPIP) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), &sla)
}
func (sla BasicTCPIP) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

// BasicDDNS
type BasicDDNS struct {
	Enable               KeyValue `json:"enable,omitempty"`
	Provider             string   `json:"provider,omitempty"`
	AddressingFormatType KeyValue `json:"addressingFormatType,omitempty"`
	HostName             string   `json:"hostName,omitempty"`
	PortNo               int64    `json:"portNo,omitempty"`
	DeviceDomainName     string   `json:"deviceDomainName,omitempty"`
	UserName             string   `json:"userName,omitempty"`
	CountryID            int64    `json:"countryID,omitempty"`
	Status               KeyValue `json:"status,omitempty"`
}

// BasicPort
type BasicPort struct {
	HTTPPort   int64 `json:"httpPort,omitempty"`
	RTSPPort   int64 `json:"rtspPort,omitempty"`
	HTTPSPort  int64 `json:"httpsPort,omitempty"`
	ServerPort int64 `json:"serverPort,omitempty"`
}

// BasicNTP
type BasicNTP struct {
	TimeSync      KeyValue `json:"timeSync,omitempty"`
	Interval      int64    `json:"interval,omitempty"`
	ServerAddress string   `json:"serverAddress,omitempty"`
	Port          int64    `json:"port,omitempty"`
}
