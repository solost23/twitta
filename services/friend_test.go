package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/http"
	"testing"
)

func TestService_FriendApplicationList(t *testing.T) {
	ginCtx, _ := gin.CreateTestContext(&http.TestResponseWriter{})
	type arg struct {
		ctx *gin.Context
	}
	type want struct {
		err error
	}
	tests := []struct {
		arg  arg
		want want
	}{
		{
			arg: arg{
				ctx: ginCtx,
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		results, err := (&Service{}).FriendApplicationList(test.arg.ctx)
		if err != nil {
			t.Errorf("error: %+v \n", err)
		}
		fmt.Printf("results: %+v \n", results)
	}
}