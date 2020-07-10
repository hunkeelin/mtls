package klinserver

import (
	"errors"
	"github.com/hunkeelin/SuperCAclient/lib"
	"github.com/hunkeelin/klinutils"
	"github.com/hunkeelin/pki/v1"
)

// The form to modify ServerConfig so it will take everything in bytes, this form is only useful when you have a SuperCA server up and running.
type InbytesForm struct {
	Ca           string        // the hostname of the SuperCa server
	Caport       string        // the port that SuperCA is being hosted on.
	Trustcert    string        // the cert you want to trust, most likely an interm cert definitely not rootca.
	ServerConfig *ServerConfig // the server config that will get modified.
	Org          string        // the org name for csr
	Rootca       string        // the rootca, for SuperCA its default as rootca.crt
}

//function only useful if you have a SuperCA server up and running and understand what it does.
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
	tb, err := klinutils.Wget(r)
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
	var cbarr, kbarr, tbarr [][]byte
	cbarr = append(cbarr, cb)
	kbarr = append(kbarr, kb)
	tbarr = append(tbarr, tb)
	in.ServerConfig.CertBytes = cbarr
	in.ServerConfig.KeyBytes = kbarr
	in.ServerConfig.TrustBytes = tbarr
	in.ServerConfig.Https = true
	return nil
}
