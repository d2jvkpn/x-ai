package stable_diffusion

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
)

var (
	_TestCtx    context.Context = context.TODO()
	_TestClient *Client
)

// default config: ../../configs/local.yaml
func TestMain(m *testing.M) {
	var (
		config string
		err    error
	)

	testFlag := flag.NewFlagSet("testFlag", flag.ExitOnError)
	flag.Parse() // must do

	testFlag.StringVar(&config, "config", "../../configs/local.yaml", "config filepath")
	testFlag.Parse(flag.Args())

	defer func() {
		if err != nil {
			fmt.Printf("!!! TestMain: %v\n", err)
			os.Exit(1)
		}
	}()

	if _TestClient, err = NewClient(config, "stable_diffusion"); err != nil {
		return
	}

	m.Run()
}
