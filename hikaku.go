package hikaku

import (
	"context"
	"fmt"
	"reflect"
)

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
	// maybe for the first time, it will have to add a context
	ctx = checkInitContext(ctx)

	// cfg := applyOpts(newConfig(), opts)

	valueA := convertPtrToValue[T](a)
	valueB := convertPtrToValue[T](b)

	if !valuesValid(valueA, valueB) {
		return ErrValuesInvalid
	}

	sig := make(chan error)

	// traverse the whole structs to create functions that will be processed later on
	// those two functions are starters of the worker that process the buffers
	switchValue(ctx, valueA)
	switchValue(ctx, valueB)

	fmt.Println("finished traversing")

	// worker that process the execution and each function will go through the switchValues again
	// it will continously accumulate functions to process to enhance the main mapping
	go func(s chan error, c context.Context) {
		var exe *executionBuffer
		var err error
		if exe, err = getExecutionBuffer(c); err != nil {
			s <- err
			close(s)
			return
		}
		finished := false
		for !finished {
			fn := exe.Pop()
			if err = fn(); err != nil {
				s <- err
				return
			}
			// put an end to it
			if exe.Len() == 0 {
				finished = true
			}
		}
		close(s)
	}(sig, ctx)

	for e := range sig {
		fmt.Println("error", e)
	}

	select {
	case <-sig:
		fmt.Println("closed")
	}

	fmt.Println("done")

	return nil
}

// main switch
func switchValue(ctx context.Context, value reflect.Value, opts ...optsValueOptions) error {
	valueOpts := applyValueOptions(newValueOptions())
	fmt.Println(value.Kind())
	// depending on the type
	switch value.Kind() {
	// we need to loop through all structfields
	case reflect.Struct:
		if err := handleStruct(ctx, value, valueOpts); err != nil {
			return err
		}
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
		break
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
