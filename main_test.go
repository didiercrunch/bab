package main

import (
	"testing"
)

func TestLoadWebAppServersFromYamlFile(t *testing.T) {
	yml := []byte(`
---

-
    name: mongs
    url: http://localhost:3333
    subdomain: mongs
    image: images/mongo.png

-
    name: ipython
    url: http://localhost:8888
    subdomain: ipython
    image: images/ipython.png
`)
	if was, err := loadWebAppServersFromYamlFile(yml); err != nil {
		t.Error(err)
	} else if len(was) != 2 {
		t.Fail()
	}
}
