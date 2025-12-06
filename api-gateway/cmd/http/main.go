package main

import (
	"context"
	"log"

	"github.com/cg-2025-crutch/backend/api-gateway/internal/run"
)

func main() {
	ctx := context.Background()

	if err := run.Run(ctx); err != nil {
		log.Fatalf("application stopped with error: %v", err)
	}
}
