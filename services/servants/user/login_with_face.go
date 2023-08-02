package user

import (
	"context"
	"errors"

	"github.com/solost23/protopb/gen/go/common"
	"github.com/solost23/protopb/gen/go/recognition"
	"twitta/global"
)

// GenerateFaceEncodings 根据用户画像生成用户脸部编码, 返回顺序与传入用户顺序一致
func GenerateFaceEncodings(ctx context.Context, faceImages [][]byte) ([]string, error) {
	if len(faceImages) <= 0 {
		return []string{}, nil
	}
	reply, err := global.FaceRecognitionSrvClient.GenerateFaceEncoding(ctx, &recognition.GenerateFaceEncodingRequest{
		Header: &common.RequestHeader{
			TraceId:    6678678,
			OperatorId: 56,
		},
		Data: faceImages,
	})
	if err != nil {
		return nil, err
	}
	if reply.GetErrorInfo().GetCode() != 0 {
		return nil, errors.New(reply.GetErrorInfo().GetMsg())
	}
	if len(faceImages) != len(reply.GetFaceEncodings()) {
		return nil, errors.New("查询出的用户编码数量比头像数少")
	}
	return reply.GetFaceEncodings(), nil
}

// CompareFace 根据用户画像对比人像库脸部编码，返回人的UserId，若没有，返回false
func CompareFace(ctx context.Context, faceImage []byte) (userId string, isFound bool, err error) {
	reply, err := global.FaceRecognitionSrvClient.CompareFaces(ctx, &recognition.CompareFacesRequest{
		Header: &common.RequestHeader{
			TraceId:    6678678,
			OperatorId: 56,
		},
		Data: faceImage,
	})
	if err != nil {
		return "", false, err
	}
	if reply.GetErrorInfo().GetCode() != 0 {
		return "", false, errors.New(reply.GetErrorInfo().GetMsg())
	}
	if !reply.GetIsFound() {
		return "", reply.GetIsFound(), nil
	}
	return reply.GetUserId(), reply.GetIsFound(), nil
}
