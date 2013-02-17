package netkit

import (
	"net/http"
	"net/url"
	"strings"
)

// serve mux
type ServeMux struct {
	routes map[string][]*Handler
}

// return new serve mux instance
func NewServeMux() *ServeMux {
	return &ServeMux{make(map[string][]*Handler)}
}

// match requests against registered handlers, serve http
func (self *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, h := range self.routes[r.Method] {
		if params, ok := h.parse(r.URL.Path); ok {
			if len(params) > 0 {
				r.URL.RawQuery = url.Values(params).Encode() + "&" + r.URL.RawQuery
			}
			h.ServeHTTP(w, r)
			return
		}
	}
	allowed := make([]string, 0, len(self.routes))
	for method, routes := range self.routes {
		if method == r.Method {
			continue
		}
		for _, h := range routes {
			if _, ok := h.parse(r.URL.Path); ok {
				allowed = append(allowed, method)
			}
		}
	}
	if len(allowed) == 0 {
		//http.NotFound(w, r)
		http.Redirect(w, r, "/error/404", 303)
		return
	}
	w.Header().Add("Allow", strings.Join(allowed, ", "))
	//http.Error(w, "Method Not Allowed", 405)
	http.Redirect(w, r, "/error/405", 303)
}

// register and http handler for a paticular method and path
func (self *ServeMux) Route(method, path string, h http.Handler) {
	self.routes[method] = append(self.routes[method], &Handler{path, h})
	n := len(path)
	if n > 0 && path[n-1] == '/' {
		self.Route(method, path[:n-1], http.RedirectHandler(path, 301))
	}
}

// wrapper for Route to allow use of handler functions
func (self *ServeMux) RouteFunc(method, path string, h http.HandlerFunc) {
	self.Route(method, path, h)
}

// register forwarder handler
func (self *ServeMux) Forward(path, newpath string) {
	self.Route("GET", path, http.RedirectHandler(newpath, 301))
}

// register static file server handler
func (self *ServeMux) Static(path, folder string) {
	n := len(path)
	if n > 0 && path[n-1] != '/' {
		path = path + "/"
	}
	h := http.StripPrefix(path, http.FileServer(http.Dir(folder)))
	self.Route("GET", path, h)
}

// route
type Handler struct {
	path 	string
	http.Handler
}

// parse registered pattern
func (self *Handler) parse(path string) (url.Values, bool) {
	p := make(url.Values)
	var i, j int
	for i < len(path) {
		switch {
		case j >= len(self.path):
			if self.path != "/" && len(self.path) > 0 && self.path[len(self.path)-1] == '/' {
				return p, true
			}
			return nil, false
		case self.path[j] == ':':
			var name, val string
			var nextc byte
			name, nextc, j = match(self.path, isBoth, j+1)
			val, _, i = match(path, byteParse(nextc), i)
			p.Add(":"+name, val)
		case path[i] == self.path[j]:
			i++
			j++
		default:
			return nil, false
		}
	}
	if j != len(self.path) {
		return nil, false
	}
	return p, true
}

// match path with registered handler
func match(s string, f func(byte) bool, i int) (matched string, next byte, j int) {
	j = i
	for j < len(s) && f(s[j]) {
		j++
	}
	if j < len(s) {
		next = s[j]
	}
	return s[i:j], next, j
}

// determine type of byte
func byteParse(b byte) func(byte) bool {
	return func(c byte) bool {
		return c != b && c != '/'
	}
}

// test for alpha byte
func isAlpha(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// test for numerical byte
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// test for alpha or numerical byte
func isBoth(ch byte) bool {
	return isAlpha(ch) || isDigit(ch)
}