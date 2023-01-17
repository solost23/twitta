package services

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"twitta/forms"
	"twitta/global"
	"twitta/pkg/constants"
	"twitta/pkg/models"
	"twitta/pkg/utils"
)

func (*Service) ChatList(c *gin.Context, id string) ([]*forms.ChatListResponse, error) {
	db := global.DB
	user := utils.GetUser(c)

	// 直接查询所有记录，并返回
	logPrivateLatters := make([]*models.LogPrivateLatter, 0)
	find := bson.M{
		"type":      3,
		"user_id":   bson.M{"$in": []string{user.ID, id}},
		"target_id": bson.M{"$in": []string{user.ID, id}},
	}
	err := models.NewLogPrivateLatter().Find(c, db, constants.Mongo, find, &logPrivateLatters)
	if err != nil {
		return nil, err
	}
	chatListResponse := make([]*forms.ChatListResponse, 0, len(logPrivateLatters))
	for _, logPrivateLatter := range logPrivateLatters {
		chatListResponse = append(chatListResponse, &forms.ChatListResponse{
			UserId:    logPrivateLatter.UserId,
			Msg:       logPrivateLatter.Content,
			CreatedAt: logPrivateLatter.CreatedAt.Format(constants.TimeFormat),
		})
	}
	return chatListResponse, nil
}
