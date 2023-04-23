package lang_chain

import (
	_ "embed"
	// "fmt"
)

var (
	//go:embed langchain_index.py
	_LangchainIndex []byte

	//go:embed langchain_query.py
	_LangchainQuery []byte
)
