REST.Is.Popular
===============

RESTful microframework for Go

## Example
```go
import (
	"github.com/smokku/rip"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"math/rand"
	"net/http"
)


// Create a resource and attach HTTP methods handlers
type randomResource struct{}

func (randomResource) Get(id string, req *http.Request) (int, interface{}, http.Header) {
	if len(id) > 0 {
		i, err := strconv.ParseUint(id, 10, 0)

		if err != nil {
			return http.StatusBadRequest, err.Error(), nil
		}

		return http.StatusOK, rand.Intn(i), nil

	} else {
		return http.StatusNotImplemented, "Listing all random numbers not implemented", nil
	}
}


func main() {
	// Create REST API handler and attach resources
	apiHandler = rip.New()
	apiHandler.Add("random", randomResource{})

	// Use Goji SubRouter to attach API at /api/*
	api := web.New()
	goji.Handle("/api/*", api)
	api.Use(middleware.SubRouter)
	api.Handle("/*", apiHandler)

	// Serve static files from public/ just for kicks
	goji.Handle("/*", http.FileServer(http.Dir("public")))

	// Go, go, go...
	goji.Serve()
}
```

    $ curl -i localhost:8000/api/random/100

