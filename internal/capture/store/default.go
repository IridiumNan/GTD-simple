package store

import (
	"log"
	"os"
	"path/filepath"

	"github.com/IridiumNan/GTD-simple/internal/models"
)

// ds is default store obj
// so that you can use this module like a package
var ds models.Store

func init() {
	df := filepath.Join(models.StoreDir, models.StoreJSONFile)
	// ensure data dir
	if err := os.MkdirAll(models.StoreDir, 0o755); err != nil {
		log.Fatalf("error while making dir for json store, dir: %s, err: %s", models.StoreDir, err.Error())
	}
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
