package models

import (
	"encoding/xml"
	"net/http"
	"time"
	"vms/internal/models/hikivision"

	uuid "github.com/google/uuid"
)

type DeviceCommand struct {
	CommandID     string          `json:"id,omitempty"`
	CameraID      string          `json:"cameraID,omitempty"`
	NVRID         string          `json:"nvrID,omitempty"`
	RequestUUID   uuid.UUID       `json:"clientID,omitempty"`
	Cmd           string          `json:"cmd,omitempty"`
	Status        string          `json:"status,omitempty"`
	EventTime     string          `json:"eventTime,omitempty"`
	EventType     string          `json:"eventType,omitempty"`
	IPAddress     string          `json:"ipAddress,omitempty"`
	StartIPAdress string          `json:"startIpAddress,omitempty"`
	EndIPAdress   string          `json:"endIpAddress,omitempty"`
	StartPort     string          `json:"startPort,omitempty"`
	EndPort       string          `json:"endPort,omitempty"`
	HttpPort      string          `json:"httpPort,omitempty"`
	OnvifPort     string          `json:"onvifPort,omitempty"`
	UserName      string          `json:"userName,omitempty"`
	Password      string          `json:"password,omitempty"`
	IndexNVR      string          `json:"indexNVR,omitempty"`
	Channel       string          `json:"channel,omitempty"`
	Track         string          `json:"track,omitempty"`
	StartTime     string          `json:"startTime,omitempty"`
	EndTime       string          `json:"endTime,omitempty"`
	ProtocolType  string          `json:"protocolType,omitempty"`
	VideoConfig   DTO_VideoConfig `json:"videoConfig,omitempty"`

	SetImageConfigCamera       VideoOverlays              `json:"setImageConfigCamera,omitempty"`
	StreamingChannelList       StreamingChannelList       `json:"streamingChannelList,omitempty"`
	StreamingChannelListNVR    StreamingChannelListNVR    `json:"streamingChannelListNVR,omitempty"`
	StreamingChannelListCamera StreamingChannelListCamera `json:"streamingChannelListCamera,omitempty"`
	StreamingChannel           StreamingChannel           `json:"streamingChannel,omitempty"`
	NetworkConfig              DTO_NetworkConfig          `json:"networkConfig,omitempty"`
	ChangePassword             DTO_ChangePassword         `json:"changePassword,omitempty"`
	ChangePasswordSeries       []DTO_ChangePassword       `json:"changePasswordSeries,omitempty"`
	ConfigCamera               map[string]ConfigCamera    `json:"configCamera,omitempty"`
	CamerasVideoConfig         []CamerasVideoConfig       `json:"camerasVideoConfig,omitempty"`
	CamerasImageConfig         []CamerasImageConfig       `json:"camerasImageConfig,omitempty"`
	CamerasNetworkConfig       []CamerasNetworkConfig     `json:"camerasNetworkConfig,omitempty"`

	ConfigNVR       map[string]ConfigNVR `json:"configNVR,omitempty"`
	SetInputProxy   interface{}          `json:"inputProxyChannel,omitempty"`  //SetInputProxy        hikivision.InputProxyChannel `json:"inputProxyChannel,omitempty"`
	SetInputProxies interface{}          `json:"inputProxyChannels,omitempty"` //SetInputProxy        hikivision.InputProxyChannel `json:"inputProxyChannel,omitempty"`

	TrackDailyParam hikivision.TrackDailyParam `json:"trackDailyParam,omitempty"`
}

type CamerasVideoConfig struct {
	StreamingChannelListCameras []StreamingChannelListCamera `json:"streamingChannelListCameras,omitempty"`
	IPAddress                   string                       `json:"ipAddress,omitempty"`
	UserName                    string                       `json:"userName,omitempty"`
	Password                    string                       `json:"password,omitempty"`
	HttpPort                    string                       `json:"httpPort,omitempty"`
	OnvifPort                   string                       `json:"onvifPort,omitempty"`
}

