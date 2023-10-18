package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tr1v3r/pkg/log"
	"github.com/tr1v3r/pkg/shutdown"
	"github.com/tr1v3r/rule"
	"github.com/tr1v3r/rule/driver"
)

const treeName = "default"

var f rule.Forest

func Serve(timeout time.Duration, handler http.Handler) {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() { // service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %s", err)
		}
	}()

	// wati shutdown signal
	<-shutdown.Cancel()
	log.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Error("timeout of %s.", timeout)
	default:
		log.Info("work done.")
	}
	log.Info("Server exiting")
}

func InitForest(builders ...rule.TreeBuilder) { f = rule.NewForest(builders...) }

func RefreshForest() { f = f.Build() }

func DefaultBuilder(rules ...rule.Rule) rule.TreeBuilder {
	return func() rule.Tree {
		tree, err := rule.NewTree(&webDriver{PathParser: driver.SlashPathParser, Modem: driver.DummyModem},
			treeName, `{}`, rules...)
		if err != nil {
			panic(fmt.Errorf("build new tree fail: %w", err))
		}
		return tree
	}
}

type webDriver struct {
	driver.Modem
	driver.PathParser

	driver.StdCalculator
}

func (webDriver) Name() string { return "default" }
