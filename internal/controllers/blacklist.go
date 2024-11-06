package controllers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
	"vms/appconfig"
	"vms/internal/models"

	"vms/comongo/minioclient"

	"vms/comongo/reposity"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// var varURLdowprivate = "https://sbs.basesystem.one/ivis/storage/api/v0/libraries/download/"

// CreateCreateBlackList		godoc
// @Summary      	Create a new blacklist
// @Description  	Takes a blacklist JSON and store in DB. Return saved JSON.
// @Tags         	blacklist
// @Produce			json
// @Param        	BlackList  body   models.DTO_BlackList_Created  true  "Event JSON"
// @Success      	200   {object}  models.JsonDTORsp[models.DTO_BlackList_Created]
// @Router       	/blacklist [post]
// @Security		BearerAuth
func CreateBlacklist(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_BlackList_Created]()
	// Call BindJSON to bind the received JSON to
	var dto models.DTO_BlackList_Created
	var arrayImage []models.ImageDetails
	var itemImage models.ImageDetails

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

	arrayImage = append(arrayImage, itemImage)

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
		}
		dto.Imgs[i].URLImage = image
		itemImage.URL = image
		itemImage.ID = data.ID
		arrayImage = append(arrayImage, itemImage)

	}

	code, _, status, err := RegisterBlackList(arrayImage, dto.ID.String(), false, dto.Type)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	if code == 200 {
		dto, err = reposity.CreateItemFromDTO[models.DTO_BlackList_Created, models.BlackList](dto)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

		jsonRsp.Data = dto
		jsonRsp.Message = "Created sussuces"
		c.JSON(http.StatusCreated, &jsonRsp)
	} else {
		// err := reposity.DeleteItemByID[models.BlackList](dto.ID.String())
		// if err != nil {
		// 	jsonRsp.Code = http.StatusInternalServerError
		// 	jsonRsp.Message = err.Error()
		// 	c.JSON(http.StatusInternalServerError, &jsonRsp)
		// 	return
		// }
		mess := ConverCodeRegisterBlackList(code)
		jsonRsp.Data = dto
		jsonRsp.Message = mess + status
		c.JSON(http.StatusInternalServerError, &jsonRsp)
	}
}

// UpdateBlackList		 	godoc
// @Summary      	Update single blacklist by id
// @Description  	Updates and returns a single event whose ID value matches the id. New data must be passed in the body.
// @Tags         	blacklist
// @Produce      	json
// @Param        	id  path  string  true  "Update blacklist by id"
// @Param        	blacklist  body      models.DTO_BlackList_Edit  true  "blacklist JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_BlackList_Edit]
// @Router       	/blacklist/{id} [put]
// @Security		BearerAuth
func UpdateName(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_BlackList_Edit]()

	var dto models.DTO_BlackList_Edit
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	// token, err := LoginIAM()
	// if err != nil {
	// 	jsonRsp.Code = http.StatusInternalServerError
	// 	jsonRsp.Message = err.Error()
	// 	c.JSON(http.StatusInternalServerError, &jsonRsp)
	// 	return
	// }

	var arrayImage []models.ImageDetails
	var itemImage models.ImageDetails
	itemImage.ID = dto.MainImageID
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

	itemImage.URL = image
	arrayImage = append(arrayImage, itemImage)
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
		itemImage.URL = image
		arrayImage = append(arrayImage, itemImage)

		itemImage.ID = data.ID

	}

	code, _, status, err := RegisterBlackList(arrayImage, c.Param("id"), false, dto.Type)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	if code == 200 {

		// Update entity from DTO
		dto, err = reposity.UpdateItemByIDFromDTO[models.DTO_BlackList_Edit, models.BlackList](c.Param("id"), dto)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}

		// Get all images for frontend
		dtoOld, err := reposity.ReadItemByIDIntoDTO[models.DTO_BlackList_Edit, models.BlackList](c.Param("id"))
		if err != nil {
			jsonRsp.Code = http.StatusNotFound
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusNotFound, &jsonRsp)
			return
		}
		StorageBucket := "ivis-storage"
		imagelink, err := GetLibraryNameByID(dtoOld.MainImageID)
		if err != nil {
			fmt.Println("Error get name of image blacklist : ", err)
		}
		//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
		image, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
		if err != nil {
			fmt.Println("Err get url file token of image blacklist")
		}

		dtoOld.MainImageURL = image

		for j, value := range dtoOld.Imgs {
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

			dtoOld.Imgs[j].URLImage = image
		}

		jsonRsp.Data = dtoOld
		c.JSON(http.StatusOK, &jsonRsp)
	} else {
		mess := ConverCodeRegisterBlackList(code)
		jsonRsp.Data = dto
		jsonRsp.Message = mess + status
		c.JSON(http.StatusInternalServerError, &jsonRsp)
	}
}

