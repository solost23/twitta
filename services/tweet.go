package services

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/solost23/protopb/gen/go/protos/common"
	es_service "github.com/solost23/protopb/gen/go/protos/es"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"twitta/forms"
	"twitta/pkg/constants"
	"twitta/pkg/domain"
	"twitta/pkg/models"
	"twitta/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)

func (*Service) TweetSend(c *gin.Context, params *forms.TweetCreateForm) error {
	db := models.NewDB()
	user := utils.GetUser(c)
	data := &models.Tweet{
		BaseModel: models.BaseModel{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		ID:           utils.UUID(),
		UserID:       user.ID,
		Title:        params.Title,
		Content:      params.Content,
		ThumbCount:   0,
		CommentCount: 0,
	}
	_, err := models.GInsertOne[models.User](c, db.GetCollection(models.NewTweet().TableName()), data)
	if err != nil {
		return err
	}

	// 推文携带创建者信息，方便后续直接搜索展示
	type Document struct {
		*models.Tweet
		Username string `json:"username"`
		Avatar   string `json:"avatar"`
	}
	postDocument := Document{Tweet: data, Username: user.Username, Avatar: user.Avatar}
	documentJson, _ := json.Marshal(postDocument)
	_, err = domain.NewESClient().Create(c, &es_service.CreateRequest{
		Header: &common.RequestHeader{
			Requester:   user.Username,
			OperatorUid: -1,
			TraceId:     -1,
		},
		Index:      constants.ESCINDEXTWEET,
		DocumentId: data.ID,
		Document:   string(documentJson),
	})
	if err != nil {
		zap.S().Errorf(err.Error())
	}
	return nil
}

func (*Service) TweetDelete(c *gin.Context, id string) error {
	db := models.NewDB()
	user := utils.GetUser(c)

	tweet, err := models.GWhereFirst[models.Tweet](c, db.GetCollection(models.NewTweet().TableName()), bson.M{"_id": id})
	if err != nil {
		return err
	}
	if user.ID != tweet.UserID {
		return errors.New("本推文所属用户不是您，无权删除")
	}
	_, err = models.GWhereDelete[models.Tweet](c, db.GetCollection(models.NewTweet().TableName()), bson.M{"_id": id})
	if err != nil {
		return err
	}

	_, err = domain.NewESClient().Delete(c, &es_service.DeleteRequest{
		Header: &common.RequestHeader{
			Requester:   user.Username,
			OperatorUid: -1,
			TraceId:     -1,
		},
		Index:      constants.ESCINDEXTWEET,
		DocumentId: id,
	})
	if err != nil {
		zap.S().Errorf(err.Error())
	}
	return nil
}

func (*Service) TweetList(c *gin.Context, params *utils.PageForm) (*forms.TweetList, error) {
	db := models.NewDB()

	tweets, err := models.GWhereFind[models.Tweet](c, db.GetCollection(models.NewTweet().TableName()), bson.M{})
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(tweets))
	for _, tweet := range tweets {
		userIds = append(userIds, tweet.UserID)
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
	// 封装数据返回
	records := make([]*forms.Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		records = append(records, &forms.Tweet{
			UserId:       tweet.UserID,
			Username:     userIdToInfoMaps[tweet.UserID].Username,
			Avatar:       userIdToInfoMaps[tweet.UserID].Avatar,
			ID:           tweet.ID,
			Title:        tweet.Title,
			Content:      tweet.Content,
			CreatedAt:    tweet.CreatedAt.Format(constants.TimeFormat),
			ThumbCount:   tweet.ThumbCount,
			CommentCount: tweet.CommentCount,
		})
	}
	result := &forms.TweetList{
		Records: records,
	}
	return result, nil
}

func (*Service) TweetOwnList(c *gin.Context) (*forms.TweetList, error) {
	db := models.NewDB()
	user := utils.GetUser(c)

	tweets, err := models.GWhereFind[models.Tweet](c, db.GetCollection(models.NewTweet().TableName()), bson.M{"user_id": user.ID})
	if err != nil {
		return nil, err
	}
	records := make([]*forms.Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		records = append(records, &forms.Tweet{
			UserId:       user.ID,
			Username:     user.Username,
			Avatar:       user.Avatar,
			ID:           tweet.ID,
			Title:        tweet.Title,
			Content:      tweet.Content,
			CreatedAt:    tweet.CreatedAt.Format(constants.TimeFormat),
			ThumbCount:   tweet.ThumbCount,
			CommentCount: tweet.CommentCount,
		})
	}

	result := &forms.TweetList{
		Records: records,
	}
	return result, nil
}

func (*Service) TweetFavoriteList(c *gin.Context) (*forms.TweetList, error) {
	db := models.NewDB()
	user := utils.GetUser(c)

	favorites, err := models.GWhereFind[models.Favorite](c, db.GetCollection(models.NewFavorite().TableName()), bson.M{"user_id": user.ID})
	if err != nil {
		return nil, err
	}
	tweetIds := make([]string, 0, len(favorites))
	for _, favorite := range favorites {
		tweetIds = append(tweetIds, favorite.TweetId)
	}
	tweets, err := models.GWhereFind[models.Tweet](c, db.GetCollection(models.NewTweet().TableName()), bson.M{"_id": bson.M{"$in": tweetIds}})
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(tweets))
	for _, tweet := range tweets {
		userIds = append(userIds, tweet.UserID)
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
	// 封装数据返回
	records := make([]*forms.Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		records = append(records, &forms.Tweet{
			UserId:       tweet.UserID,
			Username:     userIdToInfoMaps[tweet.UserID].Username,
			Avatar:       userIdToInfoMaps[tweet.UserID].Avatar,
			ID:           tweet.ID,
			Title:        tweet.Title,
			Content:      tweet.Content,
			CreatedAt:    tweet.CreatedAt.Format(constants.TimeFormat),
			ThumbCount:   tweet.ThumbCount,
			CommentCount: tweet.CommentCount,
		})
	}

	result := &forms.TweetList{
		Records: records,
	}
	return result, nil
}

func (*Service) TweetFavorite(c *gin.Context, params *forms.TweetFavoriteForm) error {
	db := models.NewDB()
	user := utils.GetUser(c)

	// 查询此用户有无收藏此文章
	_, err := models.GWhereFirst[models.Favorite](c, db.GetCollection(models.NewFavorite().TableName()), bson.M{"user_id": user.ID, "tweet_id": params.Id})
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
	_, err = models.GInsertOne[models.Favorite](c, db.GetCollection(models.NewFavorite().TableName()), data)
	if err != nil {
		return err
	}
	return nil
}

func (*Service) TweetFavoriteDelete(c *gin.Context, id string) error {
	db := models.NewDB()
	user := utils.GetUser(c)

	_, err := models.GWhereDelete[models.Favorite](c, db.GetCollection(models.NewFavorite().TableName()), bson.M{"user_id": user.ID, "_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) TweetSearch(c *gin.Context, params *forms.SearchForm) (*forms.TweetList, error) {

	searchResult, err := domain.NewESClient().Search(c, &es_service.SearchRequest{
		Header: &common.RequestHeader{
			Requester:   "search_tweet",
			OperatorUid: -1,
		},
		ShouldQuery: &es_service.Query{
			MultiMatchQueries: []*es_service.MultiMatchQuery{
				{Field: []string{"title", "content"}, Value: params.Keyword},
			},
		},
		Indices: []string{constants.ESCINDEXTWEET},
		Page:    int32(params.Page),
		Size:    int32(params.Size),
		Pretty:  true,
	})
	if err != nil {
		return nil, err
	}

	records := make([]*forms.Tweet, 0, len(searchResult.Records))
	for _, search := range searchResult.Records {
		record := &forms.Tweet{}
		_ = json.Unmarshal([]byte(search), record)
		records = append(records, record)
	}

	result := &forms.TweetList{
		Records: records,
		PageList: utils.PageList{
			Size:    int(searchResult.PageList.GetSize()),
			Pages:   searchResult.PageList.GetPages(),
			Total:   searchResult.PageList.GetTotal(),
			Current: int(searchResult.PageList.GetCurrent()),
		},
	}
	return result, nil
}
