package klinreq

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

func New() ReqBuilder {
	client := &http.Client{}
	r := ReqQ{
		Client: client,
	}

	return ReqBuilder{
		ReqQ: r,
	}
}
func (r *ReqBuilder) SetUrl(b string) *ReqBuilder {
	r.ReqQ.BodyBytes = &u
	return r
}
func (r *ReqBuilder) SetBodyBytes(b []byte) *ReqBuilder {
	r.ReqQ.Url = &b
	return r
}
func (r *ReqBuilder) SetHeaders(h map[string]string) *ReqBuilder {
	r.ReqQ.Headers = h
	return r
}
func (r *ReqBuilder) SetTimeOut(h int) *ReqBuilder {
	r.ReqQ.Client.Timeout = time.Duration(h) * time.Second
	return r
}
func (r *ReqBuilder) SetMethod(m string) *ReqBuilder {
	r.ReqQ.Method = &m
	return r
}
func (r *ReqBuilder) NoVerify() *ReqBuilder {
	r.ReqQ.NoVerify = true
	return b
}
func (r *ReqBuilder) SetJson(j interface{}) *ReqBuilder {
	r.ReqQ.Json = &j
	return r
}
func (r *ReqBuilder) Do() (*http.Response, error) {
	var client *http.Client
	if r.ReqQ.Client == nil {
		client = &http.Client{}
	} else {
		client = r.ReqQ.Client
	}
	var (
		h     *http.Response
		ebody *bytes.Reader
	)
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = r.ReqQ.NoVerify
	err := r._check()
	if err != nil {
		return h, err
	}

	// check if json exist
	if r.Json != nil {
		eJson, err := json.Marshal(*r.ReqQ.Json)
		if err != nil {
			return h, err
		}
		ebody = bytes.NewReader(eJson)
	} else {
		ebody = bytes.NewReader([]byte(""))
	}
	// check if bodyBytes exist bodybytes overwrite everything
	if r.BodyBytes != nil {
		ebody = bytes.NewReader(r.BodyBytes)
	}
	req, err := http.NewRequest(*r.ReqQ.Method, *r.ReqQ.Url, ebody)
	if err != nil {
		return h, err
	}

	if r.ReqQ.Headers != nil {
		for k, v := range r.ReqQ.Headers {
			req.Header.Set(k, v)
		}
	}
	client.Transport = &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	h, err = client.Do(req)
	if err != nil {
		return h, err
	}
	return h, nil
}

func (r *ReqBuilder) _check() error {
	// make GET as default
	if r.ReqQ.Method == nil {
		method := "GET"
		r.ReqQ.Method = &method
	}

	// check if url is valid
	_, err := url.ParseRequestURI(*r.ReqQ.Url)
	if err != nil {
		return err
	}
	if r.ReqQ.Url == nil {
		return fmt.Errorf("url not valid")
	}
	return nil
}
