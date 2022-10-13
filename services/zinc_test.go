package services

import (
	"context"
	"fmt"
	"testing"
)

func TestZinc_InsertDocument(t *testing.T) {
	z := &Zinc{"admin", "Complexpass#123"}
	type arg struct {
		ctx      context.Context
		index    string
		id       string
		document map[string]interface{}
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		arg  arg
		want want
	}{
		{
			name: "t1",
			arg: arg{
				ctx:   context.Background(),
				index: "user",
				id:    "1",
				document: map[string]interface{}{
					"id":   1,
					"name": "ty1",
					"age":  18,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "t2",
			arg: arg{
				ctx:   context.Background(),
				index: "user",
				id:    "2",
				document: map[string]interface{}{
					"id":   2,
					"name": "ty2",
					"age":  18,
				},
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "t3",
			arg: arg{
				ctx:   context.Background(),
				index: "user1",
				id:    "3",
				document: map[string]interface{}{
					"id":   3,
					"name": "ty3",
					"age":  18,
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := z.InsertDocument(test.arg.ctx, test.arg.index, test.arg.id, test.arg.document)
			if err != test.want.err {
				t.Errorf("err: %v", err)
			}
		})
	}

}

func TestZinc_SearchDocument(t *testing.T) {
	z := &Zinc{"admin", "Complexpass#123"}
	type arg struct {
		ctx         context.Context
		index       string
		queryString string
		from        int32
		size        int32
	}
	type want struct {
		result []map[string]interface{}
		err    error
	}
	tests := []struct {
		name string
		arg  arg
		want want
	}{
		{
			name: "t1",
			arg: arg{
				ctx:         context.Background(),
				index:       "user",
				queryString: "18",
				from:        0,
				size:        5,
			},
			want: want{
				result: nil,
				err:    nil,
			},
		},
		{
			name: "t2",
			arg: arg{
				ctx:         context.Background(),
				index:       "user",
				queryString: "ty1",
				from:        0,
				size:        20,
			},
			want: want{
				result: nil,
				err:    nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := z.SearchDocument(test.arg.ctx, test.arg.index, test.arg.queryString, test.arg.from, test.arg.size)
			if err != nil {
				t.Errorf("%v", err)
			}
			fmt.Println(result)
		})
	}
}

func TestZinc_DeleteDocument(t *testing.T) {
	z := &Zinc{"admin", "Complexpass#123"}
	type arg struct {
		ctx   context.Context
		index string
		id    string
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		arg  arg
		want want
	}{
		{
			name: "t1",
			arg: arg{
				ctx:   context.Background(),
				index: "user1",
				id:    "3",
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := z.DeleteDocument(test.arg.ctx, test.arg.index, test.arg.id)
			if err != test.want.err {
				t.Errorf("err: %v", err)
			}
		})
	}
}

func TestZinc_UpdateDocument(t *testing.T) {
	z := &Zinc{"admin", "Complexpass#123"}
	type arg struct {
		ctx      context.Context
		index    string
		id       string
		document map[string]interface{}
	}
	type want struct {
		err error
	}
	tests := []struct {
		name string
		arg  arg
		want want
	}{
		{
			name: "t1",
			arg: arg{
				ctx:   context.Background(),
				index: "user",
				id:    "2",
				document: map[string]interface{}{
					"name": "alex",
					"age":  20,
					"id":   "2",
				},
			},
			want: want{
				err: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := z.UpdateDocument(test.arg.ctx, test.arg.index, test.arg.id, test.arg.document)
			if err != test.want.err {
				t.Errorf("err: %v", err)
			}
		})
	}
}
