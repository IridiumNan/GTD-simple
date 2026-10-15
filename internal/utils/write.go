package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// EnsureWrite write data into the file
// if dir not exist, it will create it automatically
func EnsureWrite(data []byte, file string) error {
	dir := filepath.Dir(file)

	err := os.MkdirAll(dir, 0o755)
	if err != nil && os.IsExist(err) {
		return fmt.Errorf("EnsureWrite: error while creating dir, dir: %s, err: %s", dir, err.Error())
	}

	err = os.WriteFile(file, data, 0o644)
	if err != nil {
		return fmt.Errorf("EnsureWrite: error while write into file, file: %s, err: %s", file, err.Error())
	}

	return nil
}
