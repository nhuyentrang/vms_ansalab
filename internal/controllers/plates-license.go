package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"vms/internal/models"

	"vms/comongo/minioclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
)

// LicensePlates		godoc
// @Summary      	Create a new LicensePlates
// @Description  	Takes a license plates JSON and store in DB. Return saved JSON.
// @Tags         	LicensePlates
// @Produce			json
// @Param        	licenseplates  body   models.LicensePlates  true  "Event JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.LicensePlates]
// @Router       	/licenseplates [post]
// @Security		BearerAuth
func CreateLicensePlates(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_LicensePlates]()

	var itemImage models.ImageDetails
	// Call BindJSON to bind the received JSON to
	var dto models.DTO_LicensePlates
	if err := c.BindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	itemImage.ID = dto.MainImageID
	StorageBucket := "ivis-storage"
	imagelink, err := GetLibraryNameByID(itemImage.ID)
	if err != nil {
		fmt.Println("Error get name of image blacklist : ", err)
	}
	//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
	image, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
	if err != nil {
		fmt.Println("Err get url file token of image blacklist")
	}

	dto.MainImageURL = image

	for i, data := range dto.Imgs {
		StorageBucket := "ivis-storage"
		imagelink, err := GetLibraryNameByID(data.ID)
		if err != nil {
			fmt.Println("Error get name of image blacklist : ", err)
		}
		//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
		image, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
		if err != nil {
			fmt.Println("Err get url file token of image blacklist")
			return
		}
		dto.Imgs[i].URLImage = image

	}
	dto, err = reposity.CreateItemFromDTO[models.DTO_LicensePlates, models.LicensePlates](dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	jsonRsp.Message = "Successfully"
	c.JSON(http.StatusCreated, &jsonRsp)
}

// LicensePlates		 	godoc
// @Summary      	Update single licenseplates by id
// @Description  	Updates and returns a single event whose ID value matches the id. New data must be passed in the body.
// @Tags         	licenseplates
// @Produce      	json
// @Param        	id  path  string  true  "Update licenseplates by id"
// @Param        	licenseplates  body      models.DTO_LicensePlates  true  "licenseplates JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_LicensePlates]
// @Router       	/licenseplates/{id} [put]
// @Security		BearerAuth
func UpdateLicensePlates(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_LicensePlates]()

	var dto models.DTO_LicensePlates
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	StorageBucket := "ivis-storage"

	for _, value := range dto.Imgs {
		if value.URLImage != "" {
			// Update the main image information in the auxiliary database
			err := UpdateStorageImage(value.ID, value.ImageName, value.Type)
			if err != nil {
				fmt.Println("Error updating main image information in auxiliary database: ", err)
			}
		}

	}
	imagelink, err := GetLibraryNameByID(dto.MainImageID)
	if err != nil {
		fmt.Println("Error get name of image blacklist : ", err)
	}

	//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
	image, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
	if err != nil {
		fmt.Println("Err get url file token of image blacklist")
	}
	dto.MainImageURL = image

	for i, data := range dto.Imgs {
		if data.URLImage != "" {
			// Update the main image information in the auxiliary database
			err = UpdateStorageImage(data.ID, data.ImageName, data.Type)
			if err != nil {
				fmt.Println("Error updating main image information in auxiliary database: ", err)
			}
		}
		StorageBucket := "ivis-storage"
		imagelink, err := GetLibraryNameByID(data.ID)
		if err != nil {
			fmt.Println("Error get name of image blacklist : ", err)
		}
		//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
		image, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
		if err != nil {
			fmt.Println("Err get url file token of image blacklist")
		}

		dto.Imgs[i].URLImage = image

	}
	dto, err = reposity.UpdateItemByIDFromDTO[models.DTO_LicensePlates, models.LicensePlates](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// LicensePlates		 godoc
// @Summary      Get single licenseplates by id
// @Description  Returns the licenseplates whose ID value matches the id.
// @Tags         licenseplates
// @Produce      json
// @Param        id  path  string  true  "Search licenseplates by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_LicensePlates]
// @Router       /licenseplates/{id} [get]
// @Security		BearerAuth
func ReadLicensePlates(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_LicensePlates]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_LicensePlates, models.LicensePlates](c.Param("id"))
	if err != nil {
		jsonRsp.Code = http.StatusNotFound
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusNotFound, &jsonRsp)
		return
	}

	StorageBucket := "ivis-storage"
	imagelink, err := GetLibraryNameByID(dto.MainImageID)
	if err != nil {
		fmt.Println("Error get name of image blacklist : ", err)
	}
	//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
	image, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
	if err != nil {
		fmt.Println("Err get url file token of image blacklist")
	}

	dto.MainImageURL = image

	for j, value := range dto.Imgs {
		value.StorageBucket = "ivis-storage"
		imagelink, err := GetLibraryNameByID(value.ID)
		if err != nil {
			fmt.Println("Error get name of image blacklist : ", err)
		}
		//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
		image, err := minioclient.GetPresignedURL(value.StorageBucket, imagelink)
		if err != nil {
			fmt.Println("Err get url file token of image blacklist")
		}

		dto.Imgs[j].URLImage = image
	}

	jsonRsp.Data = dto

	c.JSON(http.StatusOK, &jsonRsp)
}

