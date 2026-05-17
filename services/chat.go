package services

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"twitta/forms"
	"twitta/global"
	"twitta/pkg/dao"
	"twitta/pkg/utils"
	"twitta/pkg/ws"

	"github.com/gin-gonic/gin"
)

// RoomID 按字典序拼接两个用户 ID，保证同一对用户的房间 ID 唯一。
// 统一剥掉 ObjectID("...") 包装，使用纯 hex。
func RoomID(a, b string) string {
	a = normalizeID(a)
	b = normalizeID(b)
	if a < b {
		return a + ":" + b
	}
	return b + ":" + a
}

func normalizeID(id string) string {
	if len(id) > 12 && id[:9] == `ObjectID(` {
		id = id[10 : len(id)-2]
	}
	return id
}

// SaveChatMessage 将 WebSocket 消息持久化到 MongoDB，并向收件人推送通知。
func SaveChatMessage(msg *ws.Message) {
	fromID := normalizeID(msg.FromID)
	targetID := normalizeID(targetIDFromRoom(msg.RoomID, msg.FromID))
	data := &dao.LogPrivateLatter{
		ID:        primitive.NewObjectID(),
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.CreatedAt,
		UserId:    fromID,
		TargetId:  targetID,
		Content:   msg.Content,
		Type:      dao.LogPrivateLatterTypePrivateLatter,
		Read:      false,
	}
	_ = dao.GInsertOne(context.Background(), global.DB, data)

	// 向收件人推通知（跨节点有效）
	if ws.DefaultNotify != nil {
		ws.DefaultNotify.Notify(&ws.Notification{
			ToUserID:  targetID,
			FromID:    fromID,
			RoomID:    msg.RoomID,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		})
	}
}

func targetIDFromRoom(roomID, fromID string) string {
	fromID = normalizeID(fromID)
	for i, ch := range roomID {
		if ch == ':' {
			a, b := roomID[:i], roomID[i+1:]
			if a == fromID {
				return b
			}
			return a
		}
	}
	return ""
}

func (s *Service) ChatList(c *gin.Context, id string, params *utils.PageForm) (*forms.ChatList, error) {
	user := utils.GetUser(c)
	myID := normalizeID(user.ID.String())
	targetID := normalizeID(id)

	query := bson.M{
		"type":      dao.LogPrivateLatterTypePrivateLatter,
		"user_id":   bson.M{"$in": []string{myID, targetID}},
		"target_id": bson.M{"$in": []string{myID, targetID}},
	}
	db := global.DB
	logPrivateLatters, total, pages, err := dao.GPaginatorOrder[*dao.LogPrivateLatter](c, db, &dao.ListPageInput{
		Page: params.Page,
		Size: params.Size,
	}, bson.M{"created_at": 1}, query)
	if err != nil {
		return nil, err
	}

	records := make([]*forms.Chat, 0, len(logPrivateLatters))
	for i := 0; i < len(logPrivateLatters); i++ {
		createdAt := logPrivateLatters[i].CreatedAt.Format(time.DateTime)
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

// GetAndMarkOfflineMessages 查询发给 userID 的未读私信，标记为已读后返回。
func GetAndMarkOfflineMessages(userID string) ([]*ws.Notification, error) {
	filter := bson.M{
		"target_id": userID,
		"type":      dao.LogPrivateLatterTypePrivateLatter,
		"read":      false,
	}
	msgs, err := dao.GWhereFind[*dao.LogPrivateLatter](context.Background(), global.DB, filter)
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, nil
	}

	ids := make([]primitive.ObjectID, 0, len(msgs))
	for _, m := range msgs {
		ids = append(ids, m.ID)
	}
	_, err = dao.GWhereUpdate[*dao.LogPrivateLatter](
		context.Background(), global.DB,
		bson.M{"_id": bson.M{"$in": ids}},
		bson.M{"$set": bson.M{"read": true}},
	)
	if err != nil {
		return nil, err
	}

	notifications := make([]*ws.Notification, 0, len(msgs))
	for _, m := range msgs {
		notifications = append(notifications, &ws.Notification{
			ToUserID:  userID,
			FromID:    m.UserId,
			RoomID:    RoomID(m.UserId, userID),
			Content:   m.Content,
			CreatedAt: m.CreatedAt,
		})
	}
	return notifications, nil
}