// ReadBlackList		 godoc
// @Summary      Get single blacklist by id
// @Description  Returns the blacklist whose ID value matches the id.
// @Tags         blacklist
// @Produce      json
// @Param        id  path  string  true  "Search blacklist by id"
// @Success      200   {object}  models.JsonDTORsp[models.DTO_BlackList]
// @Router       /blacklist/{id} [get]
// @Security		BearerAuth
func ReadBlackList(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_BlackList]()

	dto, err := reposity.ReadItemByIDIntoDTO[models.DTO_BlackList, models.BlackList](c.Param("id"))
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

// DeleteBlackList		 	godoc
// @Summary      	Delete single blacklist by id
// @Description  	Delete and returns a single event whose ID value matches the id. New data must be passed in the body.
// @Tags         	blacklist
// @Produce      	json
// @Param        	id   path    string  true  "Update blacklist by id"
// @Param        	type query   string  false "Type of the blacklist" Enums(member, blacklist, unknown)
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_BlackList]
// @Router       	/blacklist/{id} [delete]
// @Security		BearerAuth
func DeleteBlackList(c *gin.Context) {

	jsonRsp := models.NewJsonDTORsp[models.DTO_BlackList]()
	var dto models.DTO_BlackList
	faceregType := c.Query("type")
	if faceregType == "" {
		faceregType = "blacklist"
	}

	err := DeletedFacereg(c.Param("id"), faceregType)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Mark the blacklist as deleted
	if faceregType == "blacklist" {
		dto.DeleteMark = true
		dto, err = reposity.UpdateItemByIDFromDTO[models.DTO_BlackList, models.BlackList](c.Param("id"), dto)
		if err != nil {
			jsonRsp.Code = http.StatusInternalServerError
			jsonRsp.Message = err.Error()
			c.JSON(http.StatusInternalServerError, &jsonRsp)
			return
		}
	}

	jsonRsp.Data = dto
	jsonRsp.Message = "Deleted success"
	c.JSON(http.StatusOK, &jsonRsp)
}

// DeleteBlackList		 	godoc
// @Summary      	Delete single blacklist by id
// @Description  	Delete and returns a single event whose ID value matches the id. New data must be passed in the body.
// @Tags         	blacklist
// @Produce      	json
// @Param        	blacklist  body      models.DTO_BlackList_Ids  true  "blacklist JSON"
// @Success      	200  {object}  models.JsonDTORsp[models.DTO_BlackList_Ids]
// @Router       	/blacklist [delete]
// @Security		BearerAuth
func DeleteBlackLists(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[models.DTO_BlackList_Ids]()

	var dto models.DTO_BlackList_Ids
	if err := c.ShouldBindJSON(&dto); err != nil {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}

	var dtoData models.DTO_BlackList
	dtoData.DeleteMark = true

	// Update entity from DTO
	for _, m := range dto.IDs {
		_, err := reposity.UpdateItemByIDFromDTO[models.DTO_BlackList, models.BlackList](m, dtoData)
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

// GetBlackList		godoc
// @Summary      	Get all blacklist groups with query filter
// @Description  	Responds with the list of all cabin as JSON.
// @Tags         	blacklist
// @Param   		keyword			query	string	false	"blacklist name keyword"		minlength(1)  	maxlength(100)
// @Param   		sort  			query	string	false	"sort"						default(+created_at)
// @Param   		limit			query	int     false  	"limit"          			minimum(1)    	maximum(100)
// @Param   		page			query	int     false  	"page"          			minimum(1) 		default(1)
// @Produce      	json
// @Success      	200  {object}   models.JsonDTOListRsp[models.DTO_BlackList]
// @Router       	/blacklist [get]
// @Security		BearerAuth
func GetBlackList(c *gin.Context) {
	jsonRspDTOCabinsBasicInfos := models.NewJsonDTOListRsp[models.DTO_BlackList]()

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
	query := reposity.NewQuery[models.DTO_BlackList, models.BlackList]()

	// Search for keyword in name
	if keyword != "" {
		query.AddConditionOfTextField("AND", "name", "LIKE", keyword)
	}
	query.AddConditionOfTextField("AND", "deleted_mark", "=", false)

	// Exec query
	dtoBlacklist, count, err := query.ExecWithPaging(sort, limit, page)
	if err != nil {
		jsonRspDTOCabinsBasicInfos.Code = http.StatusInternalServerError
		jsonRspDTOCabinsBasicInfos.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRspDTOCabinsBasicInfos)
		return
	}

	for i, dto := range dtoBlacklist {
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

		dtoBlacklist[i].MainImageURL = image

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
			dtoBlacklist[i].Imgs[j].URLImage = image
		}
	}

	jsonRspDTOCabinsBasicInfos.Count = count
	jsonRspDTOCabinsBasicInfos.Data = dtoBlacklist
	jsonRspDTOCabinsBasicInfos.Page = int64(page)
	jsonRspDTOCabinsBasicInfos.Size = int64(len(dtoBlacklist))
	c.JSON(http.StatusOK, &jsonRspDTOCabinsBasicInfos)
}

