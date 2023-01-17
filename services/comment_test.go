package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/http"
	"testing"
	"twitta/forms"
)

func TestService_CommentList(t *testing.T) {
	ginCtx, _ := gin.CreateTestContext(&http.TestResponseWriter{})
	type arg struct {
		ctx *gin.Context
		id  string
	}
	type want struct {
		results []*forms.CommentListResponse
		err     error
	}
	tests := []struct {
		arg  arg
		want want
	}{
		{
			arg: arg{
				ctx: ginCtx,
				id:  "",
			},
			want: want{
				err: nil,
			},
		},
		{
			arg: arg{
				ctx: ginCtx,
				id:  "",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		results, err := (&Service{}).CommentList(test.arg.ctx, test.arg.id)
		if err != nil {
			t.Errorf("err: %+v \n", err)
		}
		fmt.Printf("result: %+v \n", results)
	}
}
