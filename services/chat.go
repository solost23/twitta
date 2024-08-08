package services

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"twitta/forms"
	"twitta/global"
	"twitta/pkg/constants"
	"twitta/pkg/dao"
	"twitta/pkg/utils"
)

func (s *Service) ChatList(c *gin.Context, id string, params *utils.PageForm) (*forms.ChatList, error) {
	user := utils.GetUser(c)

	// 直接查询所有记录，并返回
	query := bson.M{
		"type":      dao.LogPrivateLatterTypePrivateLatter,
		"user_id":   bson.M{"$in": []string{user.ID.String(), id}},
		"target_id": bson.M{"$in": []string{user.ID.String(), id}},
	}
	db := global.DB
	logPrivateLatters, total, pages, err := dao.GPaginatorOrder[*dao.LogPrivateLatter](c, db, &dao.ListPageInput{
		Page: params.Page,
		Size: params.Size,
	}, bson.M{"_id": 1}, query)
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
