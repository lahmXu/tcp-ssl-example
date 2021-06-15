package main

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
)

func main() {
	log.SetFlags(log.Lshortfile)
	cer, err := tls.LoadX509KeyPair("VerifyClientCert/server.crt", "VerifyClientCert/server.key")
	if err != nil {
		log.Println(err)
		return
	}
	config := &tls.Config{
		ClientCAs:    loadCA1("VerifyClientCert/ca.crt"),
		Certificates: []tls.Certificate{cer},
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	ln, err := tls.Listen("tcp", ":8000", config)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("server start, listen:8000")
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		println(msg)
		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}
func loadCA1(fileName string) *x509.CertPool {
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
