# templice

A wrapper for Go's html/template and GeertJohan's go.rice  
  
If you have used template.ParseGlob and switched from direct filesystem access to go.rice, you've might run into the problems this library tries to circumvent.

[Documentation](https://go.iondynamics.net/templice)

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

#### Prepare a FuncMap before loading/parsing templates
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

In this case templice will read from the FuncMap every time Load/LoadDir is called.

If you want your programm to evaluate FuncMap only when calling SetPrep you should use a wrapper like the following one: 
```go
tpl.SetPrep(func (f template.FuncMap) templice.Func {
	return func(templ *template.Template) *template.Template {
		return templ.Funcs(f)
	}
}(funcMap))
```