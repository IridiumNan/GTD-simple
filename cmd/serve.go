/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

func handlerSig(c <-chan os.Signal) {
	s := <-c

	switch s {
	case syscall.SIGINT, syscall.SIGTERM:
		slog.Info("receiving signal sigint or sigterm, flushing data into disk")
		stat := store.Stat()
		slog.Info("current store stat", "stat", stat.String())

		err := store.Flush()
		if err != nil {
			slog.Error("error while flushing data into disk", "error", err)
		}

		os.Exit(0)
	}

	slog.Error("signal can't be handled", "signal", s)
}

func runServe(cmd *cobra.Command, args []string) {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go handlerSig(sig)

	r := web.Router()
	addr := viper.GetString("service.addr")
	err := r.Run(addr)
	if err != nil {
		log.Fatalf("error while starting serve, err: %s", err.Error())
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
