package services

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"twitta/forms"
	"twitta/global"
	"twitta/pkg/constants"
	"twitta/pkg/models"
	"twitta/pkg/utils"
)

func (*Service) FanList(c *gin.Context) ([]*forms.FansAndWhatResponse, error) {
	db := global.DB
	user := utils.GetUser(c)

	fans := make([]*models.Fan, 0)
	err := models.NewFan().Find(c, db, constants.Mongo, bson.M{"target_id": user.ID}, &fans)
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(fans))
	for _, fan := range fans {
		userIds = append(userIds, fan.UserId)
	}
	users := make([]*models.User, 0)
	err = models.NewUser().Find(c, db, constants.Mongo, bson.M{"_id": bson.M{"$in": userIds}}, &users)
	if err != nil {
		return nil, err
	}
	fansResponse := make([]*forms.FansAndWhatResponse, 0, len(users))
	for _, user := range users {
		fansResponse = append(fansResponse, &forms.FansAndWhatResponse{
			UserId:    user.ID,
			Avatar:    user.Avatar,
			Introduce: user.Introduce,
		})
	}
	return fansResponse, nil
}

func (*Service) WhatList(c *gin.Context) ([]*forms.FansAndWhatResponse, error) {
	db := global.DB
	user := utils.GetUser(c)

	whats := make([]*models.Fan, 0)
	err := models.NewFan().Find(c, db, constants.Mongo, bson.M{"user_id": user.ID}, &whats)
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(whats))
	for _, what := range whats {
		userIds = append(userIds, what.UserId)
	}
	users := make([]*models.User, 0)
	err = models.NewUser().Find(c, db, constants.Mongo, bson.M{"_id": bson.M{"$in": userIds}}, &users)
	if err != nil {
		return nil, err
	}
	whatsResponse := make([]*forms.FansAndWhatResponse, 0, len(users))
	for _, user := range users {
		whatsResponse = append(whatsResponse, &forms.FansAndWhatResponse{
			UserId:    user.ID,
			Avatar:    user.Avatar,
			Introduce: user.Introduce,
		})
	}
	return whatsResponse, nil
}

func (*Service) WhatUser(c *gin.Context, id string) error {
	db := global.DB
	user := utils.GetUser(c)

	// 不能关注自己
	if user.ID == id {
		return errors.New(fmt.Sprintf("不能关注自己"))
	}
	// 判断，如果已关注，那么直接提示不可重复关注
	what := &models.Fan{}
	err := models.NewFan().FindOne(c, db, constants.Mongo, bson.M{"user_id": user.ID, "target_id": id}, &what)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil {
		return errors.New("已关注此人，不可重复关注")
	}
	// 关注
	data := &models.Fan{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:       utils.UUID(),
		UserId:   user.ID,
		TargetId: id,
	}
	_, err = models.NewFan().InsertOne(c, db, constants.Mongo, &data)
	if err != nil {
		return err
	}
	// 将目标用户的粉丝数 +1, 源用户的关注数量 +1
	_, err = models.NewUser().Update(c, db, constants.Mongo, bson.M{"_id": user.ID}, bson.M{"$inc": bson.M{"wechat_count": 1}})
	if err != nil {
		return err
	}
	_, err = models.NewUser().Update(c, db, constants.Mongo, bson.M{"_id": id}, bson.M{"$inc": bson.M{"fans_count": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) WhatUserDelete(c *gin.Context, id string) error {
	db := global.DB
	user := utils.GetUser(c)

	_, err := models.NewFan().Delete(c, db, constants.Mongo, bson.M{"user_id": user.ID, "target_id": id})
	if err != nil {
		return err
	}
	// 将目标用户的粉丝数量 -1, 源用户的关注数量 -1
	_, err = models.NewUser().Update(c, db, constants.Mongo, bson.M{"_id": user.ID}, bson.M{"$inc": bson.M{"wechat_count": -1}})
	if err != nil {
		return err
	}
	_, err = models.NewUser().Update(c, db, constants.Mongo, bson.M{"_id": id}, bson.M{"$inc": bson.M{"fans_count": -1}})
	if err != nil {
		return err
	}
	return nil
}
