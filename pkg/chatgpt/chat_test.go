package chatgpt

import (
	"fmt"
	"testing"
)

func TestChat_t1(t *testing.T) {
	req := &ChatReq{
		Model:       default_model(),
		Temperature: default_temperature(),
		Messages: []ChatMsg{
			{Role: "user", Content: "How to implement a http server in golang?"},
		},
	}

	ans, err := _TestClient.Chat(_TestCtx, req)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("==>", ans)
}

func TestChat_t2(t *testing.T) {
	var (
		ans string
		err error
	)

	if ans, err = _TestClient.Trans2En(_TestCtx, "如何使用 golang 进行 grpc 开发?"); err != nil {
		t.Fatal(err)
	}
	fmt.Println("==>", ans)

	ans, err = _TestClient.Trans2Cn(_TestCtx, "How to become a software architect")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("==>", ans)
}
