package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/IridiumNan/GTD-simple/internal/models"
	"github.com/IridiumNan/GTD-simple/internal/utils"
	"github.com/spf13/viper"
)

//go:embed default.toml
var defaultConfByte []byte

var searchDir = []string{
	".",
	models.ConfigDir(),
}

func LoadConfigFromDir(dirs []string) error {
	// TEST:
	// slog.SetLogLoggerLevel(slog.LevelDebug)
	// viper.SetOptions(viper.WithLogger(slog.Default()))

	if len(dirs) == 0 {
		slog.Warn("paramter dir is empty, using default configuration search dir")
		return LoadConfig()
	}

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	// viper.SetConfigFile(models.ConfigFileName)

	setDefault()

	for _, dir := range dirs {
		viper.AddConfigPath(dir)
	}

	err := viper.ReadInConfig()
	if err != nil {
		slog.Warn("no configuration file found on search path, try to creating a new configuration file")
		err = initConfig(dirs[0])
		if err != nil {
			return err
		}

		err = viper.ReadConfig(bytes.NewReader(defaultConfByte))
		if err != nil {
			return fmt.Errorf("error while reading config from default byte conf, err: %s", err.Error())
		}

	}
	slog.Info("loading configuration from file success")

	return nil
}

// LoadConfig from [searchDir]
// if configuration file not found, it will initailize a new one
// This function is called on [store] module, so there is no need to load on cmd
func LoadConfig() error {
	return LoadConfigFromDir(searchDir)
}

func setDefault() {
	viper.SetDefault("service.addr", "127.0.0.1:8000")
	viper.SetDefault("store.dir", "store")
	viper.SetDefault("store.file", "data.json")
}

func initConfig(dir string) error {
	// write the [defaultConfByte] into searchDir

	fp := filepath.Join(dir, models.ConfigFileName)
	slog.Info("try to write config into file", "file", fp)
	err := utils.EnsureWrite(defaultConfByte, fp)
	if err != nil {
		return fmt.Errorf("error while init config, file path: %s, err: %s", fp, err.Error())
	}

	return nil
}
