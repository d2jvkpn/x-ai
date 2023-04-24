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
	openai_api_key string
	path           string
	scripts        map[string]string
	env            []string
}

func NewLangChain(key, path string) (lc *LangChain, err error) {
	if err = os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}

	if path, err = filepath.Abs(path); err != nil {
		return nil, err
	}

	if key == "" {
		return nil, fmt.Errorf("key is empty")
	}

	if _, err = exec.LookPath("python3"); err != nil {
		return nil, err
	}

	lc = &LangChain{
		openai_api_key: key,
		path:           path,
		scripts:        make(map[string]string, len(_Scripts)),
		env:            append(os.Environ(), fmt.Sprintf("OPENAI_API_KEY=%s", key)),
	}

	for k, v := range _Scripts {
		if err = ioutil.WriteFile(filepath.Join(path, k), v, 0764); err != nil {
			return nil, err
		}
	}

	return lc, nil
}

func (lc *LangChain) GetPath() string {
	return lc.path
}

func (lc *LangChain) PyIndex(ctx context.Context, cf, prefix string) (err error) {
	cmd := exec.CommandContext(ctx, "python3", lc.scripts["langchain_index.py"], cf, prefix)
	cmd.Env = lc.env
	return cmd.Run()
}

func (lc *LangChain) PyQuery(ctx context.Context, prefix, query string) (ans string, err error) {
	cmd := exec.CommandContext(ctx, "python3", lc.scripts["langchain_query.py"], prefix, query)
	cmd.Env = lc.env

	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf

	if err = cmd.Run(); err != nil {
		return "", err
	}

	return buf.String(), nil
}
