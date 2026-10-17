/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/IridiumNan/GTD-simple/internal/utils"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init command unpack the front-end dir",
	Long: `init command will unpack the build-in front-end tar.xz package
	default unpacked dir is ./front`,
	Run: runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	err := os.RemoveAll("./front")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("error while removing the legacy front dir: %s", err.Error())
	}
	err = utils.UnpackFront()
	if err != nil {
		log.Fatalf("error while unpacking front dir, err: %s", err.Error())
	}
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
