package services

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"time"

	"github.com/solost23/protopb/gen/go/common"
	es_service "github.com/solost23/protopb/gen/go/elastic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"twitta/forms"
	"twitta/global"
	"twitta/pkg/constants"
	"twitta/pkg/dao"
	"twitta/pkg/utils"
	servantElastic "twitta/services/servants/elastic"

	"github.com/gin-gonic/gin"
)

func (*Service) TweetSend(c *gin.Context, params *forms.TweetCreateForm) error {
	user := utils.GetUser(c)
	data := dao.Tweet{
		ID:           primitive.NewObjectID(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		UserID:       user.ID.String(),
		Title:        params.Title,
		Content:      params.Content,
		ThumbCount:   0,
		CommentCount: 0,
	}
	db := global.DB
	err := dao.GInsertOne[*dao.Tweet](c, db, &data)
	if err != nil {
		return err
	}

	go func() {
		// 推文携带创建者信息，方便后续直接搜索展示
		type Document struct {
			*dao.Tweet
			Username string `json:"username"`
			Avatar   string `json:"avatar"`
		}

		if err = servantElastic.Save(c, constants.ESCINDEXTWEET, data.ID.String(), Document{Tweet: &data, Username: user.Username, Avatar: user.Avatar}); err != nil {
			zap.S().Error(err)
		}
	}()

	return nil
}

func (s *Service) StaticUpload(c *gin.Context, file *multipart.FileHeader) (string, error) {
	folder := "twitta.tweets.static"

	url, err := UploadImg(0, folder, file.Filename, file)
	if err != nil {
		return "", err
	}
	return utils.FulfillImageOSSPrefix(utils.TrimDomainPrefix(url)), nil
}

func (*Service) TweetDelete(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	db := global.DB
	tweet, err := dao.GWhereFirst[*dao.Tweet](c, db, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if user.ID.String() != tweet.UserID {
		return errors.New("本推文所属用户不是您，无权删除")
	}
	_, err = dao.GWhereDelete[*dao.Tweet](c, db, bson.M{"_id": id})
	if err != nil {
		return err
	}

	go func() {
		if err = servantElastic.Delete(c, constants.ESCINDEXTWEET, id); err != nil {
			zap.S().Error(err)
		}
	}()

	return nil
}

func (*Service) TweetList(c *gin.Context, params *utils.PageForm) (*forms.TweetList, error) {
	db := global.DB
	tweets, err := dao.GWhereFind[*dao.Tweet](c, db, bson.M{})
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(tweets))
	for _, tweet := range tweets {
		userIds = append(userIds, tweet.UserID)
	}
	users, err := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": userIds}})
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
	// 封装数据返回
	records := make([]*forms.Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		records = append(records, &forms.Tweet{
			UserId:       tweet.UserID,
			Username:     userIdToInfoMaps[tweet.UserID].Username,
			Avatar:       userIdToInfoMaps[tweet.UserID].Avatar,
			ID:           tweet.ID.String(),
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
	user := utils.GetUser(c)

	db := global.DB
	tweets, err := dao.GWhereFind[*dao.Tweet](c, db, bson.M{"user_id": user.ID})
	if err != nil {
		return nil, err
	}
	records := make([]*forms.Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		records = append(records, &forms.Tweet{
			UserId:       user.ID.String(),
			Username:     user.Username,
			Avatar:       user.Avatar,
			ID:           tweet.ID.String(),
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
	user := utils.GetUser(c)

	db := global.DB
	favorites, err := dao.GWhereFind[*dao.Favorite](c, db, bson.M{"user_id": user.ID})
	if err != nil {
		return nil, err
	}
	tweetIds := make([]string, 0, len(favorites))
	for _, favorite := range favorites {
		tweetIds = append(tweetIds, favorite.TweetId)
	}
	tweets, err := dao.GWhereFind[*dao.Tweet](c, db, bson.M{"_id": bson.M{"$in": tweetIds}})
	if err != nil {
		return nil, err
	}
	userIds := make([]string, 0, len(tweets))
	for _, tweet := range tweets {
		userIds = append(userIds, tweet.UserID)
	}
	users, err := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": userIds}})
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
	// 封装数据返回
	records := make([]*forms.Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		records = append(records, &forms.Tweet{
			UserId:       tweet.UserID,
			Username:     userIdToInfoMaps[tweet.UserID].Username,
			Avatar:       userIdToInfoMaps[tweet.UserID].Avatar,
			ID:           tweet.ID.String(),
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
	user := utils.GetUser(c)

	// 查询此用户有无收藏此文章
	db := global.DB
	_, err := dao.GWhereFirst[*dao.Favorite](c, db, bson.M{"user_id": user.ID, "tweet_id": params.Id})
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}
	if err == nil {
		return errors.New("您已收藏此推文")
	}
	data := &dao.Favorite{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserId:    user.ID.String(),
		TweetId:   params.Id,
	}
	err = dao.GInsertOne[*dao.Favorite](c, db, data)
	if err != nil {
		return err
	}
	return nil
}

func (*Service) TweetFavoriteDelete(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	db := global.DB
	_, err := dao.GWhereDelete[*dao.Favorite](c, db, bson.M{"user_id": user.ID, "_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (*Service) TweetSearch(c *gin.Context, params *forms.SearchForm) (*forms.TweetList, error) {

	searchResult, err := global.EsSrvClient.Search(c, &es_service.SearchRequest{
		Header: &common.RequestHeader{
			Requester:  "search_tweet",
			OperatorId: -1,
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
