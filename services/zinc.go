package services

import (
	"context"
	"fmt"
	client "github.com/zinclabs/sdk-go-zincsearch"
	"twitta/global"
)

type Zinc struct {
	Host     string
	Port     int
	Username string
	Password string
	conf     *client.Configuration
}

func NewZinc() *Zinc {
	conf := client.NewConfiguration()
	conf.Servers[0].URL = fmt.Sprintf("%s:%d", global.ServerConfig.Zinc.Host, global.ServerConfig.Zinc.Port)
	conf.Servers[1].URL = fmt.Sprintf("%s:%d", global.ServerConfig.Zinc.Host, global.ServerConfig.Zinc.Port)
	return &Zinc{
		Host:     global.ServerConfig.Zinc.Host,
		Port:     global.ServerConfig.Zinc.Port,
		Username: global.ServerConfig.Zinc.Username,
		Password: global.ServerConfig.Zinc.Password,
		conf:     conf,
	}
}

// 检查索引是否存在
func (z *Zinc) EsExists(ctx context.Context, index string) bool {
	apiClient := client.NewAPIClient(z.conf)
	auth := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{z.Username, z.Password})
	_, _, err := apiClient.Index.EsExists(auth, index).Execute()
	if err != nil {
		return false
	}
	return true
}

// 创建索引
func (z *Zinc) InsertIndex(ctx context.Context, data client.MetaIndexSimple) error {
	apiClient := client.NewAPIClient(z.conf)
	auth := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{z.Username, z.Password})
	_, _, err := apiClient.Index.Create(auth).Data(data).Execute()
	if err != nil {
		return err
	}
	return nil
}

// 添加文档
func (z *Zinc) InsertDocument(ctx context.Context, index string, id string, document map[string]interface{}) error {
	apiClient := client.NewAPIClient(z.conf)
	auth := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{z.Username, z.Password})
	_, _, err := apiClient.Document.IndexWithID(auth, index, id).Document(document).Execute()
	if err != nil {
		return err
	}
	return nil
}

// 更新文档 发现总是404 弃用，采取删除+插入的方式构建全文索引
func (z *Zinc) UpdateDocument(ctx context.Context, index string, id string, document map[string]interface{}) error {
	apiClient := client.NewAPIClient(z.conf)
	auth := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{z.Username, z.Password})
	_, _, err := apiClient.Document.Update(auth, index, id).Document(document).Execute()
	if err != nil {
		return err
	}
	return nil
}

// 删除文档
func (z *Zinc) DeleteDocument(ctx context.Context, index string, id string) error {
	apiClient := client.NewAPIClient(z.conf)
	auth := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{z.Username, z.Password})
	_, _, err := apiClient.Document.Delete(auth, index, id).Execute()
	if err != nil {
		return err
	}
	return nil
}

// 搜索文档
func (z *Zinc) SearchDocument(ctx context.Context, index string, queryString string, from int32, size int32) ([]client.MetaHit, int64, error) {
	query := client.MetaZincQuery{
		Query: &client.MetaQuery{
			Bool: &client.MetaBoolQuery{
				Must: []client.MetaQuery{
					client.MetaQuery{
						QueryString: &client.MetaQueryStringQuery{
							Query: &queryString,
						},
					},
				},
			},
		},
		From: &from,
		Size: &size,
	}
	apiClient := client.NewAPIClient(z.conf)
	auth := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{z.Username, z.Password})
	resp, _, err := apiClient.Search.Search(auth, index).Query(query).Execute()
	if err != nil {
		return nil, 0, err
	}
	var count = int64(*resp.GetHits().Total.Value)
	// 搜集查询到内容，返回数据
	hits := make([]client.MetaHit, 0, len(resp.GetHits().Hits))
	for _, hit := range resp.GetHits().Hits {
		hits = append(hits, hit)
	}
	return hits, count, nil
}
