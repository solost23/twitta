package services

import "testing"

func TestConsumeEmailMessage(t *testing.T) {
	type arg struct {
		topic       string
		name        string
		addr        string
		contentType string
		content     string
	}
	type want struct {
	}
	tests := []struct {
		arg  arg
		want want
	}{
		{
			arg: arg{
				topic:       "测试1",
				name:        "ty",
				addr:        "2805347171@qq.com",
				contentType: "text/plain",
				content:     "测试内容1",
			},
			want: want{},
		},
		{
			arg: arg{
				topic:       "测试2",
				name:        "ty",
				addr:        "2805347171@qq.com",
				contentType: "text/plain",
				content:     "测试内容2",
			},
			want: want{},
		},
	}

	for _, test := range tests {
		SendEmail(test.arg.topic, test.arg.name, test.arg.addr, test.arg.contentType, test.arg.content)
	}

}
