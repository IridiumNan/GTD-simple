package store

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/IridiumNan/GTD-simple/internal/models"
)

type JSONStore struct {
	// stuffs is all stuff on RAM
	// The key is [models.Stuff.ID]
	stuffs map[string]*models.Stuff

	// fp is the store file path
	fp string
}

// Create is the implement for [models.Store.Create]
// it just add new stuff to map
func (js *JSONStore) Create(s *models.Stuff) {
	js.stuffs[s.ID] = s
}

// Delete is the implement for [models.Store.Delete]
func (js *JSONStore) Delete(id string) bool {
	if _, ok := js.stuffs[id]; !ok {
		// It means this Stuff id not exist
		return false
	}

	js.stuffs[id].IsDel = true
	return true
}

// Fetch is the implement for [models.Store.Fetch]
// return nil if this id not found on map
func (js *JSONStore) Fetch(id string) *models.Stuff {
	s, ok := js.stuffs[id]
	if !ok {
		return nil
	}
	return s
}

// All is the implement for [models.Store.All]
// every time this function is called, all data will be Flushed into json file
func (js *JSONStore) All(f func(s *models.Stuff) bool) []*models.Stuff {
	err := js.Flush()
	if err != nil {
		slog.Error("error while flushing data into json file before list all", "err", err)
	}

	allStuffs := make([]*models.Stuff, 0, len(js.stuffs))
	if f == nil {
		return allStuffs
	}

	r := make([]*models.Stuff, 0, len(js.stuffs))
	for _, s := range allStuffs {
		if !f(s) {
			continue
		}
		r = append(r, s)
	}

	return r
}

// Flush is a implement for [models.Store.Flush]
func (js *JSONStore) Flush() error {
	data, err := json.Marshal(js.stuffs)
	if err != nil {
		return fmt.Errorf("error while marshal stuffs into byte data, err: %s", err.Error())
	}

	dir := filepath.Dir(js.fp)
	err = os.MkdirAll(dir, 0o755)
	if err != nil {
		return fmt.Errorf("error while creaing dir for json data file, err: %s", err.Error())
	}

	err = os.WriteFile(js.fp, data, 0o644)
	if err != nil {
		return fmt.Errorf("error while writing json byte data into json file, file: %s, err: %s", js.fp, err.Error())
	}

	return nil
}

// NewJSONStore creating a new json store with given file path
func NewJSONStore(file string) *JSONStore {
	// try to load all data from given file path
	js := loadFromJSON(file)
	if js != nil {
		return js
	}
	return &JSONStore{
		stuffs: make(map[string]*models.Stuff, 10),
		fp:     file,
	}
}

// loadFromJSON write all error or log by slog
// not return error
func loadFromJSON(file string) *JSONStore {
	data, err := os.ReadFile(file)
	if err != nil {
		slog.Warn("failed to load data from json file", "err", err)
		return nil
	}
	stuffs := make([]*models.Stuff, 0, 10)
	err = json.Unmarshal(data, &stuffs)
	if err != nil {
		slog.Error("error while unmarshal json data from file", "err", err)
		return nil
	}

	return &JSONStore{
		stuffs: buildStuffMap(stuffs),
		fp:     file,
	}
}

// buildStuffMap build a map for key -> ID and value -> *[models.Stuff]
func buildStuffMap(s []*models.Stuff) map[string]*models.Stuff {
	r := make(map[string]*models.Stuff, len(s))
	for _, v := range s {
		r[v.ID] = v
	}

	return r
}

// djs is default json store obj
// so that you can use this module like a package
var djs models.Store

func init() {
	df := filepath.Join(models.StoreDir, models.StoreJSONFile)
	// ensure data dir
	if err := os.MkdirAll(models.StoreDir, 0o755); err != nil {
		log.Fatalf("error while making dir for json store, dir: %s, err: %s", models.StoreDir, err.Error())
	}
	djs = NewJSONStore(df)
}

// Create call default json store Create function
func Create(s *models.Stuff) {
	djs.Create(s)
}

// Delete call default json store Delete function
func Delete(id string) bool {
	return djs.Delete(id)
}

// Fetch call default json store Fetch function
func Fetch(id string) *models.Stuff {
	return djs.Fetch(id)
}

// All call default json store All function
func All(f func(s *models.Stuff) bool) []*models.Stuff {
	return djs.All(f)
}

// Flush call default json store Flush function
func Flush() error {
	return djs.Flush()
}
