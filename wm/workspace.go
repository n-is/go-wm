package wm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	ConfigFile    = "ws.conf.json"
	HistoryFile   = "prj.history"
	ProjectConfig = "prj.conf.json"
)

type WorkSpace struct {
	Projects map[string]string `json:"projects"`
	conf     string
	projects []*Project
}

// OpenWorkspace either opens a new workspace or creates
// it if it doesn't exist and returns it.
func OpenWorkSpace(folder string) *WorkSpace {

	home := homeDir()
	wsf := filepath.Join(home, "workspaces", folder)

	os.MkdirAll(wsf, os.ModePerm)

	conf := filepath.Join(wsf, ConfigFile)
	var f *os.File
	if _, err := os.Stat(conf); os.IsNotExist(err) {
		f, err = os.Create(conf)
		if err != nil {
			log.Panic(err)
		}
	} else {
		f, err = os.Open(conf)
		if err != nil {
			log.Panic(err)
		}
	}

	ws := &WorkSpace{}
	if f != nil {
		data, _ := ioutil.ReadAll(f)
		if len(data) > 0 {
			if err := json.Unmarshal(data, ws); err != nil {
				log.Panic(err)
			}
		} else {
			ws.Projects = make(map[string]string)
		}
	}
	f.Close()

	ws.conf = conf
	return ws
}

func RemoveWorkspace(folder string) error {
	home := homeDir()
	wsf := filepath.Join(home, "workspaces", folder)

	return os.RemoveAll(wsf)
}

func (w *WorkSpace) AddNewProject(name, rootPath string) (*Project, error) {
	w.Projects[name], _ = filepath.Abs(rootPath)
	w.Update()

	p := newProject(rootPath, HistoryFile, ProjectConfig)
	p.create()
	return p, nil
}

func (w *WorkSpace) OpenProject(name string) (*Project, error) {
	rootPath, ok := w.Projects[name]
	if !ok {
		return nil, errors.New("project not found")
	}

	p := newProject(rootPath, HistoryFile, ProjectConfig)
	p.open()
	return p, nil
}

func (w *WorkSpace) RemoveProject(name string) (*Project, error) {
	rootPath, ok := w.Projects[name]
	if !ok {
		return nil, errors.New("project not found")
	}
	delete(w.Projects, name)
	w.Update()

	p := newProject(rootPath, HistoryFile, ProjectConfig)
	p.open()
	return p, nil
}

func (w *WorkSpace) Update() {
	if bts, err := json.Marshal(w); err != nil {
		log.Panic(err)
	} else {
		f, err := os.Create(w.conf)
		if err != nil {
			log.Panic(err)
		}

		f.Write(bts)
		f.Close()
	}
}