// LicensePlates		 	godoc
// @Summary      	Delete single licenseplates by id
// @Description  	Delete and returns a single event whose ID value matches the id. New data must be passed in the body.
// @Tags         	licenseplates
// @Produce      	json
// @Param        	id  path  string  true  "Update licenseplates by id"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_LicensePlates]
// @Router       	/licenseplates/{id} [delete]
// @Security		BearerAuth
func DeleteLicensePlates(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_LicensePlates]()
	var dto models.DTO_LicensePlates
	dto.DeleteMark = true

	dto, err := reposity.UpdateItemByIDFromDTO[models.DTO_LicensePlates, models.LicensePlates](c.Param("id"), dto)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	jsonRsp.Data = dto
	jsonRsp.Message = "Deleted success"
	c.JSON(http.StatusOK, &jsonRsp)
}

// LicensePlates		 	godoc
// @Summary      	Delete single licenseplates by id
// @Description  	Delete and returns a single event whose ID value matches the id. New data must be passed in the body.
// @Tags         	licenseplates
// @Produce      	json
// @Param        	licenseplates  body      models.DTO_LicensePlates_Ids  true  "licenseplates JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_LicensePlates_Ids]
// @Router       	/licenseplates [delete]
// @Security		BearerAuth
func DeleteLicensePlatess(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_LicensePlates_Ids]()

	var dto models.DTO_LicensePlates_Ids
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	var dtoData models.DTO_LicensePlates
	dtoData.DeleteMark = true

	// Update entity from DTO
	for _, m := range dto.IDs {
		_, err := reposity.UpdateItemByIDFromDTO[models.DTO_LicensePlates, models.LicensePlates](m, dtoData)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
	}
	jsonRsp.Data = dto
	c.JSON(http.StatusOK, &jsonRsp)
}

// LicensePlates		godoc
// @Summary      	Get all licenseplates groups with query filter
// @Description  	Responds with the list of all cabin as JSON.
// @Tags         	licenseplates
// @Param   		keyword			query	string	false	"licenseplates name keyword"		minlength(1)  	maxlength(100)
// @Param   		sort  			query	string	false	"sort"						default(+created_at)
// @Param   		limit			query	int     false  	"limit"          			minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          			minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_LicensePlates]
// @Router       	/licenseplates [get]
// @Security		BearerAuth
func GetLicensePlates(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_LicensePlates]()

	// Get param
	keyword := c.Query("keyword")
	sort := c.Query("sort")
	limit, _ := strconv.Atoi(c.Query("limit"))
	page, _ := strconv.Atoi(c.Query("page"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 128
	}
	if page < 1 {
		page = 1
	}
	// Build query
	query := reposity.NewQuery[models.DTO_LicensePlates, models.LicensePlates]()

	// Search for keyword in name
	if keyword != "" {
		query.AddConditionOfTextField("AND", "name", "LIKE", keyword)
	}
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)

	// Exec query
	dtoCabinBasics, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}

	for i, dto := range dtoCabinBasics {
		StorageBucket := "ivis-storage"
		imagelink, err := GetLibraryNameByID(dto.MainImageID)
		if err != nil {
			fmt.Println("Error get name of image blacklist : ", err)
		}
		//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
		image, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
		if err != nil {
			fmt.Println("Err get url file token of image blacklist")
		}

		dtoCabinBasics[i].MainImageURL = image

		for j, value := range dto.Imgs {
			value.StorageBucket = "ivis-storage"
			imagelink, err := GetLibraryNameByID(value.ID)
			if err != nil {
				fmt.Println("Error get name of image blacklist : ", err)
			}
			//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
			image, err := minioclient.GetPresignedURL(value.StorageBucket, imagelink)
			if err != nil {
				fmt.Println("Err get url file token of image blacklist")
			}
			dtoCabinBasics[i].Imgs[j].URLImage = image
		}
	}

	jsonRspDTOCabinsBasicInfos.Count = count
	jsonRspDTOCabinsBasicInfos.Data = dtoCabinBasics
	jsonRspDTOCabinsBasicInfos.Page = int64(page)
	jsonRspDTOCabinsBasicInfos.Size = int64(len(dtoCabinBasics))
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)
}
