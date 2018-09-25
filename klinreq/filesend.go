package klinreq

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func customRequest(uri, paramName, path, method string, params map[string]string) (*http.Request, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, uri, body)
	return req, err

}

// Send a file via octet-stream. The File field in ReqInfo must be completed
func SendFile(i *ReqInfo) (*http.Response, error) {
	var resp *http.Response
	req, err := customRequest("https://"+i.Dest+":"+i.Dport+"/"+i.Route, "file", i.File, i.Method, i.ExtraParams)
	for k, v := range i.Headers {
		req.Header.Set(k, v)
	}
	if err != nil {
		return resp, err
	}
	cert, err := tls.LoadX509KeyPair(i.Cert, i.Key)
	if err != nil {
		return resp, errors.New("Unable to load cert and key, check if they exist or if they are a valid pair")
	}

	// Load our CA certificate
	clientCACert, err := ioutil.ReadFile(i.Trust)
	if err != nil {
		log.Fatal("Unable to open cert", err)
	}

	clientCertPool := x509.NewCertPool()
	clientCertPool.AppendCertsFromPEM(clientCACert)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: i.InsecureSkipVerify,
		Certificates:       []tls.Certificate{cert},
		RootCAs:            clientCertPool,
	}
	tr := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{
		Timeout:   500 * time.Millisecond,
		Transport: tr,
	}
	resp, err = client.Do(req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
