package services

import (
	"Twitta/forms"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/http"
	"testing"
)

func TestService_ChatList(t *testing.T) {
	ctx, _ := gin.CreateTestContext(&http.TestResponseWriter{})
	type arg struct {
		ginCtx *gin.Context
		id     string
	}
	type want struct {
		results []*forms.ChatListResponse
		err     error
	}
	tests := []struct {
		arg  arg
		want want
	}{
		{
			arg: arg{
				ginCtx: ctx,
				id:     "",
			},
			want: want{
				err: nil,
			},
		},
		{
			arg: arg{
				ginCtx: ctx,
				id:     "",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		results, err := (&Service{}).ChatList(test.arg.ginCtx, test.arg.id)
		if err != test.want.err {
			t.Errorf("error: %+v \n", err)
		}
		fmt.Println("results: ", results)
	}
}
