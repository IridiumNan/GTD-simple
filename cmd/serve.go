/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
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

func handlerSig(c <-chan os.Signal, shutdownFunc func()) {
	s := <-c

	switch s {
	case syscall.SIGINT, syscall.SIGTERM:

		shutdownFunc()
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
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	handlerSig(sig, func() {
		// report current stat
		slog.Info("receiving signal sigint or sigterm, begin shutdown service graceful")
		stat := store.Stat()
		slog.Info("current store stat", "stat", stat.String())

		// shutdown the server
		timeoutSecond := 5
		slog.Info("try to shutdown server with timeout", "timeoutSecond", timeoutSecond)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSecond)*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Server shutdown err", "err", err)
		}

		// flush data into disk
		slog.Info("begin to flushing data into disk")
		err := store.Flush()
		if err != nil {
			slog.Error("error while flushing data into disk", "error", err)
		}
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
