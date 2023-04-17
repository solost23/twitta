package services

import (
	"encoding/json"
	"errors"
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

func (*Service) CommentList(c *gin.Context, id string, params *forms.CommentInsertForm) (*forms.CommentList, error) {
	collectionComment := models.GetCollection(global.DB, constants.Mongo, (&models.Comment{}).TableName())
	collectionUser := models.GetCollection(global.DB, constants.Mongo, (&models.User{}).TableName())

	comments, total, pages, err := models.GPaginatorOrder[models.Comment](c, collectionComment, &models.ListPageInput{
		Page: params.Page,
		Size: params.Size,
	}, "created_at ASC", bson.M{"tweet_id": id, "type": 1})
	userIds := make([]string, 0, len(comments))
	for i := 0; i < cap(userIds); i++ {
		userIds = append(userIds, comments[i].UserId)
	}

	users, err := models.GWhereFind[models.User](c, collectionUser, bson.M{"_id": bson.M{"$in": userIds}})
	if err != nil {
		return nil, err
	}

	userIdToInfoMaps := make(map[string]struct {
		Username  *string
		Avatar    *string
		Introduce *string
	}, len(users))
	for i := 0; i < len(users); i++ {
		userIdToInfoMaps[users[i].ID] = struct {
			Username  *string
			Avatar    *string
			Introduce *string
		}{Username: &users[i].Username, Avatar: &users[i].Avatar, Introduce: &users[i].Introduce}
	}

	records := make([]*forms.Comment, 0, len(comments))
	for i := 0; i < cap(records); i++ {
		records = append(records, &forms.Comment{
			UserId:    &comments[i].UserId,
			PID:       &comments[i].Parent,
			Username:  userIdToInfoMaps[comments[i].UserId].Username,
			Avatar:    userIdToInfoMaps[comments[i].UserId].Avatar,
			Introduce: userIdToInfoMaps[comments[i].UserId].Introduce,
			Id:        &comments[i].ID,
			Content:   &comments[i].Content,
			Children:  nil,
		})
	}

	nodeRecords := make([]*forms.Comment, 0)
	arrNode := make([]utils.TreeNode, 0, len(comments))
	for i := 0; i < cap(arrNode); i++ {
		arrNode[i] = records[i]
	}
	rootNodes := utils.BuildTrees(arrNode)
	rootNodesByte, err := json.Marshal(rootNodes)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(rootNodesByte, &nodeRecords)
	if err != nil {
		return nil, err
	}

	result := &forms.CommentList{
		Record: nodeRecords,
		PageList: &utils.PageList{
			Size:    params.Size,
			Pages:   pages,
			Total:   total,
			Current: params.Page,
		},
	}
	return result, nil
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
		Content: *params.Content,
		Parent:  *params.ParentId,
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
