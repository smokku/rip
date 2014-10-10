package rip

import "net/http"

type Resource interface {
	AddResource(string, Resource) Resource
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type ResourceBase struct{}

func AddResource(path string, res Resource) Resource {
	return new(ResourceBase).AddResource(path, res)
}

func (r *ResourceBase) AddResource(path string, res Resource) Resource {
	return r
}

func (r *ResourceBase) ServeHTTP(http.ResponseWriter, *http.Request) {

}