// SearchBlacklist godoc
// @Summary      Search a blacklist
// @Description  Searches the blacklist with the provided image and parameters.
// @Tags         blacklist
// @Produce      json
// @Param        searchRequest body SearchRequest true "Search Request"
// @Success      200
// @Router       /blacklist/search [post]
// @Security     BearerAuth
func SearchBlacklist(c *gin.Context) {
	var req SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	/*
		fmt.Printf("Received base64 image data of length: %d\n", len(req.Image))
		imageData, err := base64.StdEncoding.DecodeString(req.Image)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 image data"})
			return
		}
		fmt.Printf("Decoded image data length: %d\n", len(imageData))
		tempFile, err := os.CreateTemp("", "upload-*.jpg")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
			return
		}
		defer tempFile.Close()
		if _, err := tempFile.Write(imageData); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write image data"})
			return
		}
	*/
	// Login to IAM to get the token
	token, err := LoginIAM()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login to IAM"})
		return
	}

	var apiRsp models.APISearchBlacklistResponse
	apiRsp.TopkMembers = make([]models.VectorSearchMemberData, 0)
	apiRsp.TopkBlacklists = make([]models.VectorSearchMemberData, 0)

	// Perform the search
	vectorSearchResp, err := SearchVectorForBlacklist(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Loop through each TopkMembers entry and fetch member data and user details
	for _, topkMember := range vectorSearchResp.TopkMembers {
		memberData, err := GetMemberFaceImage(topkMember.MemberID, token)
		if err != nil {
			log.Printf("Failed to get face image for member ID %s: %v", topkMember.MemberID, err)

		}
		userDetail, err := GetUserDetail(topkMember.MemberID, token)
		if err != nil {
			log.Printf("Failed to get user details for member ID: %s: %v", topkMember.MemberID, err)
			//continue //Skip entry if there isn't enough information
		}
		apiRsp.TopkMembers = append(apiRsp.TopkMembers, models.VectorSearchMemberData{
			Distance:   topkMember.Distance,
			MemberID:   topkMember.MemberID,
			Data:       memberData,
			UserDetail: userDetail,
		})
	}

	// Loop through each TopkBlacklists entry and fetch member data and user details
	for _, topkBlacklist := range vectorSearchResp.TopkBlacklists {

		// Query the database to get the blacklist data
		dtoBlacklist, err := reposity.ReadItemByIDIntoDTO[models.DTO_BlackList, models.BlackList](topkBlacklist.MemberID)
		if err != nil {
			log.Printf("Failed to get blacklist data for blacklist ID %s: %v", topkBlacklist.MemberID, err)
			continue
		}

		err = minioclient.Connect(
			appconfig.Minio_Endpoint,          //Cfg.GetString("minio_endpoint"),
			appconfig.Minio_Accesskey_ID,      //Cfg.GetString("minio_accesskey_id"),
			appconfig.Minio_Secret_Access_Key, //Cfg.GetString("minio_secret_access_key"),
			appconfig.Minio_UseSSL,            //Cfg.GetBool("minio_usessl")
		)
		if err != nil {
			panic("Failed to connect to minio, err: " + err.Error())
		}

		StorageBucket := "ivis-storage"
		imagelink, err := GetLibraryNameByID(dtoBlacklist.MainImageID)
		if err != nil {
			fmt.Println("Error get name of image blacklist : ", err)
		}
		//fmt.Println("Image: ", value.Image, value.ImageResult, value.ImageObject)
		image, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
		if err != nil {
			fmt.Println("Err get url file token of image blacklist")
		}

		dtoBlacklist.MainImageURL = image
		for j, value := range dtoBlacklist.Imgs {
			imagelink, err := GetLibraryNameByID(value.ID)
			if err != nil {
				fmt.Println("Error getting name of image in blacklist: ", err)
				continue // Bỏ qua ảnh này nếu gặp lỗi
			}
			err = minioclient.Connect(
				"dev-minio-api.basesystem.one",     //Cfg.GetString("minio_endpoint"),
				"PNFIK0TCWXFZQKU0",                 //Cfg.GetString("minio_accesskey_id"),
				"SU54IHJRCR3SLH4C1GXMPAZWVJFPJOPP", //Cfg.GetString("minio_secret_access_key"),
				true,                               //Cfg.GetBool("minio_usessl")
			)
			if err != nil {
				panic("Failed to connect to minio, err: " + err.Error())
			}
			StorageBucket := "ivis-storage"
			// Lấy URL token cho từng ảnh và gán vào URLImage
			imageURL, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
			if err != nil {
				fmt.Println("Error getting URL file token of image in blacklist")
				continue // Bỏ qua ảnh này nếu gặp lỗi
			}
			dtoBlacklist.Imgs[j].URLImage = imageURL
		}
		var memberData []models.MemberFaceImageData = make([]models.MemberFaceImageData, 0)
		memberData = append(memberData, models.MemberFaceImageData{
			ImageFileID:  dtoBlacklist.MainImageID,
			ImageFileURL: dtoBlacklist.MainImageURL,
			UserID:       dtoBlacklist.ID.String(),
			FaceType:     "CENTER",
			ID:           dtoBlacklist.ID.String(),
		})

		apiRsp.TopkBlacklists = append(apiRsp.TopkBlacklists, models.VectorSearchMemberData{
			Distance: topkBlacklist.Distance,
			MemberID: topkBlacklist.MemberID,
			Data:     memberData,
			UserDetail: models.UserDetail{
				UserID:      topkBlacklist.MemberID,
				Username:    dtoBlacklist.Name,
				FullName:    dtoBlacklist.Name,
				Gender:      "",
				Email:       "",
				PhoneNumber: "",
			},
		})
	}

	// Return the response from the external API to the client
	c.JSON(http.StatusOK, apiRsp)
}

