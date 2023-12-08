package hikaku

import "context"

func has[T any](ctx context.Context, key string) bool {
	data := ctx.Value(key)
	if _, ok := data.(T); ok {
		return ok
	}
	return false
}

func get[T any](ctx context.Context, key string) (*T, error) {
	data := ctx.Value(key)
	if val, ok := data.(*T); ok {
		return val, nil
	}
	return nil, ErrContextValueNotFound
}

func set[T any](ctx context.Context, key string, data *T) context.Context {
	ctx = context.WithValue(ctx, key, data)
	return ctx
}

func addCheckProbeMap(ctx context.Context) context.Context {
	if !has[ProbeMap](ctx, keyProbeMapCtx) {
		ctx = set[ProbeMap](ctx, keyProbeMapCtx, newProbeMap())
	}
	return ctx
}

func addCheckExecutionBuffer(ctx context.Context) context.Context {
	if !has[executionBuffer](ctx, keyExecutionBufferCtx) {
		ctx = set[executionBuffer](ctx, keyExecutionBufferCtx, newExecutionBuffer())
	}
	return ctx
}

func getProbeMap(ctx context.Context) (*ProbeMap, error) {
	return get[ProbeMap](ctx, keyProbeMapCtx)
}
