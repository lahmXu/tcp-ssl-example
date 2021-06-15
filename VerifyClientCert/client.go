package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	pair, e := tls.LoadX509KeyPair("VerifyClientCert/client.crt", "VerifyClientCert/client.key")
	if e != nil {
		log.Fatal("LoadX509KeyPair:", e)
	}
	conf := &tls.Config{
		RootCAs:      loadCA("VerifyClientCert/ca.crt"),
		Certificates: []tls.Certificate{pair},
	}
	conn, err := tls.Dial("tcp", "localhost:8000", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	n, err := conn.Write([]byte("hello\n"))
	if err != nil {
		log.Println(n, err)
		return
	}
	buf := make([]byte, 100)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}
	println(string(buf[:n]))
}

func loadCA(fileName string) *x509.CertPool {
	pool := x509.NewCertPool()
	caCertPath := fileName

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return nil
	}
	pool.AppendCertsFromPEM(caCrt)

	return pool
}
