package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/http"
	"testing"
	"twitta/forms"
)

func TestService_FanList(t *testing.T) {
	ginCtx, _ := gin.CreateTestContext(&http.TestResponseWriter{})
	type arg struct {
		ctx *gin.Context
	}
	type want struct {
		results []*forms.FansAndWhatResponse
		err     error
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
		results, err := (&Service{}).FanList(test.arg.ctx)
		if err != nil {
			t.Errorf("err: %+v \n", err)
		}
		fmt.Printf("results: %+v \n", results)
	}
}
