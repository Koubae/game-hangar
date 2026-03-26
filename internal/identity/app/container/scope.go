package container

import "github.com/koubae/game-hangar/pkg/database"

type Scope struct {
	c  *AppContainer
	db database.DBTX
}
