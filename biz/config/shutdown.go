package config

import (
	"context"
	"os/signal"
	"syscall"
)

var ctx, stop = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
