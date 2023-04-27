package langchain

import (
	_ "embed"
	// "fmt"
)

var (
	//go:embed langchain_index.py
	_LangchainIndex []byte
	//go:embed langchain_query.py
	_LangchainQuery []byte
	//go:embed langchain_summarize_01.py
	_LangchainSummarize01 []byte

	_Scripts = map[string][]byte{
		"langchain_index.py":        _LangchainIndex,
		"langchain_query.py":        _LangchainQuery,
		"langchain_summarize_01.py": _LangchainSummarize01,
	}
)
