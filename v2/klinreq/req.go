package klinreq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (r *ReqBuilder) SetUrl(u string) *ReqBuilder {
	r.reqQ.Url = &u
	return r
}
func (r *ReqBuilder) SetHeaders(h map[string]string) *ReqBuilder {
	r.reqQ.Headers = h
	return r
}
func (r *ReqBuilder) SetMethod(m string) *ReqBuilder {
	r.reqQ.Method = &m
	return r
}
func (r *ReqBuilder) SetJson(j interface{}) *ReqBuilder {
	r.reqQ.Json = &j
	return r
}
func (r *ReqBuilder) Do() (*http.Response, error) {
	var (
		h     *http.Response
		ebody *bytes.Reader
	)

	err := r._check()
	if err != nil {
		return h, err
	}

	if r.Json != nil {
		eJson, err := json.Marshal(*r.reqQ.Json)
		if err != nil {
			return h, err
		}
		ebody = bytes.NewReader(eJson)
	}
	req, err := http.NewRequest(*r.reqQ.Method, *r.reqQ.Url, ebody)
	if err != nil {
		return h, err
	}

	if r.reqQ.Headers != nil {
		for k, v := range r.reqQ.Headers {
			req.Header.Set(k, v)
		}
	}
	client := &http.Client{}
	h, err = client.Do(req)
	if err != nil {
		return h, err
	}
	return h, nil
}

func (r *ReqBuilder) _check() error {
	// make GET as default
	if r.reqQ.Method == nil {
		method := "GET"
		r.reqQ.Method = &method
	}

	// check if url is valid
	_, err := url.ParseRequestURI(*r.reqQ.Url)
	if err != nil {
		return err
	}
	if r.reqQ.Url == nil {
		return fmt.Errorf("url not valid")
	}
	return nil
}
