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
	fmt.Println("testing sni")
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mainHandler(w, r)
	})
	var certs, keys [][]byte
	t2c, err := ioutil.ReadFile("/tmp/test2.klin-pro.com.crt")
	if err != nil {
		panic(err)
	}
	t2k, err := ioutil.ReadFile("/tmp/test2.klin-pro.com.key")
	if err != nil {
		panic(err)
	}
	t1c, err := ioutil.ReadFile("/tmp/test1.klin-pro.com.crt")
	if err != nil {
		panic(err)
	}
	t1k, err := ioutil.ReadFile("/tmp/test1.klin-pro.com.key")
	if err != nil {
		panic(err)
	}
	certs = append(certs, t1c, t2c)
	keys = append(keys, t1k, t2k)
	j := &ServerConfig{
		BindPort: "2018",
		BindAddr: "",
		ServeMux: con,
		//	Https:     true,
		CertBytes: certs,
		KeyBytes:  keys,
		Name2cert: map[string]Keycrt{
			"test1.klin-pro.com": Keycrt{
				Cb: t2c,
				Kb: t2k,
			},
			"test2.klin-pro.com": Keycrt{
				Cb: t1c,
				Kb: t1k,
			},
		},
		SNIoverride: true,
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
