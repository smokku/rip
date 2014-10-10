package rip

import (
	"net/http"
	"strings"
)

type Resource interface {
	AddResource(string, Resource) Resource
	GetResource(string) Resource
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type ResourceBase struct {
	resources map[string]Resource
}

func New() Resource {
	return new(ResourceBase)
}

func (r *ResourceBase) AddResource(path string, res Resource) Resource {
	if r.resources == nil {
		r.resources = make(map[string]Resource)
	}
	r.resources[path] = res
	return r
}

func (r *ResourceBase) GetResource(path string) Resource {
	return r.resources[path]
}

func (r *ResourceBase) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	path := req.URL.Path[1:]
	for needle, res := range r.resources {
		if needle == path || strings.HasPrefix(path, needle+"/") {
			if len(path) > len(needle) {
				req.URL.Path = path[len(needle):]
				res.ServeHTTP(resp, req)
			}
		}
	}
}
