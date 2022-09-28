package services

import (
	"Twitta/forms"
	"Twitta/global"
	"Twitta/pkg/constants"
	"Twitta/pkg/models"
	"Twitta/pkg/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func (*Service) CommentList(c *gin.Context, id string) ([]*forms.CommentListResponse, error) {
	db := global.DB

	comments := make([]*models.Comment, 0)
	err := models.NewComment().Find(c, db, constants.Mongo, bson.M{"tweet_id": id, "type": 1}, &comments)
	if err != nil {
		return nil, err
	}
	// 待优化排序-这里直接返回全部数据
	userIds := make([]string, 0, len(comments))
	for _, comment := range comments {
		userIds = append(userIds, comment.UserId)
	}
	users := make([]*models.User, 0)
	err = models.NewUser().Find(c, db, constants.Mongo, bson.M{"_id": bson.M{"$in": userIds}}, &users)
	if err != nil {
		return nil, err
	}
	userIdToInfoMaps := make(map[string]struct {
		Username  string
		Avatar    string
		Introduce string
	}, len(users))
	for _, user := range users {
		userIdToInfoMaps[user.ID] = struct {
			Username  string
			Avatar    string
			Introduce string
		}{Username: user.Username, Avatar: user.Avatar, Introduce: user.Introduce}
	}
	commentListResponse := make([]*forms.CommentListResponse, 0, len(comments))
	for _, comment := range comments {
		commentListResponse = append(commentListResponse, &forms.CommentListResponse{
			UserId:    comment.UserId,
			Username:  userIdToInfoMaps[comment.UserId].Username,
			Avatar:    userIdToInfoMaps[comment.UserId].Avatar,
			Introduce: userIdToInfoMaps[comment.UserId].Introduce,
			Id:        comment.ID,
			Content:   comment.Content,
		})
	}
	return commentListResponse, nil
}

func (*Service) CommentThumb(c *gin.Context, id string) error {
	db := global.DB
	user := utils.GetUser(c)

	// 查看有无点赞记录，如果无，那么创建, 文章下的点赞数 +1
	comment := &models.Comment{}
	err := models.NewComment().FindOne(c, db, constants.Mongo, bson.M{"type": 0, "user_id": user.ID, "tweet_id": id}, &comment)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil {
		return errors.New("已经点赞过此推文")
	}
	data := &models.Comment{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:      utils.UUID(),
		UserId:  user.ID,
		TweetId: id,
		Type:    0,
	}
	_, err = models.NewComment().InsertOne(c, db, constants.Mongo, &data)
	if err != nil {
		return err
	}
	_, err = models.NewTweet().Update(c, db, constants.Mongo, bson.M{"_id": id}, bson.M{"$inc": bson.M{"thumb_count": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) CommentThumbDelete(c *gin.Context, id string) error {
	db := global.DB
	user := utils.GetUser(c)

	// 查询是否存在此👍
	comment := &models.Comment{}
	err := models.NewComment().FindOne(c, db, constants.Mongo, bson.M{"type": 0, "user_id": user.ID, "tweet_id": id}, &comment)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("您已经取消赞，不可重复取消")
	}
	_, err = models.NewComment().Delete(c, db, constants.Mongo, bson.M{"type": 0, "user_id": user.ID, "tweet_id": id})
	if err != nil {
		return err
	}
	_, err = models.NewTweet().Update(c, db, constants.Mongo, bson.M{"_id": id}, bson.M{"$inc": bson.M{"thumb_count": -1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) CommentInsert(c *gin.Context, id string, params *forms.CommentInsertForm) error {
	db := global.DB
	user := utils.GetUser(c)

	// 直接插入评论记录
	data := &models.Comment{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:      utils.UUID(),
		UserId:  user.ID,
		TweetId: id,
		Content: params.Content,
		Parent:  params.ParentId,
		Type:    1,
	}
	_, err := models.NewComment().InsertOne(c, db, constants.Mongo, &data)
	if err != nil {
		return err
	}
	_, err = models.NewTweet().Update(c, db, constants.Mongo, bson.M{"_id": id}, bson.M{"$inc": bson.M{"comment_count": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) CommentDelete(c *gin.Context, id string) error {
	db := global.DB
	user := utils.GetUser(c)

	// 查询是否存在此评论
	comment := &models.Comment{}
	err := models.NewComment().FindOne(c, db, constants.Mongo, bson.M{"_id": id, "type": 1, "user_id": user.ID}, &comment)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("您已经删除此评论，不可重复删除")
	}
	_, err = models.NewComment().Delete(c, db, constants.Mongo, bson.M{"_id": id, "type": 1, "user_id": user.ID})
	if err != nil {
		return err
	}
	_, err = models.NewTweet().Update(c, db, constants.Mongo, bson.M{"_id": comment.TweetId}, bson.M{"$inc": bson.M{"comment_count": -1}})
	if err != nil {
		return err
	}
	return nil
}
