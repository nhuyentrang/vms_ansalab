package controllers

import (
	"net/http"

	"vms/internal/models"

	"github.com/gin-gonic/gin"
)

// GetCameraConfigOptions		godoc
// @Summary      	Get items of camera protocols for select box
// @Description  	Responds with the list of item for camera protocol.
// @Tags         	cameras
// @Produce      	json
// @Param   		cameraType		query	string	true	"camera type config"	Enums(network_nic_type, network_ddns_type, network_nat_port_mapping_type, video_stream_type, video_type, video_resolution_type, video_bitrate_type, video_quality_type, video_framrate_type, video_maxbirate_type, video_encoding_type, video_encoding_status, video_profle_type, image_seting_type, image_seting_day_night_type, image_seting_day_night_sensitivity,image_wdr_hlc_type, image_blc_type, image_seting_exposure_time, image_seting_exposure_model)
// @Success      	200  		{object}  		models.JsonDTORsp[[]models.KeyValue]
// @Router       	/cameras/options/combox [get]
// @Security		BearerAuth
func GetCameraCombox(c *gin.Context) {

	cameraType := c.Query("cameraType")

	jsonConfigTypeRsp := models.NewJsonDTORsp[[]models.KeyValue]()
	jsonConfigTypeRsp.Data = make([]models.KeyValue, 0)

	if cameraType == "network_nic_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "camera_nic_type_auto",
			Name: "Auto",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "camera_nic_type_10_haff",
			Name: "10M Half-dup",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "camera_nic_type_100_haff",
			Name: "100M Half-dup",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "camera_nic_type_10_full",
			Name: "10M Full-dup",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "camera_nic_type_100_full",
			Name: "100M Full-dup",
		})
	}

	if cameraType == "network_ddns_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "network_ddns_type_DynDNS",
			Name: "DynDNS",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "camera_nic_type_No-IP",
			Name: "NO-IP",
		})
	}

	if cameraType == "network_nat_port_mapping_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "network_nat_port_mapping_type_auto",
			Name: "Auto",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "network_nat_port_mapping_type_manual",
			Name: "Manual",
		})
	}

	if cameraType == "video_stream_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_stream_type_main",
			Name: "Main Stream(Normal)",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_stream_type_sub",
			Name: "Sub Stream",
		})
	}

	if cameraType == "video_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_type_1",
			Name: "Video Stream",
		})
	}

	if cameraType == "video_resolution_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_resolution_type_1920x1080",
			Name: "1920*1080P",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_resolution_type_1280x720",
			Name: "1280*720P",
		})
	}

	if cameraType == "video_bitrate_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_bitrate_type_variable",
			Name: "Variable",
		})

		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_bitrate_type_constant",
			Name: "Constant",
		})
	}

	if cameraType == "video_quality_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_quality_type_lowest",
			Name: "Lowest",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_quality_type_lower",
			Name: "Lower",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_quality_type_low",
			Name: "Low",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_quality_type_medium",
			Name: "Medium",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_quality_type_higher",
			Name: "Higher",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_quality_type_highest",
			Name: "Highest",
		})
	}

	if cameraType == "video_framrate_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_framrate_type1",
			Name: "1",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_framrate_type2",
			Name: "2",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_framrate_type4",
			Name: "4",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_framrate_type6",
			Name: "6",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_framrate_type8",
			Name: "8",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_framrate_type10",
			Name: "10",
		})
	}

	if cameraType == "video_maxbirate_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_maxbirate_type256",
			Name: "256",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_maxbirate_type512",
			Name: "512",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_maxbirate_type1024",
			Name: "1024",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_maxbirate_type2048",
			Name: "2040",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_maxbirate_type3072",
			Name: "3072",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_maxbirate_type4096",
			Name: "4096",
		})
	}

	if cameraType == "video_encoding_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_encoding_type264",
			Name: "H.264",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_encoding_type265",
			Name: "H.265",
		})
	}

	if cameraType == "video_encoding_status" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_encoding_status_on",
			Name: "On",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_encoding_status_off",
			Name: "Off",
		})
	}

	if cameraType == "video_profle_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_profle_type1",
			Name: "Main Profile",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "video_profle_type2",
			Name: "Sub Profile",
		})
	}

	if cameraType == "image_seting_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_type_auto",
			Name: "Auto switch",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_type_schedule",
			Name: "Scheduled switch",
		})
	}

	if cameraType == "image_seting_day_night_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_type_auto",
			Name: "Auto",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_type_day",
			Name: "Day",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_type_night",
			Name: "Night",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_type_schedule",
			Name: "Scheduled switch",
		})
	}

	if cameraType == "image_seting_day_night_sensitivity" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_sensitivity0",
			Name: "0",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_sensitivity1",
			Name: "1",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_sensitivity2",
			Name: "2",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_sensitivity3",
			Name: "3",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_sensitivity4",
			Name: "4",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_sensitivity5",
			Name: "5",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_sensitivity6",
			Name: "6",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_day_night_sensitivity7",
			Name: "7",
		})
	}

	if cameraType == "image_seting_exposure_model" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_exposure_model_Manual",
			Name: "Manual",
		})
	}

	if cameraType == "image_seting_exposure_time" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_exposure_time0",
			Name: "1/3",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_exposure_time1",
			Name: "1/6",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_exposure_time2",
			Name: "1/12",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_exposure_time3",
			Name: "1/25",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_exposure_time4",
			Name: "1/50",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_seting_exposure_time5",
			Name: "1/100",
		})
	}

	if cameraType == "image_blc_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_blc_type_off",
			Name: "OFF",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_blc_type_up",
			Name: "Up",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_blc_type_donw",
			Name: "Down",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_blc_type_left",
			Name: "Left",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_blc_type_right",
			Name: "Right",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_blc_type_center",
			Name: "Center",
		})
	}

	if cameraType == "image_wdr_hlc_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_wdr_hlc_type_on",
			Name: "On",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_wdr_hlc_type_off",
			Name: "Off",
		})
	}

	if cameraType == "image_wbl_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_wbl_type_awb1",
			Name: "AWB1",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_wbl_type_locked_wb",
			Name: "Locked WB",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_wbl_type_incandescent_lamp",
			Name: "Incandescent Lamp",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_wbl_type_warm_light_lamp",
			Name: "Warm Light Lamp",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_wbl_type_natural_light",
			Name: "Natural Light",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_wbl_type_fluorescent_lamp",
			Name: "Fluorescent Lamp",
		})
	}

	if cameraType == "image_ehn_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_ehn_type_off",
			Name: "Fluorescent Lamp",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_ehn_type_normal",
			Name: "Normal",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_ehn_type_expert",
			Name: "Expert",
		})
	}

	if cameraType == "image_ehn_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_ehn_type_off",
			Name: "Fluorescent Lamp",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_ehn_type_normal",
			Name: "Normal",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_ehn_type_expert",
			Name: "Expert",
		})
	}

	if cameraType == "image_mirror_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_mirror_type_off",
			Name: "OFF",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_mirror_type_lr",
			Name: "Left/Right",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_mirror_type_ud",
			Name: "Up/Down",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "image_mirror_type_ctr",
			Name: "Center",
		})
	}

	if cameraType == "osd_setting_time_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_time_type_24hr",
			Name: "OFF",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_time_type_12hr",
			Name: "Left/Right",
		})
	}

	if cameraType == "osd_setting_date_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_date_type_ymd_format1",
			Name: "YYYY-MM-DD",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_date_type_mdy_format1",
			Name: "MM-DD-YYYY",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_date_type_dmy_format1",
			Name: "DD-MM-YYYY",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_date_type_ymd_format2",
			Name: "YYYY/MM/DD",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_date_type_mdy_format2",
			Name: "MM/DD/YYYY",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_date_type_dmy_format2",
			Name: "DD/MM/YYYY",
		})
	}

	if cameraType == "osd_setting_dp_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_dp_type_tr_fl_tt",
			Name: "Transparent & Flashing",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_dp_type_tr_fl_tf",
			Name: "Transparent & Not Flashing",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_dp_type_tr_fl_ft",
			Name: "Not Transparent & Flashing",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_dp_type_tr_fl_ff",
			Name: "Not Transparent & Not Flashing",
		})
	}

	if cameraType == "osd_setting_size_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_size_type_16",
			Name: "16*16",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_size_type_32",
			Name: "32*32",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_size_type_48",
			Name: "48*48",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_size_type_64",
			Name: "64*64",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_size_type_auto",
			Name: "Auto",
		})
	}

	if cameraType == "osd_setting_font_type" {
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_font_type_bl_w",
			Name: "Black&White Self-Adaptive",
		})
		jsonConfigTypeRsp.Data = append(jsonConfigTypeRsp.Data, models.KeyValue{
			ID:   "osd_setting_font_type_ctm",
			Name: "Custom",
		})
	}

	jsonConfigTypeRsp.Code = 0
	c.JSON(http.StatusOK, &jsonConfigTypeRsp)
}
