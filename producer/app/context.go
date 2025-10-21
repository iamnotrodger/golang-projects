package app

import "context"

type ApplicationContext struct {
}

func NewApplicationContext(ctx context.Context) *ApplicationContext {
	return &ApplicationContext{}
}
