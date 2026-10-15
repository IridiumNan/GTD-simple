package config_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/IridiumNan/GTD-simple/internal/capture/config"
	"github.com/IridiumNan/GTD-simple/internal/models"
	"github.com/spf13/viper"
)

var testConf = `
[service.front]
# dir is the default front-end file path
dir = "./frontend"
addr = "127.0.0.1:8001"

[service.backend]
addr = "127.0.0.1:8081"

[store]
dir = "storage"
file = "data.toml"
`

const testDir = "config_test"

var globalTest *testing.T

func TestWithWrite(t *testing.T) {
	setTest(t)
	viper.SetConfigType("toml")

	viper.ReadConfig(bytes.NewReader([]byte(testConf)))

	check("./frontend", viper.GetString("service.front.dir"))
	check("127.0.0.1:8001", viper.GetString("service.front.addr"))
	check("127.0.0.1:8081", viper.GetString("service.backend.addr"))
	check("storage", viper.GetString("store.dir"))
	check("data.toml", viper.GetString("store.file"))
}

func TestInitConfig(t *testing.T) {
	setTest(t)

	err := os.MkdirAll(testDir, 0o755)
	if err != nil {
		t.Errorf("error while making dir for test, dir: %s, err: %s", testDir, err.Error())
	}
	testFile := filepath.Join(testDir, models.ConfigFileName)

	err = os.WriteFile(testFile, []byte(testConf), 0o644)
	if err != nil {
		t.Errorf("error while writing test configuration, file: %s, err: %s", testFile, err.Error())
	}

	err = config.LoadConfigFromDir([]string{testDir})
	if err != nil {
		t.Errorf("error while loading from test dir, file: %s, err: %s", testFile, err.Error())
	}

	check("./frontend", viper.GetString("service.front.dir"))
	check("127.0.0.1:8001", viper.GetString("service.front.addr"))
	check("127.0.0.1:8081", viper.GetString("service.backend.addr"))
	check("storage", viper.GetString("store.dir"))
	check("data.toml", viper.GetString("store.file"))
}

func setTest(t *testing.T) {
	globalTest = t
}

func check[T comparable](want, got T) {
	if want != got {
		globalTest.Errorf("configuration compare error, want: %v, got: %v", want, got)
	}
}
