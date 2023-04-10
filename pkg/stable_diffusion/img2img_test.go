package stable_diffusion

import (
	"flag"
	"fmt"
	"testing"
)

// go test -run TestImg2Img -- wk/black_cat.png "white cat"
func TestImg2Img(t *testing.T) {
	args := flag.Args()
	if len(args) < 2 {
		t.Fatal("please provide image path and prompt")
	}
	fp, prompt := args[0], args[1]

	b64img, err := Base64File(fp)
	if err != nil {
		t.Fatal(err)
	}

	req := &Img2ImgReq{
		Prompt:     prompt,
		InitImages: []string{b64img},
		BatchSize:  4,
	}

	res, err := _TestClient.Img2Img(_TestCtx, req)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Images) == 0 {
		t.Fatal("no image created")
	}
	prefix, err := Imgs2Files(res.Images, "wk/stable_diffusion", nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("==> prefix:", prefix)
	err = Info2File(res.Info, prefix)
	if err != nil {
		t.Fatal(err)
	}
}
