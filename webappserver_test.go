package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestServer struct {
	Server *httptest.Server
	Get    func(w http.ResponseWriter, request *http.Request)
	Post   func(w http.ResponseWriter, request *http.Request)
	T      *testing.T
}

func (this *TestServer) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		this.Get(w, request)
	} else if request.Method == "POST" {
		this.Post(w, request)
	} else {
		this.T.Error("cannot execute ", request.Method, " method")
	}
}

func (this *TestServer) HelloWorld(w http.ResponseWriter, request *http.Request) {
	fmt.Fprint(w, "hello world\n")
}

func (this *TestServer) ReturnPath(w http.ResponseWriter, request *http.Request) {
	fmt.Fprint(w, request.URL.Path)
}

func (this *TestServer) ReturnBogusHeader(w http.ResponseWriter, request *http.Request) {
	fmt.Fprint(w, request.Header.Get("Bogus"))
}

func (this *TestServer) ReturnBody(w http.ResponseWriter, request *http.Request) {
	io.Copy(w, request.Body)

}

func (this *TestServer) SetBogusHeader(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Bogus", "has been set on the server")
	fmt.Fprint(w, request.Header.Get("Bogus"))
}

func (this *TestServer) ReturnParamValue(w http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		this.T.Error(err)
		return
	}
	fmt.Fprint(w, request.Form.Get("search"))
}

func GetTestServer(t *testing.T) *TestServer {
	ret := new(TestServer)
	ret.T = t
	server := httptest.NewServer(ret)
	ret.Server = server
	return ret
}

func TestMarshalJSON(t *testing.T) {
	wapp := &WebAppServer{"big", "http://bigtits.com", "img.png", "some name", "tits.com:8000"}
	b, err := json.Marshal(wapp)
	if err != nil {
		t.Error(err)
	}
	if s := string(b); !strings.Contains(s, `"http://big.tits.com:8000"`) {
		t.Fail()
	}
}

func TestCanHandle(t *testing.T) {
	was := &WebAppServer{Subdomain: "foo"}
	exps := map[string]bool{"foo.bigtits.com": true,
		"foo.localhost:8000": true,
		"bar.bigtits.com":    false,
		"bigtits.com":        false}
	for host, exp := range exps {
		if was.CanHandle(host) != exp {
			t.Error("error handeling ", host)
		}
	}
}

func TestHandleGetHelloWorld(t *testing.T) {
	ts := GetTestServer(t)
	ts.Get = ts.HelloWorld
	was := &WebAppServer{Subdomain: "foo", WebAppURL: ts.Server.URL}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://foo.bigtits.com/natasha", nil)
	if err != nil {
		t.Error(err)
	}
	was.HandleGet(w, req)
	if body := string(w.Body.Bytes()); body != "hello world\n" {
		t.Error(body)
	}
}

func TestHandleGetPath(t *testing.T) {
	ts := GetTestServer(t)
	ts.Get = ts.ReturnPath
	was := &WebAppServer{Subdomain: "foo", WebAppURL: ts.Server.URL}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://foo.bigtits.com/natasha", nil)
	if err != nil {
		t.Error(err)
	}
	was.HandleGet(w, req)
	if body := string(w.Body.Bytes()); body != "/natasha" {
		t.Error(body)
	}
}

func TestHandleGetBogusHeader(t *testing.T) {
	ts := GetTestServer(t)
	ts.Get = ts.ReturnBogusHeader
	was := &WebAppServer{Subdomain: "foo", WebAppURL: ts.Server.URL}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://foo.bigtits.com/natasha", nil)
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Bogus", "something in the header!")
	was.HandleGet(w, req)
	if body := string(w.Body.Bytes()); body != "something in the header!" {
		t.Error(body)
	}
}

func TestHandleGetSetBogusHeader(t *testing.T) {
	ts := GetTestServer(t)
	ts.Get = ts.SetBogusHeader
	was := &WebAppServer{Subdomain: "foo", WebAppURL: ts.Server.URL}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://foo.bigtits.com/natasha", nil)
	if err != nil {
		t.Error(err)
	}
	was.HandleGet(w, req)
	if header := w.Header().Get("Bogus"); header != "has been set on the server" {
		t.Error(header)
	}
}

func TestHandleGetReturnParamValue(t *testing.T) {
	ts := GetTestServer(t)
	ts.Get = ts.ReturnParamValue
	was := &WebAppServer{Subdomain: "foo", WebAppURL: ts.Server.URL}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://foo.bigtits.com/natasha?search=bigtits", nil)
	if err != nil {
		t.Error(err)
	}
	was.HandleGet(w, req)
	if body := string(w.Body.Bytes()); body != "bigtits" {
		t.Error(body)
	}
}

func TestHandlePostWithBody(t *testing.T) {
	ts := GetTestServer(t)
	ts.Post = ts.ReturnBody
	was := &WebAppServer{Subdomain: "foo", WebAppURL: ts.Server.URL}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "http://foo.bigtits.com/natasha", bytes.NewBuffer([]byte(`hello world in the body`)))
	if err != nil {
		t.Error(err)
	}
	was.Handle(w, req)
	if body := string(w.Body.Bytes()); body != "hello world in the body" {
		t.Error("wrong body :", body)
	}
}

func TestHandlePostBogusHeader(t *testing.T) {
	ts := GetTestServer(t)
	ts.Post = ts.SetBogusHeader
	was := &WebAppServer{Subdomain: "foo", WebAppURL: ts.Server.URL}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "http://foo.bigtits.com/natasha", nil)
	if err != nil {
		t.Error(err)
	}
	was.HandleGet(w, req)
	if header := w.Header().Get("Bogus"); header != "has been set on the server" {
		t.Error(header)
	}
}

func TestCanHandleNoService(t *testing.T) {
	ts := GetTestServer(t)

	was := &WebAppServer{Subdomain: "foo", WebAppURL: ts.Server.URL}
	ts.Server.Close()
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://foo.bigtits.com/natasha?search=bigtits", nil)
	if err != nil {
		t.Error(err)
	}
	was.HandleGet(w, req)
	if w.Code != 500 {
		t.Error("status should be 500, service not found")
	}
}
