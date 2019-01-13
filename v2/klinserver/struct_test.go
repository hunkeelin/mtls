package klinserver

import (
	"sync"
)

type conn struct {
	regex   string
	apikey  string
	pkidir  string
	concur  int
	jobdir  string
	mu      *sync.Mutex
	monorun chan struct{}
}

type payload struct {
	C string `json:"content"`
	D bool   `json:"disabled"`
}

type payloadv2 struct {
	C string `json:"content"`
	D bool   `json:"disabled"`
}
type payloadv3 struct {
	FileBytes []byte `json:"filebytes"`
	FileName  string `json:"filename"`
}
type dowork interface {
	doit() string
}

type doworkv2 interface {
	dothat() (string, []byte)
}

func (m *payload) doit() string {
	return m.C
}
func (m *payloadv2) doit() string {
	return m.C + "noob"
}
func (m *payloadv3) dothat() (string, []byte) {
	return m.FileName, m.FileBytes
}
