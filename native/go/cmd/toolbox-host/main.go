// toolbox-host is the native-messaging host launched by the browser.
// It must never write anything but frames to stdout.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"example.com/toolbox/native/internal/app"
	"example.com/toolbox/native/internal/config"
	"example.com/toolbox/native/internal/host"
)

func main() {
	log.SetOutput(os.Stderr)
	log.SetPrefix("toolbox-host: ")

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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	srv := host.New(os.Stdout)
	a, err := app.New(cfg, srv.Emit)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()
	srv.SetRouter(a.Router)

	if err := srv.Serve(ctx, os.Stdin); err != nil {
		log.Fatal(err)
	}
}