type CamerasImageConfig struct {
	VideoOverlayCamera []VideoOverlays `json:"videoOverlayCamera,omitempty"`
	IPAddress          string          `json:"ipAddress,omitempty"`
	UserName           string          `json:"userName,omitempty"`
	Password           string          `json:"password,omitempty"`
	HttpPort           string          `json:"httpPort,omitempty"`
	OnvifPort          string          `json:"onvifPort,omitempty"`
}

type CamerasNetworkConfig struct {
	NetworkConfigCameras DTO_NetworkConfig `json:"networkConfigCameras,omitempty"`
	IPAddress            string            `json:"ipAddress,omitempty"`
	UserName             string            `json:"userName,omitempty"`
	Password             string            `json:"password,omitempty"`
	HttpPort             string            `json:"httpPort,omitempty"`
	OnvifPort            string            `json:"onvifPort,omitempty"`
}

type DeviceCommandFEceRegister struct {
	CommandID   string       `json:"commandID,omitempty"`
	Cmd         string       `json:"cmd,omitempty"`
	Status      string       `json:"status,omitempty"`
	EventTime   string       `json:"eventTime,omitempty"`
	EventType   string       `json:"eventType,omitempty"`
	FaceRegData *FaceRegData `json:"faceRegData,omitempty"`
}

type KafkaJsonVMSMessage struct {
	PayLoad    PayLoad     `json:"payload,omitempty"`
	Topic      string      `json:"topic,omitempty"`
	DeviceInfo *DeviceInfo `json:"deviceInfo,omitempty"`
}

