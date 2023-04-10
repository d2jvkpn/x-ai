package stable_diffusion

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func Base64File(fp string) (b64 string, err error) {
	var bts []byte

	if bts, err = ioutil.ReadFile(fp); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bts), nil
}

func Imgs2Files(imgs []string, dir string, id *uuid.UUID) (prefix string, err error) {
	var (
		bts []byte
		now time.Time
	)

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

	for i, v := range imgs {
		if bts, err = base64.StdEncoding.DecodeString(v); err != nil {
			return prefix, err
		}

		if err = ioutil.WriteFile(prefix+fmt.Sprintf("_%04d.png", i), bts, 0664); err != nil {
			return prefix, err
		}
	}

	return prefix, nil
}

func Info2File(info, prefix string) (err error) {
	return ioutil.WriteFile(prefix+"info.json", []byte(info), 0664)
}
