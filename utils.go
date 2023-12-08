package hikaku

import (
	"context"
	"reflect"
	"strings"
)

func valuesValid(v ...reflect.Value) bool {
	for i := 0; i < len(v); i++ {
		if !v[i].IsValid() {
			return false
		}
	}
	return true
}

func getNameField(field reflect.StructField, cfg DeepDifferenceConfig) (string, error) {
	// sorry i can't accept that
	if !field.IsExported() {
		return "", ErrFieldNotExported
	}
	value := field.Tag.Get(cfg.tag)
	if value != "" {
		if strings.Contains(value, ",") {
			return strings.Split(value, ",")[0], nil
		}
		return value, nil
	}
	// if there is not json, that's fine, we will take the name of the field
	value = field.Type.Name()
	if value != "" {
		return value, nil
	}
	// well you're not helping you know
	value = field.Name
	if value != "" {
		return value, nil
	}
	// what? you're still there?
	return "", ErrFieldHasNoName
}

var keyExecutionBufferCtx string = "buffer"
var keyDiffCtx string = "differences"
var keyAttributeMapCtx string = "mapPathPerProperty"

func checkInitContext(ctx context.Context) context.Context {
	if !has[differenceContext](ctx, keyDiffCtx) {
		ctx = set[differenceContext](ctx, keyDiffCtx, newDifferenceContext())
	}
	if !has[executionBuffer](ctx, keyExecutionBufferCtx) {
		ctx = set[executionBuffer](ctx, keyExecutionBufferCtx, newExecutionBuffer())
	}
	if !has[AttributeMap](ctx, keyAttributeMapCtx) {
		ctx = set[AttributeMap](ctx, keyAttributeMapCtx, newAttributeMap())
	}
	return ctx
}

func getExecutionBuffer(ctx context.Context) (*executionBuffer, error) {
	return get[executionBuffer](ctx, keyExecutionBufferCtx)
}

func setExecutionBuffer(ctx context.Context, exe *executionBuffer) context.Context {
	return set[executionBuffer](ctx, keyExecutionBufferCtx, exe)
}

func getAttributeMap(ctx context.Context) (*AttributeMap, error) {
	return get[AttributeMap](ctx, keyAttributeMapCtx)
}

func setAttributeMap(ctx context.Context, exe *AttributeMap) context.Context {
	return set[AttributeMap](ctx, keyAttributeMapCtx, exe)
}
