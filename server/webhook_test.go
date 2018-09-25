package klinserver

import (
	//	"encoding/json"
	"fmt"
	"github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func (c *conn) notwork(w http.ResponseWriter, r *http.Request, p dowork) {
	msg := "nothing since it's foo"
	status := 400
	fmt.Println(msg, status)
	w.WriteHeader(status)
	w.Write([]byte(msg))
	return
}
func (c *conn) getfile(w http.ResponseWriter, r *http.Request, p doworkv2) {
	msg := "Got payload"
	status := 200
	fmt.Println(msg, status)
	f, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(f, &p)
	if err != nil {
		panic(err)
	}
	filename, filebytes := p.dothat()
	towrite, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer towrite.Close()
	towrite.Write(filebytes)
	w.WriteHeader(status)
	w.Write([]byte(msg))
	fmt.Println(r.Header.Get("content-type"))
	return
}
func (c *conn) showreq(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(r.Header)
	fmt.Println(string(body))
	return
}
func (c *conn) handleWebHook(w http.ResponseWriter, r *http.Request, p dowork) {
	if strings.HasPrefix(r.Header.Get("content-type"), "multipart/form-data") {
		t, _, _ := r.FormFile("file")
		to, _ := os.Create("shit")
		io.Copy(to, t)
		to.Close()
		msg := "file transferred"
		status := 200
		fmt.Println(msg, status)
		w.WriteHeader(status)
		w.Write([]byte(msg))
		fmt.Println(r.FormValue("filename"))
		return
	} else {
		msg := "Got payload"
		status := 200
		fmt.Println(msg, status)
		f, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(f, &p)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(f))
		w.WriteHeader(status)
		w.Write([]byte(msg))
		fmt.Println(r.Header.Get("content-type"))
		return
	}
}
