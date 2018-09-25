package klinreq

import (
	"fmt"
	"github.com/hunkeelin/SuperCAclient/lib"
	"github.com/hunkeelin/klinutils"
	"github.com/hunkeelin/pki"
	"io/ioutil"
	"testing"
)

type testPayload struct {
	C string `json:"content"`
	D bool   `json:"disabled"`
}

type filejson struct {
	FileBytes []byte `json:"filebytes"`
	FileName  string `json:"filename"`
}

func TestReq(t *testing.T) {
	fmt.Println("testing req")
	payload := &testPayload{
		C: "wtf",
		D: true,
	}
	cb, _ := ioutil.ReadFile("program/test2.klin-pro.com.crt")
	kb, _ := ioutil.ReadFile("program/test2.klin-pro.com.key")
	i := &ReqInfo{
		//		Cert:    "program/test2.klin-pro.com.crt",
		//		Key:     "program/test2.klin-pro.com.key",
		CertBytes: cb,
		KeyBytes:  kb,
		Dest:      "test3.klin-pro.com",
		Dport:     "2018",
		Trust:     "program/rootca.crt",
		Method:    "POST",
		Route:     "foo",
		Payload:   payload,
		Headers: map[string]string{
			"content-type": "map json",
			"api-key":      "default123",
		},
	}
	resp, err := SendPayload(i)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body), string(resp.Status))
}

func TestInbytes(t *testing.T) {
	fmt.Println("testing inbytes")
	payload := &testPayload{
		C: "wtf",
		D: true,
	}
	i := &ReqInfo{
		Dest:  "test3.klin-pro.com",
		Dport: "2018",
		//		Trust:     "program/intermca.crt",
		Method:  "POST",
		Route:   "foo",
		Payload: payload,
		Headers: map[string]string{
			"content-type": "map json",
			"api-key":      "default123",
		},
	}
	f := InbytesForm{
		Ca:        "util3.klin-pro.com",
		Caport:    "2018",
		Trustcert: "intermca.crt",
		Rootca:    "rootca.crt",
		ReqInfo:   i,
	}
	err := Inbytes(f)
	if err != nil {
		panic(err)
	}
	resp, err := SendPayload(i)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body), string(resp.Status))
}
func TestBreq(t *testing.T) {
	fmt.Println("testing Breq")
	payload := &testPayload{
		C: "wtf",
		D: true,
	}
	w := client.WriteInfo{
		CA:     "program/rootca.crt",
		CAport: "2018",
		Chain:  true,
		CSRConfig: &klinpki.CSRConfig{
			EmailAddress:       "support@abc.com",
			RsaBits:            2048,
			Country:            "USA",
			Province:           "SHIT",
			Locality:           "NOOB",
			OrganizationalUnit: "IT",
			Organization:       "klin-pro",
		},
	}
	cb, kb, err := client.Getkeycrtbyte(w)
	if err != nil {
		panic(err)
	}
	f := klinutils.WgetInfo{
		Dest:  "util3.klin-pro.com",
		Dport: "2018",
		Route: "cacerts/intermca.crt",
	}
	b, _ := klinutils.Wget(f)
	i := &ReqInfo{
		CertBytes: cb,
		KeyBytes:  kb,
		Dest:      "test3.klin-pro.com",
		Dport:     "2018",
		//		Trust:     "program/intermca.crt",
		TrustBytes: b,
		Method:     "POST",
		Route:      "foo",
		Payload:    payload,
		Headers: map[string]string{
			"content-type": "map json",
			"api-key":      "default123",
		},
	}
	resp, err := SendPayload(i)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body), string(resp.Status))
}

func TestFuck(t *testing.T) {
	fmt.Println("testing sending file reqv2")
	f, _ := ioutil.ReadFile("program/testfile")
	payload := &filejson{
		FileBytes: f,
		FileName:  "shitname",
	}
	i := &ReqInfo{
		Cert:    "program/test2.klin-pro.com.crt",
		Key:     "program/test2.klin-pro.com.key",
		Dest:    "test3.klin-pro.com",
		Dport:   "2018",
		Trust:   "program/intermca.crt",
		Method:  "POST",
		Route:   "file",
		Payload: payload,
	}
	resp, err := SendPayload(i)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body), string(resp.Status))
}

func TestSendFile(t *testing.T) {
	fmt.Println("testing filesend")
	i := &ReqInfo{
		Cert:   "program/test2.klin-pro.com.crt",
		Key:    "program/test2.klin-pro.com.key",
		Dest:   "test3.klin-pro.com",
		Dport:  "2018",
		Trust:  "program/rootca.crt",
		Method: "POST",
		File:   "program/testfile",
		Route:  "foo",
		ExtraParams: map[string]string{
			"filename": "klinFile",
		},
	}
	resp, err := SendFile(i)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body), string(resp.Status))
}