// GetBlacklistImage godoc
// @Summary      Get blacklist images by member ID
// @Description  Returns the blacklist images associated with the provided member ID.
// @Tags         blacklist
// @Produce      json
// @Param        member_id  path  string  true  "Member ID"
// @Success      200
// @Router       /blacklist/images/{member_id} [get]
// @Security     BearerAuth
func GetBlacklistImage(c *gin.Context) {
	jsonRsp := models.NewJsonDTORsp[[]models.MemberFaceImageData]()
	ChangeImageURL()
	memberID := c.Param("member_id")
	if memberID == "" {
		jsonRsp.Code = http.StatusBadRequest
		jsonRsp.Message = "Member ID is required"
		c.JSON(http.StatusBadRequest, &jsonRsp)
		return
	}
	token, err := LoginIAM()
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}
	// Call the GetMemberID function to fetch the member data
	memberData, err := GetMemberFaceImage(memberID, token)
	if err != nil {
		jsonRsp.Code = http.StatusInternalServerError
		jsonRsp.Message = err.Error()
		c.JSON(http.StatusInternalServerError, &jsonRsp)
		return
	}

	// Return the response with the member data
	jsonRsp.Data = memberData
	jsonRsp.Message = "Successfully fetched blacklist images"
	c.JSON(http.StatusOK, &jsonRsp)
}

func GetLibraryNameByID(id string) (string, error) {
	// Connect to side DB ivi-Storage
	err := ConnectToStorageDB()
	if err != nil {
		return "", fmt.Errorf("Could not connect to storage database: %v", err)
	}

	defer CloseStorageDB()

	var library models.Library

	if err := storageDB.Where("id = ?", id).First(&library).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("Library with ID %s not found", id)
		}
		return "", fmt.Errorf("Error querying Library by ID: %v", err)
	}

	// return name of image in side database(url from MinIO)
	return library.Name, nil
}

