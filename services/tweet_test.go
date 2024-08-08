package services

import (
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/http"
	"twitta/global/initialize"
	"twitta/pkg/dao"
	"twitta/pkg/utils"
)

func TestMain(m *testing.M) {
	initialize.Initialize("../configs/configs.yml")
	os.Exit(m.Run())
}

func TestServiceTweetList(t *testing.T) {
	ginCtx, _ := gin.CreateTestContext(&http.TestResponseWriter{})
	ginCtx.Set("user", &dao.User{})
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
		results, err := (&Service{}).TweetList(test.arg.ctx, &utils.PageForm{})
		if err != test.want.err {
			t.Errorf("%v \n", err.Error())
		}
		fmt.Println("results: ", results)
	}

}