// KafkaJsonAIEventMessage struct representing the JSON data
type KafkaJsonAIEventMessage struct {
	ID              string    `json:"id"`
	MessageID       string    `json:"messageID,omitempty"`
	MsVersion       string    `json:"msVersion,omitempty"`
	SensorID        string    `json:"sensorID,omitempty"`
	Description     string    `json:"description,omitempty"`
	Timestamp       int64     `json:"timestamp,omitempty"`
	TimeStart       int64     `json:"timeStart,omitempty"`
	TimeEnd         int64     `json:"timeEnd,omitempty"`
	Image           string    `json:"image,omitempty"`
	ImageResult     string    `json:"imageResult,omitempty"`
	ImageObject     string    `json:"imageObject,omitempty"`
	Video           string    `json:"video,omitempty"`
	StorageBucket   string    `json:"storageBucket,omitempty"`
	EventType       string    `json:"eventType,omitempty"`
	CamIP           string    `json:"camIP,omitempty"`
	CamName         string    `json:"camName,omitempty"`
	CameraId        string    `json:"cameraId,omitempty"`
	MemberID        string    `json:"memberID,omitempty"`
	EventTypeString string    `json:"eventTypeString,omitempty"`
	Location        string    `json:"location,omitempty"`
	Status          string    `json:"status,omitempty"`
	CabinID         uuid.UUID `json:"cabinID,omitempty" gorm:"column:cabin_id"`
	CabinName       string    `json:"cabinName,omitempty"`
	ConverTimestamp time.Time `json:"converTimestamp,omitempty"`

	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeleteMark bool      `json:"deletedMark" gorm:"column:deleted_mark;default:false"`
	DeletedAt  time.Time `json:"deletedAt" gorm:"column:deleted_at"`
	Result     string    `json:"result"`
}
type PayLoad struct {
	CommandID              string                             `json:"commandId,omitempty"`
	Cmd                    string                             `json:"cmd,omitempty"`
	Status                 string                             `json:"status,omitempty"`
	EventTime              string                             `json:"eventTime,omitempty"`
	DeviceScan             *[]ScanDevice                      `json:"scanDevice,omitempty"`
	RequestUUID            uuid.UUID                          `json:"clientID,omitempty"`
	ListCameraStatus       []DTO_CameraInfo                   `json:"listCameraStatus"`
	ListNVRStatus          []DTO_NVRInfo                      `json:"listNVRStatus"`
	NetworkConfig          *DTO_NetworkConfig                 `json:"networkConfig,omitempty"`
	ChangePasswordSeries   []DTO_ChangePassword               `json:"changePasswordSeries,omitempty"`
	CMSearchResult         *hikivision.CMSearchResult         `json:"cmSearchResult,omitempty"`
	TrackDailyDistribution *hikivision.TrackDailyDistribution `json:"trackDailyDistribution,omitempty"`
	VideoDownLoad          VideoDownLoad                      `json:"videoDownLoad,omitempty"`
	InputProxyChannel      InputProxyChannel                  `json:"InputProxyChannel,omitempty"`
	InputProxyChannels     []InputProxyChannel                `json:"InputProxyChannels,omitempty"`
	StatusCode             int                                `xml:"statusCode,omitempty" json:"statusCode,omitempty"`

	GetOSD                  *hikivision.VideoOverlay  `json:"getOSD,omitempty"`
	ImageConfig             *DTO_ImageConfig          `json:"imageConfig,omitempty"`
	ImageConfigs            *[]DTO_ImageConfig        `json:"imageConfigs,omitempty"`
	NetworkConfigs          *[]DTO_NetworkConfig      `json:"networkConfigs,omitempty"`
	VideoConfig             *StreamingChannelList     `json:"videoConfig,omitempty"` //TODO : change to DTO_VideoCOnfig
	VideoConfig_DTO         *DTO_VideoConfig          `json:"videoConfig_dto,omitempty"`
	VideoConfigInfo         VideoConfigInfoArr        `json:"videoConfigInfo,omitempty"`
	StorageConfig           hikivision.Track          `json:"track,omitempty"`
	ResponseStatus          hikivision.ResponseStatus `json:"responseStatus,omitempty"`
	Response                http.Response             `json:"response,omitempty"`
	VideoConfigCamera       *DTO_VideoConfig          `json:"videoConfigCamera,omitempty"`
	VideoConfigNVR          *DTO_VideoConfig          `json:"videoConfigNVR,omitempty"`
	CameraStatuses          []VideoConfigEditStatuses `json:"cameraStatuses,omitempty"`
	CameraImageStatuses     []ImageConfigEditStatuses `json:"cameraImageStatuses,omitempty"`
	ChangePasswordsStatuses []ChangePasswordStatuses  `json:"changePasswordsStatuses,omitempty"`
	NetworkConfigStatuses   []ChangeNetworkStatuses   `json:"networkConfigStatuses,omitempty"`

	CameraScreenshot CameraScreenshotResponse `json:"cameraScreenshot,omitempty"`
}

type ChangeNetworkStatuses struct {
	NetworkConfig DTO_NetworkConfig `json:"networkConfig,omitempty"`
	Status        bool              `json:"status,omitempty"`
	Reason        string            `json:"reason,omitempty"`
}
type ChangePasswordStatuses struct {
	CameraPassword DTO_ChangePassword `json:"cameraPassword,omitempty"`
	Status         bool               `json:"status,omitempty"`
	Reason         string             `json:"reason,omitempty"`
}

type ImageConfigEditStatuses struct {
	VideoConfigCamera VideoOverlays `json:"videoConfigCamera,omitempty"`
	Status            bool          `json:"status,omitempty"`
	Reason            string        `json:"reason,omitempty"`
}
type VideoConfigEditStatuses struct {
	VideoConfigCamera StreamingChannelList `json:"videoConfigCamera,omitempty"`
	Status            bool                 `json:"status,omitempty"`
	Reason            string               `json:"reason,omitempty"`
}

type CameraScreenshotResponse struct {
	OnvifSnapshotUri   *OnvifSnapshotUriResponse `json:"onvifSnapshotUri,omitempty"`
	HikvisionImageData []byte                    `json:"hikvisionImageData,omitempty"`
}

