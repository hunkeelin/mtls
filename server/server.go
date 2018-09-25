package klinserver

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Server(c *ServerConfig) error {
	if len(c.CertBytes) != 0 && len(c.KeyBytes) != 0 {
		return listenB(c)
	}
	if c.Cert != "" && c.Key != "" {
		return listen(c)
	}
	if !c.Https {
		return listen(c)
	}
	return errors.New("crt,key incomplete")
}
func listenB(c *ServerConfig) error {
	var err error
	clientCertPool := x509.NewCertPool()
	if c.Verify {
		var certBytes []byte
		var err error
		if len(c.TrustBytes) != 0 {
			certBytes = c.TrustBytes
		} else {
			certBytes, err = ioutil.ReadFile(c.Trust)
			if err != nil {
				return errors.New("unable to read ca file, is the file specified? this is important if you have verify as true")
			}
		}
		if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
			return errors.New("unable to append certbytes")
		}
	}
	var certlist []tls.Certificate
	keycerts, err := tls.X509KeyPair(c.CertBytes, c.KeyBytes)
	if err != nil {
		return errors.New("certbytes and keybytes doesn't match")
	}
	certlist = append(certlist, keycerts)
	tlsconfig := &tls.Config{
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion:   tls.VersionTLS12,
		Certificates: certlist,
	}
	if c.Verify {
		tlsconfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsconfig.ClientCAs = clientCertPool
	}
	s := &http.Server{
		Handler:      c.ServeMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	fmt.Println("listening to " + c.BindAddr + ":" + c.BindPort)
	if c.Https {
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
func listen(c *ServerConfig) error {
	var err error
	clientCertPool := x509.NewCertPool()
	if c.Verify {
		certBytes, err := ioutil.ReadFile(c.Trust)
		if err != nil {
			return err
		}

		if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
			return err
		}
	}
	tlsconfig := &tls.Config{
		PreferServerCipherSuites: true,
		// Only use curves which have assembly implementations
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.X25519,
		},
		MinVersion: tls.VersionTLS12,
	}
	if c.Verify {
		tlsconfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsconfig.ClientCAs = clientCertPool
	}
	s := &http.Server{
		Addr:         c.BindAddr + ":" + c.BindPort,
		TLSConfig:    tlsconfig,
		Handler:      c.ServeMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	fmt.Println("listening to " + c.BindAddr + ":" + c.BindPort)
	if c.Https {
		err = s.ListenAndServeTLS(c.Cert, c.Key)
		if err != nil {
			return err
		}
	} else {
		err = s.ListenAndServe()
		if err != nil {
			return err
		}
	}
	return err
}
