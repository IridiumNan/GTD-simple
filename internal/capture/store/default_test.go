package store_test

import (
	"sync"
	"testing"

	"github.com/IridiumNan/GTD-simple/internal/capture/store"
	"github.com/IridiumNan/GTD-simple/internal/models"
)

const n = 100

// TestStore_CreateAndFetch tests creating a new stuff and fetching it by ID[cite: 1, 2].
func TestStore_CreateAndFetch(t *testing.T) {
	// use a invalid path to confirm no content is loaded and RAM content not flush
	store.Reload("")
	// Clean up state if necessary
	store.Clear()

	// Create a new stuff model[cite: 1]
	stuff := models.NewStuff("This is a test stuff")
	store.Create(stuff)

	if stuff.ID == "" {
		t.Fatalf("expected stuff ID to be generated, got empty string")
	}

	// Fetch the created stuff[cite: 1, 2]
	fetched := store.Fetch(stuff.ID)
	if fetched == nil {
		t.Fatalf("expected to fetch stuff with ID %s, got nil", stuff.ID)
	}

	if fetched.Raw != "This is a test stuff" {
		t.Errorf("expected raw content 'This is a test stuff', got '%s'", fetched.Raw)
	}

	store.Delete(stuff.ID)
}

// TestStore_DeleteAndClear tests soft-deleting a stuff and clearing deleted stuffs[cite: 1, 2].
func TestStore_DeleteAndClear(t *testing.T) {
	store.Reload("")
	store.Clear()

	stuff := models.NewStuff("Stuff to delete")
	store.Create(stuff)

	// Verify stat count is 1[cite: 1, 2, 3]
	statBeforeDelete := store.Stat()
	if statBeforeDelete.Count != 1 {
		t.Errorf("expected count 1, got %d", statBeforeDelete.Count)
	}

	// Delete the stuff[cite: 1, 2]
	deleted := store.Delete(stuff.ID)
	if !deleted {
		t.Errorf("expected delete to return true for existing ID %s", stuff.ID)
	}

	// Verify trying to delete non-existent ID returns false[cite: 2]
	nonExistentDelete := store.Delete("non-existent-id")
	if nonExistentDelete {
		t.Errorf("expected delete to return false for non-existent ID")
	}

	// Clear deleted stuffs[cite: 1, 2]
	store.Clear()

	statAfterClear := store.Stat()
	if statAfterClear.Count != 0 {
		t.Errorf("expected count 0 after clear, got %d", statAfterClear.Count)
	}
}

// TestStore_AllWithFilter tests retrieving all stuffs using filters[cite: 1, 3].
func TestStore_AllWithFilter(t *testing.T) {
	store.Reload("")
	store.Clear()

	s1 := models.NewStuff("First stuff")
	s2 := models.NewStuff("Second stuff")
	store.Create(s1)
	store.Create(s2)

	// Test default filter (filters out deleted stuffs)[cite: 1, 3]
	defaultFilter, ok := store.Filter["default"]
	if !ok {
		t.Fatalf("expected 'default' filter to be registered")
	}

	stuffs := store.All(defaultFilter)
	if len(stuffs) != 2 {
		t.Errorf("expected 2 active stuffs, got %d", len(stuffs))
	}

	// Soft-delete one item
	store.Delete(s1.ID)

	// Fetch with default filter should now only return 1 item[cite: 1, 3]
	stuffsAfterDelete := store.All(defaultFilter)
	if len(stuffsAfterDelete) != 1 {
		t.Errorf("expected 1 active stuff after deletion, got %d", len(stuffsAfterDelete))
	}

	// Test 'none' filter (gets all stuffs including deleted ones)[cite: 1, 3]
	noneFilter, ok := store.Filter["none"]
	if !ok {
		t.Fatalf("expected 'none' filter to be registered")
	}

	allStuffs := store.All(noneFilter)
	if len(allStuffs) != 2 {
		t.Errorf("expected 2 total stuffs with 'none' filter, got %d", len(allStuffs))
	}
}

// TestStore_Stat verifies that the store returns valid stat metrics[cite: 1, 2, 3].
func TestStore_Stat(t *testing.T) {
	store.Reload("")
	store.Clear()

	stuff := models.NewStuff("Stat test content")
	store.Create(stuff)

	stat := store.Stat()
	if stat == nil {
		t.Fatalf("expected stat not to be nil")
	}

	if stat.Count <= 0 {
		t.Errorf("expected stat count to be greater than 0, got %d", stat.Count)
	}

	if stat.Size <= 0 {
		t.Errorf("expected stat size to be greater than 0, got %d", stat.Size)
	}
}

func TestStoreRace(t *testing.T) {
	// use a invalid path to confirm no content is loaded and RAM content not flush
	store.Reload("")
	// Clean up state if necessary
	store.Clear()

	stuffs := make([]*models.Stuff, 0, n)

	for range n {
		stuffs = append(stuffs, models.NewStuff("stuff race"))
	}

	ready := make(chan struct{}, n)
	start := make(chan struct{})

	var wg sync.WaitGroup

	for i := range n {
		wg.Go(func() {
			stf := stuffs[i]

			ready <- struct{}{}
			<-start

			store.Create(stf)
		})
	}

	for range n {
		<-ready
	}
	close(start)
	wg.Wait()

	stat := store.Stat()
	if stat.Count != n {
		t.Errorf("error while testing concurrent create, want: %d, got: %d", n, stat.Count)
	}

	ready = make(chan struct{}, n)
	start = make(chan struct{})

	for i := range n {
		wg.Go(func() {
			stf := stuffs[i]
			ready <- struct{}{}
			<-start
			store.Delete(stf.ID)
		})
	}

	for range n {
		<-ready
	}
	close(start)
	wg.Wait()

	stuffs = store.All(func(s *models.Stuff) bool {
		return !s.IsDel
	})
	if len(stuffs) != 0 {
		t.Errorf("error while testinig deleting all stuffs, want: 0, got: %d", stat.Count)
	}
}

func TestMultiClose(t *testing.T) {
	store.Reload("")
	start := make(chan struct{})
	ready := make(chan struct{}, n)

	var wg sync.WaitGroup

	for range n {
		wg.Go(func() {
			ready <- struct{}{}
			<-start

			store.Close()
		})
	}

	for range n {
		<-ready
	}
	close(start)
	wg.Wait()
}
