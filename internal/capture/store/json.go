package store

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IridiumNan/GTD-simple/internal/models"
	"github.com/IridiumNan/GTD-simple/internal/utils"
)

type JSONData struct {
	Stuffs []*models.Stuff `json:"stuffs"`
}

const defaultJSONFlushInterval time.Duration = 5 * time.Second

type JSONStore struct {
	// stuffs is all stuff on RAM
	// The key is [models.Stuff.ID]
	stuffs map[string]*models.Stuff

	// fp is the store file path
	fp string

	// flushed mark if current data on RAM has been flushed on file
	// if false, [Flush] should be called when read data
	// flushed bool
	// This field is deprecated
	// you should use [JSONStore.gen] to manage and sync instead

	// change mark if stat change
	// it just consider the stat cache
	// change bool
	// this field is deprecated
	// you should use [JSONStore.gen] to manage it instead

	// cache stat
	stat *models.StuffStat

	// lock for map
	mu sync.RWMutex

	// latest write generation
	// +1 when any change on stuffs
	// NOTE:
	// remember to add Lock when read this generation
	gen atomic.Uint64

	// flushGen is the latest generation that has been flushed into disk
	// assign it as gen when flushed
	// remember to add Lock
	flushGen atomic.Uint64

	// statGen is the latest generation for cached [JSONStore.stat]
	statGen atomic.Uint64

	// interval is the auto Flush interval
	// set it before call AutoFlush Function
	interval time.Duration

	autoFlush func()

	done chan any
}

// Create is the implement for [models.Store.Create]
// it just add new stuff to map
// Overwrite the old one
func (js *JSONStore) Create(s *models.Stuff) {
	js.mu.Lock()
	defer js.mu.Unlock()

	js.stuffs[s.ID] = s

	js.gen.Add(1)
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
	js.gen.Add(1)
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
func (js *JSONStore) All(f func(s *models.Stuff) bool) []*models.Stuff {
	// WARN:
	// flush action will automaticaly exec with specific time interval
	// This code is deprecated
	//
	// err := js.Flush()
	// if err != nil {
	// 	slog.Error("error while flushing data into json file before list all", "err", err)
	// }
	js.mu.RLock()

	defer js.mu.RUnlock()

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
			js.gen.Add(1)
		}
	}
}

// Flush is a implement for [models.Store.Flush]
func (js *JSONStore) Flush() error {
	// if updated, skip
	if js.flushGen == js.gen {
		slog.Debug("all data is updated, skip flush")
		return nil
	}

	slog.Info("there are something new need to be updated, try to flush", "file", js.fp, "flushed_gen", js.flushGen.Load(), "latest_gen", js.gen.Load())
	jd, oldGen := js.prepFlush()

	data, err := json.Marshal(jd)
	if err != nil {
		slog.Error("error while marshal struct into byte data, rollback flushGen", "latest_gen", js.flushGen.Load(), "old_gen", oldGen)
		js.flushGen.Store(oldGen)
		return fmt.Errorf("error while marshal stuffs into byte data, err: %s", err.Error())
	}

	err = utils.AtomicWrite(data, js.fp)
	if err != nil {
		slog.Error("error while writing byte data into file, rollback flushGen", "latest_gen", js.flushGen.Load(), "old_gen", oldGen)
		js.flushGen.Store(oldGen)
		return fmt.Errorf("error while writing json byte data into json file, file: %s, err: %s", js.fp, err.Error())
	}

	slog.Debug("Flush: flushing data into json file success", "file", js.fp)

	return nil
}

// prepFlush get current data then update [JSONStore.flushGen]
// it get all needed status
// if fail to flush data into disk, assign oldGen to [JSONStore.flushGen]
func (js *JSONStore) prepFlush() (jd *JSONData, oldGen uint64) {
	js.mu.RLock()
	defer js.mu.RUnlock()

	// store the olad gen then update
	// if flush failed, rollback the gen
	oldGen = js.flushGen.Load()
	js.flushGen.Store(js.gen.Load())

	jd = &JSONData{Stuffs: utils.MapV(js.stuffs)}

	return
}

// Stat is implement for [models.Store.Stat]
func (js *JSONStore) Stat() *models.StuffStat {
	// return cached stat
	js.mu.RLock()
	defer js.mu.RUnlock()
	if js.statGen == js.gen && js.stat != nil {
		return js.stat
	}

	slog.Debug("Stat: stat should be update, reconstruct now")

	s := utils.MapV(js.stuffs)
	data, err := json.Marshal(s)
	if err != nil {
		return nil
	}

	// update stat, so there is no change since current time
	js.statGen.Store(js.gen.Load())

	js.stat = &models.StuffStat{
		Size:  int64(len(data)),
		Count: len(js.stuffs),
	}

	return js.stat
}

// genAutoFlush is a function that start the auto flush
func (js *JSONStore) genAutoFlush() {
	ticker := time.NewTicker(js.interval)

	slog.Info("JSONStore: starting auto flush", "interval", js.interval)
	for {
		select {

		case t := <-ticker.C:
			slog.Debug("auto flush", "time", t.String())
			err := js.Flush()
			if err != nil {
				slog.Error("error while flushing data", "err", err)
			}
		case <-js.done:
			slog.Info("stop auto flush...")
			ticker.Stop()
		}
	}
}

func (js *JSONStore) Close() error {
	js.done <- true

	return nil
}

func (js *JSONStore) AutoFlush() {
	if js.autoFlush != nil {
		js.autoFlush()
		return
	}

	js.autoFlush = sync.OnceFunc(js.genAutoFlush)
	js.autoFlush()
}

// SetInterval set the interval for auto flush
// assign this before [JSONStore.AutoFlush] else it will be [defaultJSONFlushInterval]
func (js *JSONStore) SetInterval(interval time.Duration) *JSONStore {
	if interval < 0 {
		slog.Warn("invalid negtive interval, using the default interval instead", "given", interval.String(), "default", defaultJSONFlushInterval.String())
		js.interval = defaultJSONFlushInterval
		return js
	}

	slog.Info("set the interval manually", "interval", interval.String())
	js.interval = interval

	return js
}

// NewJSONStore creating a new json store with given file path
// default flush interval is [defaultJSONFlushInterval]
func NewJSONStore(file string) *JSONStore {
	// try to load all data from given file path
	js := loadFromJSON(file)
	if js != nil {
		return js.SetInterval(defaultJSONFlushInterval)
	}

	slog.Info("load from disk failed, creating a new store", "file", file)
	return &JSONStore{
		stuffs:   make(map[string]*models.Stuff, 10),
		fp:       file,
		stat:     nil,
		interval: defaultJSONFlushInterval,
		done:     make(chan any, 1),
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
		stuffs: buildStuffMap(jd.Stuffs),
		fp:     file,
		stat:   nil,
		mu:     sync.RWMutex{},
		done:   make(chan any, 1),
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
