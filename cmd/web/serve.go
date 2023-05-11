// Package main provides ...
package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/riverchu/pkg/log"
	"github.com/riverchu/rule/biz/config"
)

func main() {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: register(gin.Default()),
	}

	go func() { // service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %s", err)
		}
	}()

	// wati shutdown signal
	<-config.Cancel()
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout())
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Error("timeout of %s.", config.ShutdownTimeout())
	}
	log.Info("Server exiting")
}
