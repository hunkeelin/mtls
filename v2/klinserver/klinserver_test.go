package klinserver

import (
	"fmt"
	"github.com/hunkeelin/klinutils"
	"io/ioutil"
	"net/http"
	"testing"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
	fmt.Println(r.Proto)
}
func TestSserver(t *testing.T) {
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mainHandler(w, r)
	})
	var certs, keys [][]byte
	t1c, err := ioutil.ReadFile("/tmp/t1.crt")
	if err != nil {
		panic(err)
	}
	t1k, err := ioutil.ReadFile("/tmp/t1.key")
	if err != nil {
		panic(err)
	}
	t2c, err := ioutil.ReadFile("/tmp/t2.crt")
	if err != nil {
		panic(err)
	}
	t2k, err := ioutil.ReadFile("/tmp/t2.key")
	if err != nil {
		panic(err)
	}
	certs = append(certs, t2c, t1c)
	keys = append(keys, t2k, t1k)
	j := &ServerConfig{
		BindPort:  "2018",
		BindAddr:  "",
		ServeMux:  con,
		Https:     true,
		CertBytes: certs,
		KeyBytes:  keys,
	}
	panic(Server(j))
}
func TestBBserver(t *testing.T) {
	finish := make(chan bool)
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mainHandler(w, r)
	})
	j := &ServerConfig{
		BindPort: "2018",
		BindAddr: "",
		ServeMux: con,
	}
	jj := &ServerConfig{
		BindPort: "2019",
		BindAddr: "",
		ServeMux: con,
	}
	g := InbytesForm{
		Ca:           "util3.klin-pro.com",
		Caport:       klinutils.Stringtoport("superca"),
		Trustcert:    "intermca.crt",
		Rootca:       "rootca.crt",
		Org:          "klin-pro",
		ServerConfig: j,
	}
	err := Inbytes(g)
	if err != nil {
		panic(err)
	}
	g.ServerConfig = jj
	err = Inbytes(g)
	if err != nil {
		panic(err)
	}
	go func() {
		panic(Server(j))
	}()
	go func() {
		panic(Server(jj))
	}()
	<-finish
}
