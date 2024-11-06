package controllers

import (
	"fmt"
	"strings"
	"vms/internal/models"
)

func UpdatePasswordCamera(dto models.DTOCamera, dtoNetWorkConfig models.DTO_NetworkConfig, VideoConfig models.VideoStreamArray) models.VideoStreamArray {
	videoStreamArray := models.VideoStreamArray{}
	for entry, channel := range dto.Streams {
		channelType := "main"
		streamID := "101"
		if channel.Name == "subStream" {
			channelType = "sub"
			streamID = "102"
		}

		rtspURL := fmt.Sprintf("rtsp://%s:%s@%s:%d/Streaming/Channels/%s", dto.Username, dto.Password, dto.IPAddress, dtoNetWorkConfig.RTSP, streamID)

		stream := models.VideoStream{
			Name:      channel.Name,
			Type:      channelType,
			URL:       rtspURL,
			IsProxied: false,
			IsDefault: true,
			Channel:   channelType,
			ID:        dto.ID.String(),
			Codec:     strings.ToLower(VideoConfig[entry].Codec),
		}
		videoStreamArray = append(videoStreamArray, stream)
	}
	return videoStreamArray
}
