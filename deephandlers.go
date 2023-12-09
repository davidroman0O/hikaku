package hikaku

import (
	"context"
	"fmt"
	"reflect"
)

// Receive a value of kind struct
// so far handling struct properly
func handleStruct(
	ctx context.Context,
	value reflect.Value,
	probe *Probe,
) error {
	var exe *executionBuffer
	var err error
	if exe, err = getExecutionBuffer(ctx); err != nil {
		fmt.Println("can't get execution buffer")
		return err
	}

	var mapProbe *ProbeMap
	if mapProbe, err = getProbeMap(ctx); err != nil {
		fmt.Println("can't get probe map")
		return err
	}
	thisType := value.Type()

	var currentProbe *Probe = probe

	// TODO @droman: we should record probes at every level, including 0
	if probe.level == 0 {
		// I am root
		rootProbe := applyProbeOptions(
			newProbe(),
			probeWithKind(probe.kind),
			probeWithTypeName(value.Type().Name()),
			probeWithValue(value),
		)
		mapProbe.Add(rootProbe)
		currentProbe = rootProbe
		// return handleStruct(ctx, value, applyProbeOptions(rootProbe, probeWithLevel(1)))
	}

	fmt.Println("struct ", thisType.Name())

	for idx := 0; idx < value.NumField(); idx++ {

		fieldValue := value.Field(idx)
		fieldType := thisType.Field(idx)

		varName := value.Type().Field(idx).Name
		varType := value.Type().Field(idx).Type.Kind()

		varValue := value.Field(idx).Interface()

		localIdx := idx

		exe.Add(func() error {
			// TODO @droman: here we do not increment the path of the parent and have no notion of real parenting, which override the same existing paths
			fmt.Println("handle struct field:", varName, varType, varValue, localIdx, probe.parentPath)

			return probeValue(
				ctx,
				fieldValue,
				probeWithParentProbe(currentProbe),
				probeWithValue(fieldValue),
				probeWithLevel(probe.level),
				probeWithFieldName(varName),
				probeWithParentPath(probe.parentPath),
				probeWithTag(fieldType.Tag),
				probeWithFieldIndex(localIdx),
				probeWithParentType(fieldType.Type),
				probeWithTypeName(varType.String()),
			)
		})
	}

	return nil
}

// so far handling strings properly
func handleString(
	ctx context.Context,
	value reflect.Value,
	probe *Probe,
) error {
	var exe *executionBuffer
	var err error
	if exe, err = getExecutionBuffer(ctx); err != nil {
		fmt.Println("can't get execution buffer")
		return err
	}

	var mapProbe *ProbeMap
	if mapProbe, err = getProbeMap(ctx); err != nil {
		fmt.Println("can't get probe map")
		return err
	}

	exe.Add(func() error {
		fmt.Println(probe.fieldIndex)
		fmt.Println("handle string: ", probe.fieldName, value.Interface(), probe.fieldIndex)

		var realValue interface{}
		if probe.value.Kind() == reflect.Pointer {
			realValue = probe.value.Elem().Interface()
		} else {
			realValue = probe.value.Interface()
		}

		opts := []optionProbe{
			probeWithParentProbe(probe.parent),
			probeWithParentPath(probe.parentPath),
			probeWithData(realValue),
			probeWithTypeName(value.Type().Name()),
			probeWithFieldName(probe.fieldName),
			probeWithTag(probe.tag),
			probeWithLevel(probe.level + 1),
			probeWithParentType(probe.parentType),
			probeWithKind(value.Kind()),
		}

		if probe.isPointer {
			// only the parent will tell me if it's a pointer
			opts = append(opts, probeWithPointer())
		}

		mapProbe.Add(
			applyProbeOptions(
				newProbe(),
				opts...,
			),
		)
		return nil
	})

	return nil
}

func handleIntInt8Int16Int32Int64(
	ctx context.Context,
	value reflect.Value,
	probe *Probe,
) error {
	var exe *executionBuffer
	var err error
	if exe, err = getExecutionBuffer(ctx); err != nil {
		fmt.Println("can't get execution buffer")
		return err
	}

	var mapProbe *ProbeMap
	if mapProbe, err = getProbeMap(ctx); err != nil {
		fmt.Println("can't get probe map")
		return err
	}

	exe.Add(func() error {
		fmt.Println(probe.fieldIndex)
		fmt.Println("handle int: ", probe.fieldName, value.Interface(), probe.fieldIndex)

		var realValue interface{}
		if probe.value.Kind() == reflect.Pointer {
			realValue = probe.value.Elem().Interface()
		} else {
			realValue = probe.value.Interface()
		}

		opts := []optionProbe{
			probeWithParentProbe(probe.parent),
			probeWithParentPath(probe.parentPath),
			probeWithData(realValue),
			probeWithTypeName(value.Type().Name()),
			probeWithFieldName(probe.fieldName),
			probeWithTag(probe.tag),
			probeWithLevel(probe.level + 1),
			probeWithParentType(probe.parentType),
			probeWithKind(value.Kind()),
		}

		if probe.isPointer {
			// only the parent will tell me if it's a pointer
			opts = append(opts, probeWithPointer())
		}

		mapProbe.Add(
			applyProbeOptions(
				newProbe(),
				opts...,
			),
		)
		return nil
	})

	return nil
}

func handleArray(ctx context.Context, value reflect.Value, opts *Probe) error {

	return nil
}
