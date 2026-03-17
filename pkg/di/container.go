package di

import (
	"github.com/koubae/game-hangar/pkg/common"
)

type Container interface {
	Logger() common.Logger
	Shutdown() error
}
