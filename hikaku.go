package hikaku

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/k0kubun/pp/v3"
)

/// TODO: i need a way to store the mapping and the composite mapping of the structs
/// TODO: define a way to represent the tree of data
/// it has to be generic enough to fit other languages aka be universal (i.e. Block, Property, etc)
/// TODO: make more todos for each type of data that need to be processed
/// TODO: add good logging (should I upgrade my `go`? could use the latest standard package for that)

// Analyze type deeply to blocks
func AnalyzeType[T any]() {

}

// Analyze value deeply to blocks
func Analyze[T any](data *T) {

}

// note: the problem with recursive, is the amount of stack we will produce, so i suggest i should have an array of functions
// to process it over and over
// Deep difference to block of differences
func DeepDifference[T any](ctx context.Context, a *T, b *T, opts ...optsConfig) error {
	now := time.Now()

	var err error
	var attrsA *AttributeMap
	var attrsB *AttributeMap

	valueA := convertPtrToValue[T](a)
	valueB := convertPtrToValue[T](b)

	if !valuesValid(valueA, valueB) {
		return ErrValuesInvalid
	}

	sigA := make(chan error)
	sigB := make(chan error)

	ctxA := checkInitContext(context.TODO())
	ctxB := checkInitContext(context.TODO())

	// traverse the whole structs to create functions that will be processed later on
	// those two functions are starters of the worker that process the buffers
	switchValue(ctxA, valueA, fromParentPath("."))
	switchValue(ctxB, valueB, fromParentPath("."))

	fmt.Println("finished traversing")

	// worker that process the execution and each function will go through the switchValues again
	// it will continously accumulate functions to process to enhance the main mapping
	go processBuffer(sigA, ctxA)
	go processBuffer(sigB, ctxB)

	// dequeue the execution buffer
	if err = deepDifferenceWait(sigA, sigB); err != nil {
		return err
	}

	// then we just need to get the maps

	if attrsA, err = getAttributeMap(ctxA); err != nil {
		return err
	}

	if attrsB, err = getAttributeMap(ctxB); err != nil {
		return err
	}

	// TODO @droman: then just need to compare each level of both attrsA/B

	pp.Println(attrsA, attrsB)

	fmt.Println("done", time.Since(now).Microseconds())

	return nil
}

// `switchValue` is the core switch of the recursive stack
// we do only one pass on the root struct and then we dequeue a constant buffer of function that will be added by other functions to avoid overflowing the stack
// TODO @droman: what's the perfs though?
func switchValue(
	ctx context.Context,
	value reflect.Value,
	opts ...optsValueOptions,
) error {
	valueOpts := applyValueOptions(newValueOptions(), opts...)
	// depending on the type
	switch value.Kind() {
	// we need to loop through all structfields
	case reflect.Struct:
		return handleStruct(ctx, value, valueOpts)
	case reflect.Slice:
	case reflect.Array:
		break
	case reflect.Float32:
	case reflect.Float64:
		break
	case reflect.Int:
	case reflect.Int8:
	case reflect.Int16:
	case reflect.Int32:
	case reflect.Int64:
		break
	case reflect.Pointer:
		return switchValue(
			ctx,
			value.Elem(),
			fromPointer(), // don't need to have a if condition with that arch
			fromPath(valueOpts.path),
			fromParentPath(valueOpts.parent),
		)
	case reflect.Uint:
	case reflect.Uint8:
	case reflect.Uint16:
	case reflect.Uint32:
	case reflect.Uint64:
		break
	case reflect.Bool:
		break
	case reflect.Complex64:
	case reflect.Complex128:
		break
	case reflect.String:
		return handleString(ctx, value, valueOpts)
	case reflect.Uintptr:
		break
	case reflect.Interface:
		break
	case reflect.Func:
		break
	case reflect.Chan:
		break
	case reflect.Map:
		break
	case reflect.UnsafePointer:
		break
	case reflect.Invalid:
		break
	default:
		return ErrUnkownKind
	}
	return nil
}

// Deep difference on blocks
func DeepDifferenceBlock[T any](data *T) {

}
