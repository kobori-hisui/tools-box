// toolboxd is the always-on local daemon (loopback HTTP + SSE). Double-click
// friendly: it reads <user config dir>/toolbox/config.json (or $TOOLBOX_CONFIG).
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/toolbox/native/internal/app"
	"example.com/toolbox/native/internal/config"
	"example.com/toolbox/native/internal/daemon"
)

func main() {
	path := os.Getenv("TOOLBOX_CONFIG")
	if path == "" {
		var err error
		if path, err = config.DefaultPath(); err != nil {
			log.Fatal(err)
		}
	}
	cfg, err := config.Load(path)
	if err != nil {
		log.Fatalf("config %s: %v", path, err)
	}

	hub := daemon.NewHub()
	a, err := app.New(cfg, hub.Publish)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()

	h, err := daemon.Handler(daemon.Options{
		Router:         a.Router,
		Hub:            hub,
		Token:          cfg.Daemon.Token,
		AllowedOrigins: cfg.Daemon.AllowedOrigins,
	})
	if err != nil {
		log.Fatalf("%v (set daemon.token in %s)", err, path)
	}

	srv := &http.Server{Addr: cfg.Daemon.Addr, Handler: h, ReadHeaderTimeout: 10 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdown)
	}()

	log.Printf("toolboxd listening on http://%s", cfg.Daemon.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
