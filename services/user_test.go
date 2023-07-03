package services

import (
	"fmt"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/http"
	"twitta/forms"
	"twitta/pkg/utils"
)

func TestServiceUserSearch(t *testing.T) {
	ginCtx, _ := gin.CreateTestContext(&http.TestResponseWriter{})
	type arg struct {
		ctx    *gin.Context
		params *forms.SearchForm
	}
	type want struct {
		results []*forms.UserSearch
		err     error
	}
	tests := []struct {
		arg  arg
		want want
	}{
		{
			arg: arg{
				ctx: ginCtx,
				params: &forms.SearchForm{
					PageForm: utils.PageForm{
						Page: 1,
						Size: 10,
					},
					Keyword: "user",
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			arg: arg{
				ctx: ginCtx,
				params: &forms.SearchForm{
					PageForm: utils.PageForm{
						Page: 1,
						Size: 10,
					},
					Keyword: "用户",
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		results, err := (&Service{}).UserSearch(test.arg.ctx, test.arg.params)
		if err != test.want.err {
			t.Errorf("%v \n", err.Error())
		}
		fmt.Printf("results: %v \n", results)
	}
}
