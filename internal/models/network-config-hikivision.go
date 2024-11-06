package models

import (
	"encoding/xml"
	"time"

	uuid "github.com/google/uuid"
)

type NetworkInterfaces struct {
	ID          uuid.UUID                     `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	NVRConfigID uuid.UUID                     `json:"nvrConfigID" gorm:"type:uuid;default:uuid_generate_v4()"`
	Network     NetworkServer                 `json:"network" gorm:"column:network;embedded;type:jsonb"`
	NTP         NTPServer                     `json:"ntp" gorm:"column:ntp;embedded;type:jsonb"`
	DDNS        DDNSServer                    `json:"ddns" gorm:"column:ddns;embedded;type:jsonb"`
	Protocol    AdminAccessProtocolServerList `json:"protocol" gorm:"column:protocol;embedded;type:jsonb"`
	Time        Time                          `json:"time" gorm:"column:time;embedded;type:jsonb"`
	Image       VideoOverlayServer            `json:"image" gorm:"column:image;embedded;type:jsonb"`
	Video       StreamingChannel              `json:"video" gorm:"column:video;embedded;type:jsonb"`
	// Timestamp
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

type DTO_NetworkInterfaces struct {
	ID          uuid.UUID                     `json:"id,omitempty"`
	NVRConfigID uuid.UUID                     `json:"nvrConfigID,omitempty"`
	Network     NetworkServer                 `json:"network,omitempty"`
	NTP         NTPServer                     `json:"ntp,omitempty"`
	DDNS        DDNSServer                    `json:"ddns,omitempty"`
	Protocol    AdminAccessProtocolServerList `json:"protocol,omitempty"`
	Time        Time                          `json:"time,omitempty" gorm:"column:time;embedded;type:jsonb"`
	Image       VideoOverlayServer            `json:"image,omitempty" gorm:"column:image;embedded;type:jsonb"`
	Video       StreamingChannel              `json:"video,omitempty" gorm:"column:video;embedded;type:jsonb"`
}

type DTO_Update_NetworkInterfaces struct {
	ID          uuid.UUID                     `json:"id,omitempty"`
	NVRConfigID uuid.UUID                     `json:"nvrConfigID,omitempty"`
	Network     NetworkServer                 `json:"network,omitempty"`
	NTP         NTPServer                     `json:"ntp,omitempty"`
	DDNS        DDNSServer                    `json:"ddns,omitempty"`
	Protocol    AdminAccessProtocolServerList `json:"protocol,omitempty"`
	Time        Time                          `json:"time,omitempty" gorm:"column:time;embedded;type:jsonb"`
	Image       VideoOverlayServer            `json:"image,omitempty" gorm:"column:image;embedded;type:jsonb"`
	Video       StreamingChannel              `json:"video,omitempty" gorm:"column:video;embedded;type:jsonb"`
}

type NetworkServer struct {
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

type DDNSServer struct {
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

type AdminAccessProtocolServerList struct {
	Version              string                `xml:"version,attr"`
	XMLName              xml.Name              `xml:"AdminAccessProtocolList,omitempty"`
	XMLNamespace         string                `xml:"xmlns,attr"`
	AdminAccessProtocols []AdminAccessProtocol `xml:"AdminAccessProtocol"`
}

type AdminAccessProtocolServer struct {
	Version  string `xml:"version,attr"`
	ID       int    `xml:"id"`
	Enabled  bool   `xml:"enabled"`
	Protocol string `xml:"protocol"`
	PortNo   int    `xml:"portNo"`
}

type VideoOverlayServer struct {
	Version              string               `xml:"version,attr"`
	XMLName              xml.Name             `xml:"VideoOverlay,omitempty"`
	XMLNamespace         string               `xml:"xmlns,attr"`
	NormalizedScreenSize NormalizedScreenSize `xml:"normalizedScreenSize"`
	Attribute            Attribute            `xml:"attribute"`
	TextOverlayList      TextOverlayList      `xml:"TextOverlayList"`
	DateTimeOverlay      DateTimeOverlay      `xml:"DateTimeOverlay"`
	ChannelNameOverlay   ChannelNameOverlay   `xml:"channelNameOverlay"`
}

type NormalizedScreenSizeServer struct {
	NormalizedScreenWidth  int `xml:"normalizedScreenWidth"`
	NormalizedScreenHeight int `xml:"normalizedScreenHeight"`
}

type AttributeServer struct {
	Transparent bool `xml:"transparent"`
	Flashing    bool `xml:"flashing"`
}

type TextOverlayListServer struct {
	Size int `xml:"size,attr"`
}

type DateTimeOverlayServer struct {
	Enabled     bool   `xml:"enabled"`
	PositionX   int    `xml:"positionX"`
	PositionY   int    `xml:"positionY"`
	DateStyle   string `xml:"dateStyle"`
	TimeStyle   string `xml:"timeStyle"`
	DisplayWeek bool   `xml:"displayWeek"`
}

type ChannelNameOverlayServer struct {
	Enabled   bool `xml:"enabled"`
	PositionX int  `xml:"positionX"`
	PositionY int  `xml:"positionY"`
}

type Time struct {
	XMLName           xml.Name `xml:"Time,omitempty"`
	XMLVersion        string   `xml:"version,attr"`
	XMLNamespace      string   `xml:"xmlns,attr"`
	TimeMode          string   `xml:"timeMode,omitempty" json:"timeMode,omitempty"`
	LocalTime         string   `xml:"localTime,omitempty" json:"localTime,omitempty"`
	TimeZone          string   `xml:"timeZone,omitempty" json:"timeZone,omitempty"`
	SatelliteInterval string   `xml:"satelliteInterval,omitempty" json:"satelliteInterval,omitempty"`
}

type StreamingChannelHikvision struct {
	Version      string   `xml:"version,attr"`
	XMLName      xml.Name `xml:"StreamingChannel,omitempty"`
	XMLNamespace string   `xml:"xmlns,attr"`
	ID           string   `xml:"id"`
	ChannelName  string   `xml:"channelName"`
	Enabled      bool     `xml:"enabled"`
	Transport    struct {
		ControlProtocolList struct {
			ControlProtocol struct {
				StreamingTransport string `xml:"streamingTransport"`
			} `xml:"ControlProtocol"`
		} `xml:"ControlProtocolList"`
	} `xml:"Transport"`
	Video struct {
		Enabled                 bool   `xml:"enabled"`
		DynVideoInputChannelID  string `xml:"dynVideoInputChannelID"`
		VideoCodecType          string `xml:"videoCodecType"`
		VideoScanType           string `xml:"videoScanType"`
		VideoResolutionWidth    int    `xml:"videoResolutionWidth"`
		VideoResolutionHeight   int    `xml:"videoResolutionHeight"`
		VideoQualityControlType string `xml:"videoQualityControlType"`
		FixedQuality            int    `xml:"fixedQuality"`
		VbrUpperCap             int    `xml:"vbrUpperCap"`
		VbrLowerCap             int    `xml:"vbrLowerCap"`
		MaxFrameRate            int    `xml:"maxFrameRate"`
		GovLength               int    `xml:"GovLength"`
		SnapShotImageType       string `xml:"snapShotImageType"`
		SmartCodec              struct {
			Enabled bool `xml:"enabled"`
		} `xml:"SmartCodec"`
	} `xml:"Video"`
	Audio struct {
		Enabled              bool   `xml:"enabled"`
		AudioInputChannelID  string `xml:"audioInputChannelID"`
		AudioCompressionType string `xml:"audioCompressionType"`
	} `xml:"Audio"`
}

// Đổi TCP/IP
type SetTCPIP struct {
	SetNetworkInterfaces      SetNetworkInterfaces      `xml:"tds:SetNetworkInterfaces"`
	SetNetworkDefaultGateways SetNetworkDefaultGateways `xml:"tds:SetNetworkDefaultGateways"`
	SetDNSServer              SetDNSServer              `xml:"tds:SetDNSServer"`
}

type SetNetworkInterfaces struct {
	Token string                               `xml:"Token"`
	Link  Links                                `xml:"Link"`
	MTU   int                                  `xml:"MTU"`
	IPv4  IPv4NetworkInterfaceSetConfiguration `xml:"IPv4"`
	IPv6  IPv6NetworkInterfaceSetConfiguration `xml:"IPv6"`
}

type Links struct {
	AutoNegotiation bool   `xml:"AutoNegotiation"`
	Speed           int    `xml:"Speed"`
	Duplex          string `xml:"Duplex"`
}

type IPv4NetworkInterfaceSetConfiguration struct {
	DHCP         bool   `xml:"DHCP"`
	Address      string `xml:"Address"`
	PrefixLength int    `xml:"PrefixLength"`
}

type IPv6NetworkInterfaceSetConfiguration struct {
	Enabled            bool   `xml:"Enabled"`
	DHCP               string `xml:"DHCP"`
	AcceptRouterAdvert bool   `xml:"AcceptRouterAdvert"`
	Address            string `xml:"Address"`
	PrefixLength       int    `xml:"PrefixLength"`
}

type SetNetworkDefaultGateways struct {
	IPv4Address string `xml:"tds:IPv4Address"`
	IPv6Address string `xml:"tds:IPv6Address"`
}

type SetDNSServer struct {
	FromDHCP     bool        `xml:"tds:FromDHCP"`
	SearchDomain string      `xml:"tds:SearchDomain"`
	DNSManual    []IPAddress `xml:"tds:DNSManual"`
}

type IPAddress struct {
	Type        string `xml:"Type"`
	IPv4Address string `xml:"IPv4Address"`
	IPv6Address string `xml:"IPv6Address"`
}

// Đổi cấu hình port
type SetNetworkProtocols struct {
	NetworkProtocols []NetworkProtocol `xml:"tds:NetworkProtocols"`
}

type NetworkProtocol struct {
	Name    string `xml:"Name"`
	Enabled bool   `xml:"Enabled"`
	Port    int    `xml:"Port"`
}

// Đổi cấu hình Time
type SetTime struct {
	DateTimeType string    `xml:"DateTimeType"`
	DNSname      string    `xml:"DNSname"`
	Time         time.Time `xml:"Time"`
}
