package models

import (
	"github.com/koubae/game-hangar/account/pkg/database/mongodb"
)

type Account struct {
	mongodb.EntityID   `bson:",inline"`
	Username           string `bson:"username"`
	Password           string `bson:"password"`
	mongodb.Timestamps `bson:",inline"`
}
