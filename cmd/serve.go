/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"log/slog"

	"github.com/IridiumNan/GTD-simple/cmd/web"
	"github.com/gin-gonic/gin"
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

func runServe(cmd *cobra.Command, args []string) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	go func() {
		backAdrr := viper.GetString("service.backend.addr")
		br := web.BackRouter()
		err := br.Run(backAdrr)
		if err != nil {
			log.Fatalf("error while running backend router, err: %s", err.Error())
		}
	}()

	frontAddr := viper.GetString("service.front.addr")
	fr := gin.Default()
	frontRoot := viper.GetString("service.front.dir")

	fr.Static("/", frontRoot)

	err := fr.Run(frontAddr)
	if err != nil {
		log.Fatalf("error while serving front end static file, root: %s, addr: %s, err: %s", frontRoot, frontAddr, err.Error())
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
