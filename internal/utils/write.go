package utils

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

// EnsureWrite write data into the file
// if dir not exist, it will create it automatically
// For important database, you should use [AtomicWrite] to save data instead
func EnsureWrite(data []byte, file string) error {
	dir := filepath.Dir(file)

	err := os.MkdirAll(dir, 0o755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("EnsureWrite: error while creating dir, dir: %s, err: %s", dir, err.Error())
	}

	err = os.WriteFile(file, data, 0o644)
	if err != nil {
		return fmt.Errorf("EnsureWrite: error while write into file, file: %s, err: %s", file, err.Error())
	}

	return nil
}

// AtomicWrite will create a tmp file for writing data
// then rename it to given file name
// it will creating the dir automatically
func AtomicWrite(data []byte, file string) error {
	dir := filepath.Dir(file)
	err := os.MkdirAll(dir, 0o755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("AtomicWrite: error while creating dir, dir: %s, err: %s", dir, err.Error())
	}

	tmp, err := os.CreateTemp(dir, "tmp-*")
	defer os.Remove(tmp.Name())
	if err != nil {
		return fmt.Errorf("AtomicWrite: error while creating a tmp file on target dir: %s, err: %s", dir, err.Error())
	}

	_, err = tmp.Write(data)
	if err != nil {
		return fmt.Errorf("AtomicWrite: error while writing data into tmp file, file: %s, err: %s", tmp.Name(), err.Error())
	}

	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("AtomicWrite: error while sync the data into disk for tmp file: %s, keep the legacy data", err.Error())
	}

	if err := tmp.Close(); err != nil {
		slog.Warn("AtomicWrite: error while closing the tmp file", "err", err.Error())
	}

	slog.Debug("AtomicWrite: renaming tmp file", "tmp_path", tmp.Name(), "target_path", file)
	err = os.Rename(tmp.Name(), file)
	if err != nil {
		return fmt.Errorf("AtomicWrite: error while renaming the tmp for persistent data, tmp file path: %s, err: %s", tmp.Name(), err.Error())
	}

	err = os.Chmod(file, 0o644)
	if err != nil {
		return fmt.Errorf("error while chmod 600 -> 644, err: %s", err)
	}

	return nil
}
