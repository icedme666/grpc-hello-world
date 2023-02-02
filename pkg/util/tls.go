package util

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"golang.org/x/net/http2"
)

// 用于获取TLS配置
func GetTLSConfig(CertPemPath, CertKeyPath string) *tls.Config {
	var CertKeyPair *tls.Certificate
	cert, _ := ioutil.ReadFile(CertPemPath)
	key, _ := ioutil.ReadFile(CertKeyPath)

	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		log.Println("TLS KeyPair err: %v\n", err)
	}

	CertKeyPair = &pair

	return &tls.Config{
		Certificates: []tls.Certificate{*CertKeyPair},
		NextProtos: []string{http2.NextProtoTLS},
	}
}