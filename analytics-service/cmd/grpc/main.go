package main

import (
	"context"
	"log"

	"github.com/cg-2025-crutch/backend/analytics-service/internal/run"
)

func main() {
	ctx := context.Background()
	if err := run.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
