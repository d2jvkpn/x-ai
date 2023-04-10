package stable_diffusion

import (
	"flag"
	"fmt"
	"testing"
)

// go test -run TestInterrogateCLIP -- wk/car01.png
func TestInterrogateCLIP(t *testing.T) {
	args := flag.Args()
	if len(args) < 1 {
		t.Fatal("please provide image path")
	}
	fp := args[0]

	res, err := _TestClient.InterrogateCLIP(fp)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("==> result:", res)
}

// go test -run TestInterrogateDeepDanBooru -- wk/car01.png
func TestInterrogateDeepDanBooru(t *testing.T) {
	args := flag.Args()
	if len(args) < 1 {
		t.Fatal("please provide image path")
	}
	fp := args[0]

	res, err := _TestClient.InterrogateDeepDanBooru(fp)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("==> result:", res)
}