type OnvifSnapshotUriResponse struct {
	Uri                 string `json:"uri"`
	InvalidAfterConnect bool   `json:"invalidAfterConnect"`
	InvalidAfterReboot  bool   `json:"invalidAfterReboot"`
	Timeout             string `json:"timeout"`
}
type DeviceInfo struct {
	CreatedAt   int64       `json:"createdAt"`
	UpdatedAt   int64       `json:"updatedAt"`
	ID          string      `json:"id"`
	DeviceType  string      `json:"deviceType"`
	ModelID     string      `json:"modelId"`
	DeviceCode  string      `json:"deviceCode"`
	Status      string      `json:"status"`
	Location    string      `json:"location"`
	IPAddress   string      `json:"ipAddress"`
	MacAddress  string      `json:"macAddress"`
	Metadata    interface{} `json:"metadata"`
	Telemetry   interface{} `json:"telemetry"`
	MqttAccount MqttAccount `json:"mqttAccount"`
}

type MqttAccount struct {
	Password string `json:"password"`
	Port     int    `json:"port"`
	Host     string `json:"host"`
	// SubscribeTopic []string `json:"subscribeTopic"`
	Uri string `json:"uri"`
	// PublishTopic   []string `json:"publishTopic"`
	Username string `json:"username"`
}

// Đổi cấu hình motiondetection
type SetMotionDetection struct {
	Sensitivity int    `json:"sensitivity,omitempty"`
	Object      string `json:"object,omitempty"`
	GridMap     string `json:"gridMap,omitempty"`
}

type DTO_CameraInfo struct {
	ID          string   `json:"id,omitempty"`
	IP          string   `json:"ip,omitempty"`
	Status      bool     `json:"status"`
	StatusValue KeyValue `json:"statusValue"`
}
type DTO_NVRInfo struct {
	ID          string   `json:"id,omitempty"`
	IP          string   `json:"ip,omitempty"`
	Status      bool     `json:"status"`
	StatusValue KeyValue `json:"statusValue"`
}
type ConfigCamera struct {
	NameCamera string                   `json:"nameCamera,omitempty"`
	IP         string                   `json:"ip,omitempty"`
	UserName   string                   `json:"userName,omitempty"`
	PassWord   string                   `json:"passWord,omitempty"`
	HTTPPort   string                   `json:"httpPort,omitempty"`
	RTSPPort   string                   `json:"rtspPort,omitempty"`
	OnvifPort  string                   `json:"onvifPort,omitempty"`
	ChannelNVR string                   `json:"channelNVR,omitempty"`
	Channels   map[string]ChannelCamera `json:"channels"`
	IDNVR      string                   `json:"idNVR"`
}
type ConfigNVR struct {
	NameCamera string        `json:"nameCamera,omitempty"`
	IP         string        `json:"ip,omitempty"`
	UserName   string        `json:"userName,omitempty"`
	PassWord   string        `json:"passWord,omitempty"`
	HTTPPort   string        `json:"httpPort,omitempty"`
	RTSPPort   string        `json:"rtspPort,omitempty"`
	OnvifPort  string        `json:"onvifPort,omitempty"`
	ChannelNVR string        `json:"channelNVR,omitempty"`
	Cameras    KeyValueArray `json:"cameras,omitempty"`
}

