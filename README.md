REST.Is.Popular.
================
[![Gitter](https://badges.gitter.im/Join Chat.svg)](https://gitter.im/smokku/rip?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

*RESTful microframework for Go*

There are many REST frameworks for Go Lang (and even more that claim they are), so why write another one?

I needed a skeleton code that simplifies parsing JSON-REST requests and servicn JSON responses. Simple one, that allows me to throw a few lines of code and get some results, but all I got from _The Web_ were behemots requiring me to write gazilions of API descriptors and covering everything including kitchen sink.

There was one candy though. A [blog article](http://dougblack.io/words/a-restful-micro-framework-in-go.html) that discussed simple matching of HTTP resource methods to Resource object methods. I tried it and it was a pleasure. But after some use I realized its defeciencies and how I could improve on this idea. For once, my API is JSON focused, so auto unmarshalling/marshaling is a must. Secondly, method signatures were a bit cumbersome. Thus RIP was born.


## Example

This example shows [Goji](http://goji.io/) integration, but it is not a must. You can use anything net/http compliant.

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

## Cross-Origin Resource Sharing

Once you put your API on a separate domain than your main site (or enable it for public use)
you will want to setup CORS handling for API requests. Thanks to rs/cors library it is
mindblowingly simple:

```go
import (
	"github.com/smokku/rip"
	//[...]
	"github.com/rs/cors"
)

func main() {
	//[...]
	api1 := web.New()
	//[...]
	api1.Use(middleware.SubRouter)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATH", "DELETE"},
		AllowedHeaders:   []string{"accept", "content-type"},
		AllowCredentials: true,
	})
	api1.Use(c.Handler)
	//[...]
}
```

## Lecture

- [Using HTTP Methods for RESTful Services](http://www.restapitutorial.com/lessons/httpmethods.html)

R.I.P. is inspired by Doug's Black ideas in [blog article](http://dougblack.io/words/a-restful-micro-framework-in-go.html)
and [sleepy code](https://github.com/dougblack/sleepy).

