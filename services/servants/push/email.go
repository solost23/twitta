package push

import (
	"context"
	"errors"

	"github.com/solost23/protopb/gen/go/protos/common"
	"github.com/solost23/protopb/gen/go/protos/push"
	"twitta/pkg/domain"
)

// SendEmail 发送邮件服务
func SendEmail(ctx context.Context, topic, username, addr, content string) error {
	if len(addr) <= 0 {
		return nil
	}
	reply, err := domain.NewPushClient().SendEmail(ctx, &push.SendEmailRequest{
		Header: &common.RequestHeader{
			TraceId:     6678677,
			OperatorUid: 55,
		},
		Email: &push.Email{
			Topic:       topic,
			Name:        username,
			Addr:        addr,
			ContentType: "text/plain",
			Content:     content,
		},
	})
	if err != nil {
		return err
	}
	if reply.ErrorInfo.GetCode() != 0 {
		return errors.New(reply.ErrorInfo.GetMsg())
	}
	return nil
}