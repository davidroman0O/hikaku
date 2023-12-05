package hikaku

import (
	"context"
	"testing"
)

type Basic struct {
	Hello     string
	Something string
	Else      string
	Property  string
}

// go test -v -count=1 -timeout 5s -run ^TestBasic$
func TestBasic(t *testing.T) {
	err := DeepDifference[Basic](context.Background(), &Basic{
		Hello:     "test",
		Something: "something",
		Else:      "else",
		Property:  "property",
	}, &Basic{
		Hello:     "ohoh",
		Something: "fasdfasdas",
		Else:      "fgfhgjgf",
		Property:  "xcxcx",
	})
	if err != nil {
		t.Error(err)
	}
}
