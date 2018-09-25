package klinreq

import (
	"io/ioutil"
)

type WgetInfo struct {
	Dest  string // Destination aka hostname/ip
	Dport string // The destination port
	Route string // the route of the file you try to get
}

func Wget(w WgetInfo) ([]byte, error) {
	var body []byte
	var err error
	j := &ReqInfo{
		Dest:               w.Dest,
		Dport:              w.Dport,
		TimeOut:            1500,
		Method:             "GET",
		Route:              w.Route,
		InsecureSkipVerify: true,
	}
	resp, err := SendPayload(j)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return body, err
	}
	return body, err
}
