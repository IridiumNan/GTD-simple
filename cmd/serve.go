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
	srv := &http.Server{
		Addr:    addr,
		Handler: r.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error while listen", "err", err)
			return
		}
	}()

	handlerSig(sig, func() error {
		var errMsg strings.Builder

		slog.Info("receiving signal sigint or sigterm, begin shutdown service graceful")

		// shutdown the server
		timeoutSecond := 5
		slog.Info("try to shutdown server with timeout", "timeoutSecond", timeoutSecond)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecond)*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Server shutdown err", "err", err)
			fmt.Fprintf(&errMsg, "error while server shutdown, err: %s", err)
		}

		// report current stat
		statDone := make(chan error, 1)
		go func() {
			stat := store.Stat()
			slog.Info("current store stat", "stat", stat.String())
			statDone <- nil
		}()
		select {
		case <-statDone:
			slog.Info("stat query done success")
		case <-time.After(3 * time.Second):
			slog.Error("stat query timeout, continue flush data to disk")
		}

		// flush data into disk
		flushDone := make(chan error, 1)
		go func() {
			slog.Info("begin to flushing data into disk")
			flushDone <- store.Flush()
		}()
		select {
		case err := <-flushDone:
			if err != nil {
				slog.Error("error while flushing data into disk", "error", err)
				fmt.Fprintf(&errMsg, "error while flushing data into disk, err: %s", err)
			}
		case <-time.After(3 * time.Second):
			slog.Error("flush data timeout, quit now, RAM data may lost")
		}

		if errMsg.Len() != 0 {
			return errors.New(errMsg.String())
		}

		return nil
	})
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
