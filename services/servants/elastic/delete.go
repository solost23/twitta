package es

import (
	"context"
	"errors"

	"twitta/global"

	"github.com/solost23/protopb/gen/go/common"
	es_service "github.com/solost23/protopb/gen/go/elastic"
)

// Delete 删除数据
func Delete(ctx context.Context, index, id string) error {
	if len(index) <= 0 {
		return errors.New("index 不为空")
	}
	if len(id) <= 0 {
		return errors.New("id 不为空")
	}
	reply, err := global.EsSrvClient.Delete(ctx, &es_service.DeleteRequest{
		Header: &common.RequestHeader{
			Requester:  "",
			OperatorId: -1,
			TraceId:    -1,
		},
		Index:      index,
		DocumentId: id,
	})
	if err != nil {
		return err
	}
	if reply.ErrorInfo.GetCode() != 0 {
		return errors.New(reply.ErrorInfo.GetMsg())
	}
	return nil
}
