package myhttp

import (
	"bytes"
	"crypto/x509"
	"testing"
)

func TestCertSerialization(t *testing.T) {
	k1, _ := GenerateRSAKey()
	c1, _ := GenerateCACert(k1, 10)
	WriteCertAndKey("test", c1, k1)

	c2, k2, _ := LoadCertAndKey("test")
	if !bytes.Equal(c1, c2) {
		t.Fatal("cert load failed")
	}

	if !bytes.Equal(
		x509.MarshalPKCS1PrivateKey(k1),
		x509.MarshalPKCS1PrivateKey(k2),
	) {
		t.Fatal("key load failed")
	}
}
