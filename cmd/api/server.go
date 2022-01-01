package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	// create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// used to receive shutdown errors from Shutdown()
	shutdownError := make(chan error)

	go func() {
		// create a quit channel that carries os.Signal
		quit := make(chan os.Signal, 1)

		// listens for SIGINT and SIGTERM signals and relay to quit channel
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		// reads signal from quit channel;
		// blocking code until signal is read
		s := <-quit

		// logs the caught signal
		app.logger.PrintInfo("shutting down the server", map[string]string{
			"signal": s.String(),
		})

		// context with 5-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Shutdown() returns nil if everything is fine, or error if there
		// was a problem; relayed to shutdownError
		shutdownError <- srv.Shutdown(ctx)
	}()

	// log the start
	app.logger.PrintInfo("starting server", map[string]string{
		"addr": srv.Addr,
		"env":  app.config.env,
	})

	// will return err if it's not http.ErrServerClosed
	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// if there was a problem with the graceful shutdown
	// return the error
	err = <-shutdownError
	if err != nil {
		return err
	}

	// graceful shutdown completed, log message
	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": srv.Addr,
	})

	return nil
}
