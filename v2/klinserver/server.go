package klinserver

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"golang.org/x/net/http2"
	"net/http"
	"time"
)

func Server(c *ServerConfig) error {
	if c.CertBytes == nil || c.KeyBytes == nil || len(c.KeyBytes) != len(c.CertBytes) {
		return errors.New("crt,key incomplete, either you didn't provide them or the number of key and cert doesn't match")
	}
	return listenB(c)
}
func listenB(c *ServerConfig) error {
	var err error
	var certlist []tls.Certificate
	clientCertPool := x509.NewCertPool()
	tlsconfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		NextProtos: []string{http2.NextProtoTLS},
	}
	if c.Https {
		for i, _ := range c.CertBytes {
			keycerts, err := tls.X509KeyPair(c.CertBytes[i], c.KeyBytes[i])
			if err != nil {
				return errors.New("certbytes and keybytes doesn't match")
			}
			certlist = append(certlist, keycerts)
		}
		tlsconfig.Certificates = certlist
	}
	if c.Verify && c.Https {
		if c.TrustBytes == nil {
			return fmt.Errorf("No trust bytes being provided")
		}
		trustlist := c.TrustBytes
		for _, trustca := range trustlist {
			if ok := clientCertPool.AppendCertsFromPEM(trustca); !ok {
				return errors.New("unable to append certbytes")
			}
		}
		tlsconfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsconfig.ClientCAs = clientCertPool
	}
	var r, w, i int
	if c.ReadTimeout == 0 {
		r = 5
	} else {
		r = c.ReadTimeout
	}
	if c.WriteTimeout == 0 {
		w = 10
	} else {
		w = c.WriteTimeout
	}
	if c.IdleTimeout == 0 {
		i = 120
	} else {
		i = c.IdleTimeout
	}
	s := &http.Server{
		Handler:      c.ServeMux,
		ReadTimeout:  time.Duration(r) * time.Second,
		WriteTimeout: time.Duration(w) * time.Second,
		IdleTimeout:  time.Duration(i) * time.Second,
	}
	fmt.Println("listening to " + c.BindAddr + ":" + c.BindPort)
	if c.Https {
		if c.Name2cert == nil {
			tlsconfig.BuildNameToCertificate()
		} else {
			tlsconfig.NameToCertificate = make(map[string]*tls.Certificate)
			for hostname, keycrt := range c.Name2cert {
				keycerts, err := tls.X509KeyPair(keycrt.Cb, keycrt.Kb)
				if err != nil {
					return fmt.Errorf("Unable to loadkeypair for name2cert %v", err)
				}
				tlsconfig.NameToCertificate[hostname] = &keycerts
				tlsconfig.Certificates = append(tlsconfig.Certificates, keycerts)
			}
		}
		l, err := tls.Listen("tcp", c.BindAddr+":"+c.BindPort, tlsconfig)
		if err != nil {
			return errors.New("unable to listen to port and address")
		}
		err = s.Serve(l)
		if err != nil {
			return errors.New("unable to serve https")
		}
	} else {
		s.Addr = c.BindAddr + ":" + c.BindPort
		err = s.ListenAndServe()
		if err != nil {
			return errors.New("unable to serve http")
		}
	}
	return err
}
