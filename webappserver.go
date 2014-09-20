package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var _ = fmt.Print

type WebAppServer struct {
	Subdomain string `yaml:"subdomain"`
	WebAppURL string `yaml:"url"`
	ImageUrl  string `yaml:"image"`
	Name      string `yaml:"name"`
	Domain    string `yaml:"default_domain,omitempty"`
}

func (this *WebAppServer) Localize(domain string) *WebAppServer {
	return &WebAppServer{this.Subdomain, this.WebAppURL, this.ImageUrl, this.Name, domain}
}

type jsonWebAppServer struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Image string `json:"image"`
}

func (this *WebAppServer) copyHeader(from, to http.Header) {
	for key, val := range from {
		to.Set(key, strings.Join(val, ";"))
	}
}

func (this *WebAppServer) GetWebAppBaseUrl() string {
	return fmt.Sprintf("http://%s.%s", this.Subdomain, this.Domain)
}

func (this *WebAppServer) MarshalJSON() ([]byte, error) {
	ret := &jsonWebAppServer{this.Name,
		this.GetWebAppBaseUrl(),
		this.ImageUrl}
	return json.Marshal(ret)
}

func (this *WebAppServer) HandleGet(w http.ResponseWriter, request *http.Request) {
	URL := this.WebAppURL + request.URL.Path
	if request.URL.RawQuery != "" {
		URL += "?" + request.URL.RawQuery
	}
	req, err := http.NewRequest(request.Method, URL, request.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "bab error :(", err.Error())
		return
	}
	this.copyHeader(request.Header, req.Header)
	resp, err := new(http.Client).Do(req)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprint(w, "bab error :(\n", err.Error())
		return
	}
	this.copyHeader(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func (this *WebAppServer) Handle(w http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "GET":
		this.HandleGet(w, request)
	case "POST":
		this.HandleGet(w, request)
	case "PUT":
		this.HandleGet(w, request)
	default:
		w.WriteHeader(500)
		fmt.Fprint(w, "method "+request.Method+" is not yet available")
	}

}

func (this *WebAppServer) CanHandle(host string) bool {
	r := regexp.MustCompile(`[:\.]`)
	domain := r.Split(host, -1)
	if len(domain) < 3 {
		return false
	}
	return domain[len(domain)-3] == this.Subdomain
}
