package store

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/IridiumNan/GTD-simple/internal/capture/config"
	"github.com/IridiumNan/GTD-simple/internal/models"
	"github.com/spf13/viper"
)

// ds is default store obj
// so that you can use this module like a package
var ds models.Store

func init() {
	// set level debug
	slog.SetLogLoggerLevel(slog.LevelDebug)

	// NOTE: Must call the config module load function first
	err := config.LoadConfig()
	if err != nil {
		slog.Error("while loading configuraion", "err", err)
	}

	dir := viper.GetString("store.dir")
	if dir == models.NonePath {
		slog.Warn("failed to get store dir from configuration, use default one")
		dir = models.StoreDir
	}
	file := viper.GetString("store.file")
	if file == "" {
		slog.Warn("failed to get store file name from configuration, use default one")
		dir = models.StoreJSONFile
	}
	df := filepath.Join(dir, file)
	// ensure data dir
	if err := os.MkdirAll(dir, 0o755); err != nil {
		log.Fatalf("error while making dir for json store, dir: %s, err: %s", models.StoreDir, err.Error())
	}

	slog.Debug("try to load stuff data from disk", "file", df)
	ds = NewJSONStore(df)
}

// Reload reload the data from another file
func Reload(file string) {
	ds = NewJSONStore(file)
}

// Create call default json store Create function
func Create(s *models.Stuff) {
	ds.Create(s)
}

// Delete call default json store Delete function
func Delete(id string) bool {
	return ds.Delete(id)
}

// Fetch call default json store Fetch function
func Fetch(id string) *models.Stuff {
	return ds.Fetch(id)
}

// All call default json store All function
// For filter all deleted stuffs, use [FilterNotDeleted]
func All(f func(s *models.Stuff) bool) []*models.Stuff {
	return ds.All(f)
}

// Flush call default json store Flush function
func Flush() error {
	return ds.Flush()
}

// Clear call default json store Clear function
func Clear() {
	ds.Clear()
}

// Stat call default json store Stat function
func Stat() *models.StuffStat {
	return ds.Stat()
}

// AutoFlush call default store AutoFlush function
func AutoFlush() {
	ds.AutoFlush()
}

// Close call default store Close funtion
func Close() error {
	return ds.Close()
}
