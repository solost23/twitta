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

func (*Service) ChatList(c *gin.Context, id string, params *utils.PageForm) (*forms.ChatList, error) {
	collection := models.GetCollection(global.DB, constants.Mongo, (&models.LogPrivateLatter{}).TableName())
	user := utils.GetUser(c)

	// 直接查询所有记录，并返回
	filter := bson.M{
		"type":      3,
		"user_id":   bson.M{"$in": []string{user.ID, id}},
		"target_id": bson.M{"$in": []string{user.ID, id}},
	}
	logPrivateLatters, total, pages, err := models.GPaginatorOrder[models.LogPrivateLatter](c, collection, &models.ListPageInput{
		Page: params.Page,
		Size: params.Size,
	}, "", filter)
	if err != nil {
		return nil, err
	}

	records := make([]*forms.Chat, 0, len(logPrivateLatters))
	for i := 0; i < cap(records); i++ {
		createdAt := logPrivateLatters[i].CreatedAt.Format(constants.TimeFormat)
		records = append(records, &forms.Chat{
			UserId:    &logPrivateLatters[i].UserId,
			Msg:       &logPrivateLatters[i].Content,
			CreatedAt: &createdAt,
		})
	}
	result := &forms.ChatList{
		Records: records,
		PageList: &utils.PageList{
			Size:    params.Size,
			Pages:   pages,
			Total:   total,
			Current: params.Page,
		},
	}
	return result, nil
}
