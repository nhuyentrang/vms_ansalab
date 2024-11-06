package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/viper"

	"vms/internal/models"

	"vms/comongo/kafkaclient"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

var (
	clientId                        string = "e7b56179-ac9d-419b-80ed-a18c41285aab"
	clientSecret                    string = "3b9fed7e-525a-472d-bf56-1aca6a4519ab"
	iamUrlLogin                     string = "https://sbs.basesystem.one/ivis/iam/api/v0/client/login"
	tenantId                        string = "ntq"
	apiSBSStorageLibraryResourceURL string = "https://sbs.basesystem.one/ivis/storage/api/v0/libraries/"
	apiRegisterBlackList            string = "https://facereg.basesystem.one/facereg/insert"
	apiDeleteBlackList              string = "https://dev-ivi-vectorsearch-api.basesystem.one/member/remove"
	apiSearchBlackList              string = "https://dev-ivi-vectorsearch-api.basesystem.one/member/search"

	apiSearchImageID string = "https://dev-ivi.basesystem.one/smc/iam/api/v0/mi-se/user-faces/image/face/"
	apiUserDetailURL string = "https://dev-ivi.basesystem.one/smc/iam/api/v0/users/"
)

type UserGroup struct {
	UserID       string  `json:"userId"`
	GroupID      string  `json:"groupId"`
	ParentID     string  `json:"parentId"`
	GroupCode    string  `json:"groupCode"`
	GroupName    string  `json:"groupName"`
	SalaryBase   float64 `json:"salaryBase"`
	IsDefault    bool    `json:"isDefault"`
	PositionName string  `json:"positionName"`
	PositionID   int     `json:"positionId"`
	PositionCode string  `json:"positionCode"`
	IsLeader     bool    `json:"isLeader"`
}

type UserAccount struct {
	ID                    string `json:"id"`
	Username              string `json:"username"`
	AccStatus             string `json:"accStatus"`
	Activated             bool   `json:"activated"`
	LastPasswordUpdatedAt int64  `json:"lastPasswordUpdatedAt"`
	IdentityProviderType  string `json:"identityProviderType"`
	Root                  bool   `json:"root"`
}

type iamLoginInfo struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type jsonStoregeGetURLRsp struct {
	Code    int             `json:"code"`
	Data    libraryResource `json:"data"`
	Message string          `json:"message"`
}

