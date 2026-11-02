package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"example.com/toolbox/native/pkg/protocol"
)

func TestDispatch(t *testing.T) {
	r := NewRouter()
	r.Register("ok", func(context.Context, json.RawMessage) (any, error) { return map[string]int{"n": 1}, nil })
	r.Register("void", func(context.Context, json.RawMessage) (any, error) { return nil, nil })
	r.Register("typed", func(context.Context, json.RawMessage) (any, error) {
		return nil, Errorf(protocol.CodeForbidden, "nope %d", 7)
	})
	r.Register("plain", func(context.Context, json.RawMessage) (any, error) { return nil, errors.New("boom") })
	r.Register("panic", func(context.Context, json.RawMessage) (any, error) { panic("kaput") })
	r.RegisterPrefix("fwd.", func(_ context.Context, m string, _ json.RawMessage) (any, error) { return m, nil })
	r.RegisterPrefix("fwd.deep.", func(_ context.Context, m string, _ json.RawMessage) (any, error) { return "deep:" + m, nil })

	cases := []struct {
		method string
		ok     bool
		code   string
		result string
	}{
		{"ok", true, "", `{"n":1}`},
		{"void", true, "", ``},
		{"typed", false, protocol.CodeForbidden, ``},
		{"plain", false, protocol.CodeInternal, ``},
		{"panic", false, protocol.CodeInternal, ``},
		{"nope", false, protocol.CodeMethodNotFound, ``},
		{"", false, protocol.CodeBadRequest, ``},
		{"fwd.x", true, "", `"fwd.x"`},
		{"fwd.deep.x", true, "", `"deep:fwd.deep.x"`},
	}
	for _, c := range cases {
		resp := r.Dispatch(context.Background(), protocol.Request{ID: "id-" + c.method, Method: c.method})
		if resp.ID != "id-"+c.method {
			t.Errorf("%q: id not echoed: %q", c.method, resp.ID)
		}
		if resp.OK != c.ok {
			t.Errorf("%q: ok=%v want %v (%+v)", c.method, resp.OK, c.ok, resp.Error)
		}
		if !c.ok && (resp.Error == nil || resp.Error.Code != c.code) {
			t.Errorf("%q: error=%+v want code %s", c.method, resp.Error, c.code)
		}
		if string(resp.Result) != c.result {
			t.Errorf("%q: result=%s want %s", c.method, resp.Result, c.result)
		}
	}
}
