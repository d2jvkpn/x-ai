package chatgpt

import (
	"fmt"
	"testing"
)

func TestImgGen(t *testing.T) {
	req := &ImgGenReq{
		Prompt: "The End of The World, Shanghai",
		N:      2,
		Size:   "512x512",
	}

	ans, err := _TestClient.ImgGen(_TestCtx, req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("==>", ans)

	_, err = ans.Save("wk/chatpt", nil)
	if err != nil {
		t.Fatal(err)
	}
}
