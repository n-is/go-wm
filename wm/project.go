package wm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Project struct {
	Variables map[string]interface{} `json:"Variables"`

	History  string `json:"-"`
	rootPath string
	config   string
}

func newProject(root, hist, config string) *Project {
	rt, _ := filepath.Abs(root)
	hist = filepath.Join(rt, hist)
	config = filepath.Join(rt, config)

	return &Project{
		rootPath: root,
		History:  hist,
		config:   config,
	}
}

func (p *Project) create() {
	root, err := filepath.Abs(p.rootPath)
	if err != nil {
		log.Panic(err)
	}
	os.MkdirAll(root, os.ModePerm)
	os.Chdir(root)

	_, err = os.Create(p.History)
	if err != nil {
		log.Panic(err)
	}
}

func (p *Project) open() {
	root, err := filepath.Abs(p.rootPath)
	if err != nil {
		log.Panic(err)
	}
	err = os.Chdir(root)

	if err != nil {
		log.Panic(err)
	}
}

func (p *Project) Run(f func(p *Project)) {
	f(p)
}

func (p *Project) Load() error {
	reader, err := os.Open(p.config)
	if err != nil {
		return err
	}

	data, _ := ioutil.ReadAll(reader)
	return json.Unmarshal(data, p)
}

func (p *Project) Save() error {
	writer, err := os.Create(p.config)
	if err != nil {
		return err
	}

	e, err := json.Marshal(p)
	if err != nil {
		return err
	}
	_, err = writer.Write(e)
	return err
}
