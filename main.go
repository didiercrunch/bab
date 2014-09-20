package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"launchpad.net/goyaml"
	"log"
	"net/http"
	"path"
)

const helloworld string = `
 ____  ____  ____  ____  ____  ____  ____  ____  ____  ____  ____  ____ 
||s ||||e ||||r ||||v ||||i ||||n ||||g ||||: ||||8 ||||0 ||||0 ||||0 ||
||__||||__||||__||||__||||__||||__||||__||||__||||__||||__||||__||||__||
|/__\||/__\||/__\||/__\||/__\||/__\||/__\||/__\||/__\||/__\||/__\||/__\|

`

type WebAppServers []*WebAppServer

func (this WebAppServers) Handle(w http.ResponseWriter, request *http.Request) {
	for _, was := range this {
		if was.Subdomain == mux.Vars(request)["subdomain"] {
			was.Handle(w, request)
			return
		}
	}
	fmt.Fprint(w, "could not find a web application for ")
}

func (this WebAppServers) Localize(domain string) WebAppServers {
	fmt.Println(">>", domain)
	ret := make(WebAppServers, len(this))
	for i, was := range this {
		ret[i] = was.Localize(domain)
	}
	return ret
}

var webAppServers WebAppServers
var port int = 8000

func init() {
	var err error
	webAppServers, err = loadWebAppServers()
	if err != nil {
		panic(err)
	}
}

func loadWebAppServers() (WebAppServers, error) {
	data, err := ioutil.ReadFile("webapps.yml")
	if err != nil {
		return nil, err
	}
	return loadWebAppServersFromYamlFile(data)

}

func loadWebAppServersFromYamlFile(yamlFile []byte) (WebAppServers, error) {
	ret := make(WebAppServers, 0, 10)
	err := goyaml.Unmarshal(yamlFile, &ret)
	return ret, err

}

func serveStatic(router *mux.Router) {
	handler := func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		filepath := "/" + vars["path"]
		w.Header().Set("Cache-Control", "public, max-age=43200")
		http.ServeFile(w, request, path.Join("client/", filepath))
	}
	router.HandleFunc("/{path:.*}", handler)
}

func ServeWebApps(router *mux.Router) {
	router.HandleFunc("/{path:.*}", webAppServers.Handle)
}

func WebappHandler(w http.ResponseWriter, request *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(webAppServers.Localize(request.Host))
}

func createMuxRouter() http.Handler {
	r := mux.NewRouter()
	ServeWebApps(r.Host(`{subdomain:\w+}.{domain:\w+}.{topleveldomain:[a-z]+}`).Subrouter())
	ServeWebApps(r.Host(`{subdomain:\w+}.localhost`).Subrouter())
	ServeWebApps(r.Host(`{subdomain:\w+}.\d{1-3}.\d{1-3}.\d{1-3}.\d{1-3}`).Subrouter())
	r.HandleFunc("/webapps", WebappHandler)
	serveStatic(r.PathPrefix("/").Subrouter())
	return r
}

func main() {
	fmt.Print(helloworld)
	address := fmt.Sprintf("0.0.0.0:%v", port)
	if err := http.ListenAndServe(address, createMuxRouter()); err != nil {
		log.Fatal(err)
	}
}
