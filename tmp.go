package netkit

import (
	"html/template"
	"io/ioutil"
	"net/http"
	"sync"
	"fmt"
)

// template
type Template struct {
	dir 	string
	base 	string
	cached 	map[string]*template.Template
	funcs 	template.FuncMap
	mu 		sync.Mutex
}

// register new base template instance
func NewTemplate(dir, base string) *Template {
	return &Template{
		dir: dir,
		base: base,
		cached: map[string]*template.Template{},
		funcs: template.FuncMap{
			"safe": Safe,
		},
	}
}

// html escaper
func Safe(html string) template.HTML {
	return template.HTML(html)
}

// load template files to be associate with base file
func (self *Template) LoadTemplates(name ...string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	for i := 0; i < len(name); i++ {
		t := template.New(self.base).Funcs(self.funcs)//Funcs(template.FuncMap{"safe": safe, "flash": flash})
		t = template.Must(t.ParseFiles(self.dir+"/"+self.base, self.dir+"/"+name[i]))	
		self.cached[name[i]] = t
	}
}

// load html file
func (self *Template) LoadFiles(name ...string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	html := template.New("main").Funcs(self.funcs)
	for i := 0; i < len(name); i++ {
		bd, err := ioutil.ReadFile(self.dir+"/"+name[i])
		if err != nil {
			fmt.Printf("File not found: %s\n", name[i])
		}
		t := template.Must(html.New(self.dir+"/"+name[i]).Parse(string(bd)))
		self.cached[name[i]] = t
	}
}

// render a template (provide name without suffix)
func (self *Template) Render(w http.ResponseWriter, name string) {
	self.SetType(w, "text/html; charset=utf-8")
	self.cached[name].Execute(w, map[string]interface{}{})
}

// render a template, supplying custom map values to be passed in (provide name without suffix)
func (self *Template) RenderCustom(w http.ResponseWriter, name string, c interface{}) {
	self.SetType(w, "text/html; charset=utf-8")
	self.cached[name].Execute(w, c)
}

// render raw data
func (self *Template) RenderRaw(w http.ResponseWriter, format string, a ...interface{}) {
	fmt.Fprintf(w, format, a...)
}

// set the header content type
func (self *Template) SetType(w http.ResponseWriter, typ string) {
	w.Header().Set("Content Type", typ)
}

// view a list of the templates that are cached/registered
func (self *Template) Templates() {
	fmt.Printf("CACHED TEMPLATES FOR '%s'\n", self.base)
	for k, v := range self.cached {
		fmt.Printf("---\nname: %s\ntmpl: %v\n", k, v)
	}
}