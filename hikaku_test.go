package hikaku

import (
	"testing"
)

type Basic struct {
	Hello     string
	Something *string
	Else      string `json:"else"`
	Property  string
}

type BasicNested struct {
	Something Basic
	Hello     string
	Else      string `json:"else"`
	Property  string
}

// go test -v -count=1 -timeout 5s -run ^TestBasic$
func TestBasic(t *testing.T) {
	somethingA := "something"
	somethingB := "fasdfasdas"
	err := DeepDifference[Basic](&Basic{
		Hello:     "test",
		Something: &somethingA,
		Else:      "else",
		Property:  "property",
	}, &Basic{
		Hello:     "ohoh",
		Something: &somethingB,
		Else:      "fgfhgjgf",
		Property:  "xcxcx",
	})
	if err != nil {
		t.Error(err)
	}
}

// go test -v -count=1 -timeout 5s -run ^TestBasicNested$
func TestBasicNested(t *testing.T) {
	somethingA := "something"
	somethingB := "fasdfasdas"
	err := DeepDifference[BasicNested](&BasicNested{
		Hello: "test",
		Something: Basic{
			Hello:     "ohoh",
			Something: &somethingA,
			Else:      "fgfhgjgf",
			Property:  "xcxcx",
		},
		Else:     "else",
		Property: "property",
	}, &BasicNested{
		Hello: "ohoh",
		Something: Basic{
			Hello:     "ohoh",
			Something: &somethingB,
			Else:      "ddd",
			Property:  "ggg",
		},
		Else:     "fgfhgjgf",
		Property: "xcxcx",
	})
	if err != nil {
		t.Error(err)
	}
}

type arrayPointer []*string
type nestedPointer *arrayPointer
type mainPointer *nestedPointer
type dataPointers struct {
	FindMe *mainPointer
}

// This test is the goal i which the probe system can support
// go test -v -count=1 -timeout 5s -run ^TestBasicNestedPointers$
func TestBasicNestedPointers(t *testing.T) {
	var valueA string = "valueA"
	var targetA arrayPointer = arrayPointer{
		&valueA,
	}
	var nestedTargetA nestedPointer = &targetA
	var valuePointerA mainPointer = &nestedTargetA

	var valueB string = "valueA"
	var targetB arrayPointer = arrayPointer{
		&valueB,
	}
	var nestedTargetB nestedPointer = &targetB
	var valuePointerB mainPointer = &nestedTargetB

	err := DeepDifference[dataPointers](&dataPointers{
		FindMe: &valuePointerA,
	}, &dataPointers{
		FindMe: &valuePointerB,
	})
	if err != nil {
		t.Error(err)
	}
}