type ChannelCamera struct {
	OnDemand bool   `json:"on_demand"`
	Url      string `json:"url"`
	Codec    string `json:"codec"`
	Name     string `json:"name"`
}
type Video struct {
	Enabled                 bool          `xml:"enabled" json:"enabled"`
	VideoInputChannelID     int           `xml:"videoInputChannelID" json:"videoInputChannelID"`
	VideoCodecType          string        `xml:"videoCodecType" json:"videoCodecType"` //H264  // H265
	VideoScanType           string        `xml:"videoScanType" json:"videoScanType"`
	VideoResolutionWidth    int           `xml:"videoResolutionWidth" json:"videoResolutionWidth"`
	VideoResolutionHeight   int           `xml:"videoResolutionHeight" json:"videoResolutionHeight"`
	VideoQualityControlType string        `xml:"videoQualityControlType" json:"videoQualityControlType"`
	ConstantBitRate         int           `xml:"constantBitRate" json:"constantBitRate"`
	FixedQuality            int           `xml:"fixedQuality" json:"fixedQuality"`
	VbrUpperCap             int           `xml:"vbrUpperCap" json:"vbrUpperCap"`
	VbrLowerCap             int           `xml:"vbrLowerCap" json:"vbrLowerCap"`
	MaxFrameRate            int           `xml:"maxFrameRate" json:"maxFrameRate"`
	KeyFrameInterval        int           `xml:"keyFrameInterval" json:"keyFrameInterval"`
	SnapShotImageType       string        `xml:"snapShotImageType" json:"snapShotImageType"`
	GovLength               int           `xml:"GovLength" json:"GovLength"`
	PacketTypes             []string      `xml:"PacketType" json:"PacketType"`
	Smoothing               int           `xml:"smoothing" json:"smoothing"`
	H265Profile             string        `xml:"H265Profile" json:"H265Profile"`
	SmartCodec              SmartCodec    `xml:"SmartCodec" json:"SmartCodec"`
	ChannelCamera           ChannelCamera `xml:"ChannelCamera" json:"ChannelCamera"`
}

type SetVideoConfig struct {
	Enabled                 bool       `xml:"enabled" json:"enabled"`
	VideoInputChannelID     int        `xml:"videoInputChannelID" json:"videoInputChannelID"`
	VideoCodecType          string     `xml:"videoCodecType" json:"videoCodecType"` //H264  // H265
	VideoScanType           string     `xml:"videoScanType" json:"videoScanType"`
	VideoResolutionWidth    int        `xml:"videoResolutionWidth" json:"videoResolutionWidth"`
	VideoResolutionHeight   int        `xml:"videoResolutionHeight" json:"videoResolutionHeight"`
	VideoQualityControlType string     `xml:"videoQualityControlType" json:"videoQualityControlType"`
	ConstantBitRate         int        `xml:"constantBitRate" json:"constantBitRate"`
	FixedQuality            int        `xml:"fixedQuality" json:"fixedQuality"`
	VbrUpperCap             int        `xml:"vbrUpperCap" json:"vbrUpperCap"`
	VbrLowerCap             int        `xml:"vbrLowerCap" json:"vbrLowerCap"`
	MaxFrameRate            int        `xml:"maxFrameRate" json:"maxFrameRate"`
	KeyFrameInterval        int        `xml:"keyFrameInterval" json:"keyFrameInterval"`
	SnapShotImageType       string     `xml:"snapShotImageType" json:"snapShotImageType"`
	GovLength               int        `xml:"GovLength" json:"GovLength"`
	PacketTypes             []string   `xml:"PacketType" json:"PacketType"`
	Smoothing               int        `xml:"smoothing" json:"smoothing"`
	H265Profile             string     `xml:"H265Profile" json:"H265Profile"`
	SmartCodec              SmartCodec `xml:"SmartCodec" json:"SmartCodec"`
}
type SmartCodec struct {
	Enabled bool `xml:"enabled"`
}