type libraryResource struct {
	Url      string `json:"url"`
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

type SearchRequest struct {
	Threshold float64 `json:"threshold" binding:"required"`
	Topk      int     `json:"topk" binding:"required"`
	Image     string  `json:"image" binding:"required"` // Base64-encoded image string
}

type ImageInfo struct {
	Base64 string `json:"base64"`
}

type MemberSearchVectorRequest struct {
	Member    int         `json:"member"`
	Blacklist int         `json:"blacklist"`
	ImageInfo []ImageInfo `json:"imageInfo"`
	Threshold float64     `json:"threshold"`
	Topk      int         `json:"topk"`
	TenantID  string      `json:"tenantId,omitempty"`
}
type MemberSearchVectorResponse struct {
	TopkBlacklists []struct {
		Distance float64 `json:"distance"`
		MemberID string  `json:"member_id"`
	} `json:"topk_blacklists"`
	TopkMembers []struct {
		Distance float64 `json:"distance"`
		MemberID string  `json:"member_id"`
	} `json:"topk_members"`
}

var c *viper.Viper

func LoginIAM() (token string, err error) {
	loginInfo := iamLoginInfo{
		ClientId:     clientId,
		ClientSecret: clientSecret,
	}
	p, err := json.Marshal(loginInfo)
	if err != nil {
		return "", errors.New("can not marshal login info, err: " + err.Error())
	}
	req, err := http.NewRequest("POST", iamUrlLogin, bytes.NewBuffer(p))
	if err != nil {
		return "", fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 5

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("can not make login request, err: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//fmt.Println("response Body:", string(body))
	if resp.StatusCode != 200 {
		return "", errors.New("fail to login, err: " + resp.Status)
	}
	value := gjson.Get(string(body), "access_token")

	//fmt.Println("==============> Token: ", value.String())

	return value.String(), nil
}

func GetLibraryResourceURL(token string, id string) (url string, err error) {
	req, err := http.NewRequest("GET", apiSBSStorageLibraryResourceURL+id, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Authorization", "bearer "+token)

	client := &http.Client{}
	client.Timeout = time.Second * 5

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("fail to download file, rsp: %v, code: %v", resp.Status, resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var urlRsp jsonStoregeGetURLRsp
	err = json.Unmarshal([]byte(body), &urlRsp)
	if err != nil {
		return "", err
	}

	return (apiSBSStorageLibraryResourceURL + "download/" + id), nil
}

func RegisterBlackList(data []models.ImageDetails, memberId string, IsForceUpdate bool, Type models.KeyValue) (codeRegisterBlackList int, topkBlacklists []models.TopkBlacklists, status string, err error) {

	var faceRegReq *models.FaceRegData
	if strings.ToLower(Type.Code) == "blacklist" || Type.Name == "Đối tượng trong danh sách đen" {
		faceRegReq = &models.FaceRegData{
			Member:        0,
			Blacklist:     1,
			ImageInfo:     data,
			MemberID:      memberId,
			Threshold:     0,
			Topk:          0,
			IsForceUpdate: IsForceUpdate,
			TenantID:      tenantId,
		}
	} else {
		faceRegReq = &models.FaceRegData{
			Member:        1,
			Blacklist:     0,
			ImageInfo:     data,
			MemberID:      memberId,
			Threshold:     0,
			Topk:          0,
			IsForceUpdate: false,
			TenantID:      tenantId,
		}
	}

	//log.Printf("======> Blacklist register data: %v", faceRegReq)

	cmd := models.DeviceCommandFEceRegister{
		CommandID:   uuid.New().String(),
		Cmd:         "face.regesiter",
		EventTime:   time.Now().Format(time.RFC3339),
		EventType:   "register",
		FaceRegData: faceRegReq,
	}
	cmsStr, _ := json.Marshal(cmd)
	kafkaclient.SendJsonMessages(string(cmsStr), M_deviceCommandFaceRegister)

	p, err := json.Marshal(faceRegReq)
	if err != nil {
		fmt.Println("can not marshal post info, err: " + err.Error())
		return http.StatusInternalServerError, []models.TopkBlacklists{}, "", errors.New("can not marshal post info, err: " + err.Error())
	}
	// Marshal the payload to JSON
	jsonData, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("Error marshaling payload: %s\n", err)
		return
	}

	// Print the JSON body to verify it before sending
	fmt.Println("Request body:", string(jsonData))

	// cfg := c
	req, err := http.NewRequest("POST", apiRegisterBlackList, bytes.NewBuffer(p))
	if err != nil {
		return http.StatusInternalServerError, []models.TopkBlacklists{}, "", fmt.Errorf("got error %s", err.Error())
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 5
	resp, err := client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, []models.TopkBlacklists{}, "", errors.New("can not make post request, err: " + err.Error())
	}

	//fmt.Println("day la response tra ve", resp)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Println("fail to post, err: " + resp.Status)
		return http.StatusInternalServerError, []models.TopkBlacklists{}, "", errors.New("fail to post, err: " + resp.Status)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return http.StatusInternalServerError, []models.TopkBlacklists{}, "", err
	}

	var code models.DataRespFaceReg
	err = json.Unmarshal([]byte(body), &code)
	if err != nil {
		return http.StatusInternalServerError, []models.TopkBlacklists{}, "", err
	}
	log.Printf("======> Code RegisterBlackList: %d", code.Code)

	return resp.StatusCode, code.TopkBlacklists, code.Status, nil
}

func DeletedFacereg(userId string, faceregType string) (err error) {
	var deleteFaceVectorReq *models.DeleteFaceVectorRequest
	if strings.ToLower(faceregType) == "blacklist" {
		deleteFaceVectorReq = &models.DeleteFaceVectorRequest{
			Member:    0,
			Blacklist: 1,
			MemberID:  userId,
			//TenantID:  tenantId,
		}
	} else if strings.ToLower(faceregType) == "member" {
		deleteFaceVectorReq = &models.DeleteFaceVectorRequest{
			Member:    1,
			Blacklist: 0,
			MemberID:  userId,
			//TenantID: tenantId,
		}
	}

	p, err := json.Marshal(deleteFaceVectorReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", apiDeleteBlackList, bytes.NewBuffer(p))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	client.Timeout = time.Second * 30

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("resp.StatusCode", resp.StatusCode)
	if resp.StatusCode != 200 && resp.StatusCode != 202 && resp.StatusCode != 204 {
		return errors.New("delete face vector response not ok, status code: " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var delFaceVectorRsp models.DeleteFaceVectorResponse
	err = json.Unmarshal(body, &delFaceVectorRsp)
	if err != nil {
		return err
	}

	// Verify the response from the external API
	if delFaceVectorRsp.Code != 0 && delFaceVectorRsp.Code != 816 {
		return fmt.Errorf("failed to delete face vector, return code: %d", delFaceVectorRsp.Code)
	}

	return nil
}

func SearchVectorForBlacklist(req SearchRequest) (rsp MemberSearchVectorResponse, err error) {
	// Read the image data from the file

	// Create the request body for the external API
	externalReq := MemberSearchVectorRequest{
		Member:    1,
		Blacklist: 1,
		ImageInfo: []ImageInfo{
			{Base64: req.Image},
		},
		Threshold: req.Threshold,
		Topk:      req.Topk,
	}

	// Convert the request body to JSON
	jsonData, err := json.Marshal(externalReq)
	if err != nil {
		return rsp, fmt.Errorf("failed to marshal request body: %w", err)
	}

	// Print the request body
	// fmt.Printf("Request body to external API:\n%s\n", string(jsonData))

	// Send the request to the external API
	resp, err := sendRequestToExternalAPI(jsonData)
	if err != nil {
		return rsp, err
	}

	// Parse the response from the external API
	var externalResp MemberSearchVectorResponse
	if err := json.Unmarshal(resp, &externalResp); err != nil {
		return rsp, fmt.Errorf("failed to unmarshal external API response: %w", err)
	}

	// Return the response from the external API
	return externalResp, nil
}

func sendRequestToExternalAPI(jsonData []byte) ([]byte, error) {
	// Log the size of the payload

	fmt.Printf("Payload size: %d bytes\n", len(jsonData))

	req, err := http.NewRequest("POST", apiSearchBlackList, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	// Log headers
	fmt.Println("Request Headers:")
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	// Increase timeout for larger requests
	client := &http.Client{Timeout: time.Second * 30}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Capture non-OK responses and log the response body
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body) // Capture the response body even in error
		return nil, fmt.Errorf("received non-OK response: %s, body: %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func GetMemberFaceImage(memberID string, token string) ([]models.MemberFaceImageData, error) {
	url := apiSearchImageID + memberID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Authorization", "bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	type findFaceImageByUserIDResponse struct {
		Code    int                          `json:"code"`
		Data    []models.MemberFaceImageData `json:"data"`
		Message string                       `json:"message"`
	}

	// Marshal respone
	var respone findFaceImageByUserIDResponse
	if err := json.Unmarshal(body, &respone); err != nil {
		return nil, fmt.Errorf("failed to unmarshal find face image by user id response: %w", err)
	}

	return respone.Data, nil
}

func GetUserDetail(memberID string, token string) (userDetail models.UserDetail, err error) {
	url := fmt.Sprintf("%s%s/detail", apiUserDetailURL, memberID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return userDetail, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("Authorization", "bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return userDetail, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return userDetail, fmt.Errorf("received non-OK response: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return userDetail, fmt.Errorf("failed to read response body: %w", err)
	}

	// Marshal respone
	type GetUserDetailResponse struct {
		Code    int               `json:"code"`
		Data    models.UserDetail `json:"data"`
		Message string            `json:"message"`
	}

	var respone GetUserDetailResponse
	if err := json.Unmarshal(body, &respone); err != nil {
		return userDetail, fmt.Errorf("failed to unmarshal get user detail response: %w", err)
	}

	return respone.Data, nil
}
