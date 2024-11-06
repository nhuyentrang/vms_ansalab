package controllers

import (
	"net/http"

	"vms/internal/models"
	"vms/statuscode"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
)

// CreatePCinfo		godoc
// @Summary      	Create a new PC Info
// @Description  	Takes a PC Info JSON and store in DB. Return saved JSON.
// @Tags         	recording-schedules
// @Produce			json
// @Param        	pc_Info  body   models.DTO_PcInfo true  "PC info JSON"
// @Success      	200   {object}  models.DTO_PcInfo
// @Router       	/pc-info [post]
// @Security		BearerAuth
func CreatePCInfo(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_PcInfo]()

	// Call BindJSON to bind the received JSON to
	var dto models.DTO_PcInfo
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = statuscode.StatusBindingInputJsonFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// Create new block
	dto, err := reposity.CreateItemFromDTO[models.DTO_PcInfo, models.PcInfo](dto)
	if err != nil {
		jsonRsp.Code = statuscode.StatusCreateItemFailed
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Response
	jsonRsp.Data = dto
	c.JSON(http.StatusCreated, &jsonRsp)
}
