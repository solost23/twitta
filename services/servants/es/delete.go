package es

import (
	"context"
	"errors"

	"github.com/solost23/protopb/gen/go/protos/common"
	es_service "github.com/solost23/protopb/gen/go/protos/es"
	"twitta/pkg/domain"
)

// Delete 删除数据
func Delete(ctx context.Context, index, id string) error {
	if len(index) <= 0 {
		return errors.New("index 不为空")
	}
	if len(id) <= 0 {
		return errors.New("id 不为空")
	}
	reply, err := domain.NewESClient().Delete(ctx, &es_service.DeleteRequest{
		Header: &common.RequestHeader{
			Requester:   "",
			OperatorUid: -1,
			TraceId:     -1,
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
