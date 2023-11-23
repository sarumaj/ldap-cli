package util

import (
	"os"
	"path/filepath"
)

const (
	appData       = "AppData"
	localAppData  = "LocalAppData"
	xdgConfigHome = "XDG_CONFIG_HOME"
)

// Get path of the executable that started current gr process.
func GetExecutablePath() string {
	executablePath, err := os.Executable()
	if err != nil {
		return "unknown"
	}

	evaluatedPath, err := filepath.EvalSymlinks(executablePath)
	if err != nil {
		return executablePath
	}

	return evaluatedPath
}
