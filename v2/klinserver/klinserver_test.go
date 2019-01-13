package klinserver

import (
	"fmt"
	"github.com/hunkeelin/klinutils"
	"net/http"
	"testing"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
	fmt.Println(r.Proto)
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
