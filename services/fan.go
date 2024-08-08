package services

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"twitta/global"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"twitta/forms"
	"twitta/pkg/dao"
	"twitta/pkg/utils"
)

func (*Service) FanList(c *gin.Context) ([]*forms.FansAndWhatResponse, error) {
	user := utils.GetUser(c)

	db := global.DB
	fans, err := dao.GWhereFind[*dao.Fan](c, db, bson.M{"target_id": user.ID})
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(fans))
	for i := 0; i != len(fans); i++ {
		userIds = append(userIds, fans[i].UserId)
	}
	users, err := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": userIds}})
	if err != nil {
		return nil, err
	}
	fansResponse := make([]*forms.FansAndWhatResponse, 0, len(users))
	for i := 0; i != len(users); i++ {
		fansResponse = append(fansResponse, &forms.FansAndWhatResponse{
			UserId:    users[i].ID.String(),
			Avatar:    utils.FulfillImageOSSPrefix(users[i].Avatar),
			Introduce: users[i].Introduce,
		})
	}
	return fansResponse, nil
}

func (*Service) WhatList(c *gin.Context) ([]*forms.FansAndWhatResponse, error) {
	user := utils.GetUser(c)

	db := global.DB
	fans, err := dao.GWhereFind[*dao.Fan](c, db, bson.M{"user_id": user.ID})
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(fans))
	for i := 0; i != len(userIds); i++ {
		userIds = append(userIds, fans[i].UserId)
	}

	users, err := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": userIds}})
	if err != nil {
		return nil, err
	}
	whatsResponse := make([]*forms.FansAndWhatResponse, 0, len(users))
	for i := 0; i != len(users); i++ {
		whatsResponse = append(whatsResponse, &forms.FansAndWhatResponse{
			UserId:    users[i].ID.String(),
			Avatar:    utils.FulfillImageOSSPrefix(users[i].Avatar),
			Introduce: users[i].Introduce,
		})
	}
	return whatsResponse, nil
}

func (*Service) WhatUser(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	// 不能关注自己
	if user.ID.String() == id {
		return errors.New(fmt.Sprintf("不能关注自己"))
	}
	// 判断，如果已关注，那么直接提示不可重复关注
	db := global.DB
	_, err := dao.GWhereFirst[*dao.Fan](c, db, bson.M{"user_id": user.ID, "target_id": id})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil {
		return errors.New("已关注此人，不可重复关注")
	}
	// 关注
	data := &dao.Fan{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    user.ID.String(),
		TargetId:  id,
	}
	err = dao.GInsertOne[*dao.Fan](c, db, data)
	if err != nil {
		return err
	}
	// 将目标用户的粉丝数 +1, 源用户的关注数量 +1
	_, err = dao.GWhereUpdate[*dao.User](c, db, bson.M{"_id": user.ID}, bson.M{"$inc": bson.M{"wechat_count": 1}})
	if err != nil {
		return err
	}
	_, err = dao.GWhereUpdate[*dao.User](c, db, bson.M{"_id": id}, bson.M{"$inc": bson.M{"fans_count": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) WhatUserDelete(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	db := global.DB
	_, err := dao.GWhereDelete[*dao.Fan](c, db, bson.M{"user_id": user.ID, "target_id": id})
	if err != nil {
		return err
	}
	// 将目标用户的粉丝数量 -1, 源用户的关注数量 -1
	_, err = dao.GWhereUpdate[*dao.User](c, db, bson.M{"_id": user.ID}, bson.M{"$inc": bson.M{"wechat_count": -1}})
	if err != nil {
		return err
	}
	_, err = dao.GWhereUpdate[*dao.User](c, db, bson.M{"_id": id}, bson.M{"$inc": bson.M{"fans_count": -1}})
	if err != nil {
		return err
	}
	return nil
}
