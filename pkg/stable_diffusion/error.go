package stable_diffusion

import (
// "fmt"
)

/*
{
  "detail": [
    {
      "loc": [
        "string",
        0
      ],
      "msg": "string",
      "type": "string"
    }
  ]
}
*/

type SDErrorDetail struct {
	Msg  string `json:"msg"`
	Type string `json:"type"`
}

type SDError struct {
	Detail []SDErrorDetail `json:"detail"`
}
