package myhttp

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"os"
	"time"
)

type CertType int

const (
	CertTypeServer CertType = iota
	CertTypeClient
)

func GenerateRSAKey() (*rsa.PrivateKey, error) {
	rsaPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}
	return rsaPrivKey, nil
}

func GenerateCACert(privKey *rsa.PrivateKey, valid int) ([]byte, error) {
	sn, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}
	ca := &x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization:  []string{"Company Inc"},
			Country:       []string{"US"},
			Province:      []string{"PP"},
			Locality:      []string{"SF"},
			StreetAddress: []string{"Golden"},
			PostalCode:    []string{"11111"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(valid, 0, 0),
		IsCA:      true,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	caBytes, err := x509.CreateCertificate(
		rand.Reader, ca, ca,
		&privKey.PublicKey, privKey,
	)
	if err != nil {
		return nil, err
	}
	return caBytes, nil
}

func GenerateCertWithCA(
	certType CertType,
	valid int,
	caPrivKey *rsa.PrivateKey,
	peerPublicKey *rsa.PublicKey,
	caCertBytes []byte,
	ipSANs []string,
) ([]byte, error) {
	sn, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}
	ipSANList := []net.IP{}
	for _, ipStr := range ipSANs {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, errors.New("Failed to parse ipSAN")
		}
		ipSANList = append(ipSANList, ip)
	}
	caCert, err := x509.ParseCertificate(caCertBytes)
	if err != nil {
		return nil, err
	}
	cert := &x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization:  []string{"Company Inc"},
			Country:       []string{"US"},
			Province:      []string{"PP"},
			Locality:      []string{"SF"},
			StreetAddress: []string{"Golden"},
			PostalCode:    []string{"22222"},
		},
		IPAddresses: ipSANList,
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(valid, 0, 0),
		//SubjectKeyId: []byte{1, 2, 3, 4, 6},
		//ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		//KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
	}
	switch certType {
	case CertTypeClient:
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
		cert.KeyUsage = x509.KeyUsageDigitalSignature
	case CertTypeServer:
		cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		cert.KeyUsage = x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment
	default:
		return nil, errors.New("Unrecognized cert type")
	}
	certBytes, err := x509.CreateCertificate(
		rand.Reader, cert, caCert,
		peerPublicKey, caPrivKey,
	)
	if err != nil {
		return nil, err
	}
	return certBytes, nil
}

func WriteCertAndKey(name string, certBytes []byte, key *rsa.PrivateKey) error {
	keyPEM := new(bytes.Buffer)
	if err := pem.Encode(keyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		return err
	}
	certPEM := new(bytes.Buffer)
	if err := pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}); err != nil {
		return err
	}
	if err := os.WriteFile(name+".key", keyPEM.Bytes(), 0600); err != nil {
		return err
	}
	if err := os.WriteFile(name+".crt", certPEM.Bytes(), 0644); err != nil {
		return err
	}
	return nil
}
