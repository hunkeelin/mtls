package klinreq

import (
	"errors"
	"github.com/hunkeelin/SuperCAclient/lib"
	"github.com/hunkeelin/klinutils"
	"github.com/hunkeelin/pki/v1"
)

type InbytesForm struct {
	Ca        string   // the hostname of the SuperCa server
	Caport    string   // the port that SuperCA is being hosted on.
	Trustcert string   // the cert you want to trust, most likely an interm cert definitely not rootca.
	ReqInfo   *ReqInfo // the reqinfo that will get modified.
	Org       string   // the org name for csr
	Rootca    string
}

func Inbytes(in InbytesForm) error {
	r := klinutils.WgetInfo{
		Dest:  in.Ca,
		Dport: in.Caport,
		Route: "cacerts/" + in.Rootca,
	}
	cab, err := klinutils.Wget(r)
	if err != nil {
		return errors.New("unable to wget rootca")
	}
	r = klinutils.WgetInfo{
		Dest:  in.Ca,
		Dport: in.Caport,
		Route: "cacerts/" + in.Trustcert,
	}
	trustb, err := klinutils.Wget(r)
	if err != nil {
		return errors.New("Unable to get " + in.Trustcert)
	}
	w := client.WriteInfo{
		CAName:  in.Ca,
		CABytes: cab,
		CAport:  in.Caport,
		Chain:   true,
		CSRConfig: &klinpki.CSRConfig{
			EmailAddress: "support@" + in.Org + ".com",
			RsaBits:      2048,
			Organization: in.Org,
		},
	}
	cb, kb, err := client.Getkeycrtbyte(w)
	if err != nil {
		errors.New("Unable to get key and crt")
	}
	in.ReqInfo.CertBytes = cb
	in.ReqInfo.KeyBytes = kb
	in.ReqInfo.Http = false
	in.ReqInfo.TrustBytes = trustb
	return nil
}
