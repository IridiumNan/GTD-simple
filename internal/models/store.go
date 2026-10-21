package models

type Store interface {
	// Create New [Stuff]
	// If this stuff id has been exist
	// it will overwrite the old one
	// so that you can update stuff with this function
	Create(s *Stuff)

	// Delete stuff by id
	// if this id not found on current data, it return false
	// else it will be marked as soft-delete
	Delete(id string) bool

	// Fetch get single stuff by it's id
	// if not found, return nil
	Fetch(id string) *Stuff

	// All return all stored stuff
	// the filter is used for filter
	// if f(s) == false, this stuff will be dropped
	All(f func(s *Stuff) bool) []*Stuff

	// Clear all stuffs that have been soft-deleted
	Clear()

	// Stat return current stuffs status
	// Count and byte size
	Stat() *StuffStat

	// Flush write all data on RAM to disk
	Flush() error

	// AutoFlush is a [sync.OnceFunc]
	// which flush all data into disk when elapse reach interval
	AutoFlush()

	// Close function will flush data automatically then stop connection
	// it's save to call this function across goroutine bcs it wrapped by a [sync.OnceFunc]
	// NOTE: This function will not flush the latest data, remeber to flush manually
	Close() error
}
