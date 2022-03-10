package halo

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func ResolveAddr(addr string) (*Addr, error) {
	addrs := strings.Split(addr, ":")
	if len(addrs) != 2 {
		return nil, errors.New("ip parse error: " + addr)
	}
	ip := net.ParseIP(addrs[0])
	if ip.To4() == nil {
		return nil, errors.New("ip parse error: " + addr)
	}
	if port, err := strconv.Atoi(addrs[1]); err != nil {
		return nil, errors.New("ip parse error: " + err.Error())
	} else {
		ipAddr := [4]byte{}
		copy(ipAddr[:], ip.To4())
		return &Addr{
			IP:   ipAddr,
			Port: port,
		}, nil
	}
}

func S2SockAddrInet4(addr string) (*unix.SockaddrInet4, error) {
	addrs := strings.Split(addr, ":")
	if len(addrs) != 2 {
		return nil, errors.New("ip parse error: " + addr)
	}
	ip := net.ParseIP(addrs[0])
	if ip.To4() == nil {
		return nil, errors.New("ip parse error: " + addr)
	}
	if port, err := strconv.Atoi(addrs[1]); err != nil {
		return nil, errors.New("ip parse error: " + err.Error())
	} else {
		sa := unix.SockaddrInet4{Port: port}
		copy(sa.Addr[:], ip.To4())
		return &sa, nil
	}
}

func GenerateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"secret"},
	}
}
