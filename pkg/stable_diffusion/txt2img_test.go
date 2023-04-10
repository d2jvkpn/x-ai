package stable_diffusion

import (
	"flag"
	"fmt"
	"testing"
)

func testControlnet(b64img string) *Controlnet {
	return &Controlnet{
		InputImage: b64img,
		Module:     "depth",
		Model:      "control_sd15_depth [fef5e48e]",
	}
}

func TestTxt2Img(t *testing.T) {
	args := flag.Args()
	prompt := "a dog wearing a hat"
	if len(args) > 0 {
		prompt = args[0]
	}

	res, err := _TestClient.Txt2ImgPromot(prompt, 4)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("==> res is nil:", res == nil)

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

func TestTxt2ImgWithControlnet(t *testing.T) {
	b64img, err := Base64File("wk/car01.png")
	if err != nil {
		t.Fatal(err)
	}

	req := &Txt2ImgReq{Prompt: "a red car"}
	cn := testControlnet(b64img)

	res, err := _TestClient.Txt2Img(_TestCtx, req, cn)
	if err != nil {
		t.Fatal(err)
	}

	prefix, err := Imgs2Files(res.Images, "wk/stable_diffusion", nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("~~~ prefix:", prefix)
}
