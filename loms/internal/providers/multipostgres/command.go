package multipostgres

import "context"

type command1[T any] struct {
	provider T
	f        func(context.Context, T) bool
}

type command2[T1 any, T2 any] struct {
	provider1 T1
	provider2 T2
	f         func(context.Context, T1, T2) bool
}

type command3[T1 any, T2 any, T3 any] struct {
	provider1 T1
	provider2 T2
	provider3 T3
	f         func(context.Context, T1, T2, T3) bool
}

func (c *command1[T]) execute(ctx context.Context) bool {
	return c.f(ctx, c.provider)
}

func (c *command2[T1, T2]) execute(ctx context.Context) bool {
	return c.f(ctx, c.provider1, c.provider2)
}

func (c *command3[T1, T2, T3]) execute(ctx context.Context) bool {
	return c.f(ctx, c.provider1, c.provider2, c.provider3)
}