type SourceInputPortDescriptor struct {
	AdminProtocol        string   `xml:"adminProtocol" json:"adminProtocol"`                         // required
	AddressingFormatType string   `xml:"addressingFormatType" json:"addressingFormatType"`           // required
	HostName             string   `xml:"hostName,omitempty" json:"hostName,omitempty"`               // optional
	IPAddress            string   `xml:"ipAddress,omitempty" json:"ipAddress,omitempty"`             // optional
	IPv6Address          string   `xml:"ipv6Address,omitempty" json:"ipv6Address,omitempty"`         // optional
	ManagePortNo         int      `xml:"managePortNo" json:"managePortNo"`                           // required
	SrcInputPort         string   `xml:"srcInputPort" json:"srcInputPort"`                           // required
	UserName             string   `xml:"userName" json:"userName"`                                   // required
	Password             string   `xml:"password" json:"password"`                                   // required
	StreamType           string   `xml:"streamType,omitempty" json:"streamType,omitempty"`           // optional
	DeviceID             string   `xml:"deviceID,omitempty" json:"deviceID,omitempty"`               // optional
	DeviceTypeName       string   `xml:"deviceTypeName,omitempty" json:"deviceTypeName,omitempty"`   // optional & read-only
	SerialNumber         string   `xml:"serialNumber,omitempty" json:"serialNumber,omitempty"`       // optional & read-only
	FirmwareVersion      string   `xml:"firmwareVersion,omitempty" json:"firmwareVersion,omitempty"` // optional & read-only
	FirmwareCode         string   `xml:"firmwareCode,omitempty" json:"firmwareCode,omitempty"`       // optional & read-only
	MacAddress           string   `json:"macAddress"`                                                // optional
	NVR                  KeyValue `json:"nvr"`                                                       // optional
	Box                  KeyValue `json:"box"`                                                       // optional
}
type NVRInfo struct {
	IPAddressNVR string `xml:"ipAddressNVR,omitempty"` // optional
	PortNVR      int    `xml:"portNVR,omitempty"`      // optional
	IPCChannelNo int    `xml:"ipcChannelNo,omitempty"` // optional
}

// ResponseStatus represents the XML_ResponseStatus and JSON_ResponseStatus resource.
type ResponseStatus struct {
	XMLName       xml.Name                     `xml:"ResponseStatus,omitempty"`
	XMLVersion    string                       `xml:"version,attr"`
	XMLNamespace  string                       `xml:"xmlns,attr"`
	RequestURL    string                       `xml:"requestURL,omitempty" json:"requestURL,omitempty"`
	StatusCode    int                          `xml:"statusCode,omitempty" json:"statusCode,omitempty"`
	StatusString  string                       `xml:"statusString,omitempty" json:"statusString,omitempty"`
	ID            int                          `xml:"id,omitempty" json:"id,omitempty"`
	SubStatusCode string                       `xml:"subStatusCode,omitempty" json:"subStatusCode,omitempty"`
	ErrorCode     int                          `xml:"errorCode,omitempty" json:"errorCode,omitempty"`
	ErrorMsg      string                       `xml:"errorMsg,omitempty" json:"errorMsg,omitempty"`
	AdditionalErr *ResponseStatusAdditionalErr `xml:"AdditionalErr,omitempty" json:"AdditionalErr,omitempty"`
}
type InputProxyChannel struct {
	XMLName         xml.Name                  `xml:"InputProxyChannel,omitempty"`
	XMLVersion      string                    `xml:"version,attr"`
	XMLNamespace    string                    `xml:"xmlns,attr"`
	ID              string                    `xml:"id"`             // required
	Name            string                    `xml:"name,omitempty"` // optional
	SourceInputPort SourceInputPortDescriptor `json:"sourceInputPortDescriptor,omitempty"`
	EnableAnr       *bool                     `xml:"enableAnr,omitempty"` // optional
	NVRInfo         NVRInfo                   `xml:"NVRInfo,omitempty"`
}

// ResponseStatusAdditionalErr represents the additional error status, which is
// valid when StatusCode is set to 9.
type ResponseStatusAdditionalErr struct {
	StatusList []ResponseStatusAdditionalErrStatus `xml:"StatusList,omitempty" json:"StatusList,omitempty"`
}

// ResponseStatusAdditionalErrStatus represents a single status information.
type ResponseStatusAdditionalErrStatus struct {
	Status string `xml:"Status,omitempty" json:"Status,omitempty"`
}

// ResponseStatusAdditionalErrStatusInfo represents information of status.
type ResponseStatusAdditionalErrStatusInfo struct {
	ID            string `xml:"id,omitempty" json:"id,omitempty"`
	StatusCode    int    `xml:"statusCode,omitempty" json:"statusCode,omitempty"`
	StatusString  string `xml:"statusString,omitempty" json:"statusString,omitempty"`
	SubStatusCode string `xml:"subStatusCode,omitempty" json:"subStatusCode,omitempty"`
	ErrorCode     int    `xml:"errorCode,omitempty" json:"errorCode,omitempty"`
	ErrorMsg      string `xml:"errorMsg,omitempty" json:"errorMsg,omitempty"`
}

