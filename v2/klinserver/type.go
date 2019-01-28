package klinserver

import (
	"net/http"
)

type Keycrt struct {
	Cb, Kb []byte
}
type ServerConfig struct {
	BindAddr     string
	BindPort     string
	Cert         string   //the location of the .crt for https
	Key          string   // the location of the .key for https
	CertBytes    [][]byte // the .crt in bytes will take preceding over Cert
	KeyBytes     [][]byte // the .key in bytes will take preceding over key
	Trust        string   // trust cert location
	TrustBytes   [][]byte // trust cert in  bytes will take preceding over Trust
	Https        bool     // whether to host in https or not
	Verify       bool
	ReadTimeout  int
	WriteTimeout int
	IdleTimeout  int
	ServeMux     *http.ServeMux    // the http.ServeMux
	Name2cert    map[string]Keycrt // key == hostname, value == cert in bytes
}
