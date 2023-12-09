package hikaku

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/k0kubun/pp/v3"
)

type DeepDifferenceConfig struct {
	tag string
}

func NewDeepDifferenceConfig() *DeepDifferenceConfig {
	return &DeepDifferenceConfig{
		tag: "json",
	}
}

func ApplyDeepDifferenceConfig(c *DeepDifferenceConfig, opts ...OptionsDeepDifferences) *DeepDifferenceConfig {
	for i := 0; i < len(opts); i++ {
		opts[i](c)
	}
	return c
}

type OptionsDeepDifferences func(c *DeepDifferenceConfig)

func DeepDifferenceWithTag(tagName string) OptionsDeepDifferences {
	return func(c *DeepDifferenceConfig) {
		c.tag = tagName
	}
}

// note: the problem with recursive, is the amount of stack we will produce, so i suggest i should have an array of functions
// to process it over and over
// Deep difference to block of differences
func DeepDifference[T any](
	a *T,
	b *T,
	opts ...OptionsDeepDifferences,
) error {
	now := time.Now()

	var err error
	var attrsA *ProbeMap
	var attrsB *ProbeMap

	valueA := reflect.ValueOf(*a) // directly the value of the A
	valueB := reflect.ValueOf(*b) // directly the value of the

	if !valuesValid(valueA, valueB) {
		return ErrValuesInvalid
	}

	sigA := make(chan error)
	sigB := make(chan error)

	// add execution buffers
	ctxA := addCheckExecutionBuffer(context.TODO())
	ctxB := addCheckExecutionBuffer(context.TODO())

	// We need probes
	ctxA = addCheckProbeMap(ctxA)
	ctxB = addCheckProbeMap(ctxB)

	// traverse the whole structs to create functions that will be processed later on
	// those two functions are starters of the worker that process the buffers
	probeValue(
		ctxA,
		valueA,
		probeWithKind(valueA.Kind()),
		probeWithTypeName(valueA.Type().Name()),
		probeWithParentPath("."),
		probeWithLevel(0),
	)
	probeValue(
		ctxB,
		valueB,
		probeWithKind(valueB.Kind()),
		probeWithTypeName(valueB.Type().Name()),
		probeWithParentPath("."),
		probeWithLevel(0),
	)

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

	if attrsA, err = getProbeMap(ctxA); err != nil {
		return err
	}

	if attrsB, err = getProbeMap(ctxB); err != nil {
		return err
	}

	// TODO @droman: then just need to compare each level of both attrsA/B

	pp.Println(attrsA, attrsB)

	fmt.Println("done", time.Since(now).Microseconds())

	return nil
}

// `probeValue` is the core switch of the recursive stack, only the first root level will be on the stack then a dequeuing will need to happen to continue the exploration of the nested structures
// We do only one pass on the root struct and then we dequeue a constant buffer of function that will be added by other functions to avoid overflowing the stack
// TODO @droman: what's the perfs though?
// context of ExecutionBuffer + ProbeMap are required
func probeValue(
	ctx context.Context,
	value reflect.Value,
	opts ...optionProbe,
) error {
	probeConfig := applyProbeOptions(newProbe(), opts...)
	// depending on the type
	switch value.Kind() {
	// we need to loop through all structfields
	case reflect.Struct:
		return handleStruct(ctx, value, probeConfig)
	case reflect.Slice:
	case reflect.Array:
		break
	case reflect.Float32:
	case reflect.Float64:
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return handleIntInt8Int16Int32Int64(ctx, value, probeConfig)
	case reflect.Pointer:
		// same probe as before but as a pointer
		// TODO @droman: I should support more complex types through other means of identification
		// simply because something can create a pointer of a pointer of an array that contain pointers of a value
		return probeValue(
			ctx,
			value.Elem(),
			probeWithProbe(probeConfig), // need a function that pass the probe with just pointer
			probeWithPointer(),          // don't need to have a if condition with that arch
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
		return handleString(ctx, value, probeConfig)
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
