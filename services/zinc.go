package services

import (
	"context"
	"fmt"

	client "github.com/zinclabs/sdk-go-zincsearch"
)

type Zinc struct {
	Username string
	Password string
}

// 添加文档
func (z *Zinc) InsertDocument(ctx context.Context, index string, id string, document map[string]interface{}) error {
	configuration := client.NewConfiguration()
	apiClient := client.NewAPIClient(configuration)
	auth := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{z.Username, z.Password})
	_, _, err := apiClient.Document.IndexWithID(auth, index, id).Document(document).Execute()
	if err != nil {
		return err
	}
	return nil
}

// 更新文档 // 弃用，总显示404 Not Found 更新采用先删除后添加的方式
func (z *Zinc) UpdateDocument(ctx context.Context, index string, id string, document map[string]interface{}) error {
	fmt.Println(index, id, document)
	configuration := client.NewConfiguration()
	apiClient := client.NewAPIClient(configuration)
	auth := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{z.Username, z.Password})
	_, _, err := apiClient.Document.Update(auth, index, id).Document(document).Execute()
	if err != nil {
		return err
	}
	return nil
}

// 删除文档
func (z *Zinc) DeleteDocument(ctx context.Context, index string, id string) error {
	configuration := client.NewConfiguration()
	apiClient := client.NewAPIClient(configuration)
	auth := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{z.Username, z.Password})
	_, _, err := apiClient.Document.Delete(auth, index, id).Execute()
	if err != nil {
		return err
	}
	return nil
}

// 搜索文档
func (z *Zinc) SearchDocument(ctx context.Context, index string, queryString string, from int32, size int32) ([]map[string]interface{}, error) {
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
	configuration := client.NewConfiguration()
	apiClient := client.NewAPIClient(configuration)
	auth := context.WithValue(ctx, client.ContextBasicAuth, client.BasicAuth{z.Username, z.Password})
	resp, _, err := apiClient.Search.Search(auth, index).Query(query).Execute()
	if err != nil {
		return nil, err
	}
	// 搜集查询到内容，返回数据
	sources := make([]map[string]interface{}, len(resp.GetHits().Hits))
	for _, hit := range resp.GetHits().Hits {
		sources = append(sources, hit.Source)
	}
	return sources, nil
}
