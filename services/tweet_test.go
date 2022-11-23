package services

import (
	"Twitta/global/initialize"
	"Twitta/pkg/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	initialize.Initialize("../configs/configs.yml")
	os.Exit(m.Run())
}

func TestService_TweetList(t *testing.T) {
	ginCtx, _ := gin.CreateTestContext(&http.TestResponseWriter{})
	ginCtx.Set("user", &models.User{})
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
		results, err := (&Service{}).TweetList(test.arg.ctx)
		if err != test.want.err {
			t.Errorf("%v \n", err.Error())
		}
		fmt.Println("results: ", results)
	}

}
