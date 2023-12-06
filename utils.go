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
var keyAttributeMapCtx string = "mapPathPerProperty"

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

// Receive a value of kind struct
func handleStruct(ctx context.Context, value reflect.Value, opts *valueOptions) error {
	var exe *executionBuffer
	var attrs *AttributeMap

	var err error
	if exe, err = getExecutionBuffer(ctx); err != nil {
		fmt.Println("can't get execution buffer")
		return err
	}

	if attrs, err = getAttributeMap(ctx); err != nil {
		fmt.Println("can't get attrs map")
		return err
	}

	thisStruct := attrs.Add(opts.parent, value)

	fmt.Println(*thisStruct, thisStruct.Name, thisStruct.ParentPath, thisStruct.Path, thisStruct.Tag)

	for idx := 0; idx < value.NumField(); idx++ {

		fieldValue := value.Field(idx)

		// if varType != reflect.Slice && varType != reflect.Struct {
		// 	// valueCompose[varName] = varValue
		// }

		varName := value.Type().Field(idx).Name
		varType := value.Type().Field(idx).Type.Kind()
		varValue := value.Field(idx).Interface()
		fmt.Println("properties")

		exe.Add(func() error {
			fmt.Println(varName, varType, varValue)
			// return switchValue(ctx, fieldValue, fromParent(fmt.Sprintf("%v.%v.%v", opts.path, value.Type().Name(), varName)))
			return switchValue(ctx, fieldValue, fromParent(thisStruct.Path))
		})
		// hashStructNames[varName] = strings.Split(valueVal.Type().Field(idx).Tag.Get("json"), ",")[0]
	}

	// ctx = setExecutionBuffer(ctx, exe)

	return nil
}

func handleString(ctx context.Context, value reflect.Value, opts *valueOptions) error {
	var exe *executionBuffer
	var err error
	if exe, err = getExecutionBuffer(ctx); err != nil {
		fmt.Println("can't get execution buffer")
		return err
	}

	exe.Add(func() error {
		fmt.Println("add that string", value.Type().Name(), value.Interface())
		return nil
	})

	return nil
}

func handleArray(ctx context.Context, value reflect.Value, opts *valueOptions) error {

	return nil
}
