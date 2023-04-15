package chatgpt

import (
	"fmt"
	"testing"
)

func TestModels(t *testing.T) {
	ans, err := _TestClient.Models(_TestCtx)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("~~~ ans: %v\n", ans)
}
