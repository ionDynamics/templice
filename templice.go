package templice

import (
	"html/template"
	"io"
	"os"
	"sync"

	"github.com/GeertJohan/go.rice"
)

type Func func(*template.Template) *template.Template

type Template struct {
	box *rice.Box
	mtx sync.RWMutex
	tpl *template.Template
	pre Func

	dev      bool
	lastRoot string
}

//New initializes a new Templice.Template
func New(bx *rice.Box) *Template {
	t := &Template{box: bx, tpl: template.New("")}
	return t
}

//SetPrep's given function is called before parsing the templates in Load / LoadDir
func (t *Template) SetPrep(pre Func) *Template {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.pre = pre
	return t
}

//Load is a shortcut for LoadDir("")
func (t *Template) Load() error {
	return t.LoadDir("")
}

//LoadDir prepares, loads and parses templates in the given directory
func (t *Template) LoadDir(root string) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	t.tpl = template.New("")

	if t.pre != nil {
		t.unsafeDo(t.pre)
	}

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

//Dev forces ExecuteTemplate to reload templates before execution
func (t *Template) Dev() *Template {
	t.dev = true
	return t
}

//ExecuteTemplate writes the given template and data to the writer
func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	if t.dev {
		t.LoadDir(t.lastRoot)
	}

	t.mtx.RLock()
	defer t.mtx.RUnlock()

	return t.tpl.ExecuteTemplate(wr, name, data)
}

//Do executes
func (t *Template) Do(f Func) *Template {
	t.mtx.Lock()
	defer t.mtx.Unlock()

	return t.unsafeDo(f)
}

func (t *Template) unsafeDo(f Func) *Template {
	tpl := f(t.tpl)
	if tpl != nil {
		t.tpl = tpl
	}
	return t
}