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
func RoomID(a, b string) string {
	if a < b {
		return a + ":" + b
	}
	return b + ":" + a
}

// SaveChatMessage 将 WebSocket 消息持久化到 MongoDB。
func SaveChatMessage(msg *ws.Message) {
	data := &dao.LogPrivateLatter{
		ID:        primitive.NewObjectID(),
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.CreatedAt,
		UserId:    msg.FromID,
		TargetId:  targetIDFromRoom(msg.RoomID, msg.FromID),
		Content:   msg.Content,
		Type:      dao.LogPrivateLatterTypePrivateLatter,
	}
	_ = dao.GInsertOne(context.Background(), global.DB, data)
}

func targetIDFromRoom(roomID, fromID string) string {
	// roomID 格式: "uid1:uid2"
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

	query := bson.M{
		"type":      dao.LogPrivateLatterTypePrivateLatter,
		"user_id":   bson.M{"$in": []string{user.ID.String(), id}},
		"target_id": bson.M{"$in": []string{user.ID.String(), id}},
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

