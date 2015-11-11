package templice

import (
	"html/template"	
	"io"
	"sync"
	"os"

	"github.com/GeertJohan/go.rice"
)

type Template struct {
	box *rice.Box
	mtx sync.RWMutex
	tpl *template.Template

	dev bool
	lastRoot string
}

func New(bx *rice.Box) *Template {
	t := &Template{box: bx, tpl: template.New("")}
	return t
}

func (t *Template) Load() error {
	return t.LoadDir("")
}

func (t *Template) LoadDir(root string) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.tpl = template.New("")

	t.lastRoot = root
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		str, err := t.box.String(path)
		if err != nil {
			return err
		}

		t.tpl, err = t.tpl.New(path).Parse(str)

		return err
	}

	return t.box.Walk(root, walkFunc)
}

func (t *Template) Dev() *Template {
	t.dev = true
	return t
}

func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	if t.dev {
		t.LoadDir(t.lastRoot)
	}

	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return t.tpl.ExecuteTemplate(wr, name, data)
}

func (t *Template) Do(f func(*template.Template)) *Template {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	f(t.tpl)
	
	return t
}