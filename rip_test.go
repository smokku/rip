package rip

import (
	"testing"

	"net/http"
	"net/http/httptest"
)

type mainResource struct {
	ResourceBase
}
type mainSubResource struct {
	ResourceBase
}
type subResource struct {
	ResourceBase
}

func TestMain(t *testing.T) {
	user := AddResource("foos", new(mainResource))
	user.AddResource("bars", new(mainSubResource))
	user.AddResource("bazes", new(subResource))

	req, err := http.NewRequest("GET", "http://example.com/foos", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	user.ServeHTTP(w, req)
	t.Logf("%d - %s", w.Code, w.Body.String())
}
