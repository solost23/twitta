package services

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"twitta/forms"
	"twitta/pkg/constants"
	"twitta/pkg/models"
	"twitta/pkg/utils"
)

func (*Service) FriendApplicationList(c *gin.Context) ([]*forms.FriendApplicationListResponse, error) {
	db := models.NewDB()
	user := utils.GetUser(c)

	logPrivateLatters, err := models.GWhereFind[models.LogPrivateLatter](c, db.GetCollection(models.NewLogPrivateLatter().TableName()),
		bson.M{"target_id": user.ID, "type": bson.M{"$in": []uint{models.LogPrivateLatterTypeAcceptOrReject, models.LogPrivateLatterTypeReject, models.LogPrivateLatterTypeAccept}}})
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(logPrivateLatters))
	for _, logPrivateLatter := range logPrivateLatters {
		userIds = append(userIds, logPrivateLatter.UserId)
	}
	users, err := models.GWhereFind[models.User](c, db.GetCollection(models.NewUser().TableName()), bson.M{"_id": bson.M{"$in": userIds}})
	if err != nil {
		return nil, err
	}
	userIdToInfoMaps := make(map[string]struct {
		Username string
		Avatar   string
	}, len(users))
	for _, user := range users {
		userIdToInfoMaps[user.ID] = struct {
			Username string
			Avatar   string
		}{Username: user.Username, Avatar: user.Avatar}
	}
	result := make([]*forms.FriendApplicationListResponse, 0, len(logPrivateLatters))
	for _, logPrivateLatter := range logPrivateLatters {
		result = append(result, &forms.FriendApplicationListResponse{
			UserId:    logPrivateLatter.UserId,
			Username:  userIdToInfoMaps[logPrivateLatter.UserId].Username,
			Avatar:    userIdToInfoMaps[logPrivateLatter.UserId].Avatar,
			Content:   logPrivateLatter.Content,
			Type:      logPrivateLatter.Type,
			CreatedAt: logPrivateLatter.CreatedAt.Format(constants.TimeFormat),
		})
	}
	return result, nil
}

func (*Service) FriendApplicationSend(c *gin.Context, params *forms.FriendApplicationSendForm) error {
	db := models.NewDB()
	user := utils.GetUser(c)

	if user.ID == params.UserId {
		return errors.New(fmt.Sprintf("不能对自己发申请"))
	}
	// 发送申请有限制
	// 如果此人已经是朋友，那么所发内容都是私信
	// 如果此人不是朋友，那么所发内容为申请信息
	_, err := models.GWhereFirst[models.Friend](c, db.GetCollection(models.NewFriend().TableName()), bson.M{"user_id": user.ID, "friend_id": params.UserId})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	fmt.Println(user.ID, params.UserId)
	var msgType uint = 0
	if err == nil {
		// 已经是朋友-此为朋友私信
		msgType = models.LogPrivateLatterTypePrivateLatter
	}
	data := &models.LogPrivateLatter{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:       utils.UUID(),
		UserId:   user.ID,
		TargetId: params.UserId,
		Content:  params.Content,
		Type:     msgType,
	}
	_, err = models.GInsertOne[models.LogPrivateLatter](c, db.GetCollection(models.NewLogPrivateLatter().TableName()), &data)
	if err != nil {
		return err
	}
	return nil
}

func (*Service) FriendApplicationAccept(c *gin.Context, id string) error {
	db := models.NewDB()
	user := utils.GetUser(c)

	if user.ID == id {
		return errors.New(fmt.Sprintf("自己和自己不可能成为朋友"))
	}
	// 查询此人是否已经是我的朋友，如果不是，添加到朋友表，否则返回错误
	_, err := models.GWhereFirst[models.Friend](c, db.GetCollection(models.NewFriend().TableName()), bson.M{"user_id": user.ID, "friend_id": id})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil {
		return errors.New(fmt.Sprintf("此人已经是您的朋友，不可重复通过申请"))
	}
	datas := []interface{}{
		&models.Friend{
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ID:       utils.UUID(),
			UserId:   user.ID,
			FriendId: id,
		}, &models.Friend{
			BaseModel: models.BaseModel{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ID:       utils.UUID(),
			UserId:   id,
			FriendId: user.ID,
		}}
	_, err = models.GInsertMany[models.Friend](c, db.GetCollection(models.NewFriend().TableName()), datas)
	if err != nil {
		return err
	}
	// 修改私信表记录状态
	_, err = models.GWhereUpdate[models.LogPrivateLatter](c, db.GetCollection(models.NewLogPrivateLatter().TableName()), bson.M{"user_id": id, "target_id": user.ID, "type": 0}, bson.M{"$set": bson.M{"type": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) FriendApplicationReject(c *gin.Context, id string) error {
	db := models.NewDB()
	user := utils.GetUser(c)

	// 直接查找到此条私信，然后状态修改为拒绝
	_, err := models.GWhereFirst[models.LogPrivateLatter](c, db.GetCollection(models.NewLogPrivateLatter().TableName()), bson.M{"user_id": id, "target_id": user.ID, "type": 0})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New(fmt.Sprintf("此私信记录不存在"))
	}
	_, err = models.GWhereUpdate[models.LogPrivateLatter](c, db.GetCollection(models.NewLogPrivateLatter().TableName()), bson.M{"user_id": id, "target_id": user.ID, "type": 0}, bson.M{"$set": bson.M{"type": 2}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) FriendDelete(c *gin.Context, id string) error {
	db := models.NewDB()
	user := utils.GetUser(c)

	if user.ID == id {
		return errors.New(fmt.Sprintf("不能删除自己"))
	}
	_, err := models.GWhereDelete[models.Friend](c, db.GetCollection(models.NewFriend().TableName()), bson.M{"user_id": user.ID, "friend_id": id})
	if err != nil {
		return err
	}
	_, err = models.GWhereDelete[models.Friend](c, db.GetCollection(models.NewFriend().TableName()), bson.M{"user_id": id, "friend_id": user.ID})
	if err != nil {
		return err
	}
	// 删除此朋友的所有申请记录以及聊天内容
	_, err = models.GWhereDelete[models.LogPrivateLatter](c, db.GetCollection(models.NewLogPrivateLatter().TableName()), bson.M{"user_id": id, "target_id": user.ID})
	if err != nil {
		return err
	}
	_, err = models.GWhereDelete[models.LogPrivateLatter](c, db.GetCollection(models.NewLogPrivateLatter().TableName()), bson.M{"target_id": id, "user_id": id})
	if err != nil {
		return err
	}
	return nil
}
