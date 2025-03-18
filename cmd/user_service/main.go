package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jozzo6/casino_loyalty_reward_system/internal/http/users"
)

func main() {
	ctx := context.Background()
	log.Println("casino_loyalty_reward_system user api: starting server...")
	server, err := users.NewServer(ctx)
	if err != nil {
		log.Fatalf("failed to initialize api: %s", err)
	}

	sCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve failed: %s", err)
		}
	}()
	log.Println("casino_loyalty_reward_system api: server started")
	<-sCtx.Done()

	err = server.Close(ctx, 10*time.Second)
	if err != nil {
		log.Fatalf("failed to stop api: %s", err)
	}
}
