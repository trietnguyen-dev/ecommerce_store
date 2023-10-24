package utils

import (
	"context"
	"time"
)

func NewCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 6*time.Second)
}
