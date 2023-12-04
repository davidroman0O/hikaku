package hikaku

import (
	"context"
	"testing"
)

type Basic struct {
	Hello string
}

func TestBasic(t *testing.T) {
	err := DeepDifference[Basic](context.Background(), &Basic{
		Hello: "test",
	}, &Basic{
		Hello: "ohoh",
	})
	if err != nil {
		t.Error(err)
	}
}
