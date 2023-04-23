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
type LCC struct {
	openai_api_key     string
	py_index, py_query string
}

func NewLCC(key, path string) (lcc *LCC, err error) {
	if err = os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	lcc = &LCC{
		openai_api_key: key,
		py_index:       filepath.Join(path, "langchain_index.py"),
		py_query:       filepath.Join(path, "langchain_query.py"),
	}

	if err = ioutil.WriteFile(lcc.py_index, _LangchainIndex, 0764); err != nil {
		return nil, err
	}

	if err = ioutil.WriteFile(lcc.py_query, _LangchainQuery, 0764); err != nil {
		return nil, err
	}

	return lcc, nil
}

func (lcc *LCC) env() []string {
	return []string{
		fmt.Sprintf("OPENAI_API_KEY=%s", lcc.openai_api_key),
	}
}

func (lcc *LCC) PyIndex(ctx context.Context, cf, prefix string) (err error) {
	cmd := exec.CommandContext(ctx, "python3", lcc.py_index, cf, prefix)
	cmd.Env = append(cmd.Env, lcc.env()...)
	return cmd.Run()
}

func (lcc *LCC) PyQuery(ctx context.Context, prefix, query string) (ans string, err error) {
	cmd := exec.CommandContext(ctx, "python3", lcc.py_index, prefix, query)
	cmd.Env = append(cmd.Env, lcc.env()...)

	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf

	if err = cmd.Run(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
