package mongodb

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type Timestamps struct {
	Created *time.Time `bson:"created,omitempty"`
	Updated *time.Time `bson:"updated,omitempty"`
}

func (model *Timestamps) OnCreate() {
	now := time.Now()
	model.Created = &now
	model.Updated = &now
}

func (model *Timestamps) OnUpdate() {
	now := time.Now()
	model.Updated = &now
}

type EntityID struct {
	ID *primitive.ObjectID `bson:"_id,omitempty"`
}

func (model *EntityID) OnCreated(insert *mongo.InsertOneResult) {
	if oid, ok := insert.InsertedID.(primitive.ObjectID); ok {
		model.ID = &oid
	}
}
