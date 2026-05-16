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

func (*Service) FriendList(c *gin.Context) ([]*forms.FriendListResponse, error) {
	user := utils.GetUser(c)

	db := global.DB
	friends, err := dao.GWhereFind[*dao.Friend](c, db, bson.M{"user_id": user.ID.String()})
	if err != nil {
		return nil, err
	}
	oids := make([]primitive.ObjectID, 0, len(friends))
	for _, f := range friends {
		oid, err := utils.ParseObjectID(f.FriendId)
		if err != nil {
			continue
		}
		oids = append(oids, oid)
	}
	users, err := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": oids}})
	if err != nil {
		return nil, err
	}
	result := make([]*forms.FriendListResponse, 0, len(users))
	for _, u := range users {
		result = append(result, &forms.FriendListResponse{
			UserId:    u.ID.String(),
			Username:  u.Username,
			Nickname:  u.Nickname,
			Avatar:    utils.FulfillImageOSSPrefix(u.Avatar),
			Introduce: u.Introduce,
		})
	}
	return result, nil
}

func (*Service) FriendApplicationList(c *gin.Context) ([]*forms.FriendApplicationListResponse, error) {
	user := utils.GetUser(c)

	db := global.DB
	logPrivateLatters, err := dao.GWhereFind[*dao.LogPrivateLatter](c, db,
		bson.M{"target_id": user.ID.String(), "type": bson.M{"$in": []uint{dao.LogPrivateLatterTypeAcceptOrReject, dao.LogPrivateLatterTypeReject, dao.LogPrivateLatterTypeAccept}}})
	if err != nil {
		return nil, err
	}
	oids := make([]primitive.ObjectID, 0, len(logPrivateLatters))
	for _, l := range logPrivateLatters {
		oid, err := utils.ParseObjectID(l.UserId)
		if err != nil {
			continue
		}
		oids = append(oids, oid)
	}
	users, err := dao.GWhereFind[*dao.User](c, global.DB, bson.M{"_id": bson.M{"$in": oids}})
	if err != nil {
		return nil, err
	}
	userIdToInfoMaps := make(map[string]struct {
		Username string
		Avatar   string
	}, len(users))
	for _, user := range users {
		userIdToInfoMaps[user.ID.String()] = struct {
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
			CreatedAt: logPrivateLatter.CreatedAt.Format(time.DateTime),
		})
	}
	return result, nil
}

func (*Service) FriendApplicationSend(c *gin.Context, params *forms.FriendApplicationSendForm) error {
	user := utils.GetUser(c)

	if user.ID.String() == params.UserId {
		return errors.New(fmt.Sprintf("不能对自己发申请"))
	}
	// 发送申请有限制
	// 如果此人已经是朋友，那么所发内容都是私信
	// 如果此人不是朋友，那么所发内容为申请信息
	db := global.DB
	_, err := dao.GWhereFirst[*dao.Friend](c, db, bson.M{"user_id": user.ID.String(), "friend_id": params.UserId})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	fmt.Println(user.ID, params.UserId)
	var msgType uint = 0
	if err == nil {
		// 已经是朋友-此为朋友私信
		msgType = dao.LogPrivateLatterTypePrivateLatter
	} else {
		// 不是朋友，检查是否已有待处理的申请
		_, err = dao.GWhereFirst[*dao.LogPrivateLatter](c, db, bson.M{"user_id": user.ID.String(), "target_id": params.UserId, "type": dao.LogPrivateLatterTypeAcceptOrReject})
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			return err
		}
		if err == nil {
			return errors.New("已发送过好友申请，请等待对方处理")
		}
	}
	data := &dao.LogPrivateLatter{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    user.ID.String(),
		TargetId:  params.UserId,
		Content:   params.Content,
		Type:      msgType,
	}
	err = dao.GInsertOne[*dao.LogPrivateLatter](c, db, data)
	if err != nil {
		return err
	}
	return nil
}

func (*Service) FriendApplicationAccept(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	if user.ID.String() == id {
		return errors.New(fmt.Sprintf("自己和自己不可能成为朋友"))
	}
	// 查询此人是否已经是我的朋友，如果不是，添加到朋友表，否则返回错误
	db := global.DB
	_, err := dao.GWhereFirst[*dao.Friend](c, db, bson.M{"user_id": user.ID.String(), "friend_id": id})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil {
		return errors.New(fmt.Sprintf("此人已经是您的朋友，不可重复通过申请"))
	}
	datas := []*dao.Friend{
		{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserId:    user.ID.String(),
			FriendId:  id,
		},
		{
			ID:        primitive.NewObjectID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			UserId:    id,
			FriendId:  user.ID.String(),
		}}
	err = dao.GInsertMany[*dao.Friend](c, db, datas)
	if err != nil {
		return err
	}
	// 修改私信表记录状态
	_, err = dao.GWhereUpdate[*dao.LogPrivateLatter](c, db, bson.M{"user_id": id, "target_id": user.ID.String(), "type": 0}, bson.M{"$set": bson.M{"type": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) FriendApplicationReject(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	// 直接查找到此条私信，然后状态修改为拒绝
	db := global.DB
	_, err := dao.GWhereFirst[*dao.LogPrivateLatter](c, db, bson.M{"user_id": id, "target_id": user.ID.String(), "type": 0})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New(fmt.Sprintf("此私信记录不存在"))
	}
	_, err = dao.GWhereUpdate[*dao.LogPrivateLatter](c, db, bson.M{"user_id": id, "target_id": user.ID.String(), "type": 0}, bson.M{"$set": bson.M{"type": 2}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) FriendDelete(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	if user.ID.String() == id {
		return errors.New(fmt.Sprintf("不能删除自己"))
	}
	db := global.DB
	_, err := dao.GWhereDelete[*dao.Friend](c, db, bson.M{"user_id": user.ID.String(), "friend_id": id})
	if err != nil {
		return err
	}
	_, err = dao.GWhereDelete[*dao.Friend](c, db, bson.M{"user_id": id, "friend_id": user.ID.String()})
	if err != nil {
		return err
	}
	// 删除此朋友的所有申请记录以及聊天内容
	_, err = dao.GWhereDelete[*dao.LogPrivateLatter](c, db, bson.M{"user_id": id, "target_id": user.ID.String()})
	if err != nil {
		return err
	}
	_, err = dao.GWhereDelete[*dao.LogPrivateLatter](c, db, bson.M{"target_id": user.ID.String(), "user_id": id})
	if err != nil {
		return err
	}
	return nil
}
