package lang_chain

import (
	// "fmt"
	"io/ioutil"
	"time"

	"github.com/google/uuid"
	yaml "gopkg.in/yaml.v3"
)

type FaissIndex struct {
	Id      uuid.UUID `json:"id" yaml:"id"`
	Created int64     `json:"created" yaml:"created"`
	Sources []Source  `json:"sources" yaml:"sources"`
}

type Source struct {
	Name   string `json:"name"   yaml:"name"`
	Type   string `json:"type"   yaml:"type"`
	Source string `json:"source" yaml:"source"`
}

func NewFaissIndex(sources []Source) (index *FaissIndex, err error) {
	var id uuid.UUID

	// if len(sources) == 0 {
	//	return nil, fmt.Errorf("empty sources")
	// }

	if id, err = uuid.NewUUID(); err != nil {
		return nil, err
	}

	return &FaissIndex{
		Id:      id,
		Created: time.Now().Unix(),
		Sources: sources,
	}, nil
}

func (index *FaissIndex) Uuid() string {
	return index.Id.String()
}

func (index *FaissIndex) SaveYaml(fp string) (err error) {
	var bts []byte

	if bts, err = yaml.Marshal(index); err != nil {
		return err
	}

	return ioutil.WriteFile(fp, bts, 0664)
}
