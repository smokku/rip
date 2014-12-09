package rip

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

const (
	GET     = "GET"
	POST    = "POST"
	PUT     = "PUT"
	PATCH   = "PATCH"
	DELETE  = "DELETE"
	HEAD    = "HEAD"
	OPTIONS = "OPTIONS"
)

// GetSupported is the interface that provides the Get
// method a resource must support to receive HTTP GETs.
type GetSupported interface {
	Get(string, *http.Request) (int, interface{}, http.Header)
}

// PostSupported is the interface that provides the Post
// method a resource must support to receive HTTP POSTs.
type PostSupported interface {
	Post(string, *http.Request) (int, interface{}, http.Header)
}

// PatchSupported is the interface that provides the Patch
// method a resource must support to receive HTTP PATCHs.
type PatchSupported interface {
	Patch(string, *http.Request) (int, interface{}, http.Header)
}

// PutSupported is the interface that provides the Put
// method a resource must support to receive HTTP PUTs.
type PutSupported interface {
	Put(string, *http.Request) (int, interface{}, http.Header)
}

// DeleteSupported is the interface that provides the Delete
// method a resource must support to receive HTTP DELETEs.
type DeleteSupported interface {
	Delete(string, *http.Request) (int, interface{}, http.Header)
}

// HeadSupported is the interface that provides the Head
// method a resource must support to receive HTTP HEADs.
type HeadSupported interface {
	Head(string, *http.Request) (int, interface{}, http.Header)
}

// HeadSupported is the interface that provides the Head
// method a resource must support to receive HTTP HEADs.
type OptionsSupported interface {
	Options(string, *http.Request) (int, interface{}, http.Header)
}

func requestHandler(resource interface{}) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {

		if req.ParseForm() != nil {
			resp.WriteHeader(http.StatusBadRequest)
			return
		}

		var handler func(string, *http.Request) (int, interface{}, http.Header)

		switch req.Method {
		case GET:
			if resource, ok := resource.(GetSupported); ok {
				handler = resource.Get
			}
		case POST:
			if resource, ok := resource.(PostSupported); ok {
				handler = resource.Post
			}
		case PUT:
			if resource, ok := resource.(PutSupported); ok {
				handler = resource.Put
			}
		case PATCH:
			if resource, ok := resource.(PatchSupported); ok {
				handler = resource.Patch
			}
		case DELETE:
			if resource, ok := resource.(DeleteSupported); ok {
				handler = resource.Delete
			}
		case HEAD:
			if resource, ok := resource.(HeadSupported); ok {
				handler = resource.Head
			}
		case OPTIONS:
			if resource, ok := resource.(OptionsSupported); ok {
				handler = resource.Options
			}
		}

		var code = http.StatusOK
		var content []byte
		var header http.Header

		if handler == nil {
			code = http.StatusMethodNotAllowed
		} else if req.URL.Path[0] == '/' {
			pathTokens := strings.Split(req.URL.Path, "/")
			var id string
			if len(pathTokens) > 2 {
				id = pathTokens[2]
			}

			var data interface{}
			code, data, header = handler(id, req)

			if header == nil {
				header = make(http.Header)
			}

			switch v := data.(type) {
			case nil:
				// pass
			case string:
				content = []byte(v)
			case int, int8, int16, int32, int64:
				content = []byte(strconv.Itoa(data.(int)))
			case uint, uint8, uint16, uint32, uint64:
				content = []byte(strconv.FormatUint(data.(uint64), 10))
			case float32, float64:
				content = []byte(strconv.FormatFloat(data.(float64), 'f', -1, 64))
			case bool:
				if data.(bool) {
					content = []byte("true")
				} else {
					content = []byte("false")
				}
			default:
				var err error
				content, err = json.MarshalIndent(data, "", "  ")
				if err != nil {
					code = http.StatusInternalServerError
				} else {
					header.Set("Content-Type", "application/json")
				}
			}

			if len(content) > 0 && len(header.Get("Content-Type")) == 0 {
				header.Set("Content-Type", "text/plain")
			}
		} else {
			code = http.StatusBadRequest
		}

		for name, values := range header {
			for _, value := range values {
				resp.Header().Add(name, value)
			}
		}

		resp.WriteHeader(code)

		if len(content) == 0 && code != http.StatusOK {
			content = []byte(http.StatusText(code))
		}
		resp.Write(content)
	}
}
