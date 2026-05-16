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
	fans, err := dao.GWhereFind[*dao.Fan](c, db, bson.M{"target_id": user.ID.String()})
	if err != nil {
		return nil, err
	}
	oids := make([]primitive.ObjectID, 0, len(fans))
	for _, f := range fans {
		oid, err := utils.ParseObjectID(f.UserId)
		if err != nil {
			continue
		}
		oids = append(oids, oid)
	}
	users, err := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": oids}})
	if err != nil {
		return nil, err
	}
	fansResponse := make([]*forms.FansAndWhatResponse, 0, len(users))
	for _, u := range users {
		fansResponse = append(fansResponse, &forms.FansAndWhatResponse{
			UserId:    u.ID.String(),
			Avatar:    utils.FulfillImageOSSPrefix(u.Avatar),
			Introduce: u.Introduce,
		})
	}
	return fansResponse, nil
}

func (*Service) WhatList(c *gin.Context) ([]*forms.FansAndWhatResponse, error) {
	user := utils.GetUser(c)

	db := global.DB
	fans, err := dao.GWhereFind[*dao.Fan](c, db, bson.M{"user_id": user.ID.String()})
	if err != nil {
		return nil, err
	}
	oids := make([]primitive.ObjectID, 0, len(fans))
	for _, f := range fans {
		oid, err := utils.ParseObjectID(f.TargetId)
		if err != nil {
			continue
		}
		oids = append(oids, oid)
	}
	users, err := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": oids}})
	if err != nil {
		return nil, err
	}
	whatsResponse := make([]*forms.FansAndWhatResponse, 0, len(users))
	for _, u := range users {
		whatsResponse = append(whatsResponse, &forms.FansAndWhatResponse{
			UserId:    u.ID.String(),
			Avatar:    utils.FulfillImageOSSPrefix(u.Avatar),
			Introduce: u.Introduce,
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
	targetOid, err := utils.ParseObjectID(id)
	if err != nil {
		return errors.New("无效用户ID")
	}
	// 判断，如果已关注，那么直接提示不可重复关注
	db := global.DB
	_, err = dao.GWhereFirst[*dao.Fan](c, db, bson.M{"user_id": user.ID.String(), "target_id": id})
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
	_, err = dao.GWhereUpdate[*dao.User](c, db, bson.M{"_id": targetOid}, bson.M{"$inc": bson.M{"fans_count": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) WhatUserDelete(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	targetOid, err := utils.ParseObjectID(id)
	if err != nil {
		return errors.New("无效用户ID")
	}
	db := global.DB
	_, err = dao.GWhereDelete[*dao.Fan](c, db, bson.M{"user_id": user.ID.String(), "target_id": id})
	if err != nil {
		return err
	}
	// 将目标用户的粉丝数量 -1, 源用户的关注数量 -1
	_, err = dao.GWhereUpdate[*dao.User](c, db, bson.M{"_id": user.ID}, bson.M{"$inc": bson.M{"wechat_count": -1}})
	if err != nil {
		return err
	}
	_, err = dao.GWhereUpdate[*dao.User](c, db, bson.M{"_id": targetOid}, bson.M{"$inc": bson.M{"fans_count": -1}})
	if err != nil {
		return err
	}
	return nil
}