// UpdateImageURL updates the URLImage of an image in the database based on its ID.
func UpdateStorageImage(imageID, newName, Type string) error {
	// Connect to the storage database
	err := ConnectToStorageDB()
	if err != nil {
		return fmt.Errorf("Could not connect to storage database: %v", err)
	}
	defer CloseStorageDB()

	// Find the image record by ID
	var library models.Library
	if err := storageDB.Where("id = ?", imageID).First(&library).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("Image with ID %s not found", imageID)
		}
		return fmt.Errorf("Error querying Library by ID: %v", err)
	}

	// Update the URLImage, updated_by, and updated_at fields
	library.Name = newName
	library.Type = Type
	library.UpdatedAt = time.Now()

	// Save the changes to the database
	if err := storageDB.Save(&library).Error; err != nil {
		return fmt.Errorf("Error updating image URL: %v", err)
	}

	return nil
}

func SaveImageToDatabase(uploadInfo minio.UploadInfo, objectName, userID, username, typeImage string, size int) (string, error) {
	// Kết nối tới cơ sở dữ liệu storage nếu chưa kết nối
	err := ConnectToStorageDB()
	if err != nil {
		return "", fmt.Errorf("Could not connect to storage database: %v", err)
	}
	defer CloseStorageDB()

	// Tạo UUID cho ID ảnh mới
	imageID, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("Failed to generate image ID: %v", err)
	}

	// Khởi tạo đối tượng Library để lưu thông tin ảnh mới
	library := models.Library{
		ID:                  imageID,
		CreatedAt:           time.Now(),
		CreatedBy:           "admin",
		UpdatedAt:           time.Now(),
		UpdatedBy:           "admin",
		IsPublic:            true,       // Giả định ảnh được công khai
		Name:                objectName, // Đường dẫn của ảnh trên MinIO
		Size:                size,       // Kích thước của ảnh
		Type:                typeImage,  // Loại ảnh
		UserID:              userID,
		Username:            username,
		GetOriginalFilename: &uploadInfo.Key, // Lưu tên file gốc
	}

	// Chèn đối tượng vào cơ sở dữ liệu
	if err := storageDB.Create(&library).Error; err != nil {
		return "", fmt.Errorf("Error inserting image into storage database: %v", err)
	}

	fmt.Println("Image record created successfully in database")
	return library.ID.String(), nil
}

// Todo: Remove following code

// Tạo một đối tượng lưu trữ kết nối phụ cho storage
var storageDB *gorm.DB

// ConnectToStorageDB: Kết nối tới cơ sở dữ liệu phụ (storage)
func ConnectToStorageDB() error {
	currentSchema := "ivis_storage"
	sqlDsn := fmt.Sprintf("host=4.194.17.112 port=5432 dbname=storage sslmode=disable user=ivissbs password=hugQRT8k5HLt4b4a2kHq7fRql9i30E51ggbM")

	// Kết nối với PostgreSQL cho cơ sở dữ liệu phụ
	database, err := gorm.Open(postgres.New(postgres.Config{
		DSN: sqlDsn,
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   currentSchema + ".", // schema name
			SingularTable: true,                // use singular table name
		},
	})

	if err != nil {
		return fmt.Errorf("Failed to connect to storage database: %w", err)
	}

	// Lưu đối tượng kết nối vào biến global storageDB
	storageDB = database

	// Thiết lập pool kết nối cho storageDB nếu cần
	sqlDB, err := storageDB.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	fmt.Println("Success connect to Storage DB")
	return nil
}

