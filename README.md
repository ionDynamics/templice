# templice

A wrapper for Go's html/template and GeertJohan's go.rice

### Examples
Import the package: `import "go.iondynamics.net/templice"`  
most of the time you'll also need at least `import "github.com/GeertJohan/go.rice"`

#### Basic Usage
```go
tpl := templice.New(rice.MustFindBox("template)).Load()
tpl.ExecuteTemplate(os.Stdout, "hello.tpl", "world")
```

#### Usage with subdirectories
```go
tpl := templice.New(rice.MustFindBox("http-files"))
tpl.LoadDir("templates")
tpl.ExecuteTemplate(os.Stdout, "templates/hello.tpl", "world")
```

#### Load a FuncMap
```go
tpl := templice.New(rice.MustFindBox("template))
funcMap := template.FuncMap{
	"up": strings.ToUpper,
}

tpl.SetPrep(func(templ *template.Template) *template.Template {
	return templ.Funcs(funcMap)
})
tpl.Load()
tpl.ExecuteTemplate(os.Stdout, "hello.tpl", "world")
```
