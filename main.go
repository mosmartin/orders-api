package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"

	"github.com/mosmartin/orders-api/app"
)

func main() {
	a := app.New()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := a.Start(ctx)
	if err != nil {
		slog.Error(err.Error())
		fmt.Println("failed to start the app:", err)
	}
}
