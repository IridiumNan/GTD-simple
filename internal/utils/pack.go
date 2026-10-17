package utils

import (
	"bytes"
	_ "embed"
	"log/slog"
	"os"
	"os/exec"
)

//go:embed front.tar.xz
var frontXZPack []byte

// UnpackFront to current dir
func UnpackFront() error {
	slog.Info("unpacking front package")
	unPackCmd := exec.Command("tar", "-xJvf", "-")

	unPackCmd.Stdin = bytes.NewReader(frontXZPack)
	unPackCmd.Stdout = os.Stdout
	unPackCmd.Stderr = os.Stdout

	return unPackCmd.Run()
}
