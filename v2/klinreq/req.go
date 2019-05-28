package klinreq

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (r *ReqBuilder) SetUrl(u string) *ReqBuilder {
	r.ReqQ.Url = &u
	return r
}
func (r *ReqBuilder) SetHeaders(h map[string]string) *ReqBuilder {
	r.ReqQ.Headers = h
	return r
}
func (r *ReqBuilder) SetMethod(m string) *ReqBuilder {
	r.ReqQ.Method = &m
	return r
}
func (r *ReqBuilder) NoVerify() *ReqBuilder {
	r.ReqQ.NoVerify = true
	return r
}
func (r *ReqBuilder) SetJson(j interface{}) *ReqBuilder {
	r.ReqQ.Json = &j
	return r
}
func (r *ReqBuilder) Do() (*http.Response, error) {
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

	if r.Json != nil {
		eJson, err := json.Marshal(*r.ReqQ.Json)
		if err != nil {
			return h, err
		}
		ebody = bytes.NewReader(eJson)
	}
	if ebody == nil {
		ebody = []byte("")
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
	client := &http.Client{}
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
