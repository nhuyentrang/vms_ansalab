package models

import (
	"time"

	uuid "github.com/google/uuid"
)

// TableName overrides the table name
func (AIEngine) TableName() string {
	return "ivis_vms.ai_engine"
}

// AIEngine model
type AIEngine struct {
	ID        uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	MachineID string    `json:"machineID" gorm:"column:machine_id"`

	// AIEngine agent info
	AgentVersion  string    `json:"agentVersion" gorm:"column:agent_version"`
	Online        bool      `json:"online" gorm:"column:online;default:false"`
	ConnectedAt   time.Time `json:"connectedAt" gorm:"column:connected_at"`
	LastMessageAt time.Time `json:"lastMessageAt" gorm:"column:last_message_at"`

	Hostname   string `json:"hostname" gorm:"column:hostname"`
	IPAddress  string `json:"ipAddress" gorm:"column:ip_address"`
	MACAddress string `json:"macAddress" gorm:"column:mac_address"`
	Uptime     int64  `json:"uptime" gorm:"column:uptime"`

	OSPlatform string        `json:"osPlatform" gorm:"column:os_platform"`
	OSExtInfo  KeyValueArray `json:"osExtInfo" gorm:"column:os_ext_info;embedded;type:jsonb"`
	/*
		// Following fields are in the OSExtInfo keyvalue array
		OSArch               string    `json:"osArch" gorm:"column:os_arch"`
		OSPlatform           string    `json:"osPlatform" gorm:"column:os_platform"`
		OSFamily             string    `json:"osFamily" gorm:"column:os_family"`
		OSVersion            string    `json:"osVersion" gorm:"column:os_version"`
	*/

	CPUUtilization float64       `json:"cpu" gorm:"column:cpu_utilization"`
	CPUExtInfo     KeyValueArray `json:"cpuExtInfo" gorm:"column:cpu_ext_info;embedded;type:jsonb"`
	/*
		// Following fields are in the CPUExtInfo keyvalue array
		CPUCore      int     `json:"cpuCore" gorm:"column:cpu_core"`
		CPUFrequency float64 `json:"cpuFrequency" gorm:"column:cpu_frequency"`
		CPUModel     string  `json:"cpuModel" gorm:"column:cpu_model"`
	*/

	MEMUtilization float64       `json:"mem" gorm:"column:mem_utilization"`
	MEMExtInfo     KeyValueArray `json:"memExtInfo" gorm:"column:mem_ext_info;embedded;type:jsonb"`
	/*
		// Following fields are in the MEMExtInfo keyvalue array
		MemTotal float64 `json:"memTotal" gorm:"column:mem_total"`
		MemUnit  string  `json:"memUnit" gorm:"column:mem_unit"`
	*/

	StorageUtilization float64       `json:"storage" gorm:"column:storage_utilization"`
	StorageExtInfo     KeyValueArray `json:"storageExtInfo" gorm:"column:storage_ext_info;embedded;type:jsonb"`
	/*
		// Following fields are in the StorageExtInfo keyvalue array
		NumberOfDisks  int     `json:"numberOfDisks"`
		Disk0Model     string  `json:"disk0Model"`
		Disk0Size      float64 `json:"disk0Size"`
		Disk0Percent   float64 `json:"disk0Percent"`
		Disk0Unit      string  `json:"disk0Unit"`
		Disk0MountPath string  `json:"disk0MountPath"`
	*/

	GPUUtilization float64       `json:"gpu" gorm:"column:gpu_utilization"`
	GPUExtInfo     KeyValueArray `json:"gpuExtInfo" gorm:"column:gpu_ext_info;embedded;type:jsonb"`
	/*
		// Following fields are in the GPUExtInfo keyvalue array
		NumberOfGPUs      int     `json:"numberOfGPUs" gorm:"column:number_of_gpus"`
		GPUDriverVersion  string  `json:"gpu0DriverVersion" gorm:"column:gpu_driver_version"`
		GPU0Name          string  `json:"gpu0Name" gorm:"column:gpu_name"`
		GPU0MemoryPercent float64 `json:"gpu0MemoryPercent" gorm:"column:gpu_memory_percent"`
		GPU0FanSpeed      float64 `json:"gpu0FanSpeed" gorm:"column:gpu_fan_speed"`
		GPU0Temperature   float64 `json:"gpu0Temperature" gorm:"column:gpu_temperature"`
	*/

	// Store info about available AI models (model name, model type, model characteristic, model capability)
	AIModels KeyValueArray `json:"aiModels" gorm:"column:ai_models;embedded;type:jsonb"`

	// Camera back-referencer table "cam_ai_engine_planning"
	Cameras []*Camera `gorm:"many2many:cam_ai_engine_planning;"`

	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime:true"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;autoUpdateTime:true"`
	DeletedAt time.Time `json:"deletedAt" gorm:"column:deleted_at"`
}

// WS message between AIEngine and Management Service
type WSMsgAIEngine struct {
	ID        uuid.UUID `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	SenderID  string    `json:"senderID"`
	//Type      string    `json:"type"`
	Action string `json:"action"`
	Data   string `json:"data"`
}

