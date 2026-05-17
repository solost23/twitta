package services

import (
	"context"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/solost23/protopb/gen/go/common"
	"github.com/solost23/protopb/gen/go/oss"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"twitta/forms"
	"twitta/global"
	"twitta/pkg/dao"
	"twitta/pkg/utils"
)

// parseObjectIDs converts a slice of ObjectID("...") strings to []primitive.ObjectID,
// silently skipping any that fail to parse.
func parseObjectIDs(ids []string) []primitive.ObjectID {
	oids := make([]primitive.ObjectID, 0, len(ids))
	for _, id := range ids {
		if oid, err := utils.ParseObjectID(id); err == nil {
			oids = append(oids, oid)
		}
	}
	return oids
}

// collectStrings extracts a string field from each element of a slice.
func collectStrings[T any](items []T, fn func(T) string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		out = append(out, fn(item))
	}
	return out
}

type userInfo struct{ Username, Avatar string }

// buildUserMap builds a map from ObjectID("...") string → userInfo.
func buildUserMap(users []*dao.User) map[string]userInfo {
	m := make(map[string]userInfo, len(users))
	for _, u := range users {
		m[u.ID.String()] = userInfo{Username: u.Username, Avatar: u.Avatar}
	}
	return m
}

// fillOriginTweets populates OriginTweet for any retweet records in the list.
func fillOriginTweets(c *gin.Context, db *mongo.Database, records []*forms.Tweet) {
	originIds := make([]string, 0)
	for _, r := range records {
		if r.RetweetOf != "" {
			originIds = append(originIds, r.RetweetOf)
		}
	}
	if len(originIds) == 0 {
		return
	}
	originOids := parseObjectIDs(originIds)
	origins, err := dao.GWhereFind[*dao.Tweet](c, db, bson.M{"_id": bson.M{"$in": originOids}})
	if err != nil || len(origins) == 0 {
		return
	}
	authorOids := parseObjectIDs(collectStrings(origins, func(t *dao.Tweet) string { return t.UserID }))
	authors, _ := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": authorOids}})
	authorMap := buildUserMap(authors)

	originMap := make(map[string]*dao.Tweet, len(origins))
	for _, o := range origins {
		originMap[o.ID.String()] = o
	}
	for _, r := range records {
		if r.RetweetOf == "" {
			continue
		}
		o, ok := originMap[r.RetweetOf]
		if !ok {
			continue
		}
		u := authorMap[o.UserID]
		r.OriginTweet = tweetToForm(o, u.Username, u.Avatar)
	}
}


func UploadImg(userId uint, folderName string, postFilename string, file *multipart.FileHeader) (string, error) {
	fileHandle, err := file.Open()
	if err != nil {
		return "", err
	}

	defer func() { _ = fileHandle.Close() }()

	b, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		return "", err
	}

	return uploadImgOrVidBytes(userId, folderName, postFilename, b, "image")
}

func UploadVid(userId uint, folderName string, postFilename string, file *multipart.FileHeader) (string, error) {
	fileHandle, err := file.Open()
	if err != nil {
		return "", err
	}

	defer func() { _ = fileHandle.Close() }()

	b, err := ioutil.ReadAll(fileHandle)
	if err != nil {
		return "", err
	}

	return uploadImgOrVidBytes(userId, folderName, postFilename, b, "video")
}

func uploadImgOrVidBytes(userId uint, folderName string, postFileName string, fileBytes []byte, fileType string) (string, error) {
	if len(fileBytes) == 0 {
		return "", fmt.Errorf("upload image or video file is empty")
	}
	mime := strings.Split(mimetype.Detect(fileBytes).String(), " ")[0]
	if !strings.HasPrefix(mime, fileType) {
		return "", fmt.Errorf("invalid mime type: %s", mime)
	}

	filename := utils.NewMd5(
		time.Now().Format(time.DateOnly)+
			fmt.Sprintf("%d", userId)+
			utils.NewMd5(string(fileBytes))+
			postFileName) + path.Ext(postFileName)
	url, err := upload(userId, folderName, filename, fileBytes)
	if err != nil {
		return "", err
	}

	return url, nil
}

func upload(userId uint, folder, filename string, fileBytes []byte) (url string, err error) {
	uploadResponse, err := global.OssSrvClient.Upload(context.Background(), &oss.UploadRequest{
		Header: &common.RequestHeader{
			Requester: strconv.Itoa(int(userId)),
			TraceId:   10000,
		},
		Folder:     folder,
		Key:        filename,
		Data:       fileBytes,
		UploadType: "static",
	})
	if err != nil {
		return "", err
	}
	return uploadResponse.Url, nil
}
