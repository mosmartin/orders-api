package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/mosmartin/orders-api/app"
)

func main() {
	a := app.New()

	err := a.Start(context.TODO())
	if err != nil {
		slog.Error(err.Error())
		fmt.Println("failed to start the app:", err)
	}
}
