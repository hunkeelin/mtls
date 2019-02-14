package klinreq

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"golang.org/x/net/http2"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ReqInfo struct {
	Cert               string // The cert for mtls
	Key                string // The key for mtls
	CertBytes          []byte // same as cert but in bytes will overwrite Cert
	KeyBytes           []byte // same as key but in bytes will overwrite Key
	Dest               string // The destination address. It has to be hostname
	Dport              string // The destination address port
	TimeOut            int
	Trust              string // The trusted CA cert
	TrustBytes         []byte
	Method             string // The req method, POST/PATCH etc...
	Route              string // The route, by default its "/" it can be "/api"
	File               string // If you are sending file specify the file you are sending.
	Http               bool
	Headers            map[string]string
	ExtraParams        map[string]string
	Payload            interface{}
	Xml                bool
	InsecureSkipVerify bool
	BodyBytes          []byte // raw bytes of the body, will overwrite what' sin payload.
	HttpVersion        int
	Rawreq             *http.Request
}

// Send a json payload. payload should be a struct where you define your json
func SendPayload(i *ReqInfo) (*http.Response, error) {
	client := &http.Client{}
	if i.Rawreq != nil {
		resp, err := client.Do(i.Rawreq)
		if err != nil {
			return resp, fmt.Errorf("Client do error pkg-klinreq %v", err)
		}
		return resp, nil
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var resp *http.Response
	var cert tls.Certificate
	var certlist []tls.Certificate
	var err error
	if i.Cert != "" && i.Key != "" {
		cert, err = tls.LoadX509KeyPair(i.Cert, i.Key)
		if err != nil {
			return resp, err
		}
	}
	if len(i.CertBytes) != 0 && len(i.KeyBytes) != 0 {
		cert, err = tls.X509KeyPair(i.CertBytes, i.KeyBytes)
		if err != nil {
			return resp, err
		}
	}
	certlist = append(certlist, cert)

	// Load our CA certificate
	var clientCACert []byte
	if i.Trust != "" {
		var err error
		clientCACert, err = ioutil.ReadFile(i.Trust)
		if err != nil {
			return resp, err
		}
	}
	if len(i.TrustBytes) != 0 {
		clientCACert = i.TrustBytes
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: i.InsecureSkipVerify,
		Certificates:       certlist,
	}
	if i.Trust != "" {
		tlsConfig.RootCAs = clientCertPool
	}
	if len(i.TrustBytes) != 0 {
		tlsConfig.RootCAs = clientCertPool
	}
	if i.HttpVersion == 0 {
		i.HttpVersion = 2
	}
	switch i.HttpVersion {
	case 1:
		client.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	case 2:
		client.Transport = &http2.Transport{
			TLSClientConfig: tlsConfig,
		}
	default:
		return resp, errors.New("Wrong https version please specify 1 or 2 you specified" + strconv.Itoa(i.HttpVersion))
	}
	if i.TimeOut == 0 {
		client.Timeout = time.Duration(20000) * time.Millisecond
	} else {
		client.Timeout = time.Duration(i.TimeOut) * time.Millisecond
	}
	var encodepayload []byte
	if i.Xml {
		encodepayload, err = xml.Marshal(i.Payload)
	} else {
		encodepayload, err = json.Marshal(i.Payload)
	}
	if err != nil {
		panic(err)
	}
	var addr string
	if len(i.Route) > 0 {
		if string(i.Route[0]) != "/" {
			i.Route = "/" + i.Route
		}
	}
	var portinfo string
	if i.Dport == "" {
		portinfo = ""
	} else {
		portinfo = ":" + i.Dport
	}
	if i.Http {
		addr = "http://" + i.Dest + portinfo + i.Route
	} else {
		addr = "https://" + i.Dest + portinfo + i.Route
	}
	var ebody *bytes.Reader
	if len(i.BodyBytes) > 0 {
		ebody = bytes.NewReader(i.BodyBytes)
	} else {
		ebody = bytes.NewReader(encodepayload)
	}
	req, err := http.NewRequest(i.Method, addr, ebody)
	if err != nil {
		return resp, fmt.Errorf("Error making new request pkg-klinreq%v", err)
	}
	for k, v := range i.Headers {
		req.Header.Set(k, v)
	}
	req.Header.Del("Accept-Encoding")
	resp, err = client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("client do error pkg-klinreq%v", err)
	}
	return resp, nil
}
