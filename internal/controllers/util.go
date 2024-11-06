package controllers

import (
	"context"
	"fmt"
	"net"
	"time"

	"vms/internal/models"
	"vms/wssignaling"

	"vms/comongo/reposity"

	"github.com/google/uuid"
)

var (
	cmd_ScanDevice                  string = "scan_device"
	cmd_ScanDeviceIP                string = "scan_device_ip"
	cmd_ScanDeviceListIP            string = "scan_device_list_ip"
	cmd_GetDataConfig               string = "get_dataconfig"
	cmd_GetDataStoragePlayback      string = "get_data_storage_playback"
	cmd_AddDataConfig               string = "add_dataconfig"
	cmd_UpdateNetworkConfig         string = "update_networkconfig"
	cmd_GetNetWorkConfig            string = "get_networkconfig"
	cmd_UpdateNetWorkConfigSeries   string = "get_networkconfigs"
	cmd_SetVideoConfigOfMultiCamera string = "set_video_config_of_multi_Camera"
	cmd_GetVideoConfig              string = "get_video_config"
	cmd_SetVideoConfigCamera        string = "set_video_config_of_Camera"
	cmd_UpdateOSDConfig             string = "set_osd_config"
	cmd_UpdateOSDConfigs            string = "set_osd_configs"
	cmd_UpdateOSDConfigNVR          string = "set_osd_config_of_NVR"
	cmd_GetOSDConfig                string = "get_osd_config"
	cmd_GetOSDConfigNVR             string = "get_osd_config_of_NVR"
	cmd_GetImageConfig              string = "get_image_config"
	cmd_UpdateOTA                   string = "update_OTA"
	cmd_AddCameratoNVR              string = "add_camera"
	cmd_ChangePassWord              string = "change_password"
	cmd_ChangePassWordSeries        string = "change_password_series"
	cmd_PingCamera                  string = "PingCamera"
	cmd_UpdateIPandPortHTTP         string = "update_ip_and_porthttp"
	cmd_GetCalendelPlayback         string = "get_data_calender_playback"
	cmd_DeleteFileStream            string = "delete_file_stream"
	cmd_RemoveCameraFromNVR         string = "remove_camera"
	cmd_RemoveCamerasFromNVR        string = "remove_cameras"
	cmd_GetCameraSnapshot           string = "get_camera_snapshot"
	cmd_DownLoadClip                string = "download_clip"
	cmd_ExtractClip                 string = "extract_clip"
	cmd_GetAttachedCameras          string = "get_attached_cameras"
	cmd_GetOnSiteVideoConfig        string = "get_onsite_video_config"
	cmd_SetVideoConfigOfNVR         string = "set_video_config_of_NVR"
	cmd_GetVideoConfigNVR           string = "get_video_config_of_NVR"
)

func ConvertAIEventString(eventAI string) string {
	switch eventAI {
	case "AI_EVENT_SABOTAGE_DETECTION":
		return "Phát hiện camera bị phá hoại"
	case "AI_EVENT_DANGEROUS_OBJECT_DETECTION":
		return "Phát hiện đối tượng nguy hiểm"
	case "AI_EVENT_ABNORMAL_BEHAVIOR_DETECTION":
		return "Phát hiện hành vi bất thường"
	case "AI_EVENT_BLACKLIST_FACE_RECOGNITION":
		return "Phát hiện đối tượng danh sách đen"
	case "AI_EVENT_MEMBER_RECOGNITION":
		return "Nhận dạng thành viên"
	case "AI_EVENT_UNKNOWN_FACE_DETECTION":
		return "Phát hiện người lạ"
	case "LICENSE_PLATE_RECOGNITION":
		return "Nhận dạng biển số xe"
	default:
		return "Sự kiện AI chưa phân loại"
	}
}

func Get10MinuteIntervalsInRangeStamp(startTimestamp, endTimestamp int64) []int64 {
	start := time.Unix(0, startTimestamp*int64(time.Millisecond))
	end := time.Unix(0, endTimestamp*int64(time.Millisecond))

	var timestamps []int64

	for t := end; t.After(start); t = t.Add(-10 * time.Minute) {
		timestamps = append(timestamps, t.UnixNano()/int64(time.Millisecond))
	}

	return timestamps
}

