package hikaku

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

func convertPtrToValue[T any](a *T) reflect.Value {
	return reflect.ValueOf(a)
}

func valuesValid(v ...reflect.Value) bool {
	for i := 0; i < len(v); i++ {
		if !v[i].IsValid() {
			return false
		}
	}
	return true
}

func getNameField(field reflect.StructField, cfg config) (string, error) {
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

func checkInitContext(ctx context.Context) context.Context {
	if !has[differenceContext](ctx, keyDiffCtx) {
		ctx = set[differenceContext](ctx, keyDiffCtx, newDifferenceContext())
	}
	if !has[executionBuffer](ctx, keyExecutionBufferCtx) {
		ctx = set[executionBuffer](ctx, keyExecutionBufferCtx, newExecutionBuffer())
	}
	return ctx
}

func getExecutionBuffer(ctx context.Context) (*executionBuffer, error) {
	return get[executionBuffer](ctx, keyExecutionBufferCtx)
}

func setExecutionBuffer(ctx context.Context, exe *executionBuffer) context.Context {
	return set[executionBuffer](ctx, keyExecutionBufferCtx, exe)
}

// Receive a value of kind struct
func handleStruct(ctx context.Context, value reflect.Value, opts *valueOptions) error {
	var exe *executionBuffer
	var err error
	if exe, err = getExecutionBuffer(ctx); err != nil {
		fmt.Println("can't get execution buffer")
		return err
	}

	for idx := 0; idx < value.NumField(); idx++ {

		fieldValue := value.Field(idx)
		// varName := valueVal.Type().Field(idx).Name
		// varType := valueVal.Type().Field(idx).Type.Kind()
		// varValue := valueVal.Field(idx).Interface()
		// if varType != reflect.Slice && varType != reflect.Struct {
		// 	// valueCompose[varName] = varValue
		// }
		fmt.Println("properties")
		exe.Add(func() error {
			return switchValue(ctx, fieldValue)
		})
		// hashStructNames[varName] = strings.Split(valueVal.Type().Field(idx).Tag.Get("json"), ",")[0]
	}

	// ctx = setExecutionBuffer(ctx, exe)

	return nil
}
