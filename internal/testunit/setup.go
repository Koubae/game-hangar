package testunit

import "github.com/koubae/game-hangar/pkg/common"

func Setup() {
	common.CreateLogger("ERROR", "/tmp/")
}
