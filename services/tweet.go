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
		Images:       params.Images,
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

	oid, err := utils.ParseObjectID(id)
	if err != nil {
		return errors.New("无效推文ID")
	}
	db := global.DB
	tweet, err := dao.GWhereFirst[*dao.Tweet](c, db, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	if user.ID.String() != tweet.UserID {
		return errors.New("本推文所属用户不是您，无权删除")
	}
	_, err = dao.GWhereDelete[*dao.Tweet](c, db, bson.M{"_id": oid})
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
	tweets, total, pages, err := dao.GPaginatorOrder[*dao.Tweet](c, db, &dao.ListPageInput{
		Page: params.Page,
		Size: params.Size,
	}, bson.M{"created_at": -1}, bson.M{})
	if err != nil {
		return nil, err
	}
	userOids := parseObjectIDs(collectStrings(tweets, func(t *dao.Tweet) string { return t.UserID }))
	users, err := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": userOids}})
	if err != nil {
		return nil, err
	}
	userMap := buildUserMap(users)
	records := make([]*forms.Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		u := userMap[tweet.UserID]
		records = append(records, tweetToForm(tweet, u.Username, u.Avatar))
	}
	fillOriginTweets(c, db, records)
	result := &forms.TweetList{
		PageList: utils.PageList{
			Size:    params.Size,
			Pages:   pages,
			Total:   total,
			Current: params.Page,
		},
		Records: records,
	}
	return result, nil
}

func (*Service) TweetOwnList(c *gin.Context) (*forms.TweetList, error) {
	user := utils.GetUser(c)

	db := global.DB
	tweets, err := dao.GWhereFind[*dao.Tweet](c, db, bson.M{"user_id": user.ID.String()})
	if err != nil {
		return nil, err
	}
	records := make([]*forms.Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		records = append(records, tweetToForm(tweet, user.Username, user.Avatar))
	}
	fillOriginTweets(c, db, records)

	result := &forms.TweetList{
		Records: records,
	}
	return result, nil
}

func (*Service) TweetFavoriteList(c *gin.Context) (*forms.TweetList, error) {
	user := utils.GetUser(c)

	db := global.DB
	favorites, err := dao.GWhereFind[*dao.Favorite](c, db, bson.M{"user_id": user.ID.String()})
	if err != nil {
		return nil, err
	}
	tweetOids := parseObjectIDs(collectStrings(favorites, func(f *dao.Favorite) string { return f.TweetId }))
	tweets, err := dao.GWhereFind[*dao.Tweet](c, db, bson.M{"_id": bson.M{"$in": tweetOids}})
	if err != nil {
		return nil, err
	}
	userOids := parseObjectIDs(collectStrings(tweets, func(t *dao.Tweet) string { return t.UserID }))
	users, err := dao.GWhereFind[*dao.User](c, db, bson.M{"_id": bson.M{"$in": userOids}})
	if err != nil {
		return nil, err
	}
	userMap := buildUserMap(users)
	records := make([]*forms.Tweet, 0, len(tweets))
	for _, tweet := range tweets {
		u := userMap[tweet.UserID]
		records = append(records, tweetToForm(tweet, u.Username, u.Avatar))
	}
	fillOriginTweets(c, db, records)

	result := &forms.TweetList{
		Records: records,
	}
	return result, nil
}

func (*Service) TweetFavorite(c *gin.Context, params *forms.TweetFavoriteForm) error {
	user := utils.GetUser(c)

	// 查询此用户有无收藏此文章
	db := global.DB
	_, err := dao.GWhereFirst[*dao.Favorite](c, db, bson.M{"user_id": user.ID.String(), "tweet_id": params.Id})
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

	oid, err := utils.ParseObjectID(id)
	if err != nil {
		return errors.New("无效收藏ID")
	}
	db := global.DB
	_, err = dao.GWhereDelete[*dao.Favorite](c, db, bson.M{"user_id": user.ID.String(), "_id": oid})
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

func tweetToForm(tweet *dao.Tweet, username, avatar string) *forms.Tweet {
	return &forms.Tweet{
		UserId:       tweet.UserID,
		Username:     username,
		Avatar:       avatar,
		ID:           tweet.ID.String(),
		Title:        tweet.Title,
		Content:      tweet.Content,
		Images:       tweet.Images,
		CreatedAt:    tweet.CreatedAt.Format(time.DateTime),
		ThumbCount:   tweet.ThumbCount,
		CommentCount: tweet.CommentCount,
		RetweetCount: tweet.RetweetCount,
		RetweetOf:    tweet.RetweetOf,
	}
}

func (*Service) TweetDetail(c *gin.Context, id string) (*forms.Tweet, error) {
	oid, err := utils.ParseObjectID(id)
	if err != nil {
		return nil, errors.New("无效推文ID")
	}
	db := global.DB
	tweet, err := dao.GWhereFirst[*dao.Tweet](c, db, bson.M{"_id": oid})
	if err != nil {
		return nil, err
	}
	authorOid, err := utils.ParseObjectID(tweet.UserID)
	if err != nil {
		return nil, errors.New("无效用户ID")
	}
	author, err := dao.GWhereFirst[*dao.User](c, db, bson.M{"_id": authorOid})
	if err != nil {
		return nil, err
	}
	result := tweetToForm(tweet, author.Username, author.Avatar)

	if tweet.RetweetOf != "" {
		originOid, err := utils.ParseObjectID(tweet.RetweetOf)
		if err == nil {
			origin, err := dao.GWhereFirst[*dao.Tweet](c, db, bson.M{"_id": originOid})
			if err == nil {
				originAuthorOid, err := utils.ParseObjectID(origin.UserID)
				if err == nil {
					originAuthor, err := dao.GWhereFirst[*dao.User](c, db, bson.M{"_id": originAuthorOid})
					if err == nil {
						result.OriginTweet = tweetToForm(origin, originAuthor.Username, originAuthor.Avatar)
					}
				}
			}
		}
	}
	return result, nil
}

func (*Service) TweetRetweet(c *gin.Context, id string) error {
	user := utils.GetUser(c)

	oid, err := utils.ParseObjectID(id)
	if err != nil {
		return errors.New("无效推文ID")
	}
	db := global.DB
	origin, err := dao.GWhereFirst[*dao.Tweet](c, db, bson.M{"_id": oid})
	if err != nil {
		return errors.New("推文不存在")
	}
	// 不允许重复转发同一条
	_, err = dao.GWhereFirst[*dao.Tweet](c, db, bson.M{"user_id": user.ID.String(), "retweet_of": id})
	if err == nil {
		return errors.New("您已转发过这条推文")
	}
	if !errors.Is(err, mongo.ErrNoDocuments) {
		return err
	}

	data := &dao.Tweet{
		ID:        primitive.NewObjectID(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID.String(),
		RetweetOf: id,
	}
	if err = dao.GInsertOne[*dao.Tweet](c, db, data); err != nil {
		return err
	}
	_, err = dao.GWhereUpdate[*dao.Tweet](c, db, bson.M{"_id": origin.ID}, bson.M{"$inc": bson.M{"retweet_count": 1}})
	return err
}

