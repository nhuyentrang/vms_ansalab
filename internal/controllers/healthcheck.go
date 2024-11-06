package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"vms/internal/models"

	"vms/comongo/kafkaclient"

	"github.com/gin-gonic/gin"
)

// TODO check for status online / offline of NVR
// Health		 godoc
// @Summary      Health is used to handle HTTP Health requests to this service.
// @Description  Use this for liveness probes or any other checks which only validate if the services is running.
// @Tags         healthcheck
// @Success      200
// @Router       /health [get]
func Health(c *gin.Context) {
	cmd := models.DeviceCommand{
		CommandID:    "799be4b6-6492-47a9-8a6e-7af19b3858cb",
		Cmd:          cmd_GetOnSiteVideoConfig,
		EventTime:    time.Now().Format(time.RFC3339),
		IPAddress:    "192.168.2.204",
		UserName:     "admin",
		Password:     "123456",
		HttpPort:     "80",
		ProtocolType: "ONVIF",
		CameraID:     "1dc981ff-6c81-4bd6-982b-b053dafec56d",
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommand)
	c.String(http.StatusOK, "Success")
}

// Ready		 godoc
// @Summary      Ready is used to handle HTTP Ready requests to this service.
// @Description  Use this for readiness probes or any checks that validate the service is ready to accept traffic.
// @Tags         healthcheck
// @Success      200
// @Router       /ready [get]
func Ready(c *gin.Context) {
	c.String(http.StatusOK, "Success")
}
