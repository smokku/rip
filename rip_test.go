package rip

import (
	"testing"

	"net/http"
	"net/http/httptest"
)

type fooResource struct {
}
type fooSubResource struct {
}
type bazResource struct {
}

func TestMain(t *testing.T) {
	root := New().
		Add("foos", fooResource{}).
		Add("bazes", bazResource{})
	//root.Get("foos").
	//	Add("bars", fooSubResource{})

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

		rsp := httptest.NewRecorder()
		root.ServeHTTP(rsp, req)
		t.Logf("%d - %s", rsp.Code, rsp.Body.String())
	}
}
