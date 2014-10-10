package rip

import (
	"testing"

	"net/http"
	"net/http/httptest"
)

type fooResource struct {
	ResourceBase
}
type fooSubResource struct {
	ResourceBase
}
type bazResource struct {
	ResourceBase
}

func TestMain(t *testing.T) {
	root := New().
		AddResource("foos", new(fooResource)).
		AddResource("bazes", new(bazResource))
	root.GetResource("foos").
		AddResource("bars", new(fooSubResource))

	type MainRequestTest struct {
		method  string
		url     string
		success bool
	}
	var mainRequestTests = []MainRequestTest{
		{GET, "/foos", true},
		{GET, "/foos/123", true},
		{POST, "/foos", true},
		{GET, "/notthere", false},
		{GET, "/foos/none", false},
		{DELETE, "/foos/123", true},
		{GET, "/foos/bars", true},
		{GET, "/foos/123/bazes", true},
		{POST, "/foos/123/bazes", true},
	}

	for _, d := range mainRequestTests {
		req, err := http.NewRequest(d.method, d.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()
		root.ServeHTTP(w, req)
		t.Logf("%d - %s", w.Code, w.Body.String())
	}
}
