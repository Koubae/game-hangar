package vars

import (
	"path/filepath"
	"runtime"
)

var RootDir string
var ConfigDir string

// Get the root directory of the project
func init() {
	RootDir = getRootProjectDirectory()
	ConfigDir = filepath.Join(RootDir, "conf")
}

func getRootProjectDirectory() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	rootDir, _ := filepath.Abs(dir)
	return rootDir
}
