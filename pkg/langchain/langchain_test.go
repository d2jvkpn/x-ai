package langchain

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestLangChain(t *testing.T) {
	lc, err := NewLangChain("sk-xxxx", "./wk")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("LangChain: %v\n", lc)

	_, err = os.Stat("./wk/langchain_index.py")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat("./wk/langchain_query.py")
	if err != nil {
		t.Fatal(err)
	}
}

func TestExec(t *testing.T) {
	cmd := exec.CommandContext(context.TODO(), "ls")
	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}

	fmt.Println(buf.String())
}
