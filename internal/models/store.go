package models

type Store interface {
	// Create New [Stuff]
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

	// Flush write all data on RAM to disk
	Flush() error
}