func Get4HoursIntervalsInRangeStamp(startTimestamp, endTimestamp int64) []int64 {
	start := time.Unix(0, startTimestamp*int64(time.Millisecond))
	end := time.Unix(0, endTimestamp*int64(time.Millisecond))

	var timestamps []int64

	for t := end; t.After(start); t = t.Add(-4 * time.Hour) {
		timestamps = append(timestamps, t.UnixNano()/int64(time.Millisecond))
	}

	return timestamps
}

func ConvertConnectToAlertTypeString(IS string) string {
	switch IS {
	case "SensorConnected":
		return "Đã kết nối"
	case "SensorDisconnected":
		return "Mất kết nối"
	case "not connected":
		return "Không thể kết nối"
	case "not availble":
		return "Không thể truyền tải hình ảnh"
	default:
		return "Cảnh báo hệ thống"
	}
}
func SaveNewAlertReport(cabinID uuid.UUID, sensorID uuid.UUID, alertType string, cabinName string, deviceName string, location string, isWss bool) error {
	if isWss {
		data := models.DTO_Report{
			ID:         uuid.New(),
			CabinID:    cabinID,
			SensorID:   sensorID,
			AlertType:  alertType,
			CabinName:  cabinName,
			Location:   location,
			DeviceName: deviceName,
			Deleted:    false,
			Status:     "NEW",
		}
		if deviceName != "AlarmSensor" {
			_, err := reposity.CreateItemFromDTO[models.DTO_Report, models.Report](data)
			return err
		}
	}
	return nil
}

func SaveNewSystemWaring(cabinID uuid.UUID, sensorID uuid.UUID, alertType string, cabinName string, deviceName string, location string, isWss bool) error {
	if isWss {
		data := models.DTO_SystemWaring{
			ID:         uuid.New(),
			CabinID:    cabinID,
			SensorID:   sensorID,
			AlertType:  alertType,
			CabinName:  cabinName,
			Location:   location,
			DeviceName: deviceName,
			Deleted:    false,
			Status:     "NEW",
		}
		_, err := reposity.CreateItemFromDTO[models.DTO_SystemWaring, models.SystemWaring](data)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveDevicelog(deviceCode, deviceType, action, detail, protocol string) error {
	log := models.DTO_DeviceLog{
		ID:         uuid.New(),
		Atts:       time.Now().UnixMilli(),
		Command:    "",
		DeviceCode: deviceCode,
		DeviceType: deviceType,
		Action:     action,
		Detail:     detail,
		Protocol:   protocol,
	}
	_, err := reposity.CreateItemFromDTO[models.DTO_DeviceLog, models.DeviceLog](log)

	// Also send device logs to websocket
	wssignaling.SendNotifyMessage("logs", "devicelog", log)
	return err
}

func GetWeeksInRange(startTimestamp, endTimestamp int64) []time.Time {
	start := time.Unix(0, startTimestamp*int64(time.Millisecond))
	end := time.Unix(0, endTimestamp*int64(time.Millisecond))

	var weeks []time.Time

	for t := start; t.Before(end) || t.Equal(end) || t.After(end); {
		year, week := t.ISOWeek()
		firstDayOfISOWeek := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, (week-1)*7)
		weeks = append(weeks, firstDayOfISOWeek)

		t = t.AddDate(0, 0, 7)
	}

	fmt.Println("======> weeks: ", weeks, start, end)
	return weeks
}

func GetHoursInRange(startTimestamp, endTimestamp int64) []time.Time {
	start := time.Unix(0, startTimestamp*int64(time.Millisecond))
	end := time.Unix(0, endTimestamp*int64(time.Millisecond))

	var hours []time.Time

	// Lặp qua từng giờ trong phạm vi thời gian và thêm vào danh sách
	for t := start; t.Before(end); t = t.Add(time.Hour) {
		hours = append(hours, t)
	}

	return hours
}

func GetDaysInRange(startTimestamp int64, endTimestamp int64) []time.Time {
	start := time.Unix(0, startTimestamp*int64(time.Millisecond))
	end := time.Unix(0, endTimestamp*int64(time.Millisecond))

	var days []time.Time

	// loop through days between start and end dates
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		days = append(days, d)
	}

	return days
}

