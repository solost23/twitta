package services

import (
	"Twitta/forms"
	"Twitta/global"
	"Twitta/pkg/constants"
	"Twitta/pkg/models"
	"Twitta/pkg/utils"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

func (*Service) TweetSend(c *gin.Context, params *forms.TweetCreateForm) error {
	db := global.DB
	user := utils.GetUser(c)
	data := &models.Tweet{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:      utils.UUID(),
		UserID:  user.ID,
		Title:   params.Title,
		Content: params.Content,
	}
	_, err := models.NewTweet().InsertOne(c, db, constants.Mongo, data)
	if err != nil {
		return err
	}
	return nil
}

func (*Service) TweetDelete(c *gin.Context, id string) error {
	db := global.DB
	user := utils.GetUser(c)

	tweet := &models.Tweet{}
	err := models.NewTweet().FindOne(c, db, constants.Mongo, bson.M{"_id": id}, tweet)
	if err != nil {
		return err
	}
	if user.ID != tweet.UserID {
		return errors.New("本推文所属用户不是您，无权删除")
	}
	_, err = models.NewTweet().Delete(c, db, constants.Mongo, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) TweetList(c *gin.Context) ([]*forms.TweetListResponse, error) {
	db := global.DB

	tweets := make([]*models.Tweet, 0)
	err := models.NewTweet().Find(c, db, constants.Mongo, bson.M{}, &tweets)
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(tweets))
	for _, tweet := range tweets {
		userIds = append(userIds, tweet.UserID)
	}
	users := make([]*models.User, 0)
	err = models.NewUser().Find(c, db, constants.Mongo, bson.M{"_id": bson.M{"$in": userIds}}, &users)
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
	// 封装数据返回
	tweetListResponse := make([]*forms.TweetListResponse, 0, len(tweets))
	for _, tweet := range tweets {
		tweetListResponse = append(tweetListResponse, &forms.TweetListResponse{
			UserId:       tweet.UserID,
			Username:     userIdToInfoMaps[tweet.UserID].Username,
			Avatar:       userIdToInfoMaps[tweet.UserID].Avatar,
			TweetId:      tweet.ID,
			Title:        tweet.Title,
			Content:      tweet.Content,
			TweetTime:    tweet.CreatedAt.Format(constants.TimeFormat),
			ThumbCount:   tweet.ThumbCount,
			CommentCount: tweet.CommentCount,
		})
	}
	return tweetListResponse, nil
}

func (*Service) TweetOwnList(c *gin.Context) ([]*forms.TweetListResponse, error) {
	db := global.DB
	user := utils.GetUser(c)

	tweets := make([]*models.Tweet, 0)
	err := models.NewTweet().Find(c, db, constants.Mongo, bson.M{"user_id": user.ID}, &tweets)
	if err != nil {
		return nil, err
	}
	tweetOwnList := make([]*forms.TweetListResponse, 0, len(tweets))
	for _, tweet := range tweets {
		tweetOwnList = append(tweetOwnList, &forms.TweetListResponse{
			UserId:       user.ID,
			Username:     user.Username,
			Avatar:       user.Avatar,
			TweetId:      tweet.ID,
			Title:        tweet.Title,
			Content:      tweet.Content,
			TweetTime:    tweet.CreatedAt.Format(constants.TimeFormat),
			ThumbCount:   tweet.ThumbCount,
			CommentCount: tweet.CommentCount,
		})
	}
	return tweetOwnList, nil
}

func (*Service) TweetFavoriteList(c *gin.Context) ([]*forms.TweetListResponse, error) {
	db := global.DB
	user := utils.GetUser(c)

	favorites := make([]*models.Favorite, 0)
	err := models.NewFavorite().Find(c, db, constants.Mongo, bson.M{"user_id": user.ID}, &favorites)
	if err != nil {
		return nil, err
	}
	tweetIds := make([]string, 0, len(favorites))
	for _, favorite := range favorites {
		tweetIds = append(tweetIds, favorite.TweetId)
	}
	tweets := make([]*models.Tweet, 0)
	err = models.NewTweet().Find(c, db, constants.Mongo, bson.M{"_id": bson.M{"$in": tweetIds}}, &tweets)
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(tweets))
	for _, tweet := range tweets {
		userIds = append(userIds, tweet.UserID)
	}
	users := make([]*models.User, 0)
	err = models.NewUser().Find(c, db, constants.Mongo, bson.M{"_id": bson.M{"$in": userIds}}, &users)
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
	// 封装数据返回
	tweetFavoriteResponse := make([]*forms.TweetListResponse, 0, len(tweets))
	for _, tweet := range tweets {
		tweetFavoriteResponse = append(tweetFavoriteResponse, &forms.TweetListResponse{
			UserId:       tweet.UserID,
			Username:     userIdToInfoMaps[tweet.UserID].Username,
			Avatar:       userIdToInfoMaps[tweet.UserID].Avatar,
			TweetId:      tweet.ID,
			Title:        tweet.Title,
			Content:      tweet.Content,
			TweetTime:    tweet.CreatedAt.Format(constants.TimeFormat),
			ThumbCount:   tweet.ThumbCount,
			CommentCount: tweet.CommentCount,
		})
	}
	return tweetFavoriteResponse, nil
}

func (*Service) TweetFavorite(c *gin.Context, params *forms.TweetFavoriteForm) error {
	db := global.DB
	user := utils.GetUser(c)

	// 查询此用户有无收藏此文章
	favorite := &models.Favorite{}
	err := models.NewFavorite().FindOne(c, db, constants.Mongo, bson.M{}, favorite)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil {
		return errors.New("您已收藏此推文")
	}
	data := &models.Favorite{
		ID: utils.UUID(),
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		UserId:  user.ID,
		TweetId: params.Id,
	}
	_, err = models.NewFavorite().InsertOne(c, db, constants.Mongo, data)
	if err != nil {
		return err
	}
	return nil
}

func (*Service) TweetFavoriteDelete(c *gin.Context, id string) error {
	db := global.DB
	user := utils.GetUser(c)

	_, err := models.NewFavorite().Delete(c, db, constants.Mongo, bson.M{"user_id": user.ID, "_id": id})
	if err != nil {
		return err
	}
	return nil
}
