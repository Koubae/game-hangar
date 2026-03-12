package database

import "context"

type Connector interface {
	String() string
	Ping(ctx context.Context) error
	Shutdown() error
}
