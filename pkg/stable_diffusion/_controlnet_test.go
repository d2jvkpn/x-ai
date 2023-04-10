package stable_diffusion

import (
	"context"
	"fmt"
	"testing"
)

func testControlnetUnit(b64img string) ControlnetUnit {
	return ControlnetUnit{
		InputImage: b64img,
		Module: "depth",
		Model: "control_sd15_depth [fef5e48e]",
		Guidance: 1.0,
	}
}

func TestCNTxt2Img(t *testing.T) {
	b64img, err := Base64File("wk/car01.png")
	if err != nil {
		t.Fatal(err)
	}

	req := &CNTxt2ImgReq{
		Txt2ImgReq: Txt2ImgReq{Prompt: "a red car"},
		ControlnetUnits: []ControlnetUnit{
			testControlnetUnit(b64img),
		}}

	res, err := _TestClient.CNTxt2Img(context.TODO(), req)
	if err != nil {
		t.Fatal(err)
	}

	prefix, err := _TestClient.SaveImgs(res.Images, "wk/stable_diffusion", nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("~~~ prefix:", prefix)
}