func ConverStructUpdateNetworkConfig(reqDTO models.DTO_NetworkConfig) models.SetTCPIP {
	// AutoNegotiation, Speed, Duplex := ConvertNetWorkconfigNICType(reqDTO.NicType)
	dtoSetTCPIP := models.SetTCPIP{
		SetNetworkInterfaces: models.SetNetworkInterfaces{
			Token: reqDTO.TokenSetNetworkInterfaces,
			Link: models.Links{
				AutoNegotiation: true,
				Speed:           0,
				Duplex:          "", //Full, Half
			},
			MTU: reqDTO.MTU,
			IPv4: models.IPv4NetworkInterfaceSetConfiguration{
				DHCP:         reqDTO.DHCP,
				Address:      reqDTO.IPv4DefaultGateway,
				PrefixLength: 24,
			},
			// IPv6: models.IPv6NetworkInterfaceSetConfiguration{
			// 	Enabled:            false,
			// 	DHCP:               "",
			// 	AcceptRouterAdvert: true,
			// 	Address:            reqDTO.IPv6DefaultGateway,
			// 	PrefixLength:       24,
			// },
		},

		SetNetworkDefaultGateways: models.SetNetworkDefaultGateways{
			IPv4Address: reqDTO.IPv4DefaultGateway,
			// IPv6Address: reqDTO.IPv6DefaultGateway,
		},
		SetDNSServer: models.SetDNSServer{
			FromDHCP: reqDTO.AutoDNS,
			// SearchDomain: "test",
			DNSManual: []models.IPAddress{
				{
					Type:        "IPv4", // IPv4,IPv6
					IPv4Address: reqDTO.PrefDNS,
					IPv6Address: "",
				},
				{
					Type:        "IPv4",
					IPv4Address: reqDTO.AlterDNS,
					IPv6Address: "",
				},
			},
		},
	}
	return dtoSetTCPIP
}

func ConverStructUpdatePort(reqDTO models.DTO_NetworkConfig) models.SetNetworkProtocols {
	dtoSetTCPIP := models.SetNetworkProtocols{
		NetworkProtocols: []models.NetworkProtocol{
			{
				Name:    "HTTP",
				Enabled: true,
				Port:    reqDTO.HTTP,
			},
			// {
			// 	Name:    "HTTPS",
			// 	Enabled: true,
			// 	Port:    reqDTO.HTTPS,
			// }, //ServiceNotSupported
			{
				Name:    "RTSP",
				Enabled: true,
				Port:    reqDTO.RTSP,
			},
		},
	}
	return dtoSetTCPIP
}

func ConverStructUpdateDDNSHik(reqDTO models.DTO_NetworkConfig) models.DDNSServer {
	ddns := models.DDNSServer{
		Version:      "1.0",
		XMLNamespace: "http://www.hikvision.com/ver20/XMLSchema",
		ID:           1,
		Enabled:      true,
		Provider:     "",
		ServerAddress: struct {
			AddressingFormatType string `xml:"addressingFormatType"`
			HostName             string `xml:"hostName"`
		}{
			AddressingFormatType: "dynupdate.no-ip.com",
			HostName:             reqDTO.Port,
		},
		PortNo:           0,
		DeviceDomainName: reqDTO.Domain,
		UserName:         reqDTO.UserName,
		// CountryID:        123,
		Status: "connServerfail",
	}
	return ddns
}

func ConverStructUpdateTime(reqDTO models.DTO_NetworkConfig) models.SetTime {
	var timeType = ""
	if reqDTO.ManualTime {
		timeType = "Manual"
	} else {
		timeType = "NTP"
	}
	data := models.SetTime{
		DateTimeType: timeType,
		DNSname:      reqDTO.ServerAddressNTP,
		Time:         reqDTO.SetTime,
	}
	return data
}

