package vars

import (
	"path/filepath"
	"runtime"
)

var RootDir string

// Get the root directory of the project
func init() {
	RootDir = getRootProjectDirectory()
}

func getRootProjectDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	rootDir, _ := filepath.Abs(dir)
	return rootDir
}