// Đóng kết nối cơ sở dữ liệu storage
func CloseStorageDB() error {
	sqlDB, err := storageDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// UploadImage godoc
// @Summary      Upload an image to MinIO and return metadata
// @Description  Uploads an image to MinIO, saves the metadata in the auxiliary table, and returns the image metadata.
// @Tags         images
// @Accept       multipart/form-data
// @Produce      json
// @Param        files formData file true "Image file"
// @Success      200   {object} models.ListImgKeyValue
// @Router       /images/upload [post]
// @Security     BearerAuth
func UploadImage(c *gin.Context) {
	// Lấy file từ form-data
	file, err := c.FormFile("files")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload error: " + err.Error()})
		return
	}

	// Kết nối với MinIO
	err = minioclient.Connect(
		"dev-minio-api.basesystem.one",     //Cfg.GetString("minio_endpoint"),
		"PNFIK0TCWXFZQKU0",                 //Cfg.GetString("minio_accesskey_id"),
		"SU54IHJRCR3SLH4C1GXMPAZWVJFPJOPP", //Cfg.GetString("minio_secret_access_key"),
		true,                               //Cfg.GetBool("minio_usessl")
	)
	if err != nil {
		panic("Failed to connect to minio, err: " + err.Error())
	}

	// Mở file để đọc nội dung
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file: " + err.Error()})
		return
	}
	defer src.Close()

	// Đọc toàn bộ nội dung của file vào buffer
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file: " + err.Error()})
		return
	}

	// Chuyển đổi nội dung file thành chuỗi base64
	imageBase64 := base64.StdEncoding.EncodeToString(buffer.Bytes())

	// Định nghĩa tên object trong MinIO
	objectName := fmt.Sprintf("IAM/%s", file.Filename)
	storageBucket := "ivis-storage"

	// Upload ảnh lên MinIO sử dụng hàm UploadImageBase64
	uploadInfo, err := minioclient.UploadImageBase64(storageBucket, objectName, imageBase64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to MinIO: " + err.Error()})
		return
	}

	// Lấy URL của file đã upload
	imageURL, err := minioclient.GetPresignedURL(storageBucket, objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get image URL: " + err.Error()})
		return
	}

	imageID, err := SaveImageToDatabase(uploadInfo, objectName, "", "", "jpg", buffer.Len())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image to database: " + err.Error()})
		return
	}

	// Phản hồi cho frontend
	imgData := models.ListImgKeyValue{
		ID:            imageID, // Sử dụng ID từ cơ sở dữ liệu
		ImageName:     objectName,
		StorageBucket: storageBucket, // bucket name
		URLImage:      imageURL,      // đường dẫn đến ảnh trong minio
		Type:          "jpeg",        // Định dạng mặc định
	}

	c.JSON(http.StatusOK, imgData)
}

// Change all Main iammge url

func ChangeImageURL() {

	query := reposity.NewQuery[models.DTO_BlackList, models.BlackList]()
	dtoBlacklist, _, err := query.ExecNoPaging("-created_at")
	if err != nil {
		fmt.Println("Can't get all blacklist from DB")
		return
	}

	for i := range dtoBlacklist {
		StorageBucket := "ivis-storage"

		// Lấy URL của MainImageID
		imagelink, err := GetLibraryNameByID(dtoBlacklist[i].MainImageID)
		if err != nil {
			fmt.Println("Error getting name of main image blacklist: ", err)
			continue // Bỏ qua phần tử này nếu gặp lỗi
		}

		// Lấy URL token cho MainImage và gán vào MainImageURL
		mainImageURL, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
		if err != nil {
			fmt.Println("Error getting URL file token of main image blacklist")
			continue // Bỏ qua phần tử này nếu gặp lỗi
		}
		dtoBlacklist[i].MainImageURL = mainImageURL

		// Xử lý các hình ảnh trong Imgs
		for j := range dtoBlacklist[i].Imgs {
			imagelink, err := GetLibraryNameByID(dtoBlacklist[i].Imgs[j].ID)
			if err != nil {
				fmt.Println("Error getting name of image in blacklist: ", err)
				continue // Bỏ qua ảnh này nếu gặp lỗi
			}

			// Lấy URL token cho từng ảnh và gán vào URLImage
			imageURL, err := minioclient.GetPresignedURL(StorageBucket, imagelink)
			if err != nil {
				fmt.Println("Error getting URL file token of image in blacklist")
				continue // Bỏ qua ảnh này nếu gặp lỗi
			}
			dtoBlacklist[i].Imgs[j].URLImage = imageURL
		}

		// Cập nhật vào cơ sở dữ liệu
		if _, err := reposity.UpdateItemByIDFromDTO[models.DTO_BlackList, models.BlackList](dtoBlacklist[i].ID.String(), dtoBlacklist[i]); err != nil {
			log.Fatal("Failed to update DB of blacklist: ", err)
			return
		}
		fmt.Println("success uppdate ImgaeURL")
	}

}
