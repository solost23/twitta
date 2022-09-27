package services

import (
	"Twitta/global"
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"path"
	"time"

	"github.com/solost23/tools/minio_storage"

	"Twitta/pkg/constants"
	"Twitta/pkg/models"
	"Twitta/pkg/utils"
)

func UploadImg(user *models.User, folderName string, file *multipart.FileHeader) (string, error) {
	return uploadImage(user, folderName, file, "image")
}

func UploadVid(user *models.User, folderName string, file *multipart.FileHeader) (string, error) {
	return uploadImage(user, folderName, file, "video")
}

func uploadImage(user *models.User, folderName string, file *multipart.FileHeader, uploadType string) (string, error) {
	fileHandle, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = fileHandle.Close() }()

	fileByte, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		return "", err
	}

	postFileName := file.Filename
	fileName := utils.GetMd5Hash(
		time.Now().Format(constants.TimeFormat)+
			user.ID+
			utils.GetMd5Hash(string(fileByte))+
			postFileName) + path.Ext(postFileName)
	url, err := upload(folderName, fileName, fileHandle, uploadType, file.Size)
	if err != nil {
		return "", err
	}
	return url, nil
}

func upload(folderName string, fileName string, fileHandle multipart.File, uploadType string, fileSize int64) (result string, err error) {
	client, err := minio_storage.NewMinio(&minio_storage.Config{
		EndPoint:        global.ServerConfig.MinioConfig.EndPoint,
		AccessKeyID:     global.ServerConfig.MinioConfig.AccessKeyId,
		SecretAccessKey: global.ServerConfig.MinioConfig.SecretAccessKey,
		UserSSL:         global.ServerConfig.MinioConfig.UserSsl,
	})
	if err != nil {
		return "", err
	}
	ctx := context.Background()
	if err = minio_storage.CreateBucket(ctx, client, folderName); err != nil {
		return "", err
	}
	// 设置链接可永久下载
	policy := `
{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS": 
          ["*"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource": 
          ["arn:aws:s3:::%s"]},{"Effect":"Allow","Principal":{"AWS":["*"]},"Action": 
          ["s3:GetObject"],"Resource":["arn:aws:s3:::%s/*"]}]}
`
	if err = client.SetBucketPolicy(folderName, fmt.Sprintf(policy, folderName, folderName)); err != nil {
		return "", err
	}
	if err = minio_storage.StreamUpload(ctx, client, folderName, fileName, fileHandle, fileSize, fmt.Sprintf("Application/%s", uploadType)); err != nil {
		return "", err
	}
	requestParams := make(url.Values)
	fileUrl, err := minio_storage.GetFileUrl(ctx, client, folderName, fileName, 168*time.Hour, requestParams)
	if err != nil {
		return "", err
	}
	return fileUrl, nil
}
