package es

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/solost23/protopb/gen/go/protos/common"
	es_service "github.com/solost23/protopb/gen/go/protos/es"
	"twitta/global"
)

// Save 将数据存入ES
func Save(ctx context.Context, index, id string, data any) error {
	if len(index) <= 0 {
		return errors.New("index 不为空")
	}
	if len(id) <= 0 {
		return errors.New("id 不为空")
	}
	buf, err := json.Marshal(data)
	if err != nil {
		return err
	}
	reply, err := global.EsSrvClient.Create(ctx, &es_service.CreateRequest{
		Header: &common.RequestHeader{
			Requester:   "",
			OperatorUid: -1,
			TraceId:     -1,
		},
		Index:      index,
		DocumentId: id,
		Document:   string(buf),
	})
	if err != nil {
		return err
	}
	if reply.ErrorInfo.GetCode() != 0 {
		return errors.New(reply.ErrorInfo.GetMsg())
	}
	return nil
}
