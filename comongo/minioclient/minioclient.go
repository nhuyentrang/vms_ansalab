package minioclient

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var mc *minio.Client = nil
var isConnected bool = false

func Connect(endpoint, accessKeyID, secretAccessKey string, useSSL bool) error {
	// Initialize minio client object.
	fmt.Printf("[MinioClient] endpoint: %s, useSSL: %t\n", endpoint, useSSL)

	var err error = nil
	mc, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
		fmt.Println("[MinioClient] failed to setup client, error: ", err.Error())
		return err
	}

	isConnected = true
	fmt.Println("[MinioClient] setup completed !")
	return nil
}

func CheckObjExist(bucketName, objectName string) (bool, error) {
	if !isConnected || mc == nil {
		return false, errors.New("minio client is not valid")
	}
	_, err := mc.StatObject(context.Background(), bucketName, objectName, minio.StatObjectOptions{})
	return err == nil, err
}

func MakeBucket(bucketName string) (err error) {
	if !isConnected || mc == nil {
		return errors.New("minio client is not valid")
	}
	err = mc.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "asia-southeast1-b", ObjectLocking: false})
	return err
}

func CheckBucketExist(bucketName string) (found bool, err error) {
	if !isConnected || mc == nil {
		return false, errors.New("minio client is not valid")
	}
	found, err = mc.BucketExists(context.Background(), bucketName)
	return found, err
}

func RemoveBucket(bucketName string) (err error) {
	if !isConnected || mc == nil {
		return errors.New("minio client is not valid")
	}
	err = mc.RemoveBucket(context.Background(), bucketName)
	return err
}

func UploadImageBase64(bucketName, objectName, imageBase64 string) (uploadInfo minio.UploadInfo, err error) {
	if !isConnected || mc == nil {
		return minio.UploadInfo{}, errors.New("minio client is not valid")
	}

	// Decode the base64 string into binary data
	imageData, err := base64.StdEncoding.DecodeString(imageBase64)
	if err != nil {
		return minio.UploadInfo{}, err
	}

	// Create a reader from the binary data
	imageReader := strings.NewReader(string(imageData))

	// Upload the image to Minio
	uploadInfo, err = mc.PutObject(
		context.Background(),
		bucketName,
		objectName,
		imageReader,
		int64(len(imageData)), // Specify the size of the data
		minio.PutObjectOptions{ContentType: "image/jpeg"}, // Set the content type as needed
	)

	if err != nil {
		fmt.Println(err)
		return minio.UploadInfo{}, err
	}

	return uploadInfo, nil

	/*
		// How to convert imageBase64 into jpg image file object and upload to mino
		uploadInfo, err = mc.PutObject(context.Background(), "xuka-vms-ai-edgedevice-event", "myobject", ??? , ???.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Successfully uploaded bytes: ", uploadInfo)
		return uploadInfo, nil
	*/
}

func UploadFileBuffer(bucketName, objectName string, bufferReader io.Reader, bufferLen int64) (uploadInfo minio.UploadInfo, err error) {
	if !isConnected || mc == nil {
		return minio.UploadInfo{}, errors.New("minio client is not valid")
	}

	// Upload the image to Minio
	uploadInfo, err = mc.PutObject(
		context.Background(),
		bucketName,
		objectName,
		bufferReader,
		bufferLen, // Specify the size of the data
		minio.PutObjectOptions{ContentType: "image/jpeg"}, // Set the content type as needed
	)

	if err != nil {
		fmt.Println(err)
		return minio.UploadInfo{}, err
	}

	return uploadInfo, nil
}

func GetPresignedURL(bucketName, objectName string) (string, error) {
	if !isConnected || mc == nil {
		return "", errors.New("minio client is not valid")
	}
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	// Set the object name in the request parameters
	reqParams.Set("response-content-disposition", "attachment; filename=\""+objectName+"\"")
	// Generates a presigned url which expires in a day.
	presignedURL, err := mc.PresignedGetObject(context.Background(), bucketName, objectName, time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return presignedURL.String(), nil
}

func RemoveObject(bucketName, objectName string) (err error) {
	if !isConnected || mc == nil {
		return errors.New("minio client is not valid")
	}
	err = mc.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	return err
}

func TestBucketAccess(bucketName string) error {
	// Check minio bucket is existed
	bucketExisted, err := CheckBucketExist(bucketName)
	if err != nil {
		return err
	}
	if !bucketExisted {
		return errors.New("bucket is not existed")
	}

	if bucketExisted {
		// Upload a dummy image to bucket
		testImgB64 := "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"
		// Get the current date and time
		filename := "test_image_" + time.Now().Format("20060102150405") + ".jpg"
		// Upload the image with the generated filename
		_, err := UploadImageBase64(bucketName, filename, testImgB64)
		if err != nil {
			return err
		}

		// Get the presigned URL for the uploaded object
		_, err = GetPresignedURL(bucketName, filename)
		if err != nil {
			return err
		}

		// Remove the test image
		err = RemoveObject(bucketName, filename)
		if err != nil {
			return err
		}
	}
	return nil
}
