package chatgpt

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ImgGenReq struct {
	Prompt string `json:"prompt"`
	N      uint32 `json:"n,omitempty"`
	Size   string `json:"size,omitempty"`
}

type ImgGenRes struct {
	Created int64 `json:"created"`
	Data    []struct {
		Url string `json:"url"`
	} `json:"data,omitempty"`
}

func (req *ImgGenReq) Validate() (err error) {
	if req.Prompt == "" {
		return fmt.Errorf("empty prompt")
	}

	if req.N == 0 {
		req.N = 1
	}

	if req.Size == "" {
		req.Size = "1024x1024"
	}
	if !_RE_ImageSize.MatchString(req.Size) {
		return fmt.Errorf("invalid size")
	}

	return nil
}

func (client *Client) ImgGen(ctx context.Context, req *ImgGenReq) (res *ImgGenRes, err error) {
	var (
		request  *http.Request
		response *http.Response
	)

	if err = req.Validate(); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	_ = encoder.Encode(req)

	request, _ = http.NewRequest("POST", client.img_gen_url(), buf)
	request.WithContext(ctx)
	client.setAuth(request, true)

	if response, err = client.cli.Do(request); err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(response.Status)
	}

	res = new(ImgGenRes)
	decoder := json.NewDecoder(response.Body)
	if err = decoder.Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func (res *ImgGenRes) List() []string {
	list := make([]string, 0, len(res.Data))
	for _, v := range res.Data {
		list = append(list, v.Url)
	}

	return list
}

func (res *ImgGenRes) Save(dir string, id *uuid.UUID) (prefix string, err error) {
	var (
		bts []byte
		now time.Time
	)

	if len(res.Data) == 0 {
		return "", nil
	}
	if err = os.MkdirAll(dir, 0750); err != nil {
		return "", err
	}
	if id == nil {
		newId := uuid.New()
		id = &newId
	}

	now = time.Now()
	prefix = filepath.Join(dir, now.Format("2006-01-02T15-03-04MST")+"_"+id.String())
	// fmt.Println("~~~", prefix, res.Info)

	bts, _ = json.Marshal(res)
	if err = ioutil.WriteFile(prefix+"_info.json", bts, 0664); err != nil {
		return "", err
	}

	errs := make([]error, 0, len(res.Data))
	errch := make(chan error, 0)
	wg := new(sync.WaitGroup)
	cli := new(http.Client)

	download := func(i int, url string) {
		var (
			e    error
			res  *http.Response
			file *os.File
		)

		defer func() {
			if e != nil {
				errch <- e
			}
			wg.Done()
		}()

		if res, e = cli.Get(url); e != nil {
			return
		}
		defer res.Body.Close()
		if res.StatusCode != http.StatusOK {
			e = fmt.Errorf(res.Status)
			return
		}

		if file, e = os.Create(prefix + fmt.Sprintf("%04d.png", i)); e != nil {
			return
		}
		defer file.Close()
		_, e = io.Copy(file, res.Body)
	}

	go func() {
		for e := range errch {
			errs = append(errs, e)
		}
	}()

	for i, v := range res.Data {
		wg.Add(1)
		go download(i, v.Url)
	}

	wg.Wait()
	close(errch)

	if len(errs) > 0 {
		err = errors.Join(errs...)
	}

	return prefix, err
}
