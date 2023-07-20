package stable_diffusion

import (
	"flag"
	"fmt"
	"testing"
)

// go test -run TestControlnet -timeout 2m -- wk_02_controlnet/emptyroom6_bedroom.jpg wk_02_controlnet/style2_living.jpg
func TestControlnet(t *testing.T) {
	var err error

	args := flag.Args()
	cimg, simg := args[0], args[1]

	if cimg, err = Base64File(cimg); err != nil {
		t.Fatal(err)
	}
	if simg, err = Base64File(simg); err != nil {
		t.Fatal(err)
	}

	req := &Img2ImgReq{
		Prompt:     "",
		InitImages: []string{simg},
		BatchSize:  4,
	}

	conf := _TestClient.config
	cn := &Controlnet{
		InputImage: cimg,
		Module:     conf.Controlnet.Module,
		Model:      conf.Controlnet.Models[0],
	}
	if err := _TestClient.ValidateControlnet(cn); err != nil {
		t.Fatal(err)
	}

	res, err := _TestClient.Img2Img(_TestCtx, req, cn)
	if err != nil {
		t.Fatal(err)
	}

	if len(res.Images) == 0 {
		t.Fatal("no image created")
	}
	prefix, err := Imgs2Files(res.Images, "wk_02_controlnet", nil)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("==> prefix:", prefix)
	err = Info2File(res.Info, prefix)
	if err != nil {
		t.Fatal(err)
	}
}
