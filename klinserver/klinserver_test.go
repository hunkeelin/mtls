package klinserver

import (
	"fmt"
	"github.com/hunkeelin/SuperCAclient/lib"
	"github.com/hunkeelin/klinutils"
	"github.com/hunkeelin/pki"
	"net/http"
	"testing"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
	fmt.Println(r.Proto)
}
func TestHttp(t *testing.T) {
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tListen(w, r)
	})
	j := &ServerConfig{
		BindPort: "2018",
		BindAddr: "",
		ServeMux: con,
	}
	err := Server(j)
	if err != nil {
		panic(err)
	}
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
		Caport:       "2018",
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
func TestBserver(t *testing.T) {
	fmt.Println("testing Bserver")
	c := new(conn)
	sema := make(chan struct{}, 1)
	c.monorun = sema
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var body payload
		c.notwork(w, r, &body)
	})
	con.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		var body payloadv2
		c.handleWebHook(w, r, &body)
	})
	con.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		var body payloadv3
		c.getfile(w, r, &body)
	})
	con.HandleFunc("/show", func(w http.ResponseWriter, r *http.Request) {
		c.showreq(w, r)
	})
	r := klinutils.WgetInfo{
		Dest:  "util3.klin-pro.com",
		Dport: "2018",
		Route: "cacerts/intermca.crt",
	}
	cab, err := klinutils.Wget(r)
	if err != nil {
		panic(err)
	}
	w := client.WriteInfo{
		CABytes: cab,
		CAport:  "2018",
		Chain:   true,
		CSRConfig: &klinpki.CSRConfig{
			EmailAddress:       "support@abc.com",
			RsaBits:            2048,
			Country:            "USA",
			Province:           "SHIT",
			Locality:           "NOOB",
			OrganizationalUnit: "IT",
			Organization:       "klin-pro",
		},
		CAName: "util3.klin-pro.com",
	}
	cb, kb, err := client.Getkeycrtbyte(w)
	if err != nil {
		panic(err)
	}
	s := &ServerConfig{
		BindPort:  "2018",
		CertBytes: cb,
		KeyBytes:  kb,
		Trust:     "program/intermca.crt",
		Https:     true,
		//		Verify:    true,
		ServeMux: con,
	}
	err = Server(s)
	if err != nil {
		panic(err)
	}
}
func TestServer(t *testing.T) {
	fmt.Println("testing server")
	c := new(conn)
	sema := make(chan struct{}, 1)
	c.monorun = sema
	con := http.NewServeMux()
	con.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var body payload
		c.notwork(w, r, &body)
	})
	con.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		var body payloadv2
		c.handleWebHook(w, r, &body)
	})
	con.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		var body payloadv3
		c.getfile(w, r, &body)
	})
	con.HandleFunc("/show", func(w http.ResponseWriter, r *http.Request) {
		c.showreq(w, r)
	})
	s := &ServerConfig{
		BindPort: "2018",
		Cert:     "program/test3.klin-pro.com.crt",
		Key:      "program/test3.klin-pro.com.key",
		Trust:    "program/intermca.crt",
		Https:    true,
		Verify:   true,
		ServeMux: con,
	}
	Server(s)
}
