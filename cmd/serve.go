/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/IridiumNan/GTD-simple/cmd/web"
	"github.com/IridiumNan/GTD-simple/internal/capture/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the backend serve",
	Long: `It will start a backend serve for GTD-simple
	if configuration file not found, it will create a new one.
	default configuration path is ./config.toml`,
	Run: runServe,
}

func handlerSig(c <-chan os.Signal, shutdownFunc func() error) {
	s := <-c

	switch s {
	case syscall.SIGINT, syscall.SIGTERM:

		err := shutdownFunc()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	slog.Error("signal can't be handled", "signal", s)
}

func runServe(cmd *cobra.Command, args []string) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	r := web.Router()
	addr := viper.GetString("service.addr")
	timeout := 5 * time.Second
	srv := &http.Server{
		Addr:              addr,
		Handler:           r.Handler(),
		ReadTimeout:       timeout,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
		IdleTimeout:       6 * timeout,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error while listen", "err", err)
			return
		}
	}()

	slog.Info("starting auto flushing data into disk")
	go store.AutoFlush()

	// a closure function for handle the [syscall.SIGINT] and [syscall.SIGTERM] signal
	// for graceful shutdown
	handlerSig(sig, func() error {
		var errMsg strings.Builder
		errMsg.Grow(1 << 6)

		slog.Info("receiving signal sigint or sigterm, begin shutdown service graceful")

		// shutdown the http server with timeout 5s
		timeoutSecond := 5
		slog.Info("try to shutdown server with timeout", "timeoutSecond", timeoutSecond)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecond)*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Server shutdown err", "err", err)
			fmt.Fprintf(&errMsg, "error while server shutdown, err: %s\n", err)
		}

		// report latest stuff stat
		err := runWithTimeout(3*time.Second, func() {
			stat := store.Stat()
			if stat == nil {
				fmt.Fprintf(&errMsg, "invalid stat: %v", stat)
				return
			}
			slog.Info("current stat", "stat", stat.String)
		})
		if err != nil {
			fmt.Fprintln(&errMsg, err.Error())
		}

		// Close store
		err = store.Close()
		if err != nil {
			fmt.Fprintf(&errMsg, "error while closing store: %s\n", err)
		}

		// flushing latest data into disk
		err = runWithTimeout(3*time.Second, func() {
			slog.Info("starting flusing data before program exiting")
			err := store.Flush()
			if err != nil {
				slog.Error("error while flushing before exit", "err", err)
			}
		})
		if err != nil {
			fmt.Fprintf(&errMsg, "error while flushing latest data, data may lost err: %s", err)
		}

		if errMsg.Len() != 0 {
			return errors.New(errMsg.String())
		}

		return nil
	})
}

// runWithTimeout run this function and wait for wait time
func runWithTimeout(waitTime time.Duration, f func()) error {
	done := make(chan int, 1)
	go func() {
		f()
		done <- 0
	}()

	select {
	case <-done:
		return nil
	case <-time.After(waitTime):
		return errors.New("func exec time out, force quit")
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
