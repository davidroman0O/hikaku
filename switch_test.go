package hikaku

import (
	"reflect"
	"testing"
)

func TestSwitchStruct(t *testing.T) {
	testCases := []struct {
		Value  interface{}
		Expect error
	}{
		struct {
			Value  interface{}
			Expect error
		}{
			Value: struct {
				Hello string
			}{
				Hello: "world",
			},
			Expect: nil,
		},
	}

	for i := 0; i < len(testCases); i++ {
		err := probeValue(nil, reflect.ValueOf(testCases[i].Value))
		if err != testCases[i].Expect {
			t.Error("switch value doesn't give the correct error ", err, testCases[i].Expect)
		}
	}
}