func ConverCodeRegisterBlackList(code int) string {
	switch code {
	case 200:
		return "NO ERR"
	case 400:
		return "ERR IMAGE INFO NOT FOUND"
	case 500:
		return "ERR INTERNAL SERVER"
	case 0:
		return "DELETE OK"
	case 801:
		return "ERR FACE DEVIATTON TOO MUCH"
	case 802:
		return "ERR SIMILAR MEMBER FACE EXISTED"
	case 803:
		return "ERR BLACKLIST EXISTED"
	case 804:
		return "ERR NO FACE FEATURE FOUND"
	case 805:
		return "ERR MEMBER NOT FOUND"
	case 806:
		return "ERR NUMBER FEATURE NOT ENOUGHT"
	case 807:
		return "ERR NO FEATURE FOUND"
	case 808:
		return "ERR NOT ENOUGH REQUIRED FACES"
	case 809:
		return "ERR NO FACE WERE FOUND"
	case 810:
		return "ERR IMAGE QUALITY NG"
	case 811:
		return "ERR FACE WEARING MASK"
	case 812:
		return "ERR FAIL TO DOWLOAD IMAGE"
	case 813:
		return "ERR FAIL TO DECOD B64 IMAGE"
	case 814:
		return "ERR FAIL TO REMOVE MEMBER"
	case 815:
		return "ERR FAIL TO INSERT MEMBER"
	default:
		return "UNKNOWN"
	}
}

func ConvertNetWorkconfigNICType(code string) (AutoNegotiation bool, Speed int, Duplex string) {
	switch code {
	case "Auto":
		return true, 0, "Full"
	case "10M_Half_Dup":
		return false, 10, "Half"
	case "10M_Full_Dup":
		return false, 10, "Full"
	case "100M_Half_Dup":
		return false, 10, "Half"
	case "100M_Full_Dup":
		return false, 100, "Full"
	case "1000M_Full_Dup":
		return false, 1000, "Full"
	case "1000M_Half_Dup":
		return false, 1000, "Half"
	default:
		return true, 100, "Full"
	}
}

func CreatedNetWorkConfig(IDConfigNVR string) error {
	var dtoNetworkConfig models.DTO_NetworkConfig
	// Create NVR
	dtoNetworkConfig, err := reposity.CreateItemFromDTO[models.DTO_NetworkConfig, models.NetworkConfig](dtoNetworkConfig)
	if err != nil {
		fmt.Println("======> Add network configuration failed !!!")
		return err
	}

	// Update entity from DTO
	err = reposity.UpdateSingleColumn[models.NVRConfig](IDConfigNVR, "network_config_id", dtoNetworkConfig.ID)
	if err != nil {
		fmt.Println("======> Updated NVR configuration failed !!!")
		return err
	}
	return nil
}

func kiemTraDanhSachIPPortTrongDai(listIPPort []net.Addr, ipA, ipB string, portA, portB int) []net.Addr {
	var ketQua []net.Addr

	ipAInt := ipToInt(ipA)
	ipBInt := ipToInt(ipB)

	for _, addr := range listIPPort {
		ip, port, err := parseAddr(addr.String())
		if err != nil {
			continue
		}

		ipInt := ipToInt(ip)
		if ipAInt <= ipInt && ipInt <= ipBInt && portA <= port && port <= portB {
			// Thêm addr vào slice nếu thỏa mãn điều kiện
			ketQua = append(ketQua, addr)
		}
	}

	return ketQua
}

func ipToInt(ip string) int {
	ipAddr := net.ParseIP(ip)
	ipInt := 0
	if ipAddr != nil {
		ipInt = int(ipAddr[12])<<24 + int(ipAddr[13])<<16 + int(ipAddr[14])<<8 + int(ipAddr[15])
	}
	return ipInt
}

func parseAddr(addr string) (string, int, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, err
	}

	portInt := 0
	if port != "" {
		// Thêm context.Background() vào hàm LookupPort
		portInt, err = net.DefaultResolver.LookupPort(context.Background(), "tcp", port)
		if err != nil {
			return "", 0, err
		}
	}

	return host, portInt, nil
}

func mapVideoQuality(quality int) string {
	switch {
	case quality <= 1:
		return "Lowest"
	case quality <= 20:
		return "Lower"
	case quality <= 40:
		return "Low"
	case quality <= 60:
		return "Medium"
	case quality <= 80:
		return "High"
	default:
		return "Highest"
	}
}
