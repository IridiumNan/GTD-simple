package store

import "github.com/IridiumNan/GTD-simple/internal/models"

var Filter = map[string]func(s *models.Stuff) bool{
	"default": filterNotDeleted,
	"none":    nil,
}

// filterNotDeleted filter all deleted stuff
func filterNotDeleted(s *models.Stuff) bool {
	return !s.IsDel
}