const (
	// Message type (direction)
	//WSMsgTypeRequestFromAIEngine = "RequestFromAIEngine"
	//WSMsgTypeResponseFromVMS     = "ResponseFromVMS"

	// Command type
	AIEngineActionRegisterAgent  = "RegisterAgent"
	AIEngineActionSyncAIConfig   = "SyncAIConfig"
	AIEngineActionUpdateAIConfig = "UpdateAIConfig"
	AIEngineActionUpdateAIModel  = "UpdateAIModel"
	AIEngineActionReportUsage    = "ReportUsage"
)

// Draft...
// Command for config AI: SyncAIConfig, UpdateAIConfig, ReadAIConfig, DeleteAIConfig, ListAIConfig
// Command for control: StartAI, StopAI, RestartAI, CheckAI, ListAI
// Command for update: UpdateAIEngine, UpdateAIEngineConfig
// Command for model: GetAIModel, ListAIModel
// Command for camera: GetAICamera, ListAICamera
// Command for error: GetAIError, ListAIError

// CamAIProperty model
type WSAICamProperty struct {
	CameraName       string `json:"cameraName"`
	CameraMacAddress string `json:"cameraMacAddress"`
	Description      string `json:"description"`

	// CamAIProperty contain model name, region config...
	AICamProperty DTO_Camera_AI_Property_BasicInfo `json:"aiCamProperty"`

	// First stream is main stream, second stream is sub stream, and custom streams are after that
	// When retrive stream, should check if channel contain "main" or "Main", if not take the first stream
	VideoStream VideoStream `json:"videoStream"`

	// Camera location
	Location   string `json:"location"`
	Coordinate string `json:"coordinate"`
}

// AIEngineSyncAIConfig model
// This data is for sync all, any camera or model is not in this message will be remove
type WSAICamProperties struct {
	WSAICamProperties []WSAICamProperty `json:"wsAICamProperties"`
}

// AIEngineUsageReport model
type WSAIEngineUsageReport struct {
	HostName string `json:"hostName"`

	// Common usage info
	CPUUsage     float64 `json:"cpu"`
	MEMUsage     float64 `json:"mem"`
	StorageUsage float64 `json:"storage"`
	GPUUsage     float64 `json:"gpu"`

	// Extended cpu usage info
	CPUExtInfo []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"cpuExtInfo"`

	// Extended memory usage info
	MEMExtInfo []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"memExtInfo"`

	// Extended storage usage info
	StorageExtInfo []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"storageExtInfo"`

	// Extended gpu usage info
	GPUExtInfo []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"gpuExtInfo"`

	AIEProcState []struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"aieProcState"`
}
