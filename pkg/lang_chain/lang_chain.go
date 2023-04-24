package lang_chain

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// Lang Chain Client
type LangChain struct {
	openai_api_key     string
	path               string
	py_index, py_query string
}

func NewLangChain(key, path string) (lc *LangChain, err error) {
	if err = os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	if path, err = filepath.Abs(path); err != nil {
		return nil, err
	}

	lc = &LangChain{
		openai_api_key: key,
		path:           path,
		py_index:       filepath.Join(path, "langchain_index.py"),
		py_query:       filepath.Join(path, "langchain_query.py"),
	}

	if err = ioutil.WriteFile(lc.py_index, _LangchainIndex, 0764); err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile(lc.py_query, _LangchainQuery, 0764); err != nil {
		return nil, err
	}

	return lc, nil
}

func (lc *LangChain) SubPath(elem ...string) string {
	list := append([]string{lc.path}, elem...)
	return filepath.Join(list...)
}

func (lc *LangChain) env() []string {
	return []string{
		fmt.Sprintf("OPENAI_API_KEY=%s", lc.openai_api_key),
	}
}

func (lc *LangChain) PyIndex(ctx context.Context, cf, prefix string) (err error) {
	cmd := exec.CommandContext(ctx, "python3", lc.py_index, cf, prefix)
	cmd.Env = append(cmd.Env, lc.env()...)
	return cmd.Run()
}

func (lc *LangChain) PyQuery(ctx context.Context, prefix, query string) (ans string, err error) {
	cmd := exec.CommandContext(ctx, "python3", lc.py_index, prefix, query)
	cmd.Env = append(cmd.Env, lc.env()...)

	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf

	if err = cmd.Run(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
