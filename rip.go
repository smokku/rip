package rip

import (
	"net/http"
	//"strings"
)

/*
type Resource struct {
	Add(string, interface{}) Resource
	Get(string) Resource
	ServeHTTP(http.ResponseWriter, *http.Request)
}
*/

type Handler struct {
	*http.ServeMux
	resources map[string]interface{}
}

func New() *Handler {
	return &Handler{http.NewServeMux(), make(map[string]interface{})}
}

func (h *Handler) Add(name string, resource interface{}) *Handler {
	h.HandleFunc("/"+name+"/", requestHandler(resource))
	return h
}

/*
func (r *Handler) Get(name string) *Handler {
	return r.resources[name]
}

func (r *Handler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	name := req.URL.Path[1:]
	for needle, res := range r.resources {
		if needle == name || strings.HasPrefix(name, needle+"/") {
			if len(name) > len(needle) {
				req.URL.Path = name[len(needle):]
				res.ServeHTTP(resp, req)
			}
		}
	}
}
*/
