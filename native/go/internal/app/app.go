// Package app wires every protocol method to its implementation. Both the
// native-messaging host and the local daemon use the same Router.
package app

import (
	"context"
	"encoding/json"
	"errors"
	"runtime"
	"strings"
	"time"

	"example.com/toolbox/native/internal/config"
	"example.com/toolbox/native/internal/notify"
	"example.com/toolbox/native/internal/policy"
	"example.com/toolbox/native/internal/rpc"
	"example.com/toolbox/native/internal/sidecar"
	"example.com/toolbox/native/internal/watcher"
	"example.com/toolbox/native/pkg/protocol"
)

// App owns the long-lived resources behind the router.
type App struct {
	Router   *rpc.Router
	policy   *policy.Policy
	watchers *watcher.Manager
	sidecar  *sidecar.Client
}

// New builds the app. emit is called for every pushed event (may be nil).
func New(cfg config.Config, emit func(protocol.Event)) (*App, error) {
	pol, err := policy.New(cfg.AllowedRoots)
	if err != nil {
		return nil, err
	}
	a := &App{Router: rpc.NewRouter(), policy: pol}
	a.watchers = watcher.NewPolling(time.Second, func(e watcher.Event) {
		if emit == nil {
			return
		}
		data, _ := json.Marshal(e)
		emit(protocol.Event{Event: "fs.changed", Data: data})
	})
	if cfg.Sidecar != "" {
		// A missing sidecar is not fatal: sidecar.* then answers "unavailable".
		if c, err := sidecar.Start(cfg.Sidecar); err == nil {
			a.sidecar = c
		}
	}
	a.register()
	return a, nil
}

// Close releases watchers and the sidecar.
func (a *App) Close() {
	a.watchers.Close()
	if a.sidecar != nil {
		_ = a.sidecar.Close()
	}
}

func decode(params json.RawMessage, v any) error {
	if len(params) == 0 {
		return rpc.Errorf(protocol.CodeBadRequest, "missing params")
	}
	if err := json.Unmarshal(params, v); err != nil {
		return rpc.Errorf(protocol.CodeBadRequest, "invalid params: %v", err)
	}
	return nil
}

func (a *App) register() {
	r := a.Router

	r.Register("ping", func(context.Context, json.RawMessage) (any, error) {
		return map[string]any{"pong": true, "version": protocol.Version, "os": runtime.GOOS, "arch": runtime.GOARCH}, nil
	})

	r.Register("system.methods", func(context.Context, json.RawMessage) (any, error) {
		return r.Methods(), nil
	})

	r.Register("fs.watch", func(_ context.Context, p json.RawMessage) (any, error) {
		var in struct {
			Path      string `json:"path"`
			Recursive bool   `json:"recursive"`
		}
		if err := decode(p, &in); err != nil {
			return nil, err
		}
		path, err := a.policy.Check(in.Path)
		if err != nil {
			return nil, rpc.Errorf(protocol.CodeForbidden, "%v", err)
		}
		id, err := a.watchers.Add(path, in.Recursive)
		if err != nil {
			return nil, rpc.Errorf(protocol.CodeBadRequest, "watch: %v", err)
		}
		return map[string]string{"watchId": id}, nil
	})

	r.Register("fs.unwatch", func(_ context.Context, p json.RawMessage) (any, error) {
		var in struct {
			WatchID string `json:"watchId"`
		}
		if err := decode(p, &in); err != nil {
			return nil, err
		}
		return map[string]bool{"removed": a.watchers.Remove(in.WatchID)}, nil
	})

	r.Register("notify.show", func(ctx context.Context, p json.RawMessage) (any, error) {
		var in struct {
			Title string `json:"title"`
			Body  string `json:"body"`
		}
		if err := decode(p, &in); err != nil {
			return nil, err
		}
		if strings.TrimSpace(in.Title) == "" {
			return nil, rpc.Errorf(protocol.CodeBadRequest, "title is required")
		}
		if err := notify.Show(ctx, in.Title, in.Body); err != nil {
			if errors.Is(err, notify.ErrUnsupported) {
				return nil, rpc.Errorf(protocol.CodeUnavailable, "%v", err)
			}
			return nil, err
		}
		return map[string]bool{"shown": true}, nil
	})

	r.RegisterPrefix("sidecar.", func(ctx context.Context, method string, p json.RawMessage) (any, error) {
		if a.sidecar == nil {
			return nil, rpc.Errorf(protocol.CodeUnavailable, "sidecar is not configured")
		}
		res, err := a.sidecar.Call(ctx, method, p)
		if err != nil {
			var re *sidecar.RemoteError
			if errors.As(err, &re) {
				return nil, rpc.Errorf(re.Code, "%s", re.Message)
			}
			if errors.Is(err, sidecar.ErrNotRunning) {
				return nil, rpc.Errorf(protocol.CodeUnavailable, "%v", err)
			}
			if errors.Is(err, context.DeadlineExceeded) {
				return nil, rpc.Errorf(protocol.CodeTimeout, "sidecar timed out")
			}
			return nil, err
		}
		return json.RawMessage(res), nil
	})
}
