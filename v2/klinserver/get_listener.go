package klinserver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"golang.org/x/net/http2"
	"net"
)

func GetListener(c *ServerConfig) (net.Listener, error) {
	var err error
	var certlist []tls.Certificate
	var l net.Listener
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
				return l, fmt.Errorf("certbytes and keybytes doesn't match %v", err)
			}
			certlist = append(certlist, keycerts)
		}
		tlsconfig.Certificates = certlist
	}
	if c.Verify && c.Https {
		if c.TrustBytes == nil {
			return l, fmt.Errorf("No trust bytes being provided")
		}
		trustlist := c.TrustBytes
		for _, trustca := range trustlist {
			if ok := clientCertPool.AppendCertsFromPEM(trustca); !ok {
				return l, fmt.Errorf("unable to append certbytes")
			}
		}
		tlsconfig.ClientAuth = tls.RequireAndVerifyClientCert
		tlsconfig.ClientCAs = clientCertPool
	}
	tlsconfig.BuildNameToCertificate()
	if c.Name2cert != nil {
		for hostname, keycrt := range c.Name2cert {
			keycerts, err := tls.X509KeyPair(keycrt.Cb, keycrt.Kb)
			if err != nil {
				return l, fmt.Errorf("Unable to loadkeypair for name2cert key value of "+hostname+"%v", err)
			}
			// if key does have a value already
			if _, haveVal := tlsconfig.NameToCertificate[hostname]; haveVal {
				if c.SNIoverride {
					tlsconfig.NameToCertificate[hostname] = &keycerts
					tlsconfig.Certificates = append(tlsconfig.Certificates, keycerts)
				} else {
					return l, fmt.Errorf("Overiding an sni existing key value " + hostname + " if you are sure what you are doing enable SNIoverride")
				}
			}
		}
	}
	l, err = tls.Listen("tcp", c.BindAddr+":"+c.BindPort, tlsconfig)
	if err != nil {
		return l, fmt.Errorf("unable to listen to port and address %v", err)
	}
	return l, err
}
