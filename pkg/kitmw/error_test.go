package kitmw

import (
	"context"
	"testing"

	"github.com/go-kit/kit/endpoint"
	"github.com/DoNewsCode/std/pkg/srverr"
	"github.com/pkg/errors"
)

func TestMakeErrorMarshallerMiddleware(t *testing.T) {
	mw := MakeErrorMarshallerMiddleware()
	e1 := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, errors.New("foo")
	}
	e2 := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, srverr.NotFoundErr(errors.New("bar"), "")
	}
	e3 := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, srverr.NotFoundErr(srverr.InvalidArgumentErr(errors.New("bar"), ""), "")
	}
	e4 := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, errors.Wrap(srverr.NotFoundErr(errors.New("foo"), ""), "bar")
	}
	cases := []endpoint.Endpoint{e1, e2, e3, e4}
	for _, c := range cases {
		cc := c
		t.Run("", func(t *testing.T) {
			_, err := mw(cc)(nil, nil)
			if _, ok := err.(srverr.ServerError); !ok {
				t.Fail()
			}
		})
	}
}
