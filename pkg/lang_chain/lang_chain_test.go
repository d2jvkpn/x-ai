package lang_chain

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestLCC(t *testing.T) {
	lcc, err := NewLLC("sk-xxxx", "./wk")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("LCC: %v\n", lcc)

	_, err = os.Stat("./wk/lang_chain/langchain_index.py")
	if err != nil {
		t.Fatal(err)
	}

	_, err = os.Stat("./wk/lang_chain/langchain_query.py")
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
