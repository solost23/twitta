package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/http"
	"testing"
	"twitta/forms"
	"twitta/pkg/utils"
)

func TestService_ChatList(t *testing.T) {
	ctx, _ := gin.CreateTestContext(&http.TestResponseWriter{})
	type arg struct {
		ginCtx *gin.Context
		id     string
	}
	type want struct {
		results *forms.ChatList
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
		results, err := (&Service{}).ChatList(test.arg.ginCtx, test.arg.id, &utils.PageForm{Page: 1, Size: 10})
		if err != test.want.err {
			t.Errorf("error: %+v \n", err)
		}
		fmt.Println("results: ", results)
	}
}
