package services

import (
	"encoding/json"
	"errors"
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

func (*Service) CommentList(c *gin.Context, id string, params *forms.CommentInsertForm) (*forms.CommentList, error) {
	db := global.DB
	comments, total, pages, err := dao.GPaginatorOrder[*dao.Comment](c, db, &dao.ListPageInput{
		Page: params.Page,
		Size: params.Size,
	}, bson.M{"created_at": 1}, bson.M{"tweet_id": id, "type": 1})
	userIds := make([]string, 0, len(comments))
	for i := 0; i < cap(userIds); i++ {
		userIds = append(userIds, comments[i].UserId)
	}

	users, err := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": userIds}})
	if err != nil {
		return nil, err
	}

	userIdToInfoMaps := make(map[string]struct {
		Username  *string
		Avatar    *string
		Introduce *string
	}, len(users))
	for i := 0; i < len(users); i++ {
		userIdToInfoMaps[users[i].ID.String()] = struct {
			Username  *string
			Avatar    *string
			Introduce *string
		}{Username: &users[i].Username, Avatar: &users[i].Avatar, Introduce: &users[i].Introduce}
	}

	records := make([]*forms.Comment, 0, len(comments))
	for i := 0; i < cap(records); i++ {
		idStr := comments[i].ID.String()
		records = append(records, &forms.Comment{
			UserId:    &comments[i].UserId,
			PID:       &comments[i].ParentId,
			Username:  userIdToInfoMaps[comments[i].UserId].Username,
			Avatar:    userIdToInfoMaps[comments[i].UserId].Avatar,
			Introduce: userIdToInfoMaps[comments[i].UserId].Introduce,
			Id:        &idStr,
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
	user := utils.GetUser(c)

	// 查看有无点赞记录，如果无，那么创建, 文章下的点赞数 +1
	db := global.DB
	comment, err := dao.GWhereFirst[*dao.Comment](c, db, bson.M{"type": 0, "user_id": user.ID, "tweet_id": id})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if comment != nil {
		return errors.New("已经点赞过此推文")
	}
	data := &dao.Comment{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    user.ID.String(),
		TweetId:   id,
		Type:      dao.CommentTypeThumb,
	}
	err = dao.GInsertOne[*dao.Comment](c, db, data)
	if err != nil {
		return err
	}
	_, err = dao.GWhereUpdate[*dao.Tweet](c, db, bson.M{"_id": id}, bson.M{"$inc": bson.M{"thumb_count": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) CommentThumbDelete(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	// 查询是否存在此👍
	db := global.DB
	_, err := dao.GWhereFirst[*dao.Comment](c, db, bson.M{"type": 0, "user_id": user.ID, "tweet_id": id})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("您已经取消赞，不可重复取消")
	}
	_, err = dao.GWhereDelete[*dao.Comment](c, db, bson.M{"type": 0, "user_id": user.ID, "tweet_id": id})
	if err != nil {
		return err
	}
	_, err = dao.GWhereUpdate[*dao.Tweet](c, db, bson.M{"_id": id}, bson.M{"$inc": bson.M{"thumb_count": -1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) CommentInsert(c *gin.Context, id string, params *forms.CommentInsertForm) error {
	user := utils.GetUser(c)

	// 直接插入评论记录
	data := &dao.Comment{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    user.ID.String(),
		TweetId:   id,
		Content:   *params.Content,
		ParentId:  *params.ParentId,
		Type:      dao.CommentTypeComment,
	}

	db := global.DB
	if err := dao.GInsertOne[*dao.Comment](c, db, data); err != nil {
		return err
	}
	_, err := dao.GWhereUpdate[*dao.Tweet](c, db, bson.M{"_id": id}, bson.M{"$inc": bson.M{"comment_count": 1}})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) CommentDelete(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	// 查询是否存在此评论
	db := global.DB
	comment, err := dao.GWhereFirst[*dao.Comment](c, db, bson.M{"_id": id, "type": 1, "user_id": user.ID})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
		return errors.New("您已经删除此评论，不可重复删除")
	}
	_, err = dao.GWhereDelete[*dao.Comment](c, db, bson.M{"_id": id, "type": 1, "user_id": user.ID})
	if err != nil {
		return err
	}
	_, err = dao.GWhereUpdate[*dao.Tweet](c, db, bson.M{"_id": comment.TweetId}, bson.M{"$inc": bson.M{"comment_count": -1}})
	if err != nil {
		return err
	}
	return nil
}
