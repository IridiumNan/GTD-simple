package models

import (
	"os"
	"path"
)

const (
	APPName       = "gtd-cli"
	StoreDir      = "store"
	StoreJSONFile = "data.json"

	ConfigFileName = "config.toml"

	NonePath = ""
)

// ConfigDir return this app configuration dir
// for linux, it will be ~/.config/[APPName]
func ConfigDir() string {
	dir, err := os.UserCacheDir()
	if err != nil {
		return NonePath
	}

	return path.Join(dir, APPName)
}
