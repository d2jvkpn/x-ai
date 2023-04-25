package langchain

import (
	"fmt"
	"testing"
)

func TestFaissIndex(t *testing.T) {
	var (
		err   error
		index *FaissIndex
	)

	if index, err = NewFaissIndex([]Source{}); err != nil {
		t.Fatal(err)
	}

	if err = index.SaveYaml("wk/" + index.Uuid() + ".yaml"); err != nil {
		t.Fatal(err)
	}

	fmt.Println("==>", index)
}
