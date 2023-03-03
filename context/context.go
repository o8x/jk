package context

import (
	"context"
	"time"
)

type Context struct {
	Parent     context.Context
	Ctx        context.Context
	CancelFunc context.CancelFunc
}

func Background() context.Context {
	return context.Background()
}

func WithCancel(parent context.Context) *Context {
	ctx, cancelFunc := context.WithCancel(parent)

	return &Context{
		Parent:     parent,
		Ctx:        ctx,
		CancelFunc: cancelFunc,
	}
}

func WithTimeout(parent context.Context, timeout time.Duration) *Context {
	return WithDeadline(parent, time.Now().Add(timeout))
}

func WithDeadline(parent context.Context, deadline time.Time) *Context {
	ctx, cancelFunc := context.WithDeadline(parent, deadline)

	return &Context{
		Parent:     parent,
		Ctx:        ctx,
		CancelFunc: cancelFunc,
	}
}

func (c *Context) Wait() {
	<-c.Ctx.Done()
}

func (c *Context) Cancel() {
	c.CancelFunc()
}
