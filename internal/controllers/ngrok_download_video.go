package controllers

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DownloadRequest struct {
	XMLName      xml.Name `xml:"downloadRequest"`
	XMLVersion   string   `xml:"version,attr"`
	XMLNamespace string   `xml:"xmlns,attr"`
}

// GetEvents		godoc
// @Summary      	Get video with query
// @Description  	Responds with the list of all Event as JSON.
// @Tags         	Download video
// @Produce      	json
// @Param   		endTime		query	string	false	"endTime"	Enums(nicType,dhcp,mtuPlaceHolder,multicastDiscovery)
// @Param   		startTime		query	string	false	"startTime"		Enums(dnsType, portPlaceHolder)
// @Param   		ipDevice		query	string	false	"ipDevicep"		Enums(httpPortPlaceHolder, httpsPortPlaceHolder, rtspPortPlaceHolder, servicePortPlaceHolder)
// @Param   		trackID		query	string	false	"trackID"	Enums(streamType,videoType,resolution,bitrateType,videoQuality,frameRatePlaceHolder,maxBitRatePlaceHolder,videoEncoding,h264Plus,profile,iframeIntervalPlaceHolder,smoothing)
// @Param   		userName		query	string	false	"userName"	Enums(audioEncoding,audioInput,inputVolumePlaceHolder,environmentNoiseFilter)
// @Param   		passWord		query	string	false	"passWord"	Enums(dayNightMode,exposure,areaBLC,whiteBalance,gain,shutter,wdr,noiseReduction,hlc,mirror)
// @Success      	200  {object}   models.DTO_Event_Read_BasicInfo
// @Router       	/downloadVideo [get]
// @Security		BearerAuth
func DownloadVideo(c *gin.Context) {
	// Parse query parameters

	ipDevice := c.Query("ipDevice")
	userName := c.Query("userName")
	passWord := c.Query("passWord")
	trackID := c.Query("trackID")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")
	// Simulate downloading video
	videoURL := fmt.Sprintf("https://13aa-123-24-205-208.ngrok-free.app/downloadVideo?endTime=%s&startTime=%s&ipDevice=%s&trackID=%s&userName=%s&passWord=%s",
		endTime, startTime, ipDevice, trackID, userName, passWord)

	// Create a new HTTP client
	client := &http.Client{}
	// Create the download request
	downloadReq := DownloadRequest{

		XMLNamespace: "http://www.isapi.org/ver20/XMLSchema",
	}

	// Encode the request to XML
	reqBody, err := xml.Marshal(downloadReq)
	if err != nil {
		http.Error(c.Writer, fmt.Sprintf("Error encoding request: %v", err), http.StatusInternalServerError)
		return
	}

	// Create a new POST request
	req, err := http.NewRequest("GET", videoURL, bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(c.Writer, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError)
		return
	}

	// Set basic authentication
	req.SetBasicAuth(userName, passWord)
	req.Header.Set("Content-Type", "application/xml")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(c.Writer, fmt.Sprintf("Error sending request: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		http.Error(c.Writer, fmt.Sprintf("Failed to download video: %s", resp.Status), http.StatusInternalServerError)
		return
	}

	// Read the response body
	videoData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(c.Writer, fmt.Sprintf("Error reading response: %v", err), http.StatusInternalServerError)
		return
	}

	// Đặt tiêu đề cho phản hồi HTTP để gửi video về client
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=\"video.mp4\"")
	c.Writer.Header().Set("Content-Type", "video/mp4")
	c.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", len(videoData)))

	// Gửi dữ liệu video về client
	if _, err := c.Writer.Write(videoData); err != nil {
		http.Error(c.Writer, fmt.Sprintf("Error sending video: %v", err), http.StatusInternalServerError)
		return
	}

	log.Println("Video sent successfully")
}
