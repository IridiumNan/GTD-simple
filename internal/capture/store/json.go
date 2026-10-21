package store

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"github.com/IridiumNan/GTD-simple/internal/models"
	"github.com/IridiumNan/GTD-simple/internal/utils"
)

type JSONData struct {
	Stuffs []*models.Stuff `json:"stuffs"`
}

type JSONStore struct {
	// stuffs is all stuff on RAM
	// The key is [models.Stuff.ID]
	stuffs map[string]*models.Stuff

	// fp is the store file path
	fp string

	// flushed mark if current data on RAM has been flushed on file
	// if false, [Flush] should be called when read data
	flushed bool

	// change mark if stat change
	// it just consider the stat cache
	change bool

	// cache stat
	stat *models.StuffStat

	// lock for map
	mu sync.RWMutex
}

// markWrite update [JSONStore.change] and [JSONStore.flushed] as data updated
// call it when any write action happend on data
func (js *JSONStore) markWrite() {
	js.change = true
	js.flushed = false
}

// Create is the implement for [models.Store.Create]
// it just add new stuff to map
// Overwrite the old one
func (js *JSONStore) Create(s *models.Stuff) {
	js.mu.Lock()
	defer js.mu.Unlock()

	js.stuffs[s.ID] = s

	js.markWrite()
}

// Delete is the implement for [models.Store.Delete]
func (js *JSONStore) Delete(id string) bool {
	js.mu.Lock()
	defer js.mu.Unlock()
	if _, ok := js.stuffs[id]; !ok {
		// It means this Stuff id not exist
		return false
	}

	js.stuffs[id].IsDel = true
	js.markWrite()
	return true
}

// Fetch is the implement for [models.Store.Fetch]
// return nil if this id not found on map
func (js *JSONStore) Fetch(id string) *models.Stuff {
	js.mu.RLock()
	defer js.mu.RUnlock()
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

	js.mu.Lock()
	defer js.mu.Unlock()

	if f == nil {
		return utils.MapV(js.stuffs)
	}

	r := make([]*models.Stuff, 0, len(js.stuffs))
	for _, s := range js.stuffs {
		if !f(s) {
			continue
		}
		r = append(r, s)
	}
	return r
}

// Clear is an implement for [models.Store.Clear]
// remove all stuffs that has been deleted
func (js *JSONStore) Clear() {
	js.mu.Lock()
	defer js.mu.Unlock()
	for id, s := range js.stuffs {
		if s.IsDel {
			delete(js.stuffs, id)
			js.markWrite()
		}
	}
}

// Flush is a implement for [models.Store.Flush]
func (js *JSONStore) Flush() error {
	js.mu.RLock()
	defer js.mu.RUnlock()
	if js.flushed {
		slog.Debug("all data is updated, skip flush request")
		return nil
	}

	slog.Debug("there are something new need to be updated, try to flush", "file", js.fp)

	jd := JSONData{Stuffs: utils.MapV(js.stuffs)}

	data, err := json.Marshal(jd)
	if err != nil {
		return fmt.Errorf("error while marshal stuffs into byte data, err: %s", err.Error())
	}

	err = utils.AtomicWrite(data, js.fp)
	if err != nil {
		return fmt.Errorf("error while writing json byte data into json file, file: %s, err: %s", js.fp, err.Error())
	}

	slog.Debug("Flush: begin flushing data into json file success", "file", js.fp)
	js.flushed = true

	return nil
}

// Stat is implement for [models.Store.Stat]
func (js *JSONStore) Stat() *models.StuffStat {
	// return cached stat
	js.mu.RLock()
	defer js.mu.RUnlock()
	if !js.change && js.stat != nil {
		return js.stat
	}

	slog.Debug("Stat: stat should be update, reconstruct now")

	s := utils.MapV(js.stuffs)
	data, err := json.Marshal(s)
	if err != nil {
		return nil
	}

	// update stat, so there is no change since current time
	js.change = false

	js.stat = &models.StuffStat{
		Size:  int64(len(data)),
		Count: len(js.stuffs),
	}

	return js.stat
}

// NewJSONStore creating a new json store with given file path
func NewJSONStore(file string) *JSONStore {
	// try to load all data from given file path
	js := loadFromJSON(file)
	if js != nil {
		return js
	}

	slog.Info("load from disk failed, creating a new store", "file", file)
	return &JSONStore{
		stuffs:  make(map[string]*models.Stuff, 10),
		fp:      file,
		flushed: true,
		change:  false,
		stat:    nil,
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
	var jd JSONData
	err = json.Unmarshal(data, &jd)
	if err != nil {
		slog.Error("error while unmarshal json data from file", "err", err)
		return nil
	}

	slog.Debug("load stuffs from disk success", "file", file, "size", len(data))

	return &JSONStore{
		stuffs:  buildStuffMap(jd.Stuffs),
		fp:      file,
		flushed: true,
		change:  false,
		stat:    nil,
		mu:      sync.RWMutex{},
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
