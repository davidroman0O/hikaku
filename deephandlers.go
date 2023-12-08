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

	thisType := value.Type()

	// TODO @droman: we should record probes at every level, including 0

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
			fmt.Println("handle struct field:", varName, varType, varValue, localIdx)
			return switchValue(
				ctx,
				fieldValue,
				probeWithLevel(probe.level),
				probeWithFieldName(varName),
				probeWithParentPath(probe.parent),
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

	// var attrs *AttributeMap
	// if attrs, err = getAttributeMap(ctx); err != nil {
	// 	fmt.Println("can't get attrs map")
	// 	return err
	// }

	exe.Add(func() error {
		fmt.Println(probe.fieldIndex)
		fmt.Println("handle string: ", probe.fieldName, value.Interface(), probe.fieldIndex)
		// attrOpts := []optsAttributeData{
		// 	optAttrParent(probe.parent),
		// 	optName(probe.fieldName),
		// 	optTypeInfo(probe.typeInfo),
		// }
		// // TODO @droman: use options pattern for that one or accumulate them all?!
		// if probe.tag.Get("json") != "" {
		// 	attrOpts = append(attrOpts, withTag(
		// 		probe.tag.Get("json"),
		// 	))
		// }
		// attrs.Add(probe.parent, value, attrOpts...)

		opts := []optionProbe{
			probeWithFieldName(probe.fieldName),
			probeWithTag(probe.tag),
			probeWithLevel(probe.level + 1),
			probeWithParentPath(probe.path),
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
