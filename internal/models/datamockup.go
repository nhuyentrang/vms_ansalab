package models

import (
	"encoding/xml"
	"time"
)

// #region TCP/IP
type NetworkInterface struct {
	Token        string `json:"token"`
	Link         Link   `json:"Link"`
	DHCP         bool   `json:"DHCP"`
	IPv4Address  string `json:"IPv4Address"`
	IPv4Branch   string `json:"IPv4Branch"`
	IPv4Default  string `json:"IPv4Default"`
	IPv6Address  string `json:"IPv6Address"`
	Length       string `json:"Leng"`
	IPv6Default  string `json:"IPv6Default"`
	MAC          string `json:"MAC"`
	MUT          int    `json:"MUT"`
	DNS          bool   `json:"DNS"`
	PreferredDNS string `json:"PreferredDNS"`
	AlternateDNS string `json:"AlternateDNS"`
}

type Link struct {
	AutoNegotiation bool   `json:"autoNegotiation"`
	Speed           int    `json:"speed"`
	Duplex          string `json:"duplex"`
}

//#endregion

// #region DDNS
type DDNS struct {
	Token         string `json:"token"`
	Enabled       bool   `json:"Enabled"`
	ServerAddress string `json:"ServerAddress"`
	Domain        string `json:"Domain"`
	UserName      string `json:"UserName"`
	Password      int    `json:"Password"`
	Confirm       string `json:"Confirm"`
}

//#endregion

// #region Port
type AdminAccessProtocolList struct {
	Token     string                `json:"token"`
	Protocols []AdminAccessProtocol `json:"AdminAccessProtocol"`
}

type AdminAccessProtocol struct {
	Id       int    `json:"id"`
	Enabled  bool   `json:"enabled,omitempty"`
	Protocol string `json:"protocol"`
	PortNo   int    `json:"portNo"`
}

//#endregion

// #region NTP
type NTP struct {
	Token         string `json:"token"`
	NTP           bool   `json:"ntp,omitempty"`
	DateTimeType  string `json:"DateTimeType,omitempty"`
	ServerAddress string `json:"ServerAddress,omitempty"`
	NTPPort       string `json:"NTPPort,omitempty"`
	Period        string `json:"Period,omitempty"`
}

type SystemDateAndTime struct {
	DateTimeType  string    `json:"DateTimeType,omitempty"`
	TimeZone      string    `json:"TimeZone,omitempty"`
	UTCDateTime   time.Time `json:"UTCDateTime,omitempty"`
	LocalDateTime time.Time `json:"LocalDateTime,omitempty"`
}

type NTPList struct {
	NTP               NTP               `json:"ntp"`
	SystemDateAndTime SystemDateAndTime `json:"SystemDateAndTime"`
}

//#endregion

// #region Image
type ImageChannel struct {
	Id                       int    `json:"id"`
	InputPort                int    `json:"inputPort"`
	WDRMode                  string `json:"WDRMode"`
	WDRLevel                 int    `json:"WDRLevel"`
	BLCEanbled               bool   `json:"BLCEanbled"`
	ColorBrightnessLevel     int    `json:"colorBrightnessLevel"`
	ColorContrastLevel       int    `json:"colorContrastLevel"`
	ColorSaturationLevel     int    `json:"colorSaturationLevel"`
	HLCEanbled               bool   `json:"HLCEanbled"`
	HLCLevel                 int    `json:"HLCLevel"`
	EnableImageLossDetection bool   `json:"enableImageLossDetection"`
}

//#endregion

// #region OSD
type NormalizedScreenSize struct {
	NormalizedScreenWidth  int `json:"normalizedScreenWidth"`
	NormalizedScreenHeight int `json:"normalizedScreenHeight"`
}

type Attribute struct {
	Transparent bool `json:"transparent"`
	Flashing    bool `json:"flashing"`
}

type TextOverlays struct {
	ID          int    `json:"id"`
	Enabled     bool   `json:"enabled"`
	PositionX   int    `json:"positionX"`
	PositionY   int    `json:"positionY"`
	DisplayText string `json:"displayText"`
}

type DateTimeOverlay struct {
	Enabled     bool   `json:"enabled"`
	PositionX   int    `json:"positionX"`
	PositionY   int    `json:"positionY"`
	DateStyle   string `json:"dateStyle"`
	TimeStyle   string `json:"timeStyle"`
	DisplayWeek bool   `json:"displayWeek"`
}

type ChannelNameOverlay struct {
	Enabled   bool `json:"enabled"`
	PositionX int  `json:"positionX"`
	PositionY int  `json:"positionY"`
}

type TextOverlayList struct {
	TextOverlays []TextOverlays `json:"TextOverlay"`
}

type VideoOverlay struct {
	NormalizedScreenSize []NormalizedScreenSize `json:"normalizedScreenSize"`
	Attribute            []Attribute            `json:"attribute"`
	TextOverlayList      []TextOverlayList      `json:"textOverlayList"`
	DateTimeOverlay      []DateTimeOverlay      `json:"dateTimeOverlay"`
	ChannelNameOverlay   []ChannelNameOverlay   `json:"channelNameOverlay"`
}

type VideoOverlays struct {
	Version              string               `xml:"version,attr"`
	XMLName              xml.Name             `xml:"VideoOverlay,omitempty"`
	XMLNamespace         string               `xml:"xmlns,attr"`
	NormalizedScreenSize NormalizedScreenSize `xml:"normalizedScreenSize"`
	Attribute            Attribute            `xml:"attribute"`
	TextOverlayList      TextOverlayList      `xml:"TextOverlayList"`
	DateTimeOverlay      DateTimeOverlay      `xml:"DateTimeOverlay"`
	ChannelNameOverlay   ChannelNameOverlay   `xml:"channelNameOverlay"`
	FontSize             string               `xml:"fontSize,attr,omitempty"`
	FrontColorMode       string               `xml:"frontColorMode,attr,omitempty"`
	FrontColor           string               `xml:"frontColor,omitempty"`
	Alignment            string               `xml:"alignment,omitempty"`
	BatteryPowerOverlay  string               `xml:"BatteryPowerOverlay,omitempty"`
}

type NormalizedScreenSizes struct {
	NormalizedScreenWidth  int `xml:"normalizedScreenWidth"`
	NormalizedScreenHeight int `xml:"normalizedScreenHeight"`
}

type Attributes struct {
	Transparent bool `xml:"transparent"`
	Flashing    bool `xml:"flashing"`
}

type TextOverlayLists struct {
	Size int `xml:"size,attr"`
}

type DateTimeOverlays struct {
	Enabled     bool   `xml:"enabled"`
	PositionX   int    `xml:"positionX"`
	PositionY   int    `xml:"positionY"`
	DateStyle   string `xml:"dateStyle"`
	TimeStyle   string `xml:"timeStyle"`
	DisplayWeek bool   `xml:"displayWeek"`
}

type ChannelNameOverlays struct {
	Enabled   bool `xml:"enabled"`
	PositionX int  `xml:"positionX"`
	PositionY int  `xml:"positionY"`
}

//#endregion

//#region DDNS
//#endregion

//#region DDNS
//#endregion

//#region DDNS
//#endregion