type TrackDailyParam struct {
	XMLName     xml.Name `xml:"trackDailyParam"`
	Year        int      `xml:"year"`
	MonthOfYear int      `xml:"monthOfYear"`
}
type StreamingChannelList struct {
	Version          string             `xml:"version,attr"`
	XMLName          xml.Name           `xml:"StreamingChannelList,omitempty"`
	XMLNamespace     string             `xml:"xmlns,attr"`
	StreamingChannel []StreamingChannel `xml:"StreamingChannel"`
}
type StreamingChannelListCamera struct {
	Version          string                         `xml:"version,attr"`
	XMLName          xml.Name                       `xml:"StreamingChannelListCamera,omitempty"`
	XMLNamespace     string                         `xml:"xmlns,attr"`
	StreamingChannel []StreamingChannelUpdateCamera `xml:"StreamingChannel"`
}
type StreamingChannel struct {
	// Version      string    `xml:"version,attr"`
	// XMLName      xml.Name  `xml:"StreamingChannel,omitempty"`
	// XMLNamespace string    `xml:"xmlns,attr"`
	ID          int       `xml:"id"`
	ChannelName string    `xml:"channelName"`
	Enabled     bool      `xml:"enabled"`
	Transport   Transport `xml:"Transport"`
	Video       Video     `xml:"Video"`
	URI         string    `xml:"uri"`
}
type StreamingChannelUpdateCamera struct {
	// Version      string    `xml:"version,attr"`
	// XMLName      xml.Name  `xml:"StreamingChannel,omitempty"`
	// XMLNamespace string    `xml:"xmlns,attr"`
	ID          int            `xml:"id"`
	ChannelName string         `xml:"channelName"`
	Enabled     bool           `xml:"enabled"`
	Transport   Transport      `xml:"transport"`
	Video       SetVideoConfig `xml:"video"`
}

type Transport struct {
	MaxPacketSize       int                 `xml:"maxPacketSize"`
	ControlProtocolList ControlProtocolList `xml:"controlProtocolList"`
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

type StreamingChannelListNVR struct {
	XMLName           xml.Name                    `xml:"StreamingChannelList" json:"-"`
	Version           string                      `xml:"version,attr" json:"version"`
	XMLNamespace      string                      `xml:"xmlns,attr,omitempty" json:"xmlns"`
	StreamingChannels []StreamingChannelNVRUpdate `xml:"StreamingChannel" json:"streamingChannel"`
}

type StreamingChannelNVRUpdate struct {
	Version      string       `xml:"version,attr" json:"version,omitempty"`
	XMLNamespace string       `xml:"xmlns,attr,omitempty" json:"xmlns"`
	ID           int          `xml:"id" json:"id"`
	ChannelName  string       `xml:"channelName" json:"channelName"`
	Enabled      bool         `xml:"enabled" json:"enabled"`
	TransportNVR TransportNVR `xml:"Transport" json:"Transport"`
	Video        Video        `xml:"Video" json:"video"`
}
type TransportNVR struct {
	XMLName             xml.Name               `xml:"Transport" json:"-"`
	RtspPortNo          int                    `xml:"rtspPortNo,omitempty" json:"rtspPortNo,omitempty"`
	MaxPacketSize       int                    `xml:"maxPacketSize,omitempty" json:"maxPacketSize,omitempty"`
	ControlProtocolList ControlProtocolListNVR `xml:"ControlProtocolList,omitempty" json:"ControlProtocolList,omitempty"`
}
type ControlProtocolListNVR struct {
	ControlProtocols []ControlProtocol `xml:"ControlProtocol,omitempty" json:"ControlProtocol,omitempty"`
}

type ControlProtocolNVR struct {
	StreamingTransport string `xml:"streamingTransport" json:"streamingTransport"`
}
