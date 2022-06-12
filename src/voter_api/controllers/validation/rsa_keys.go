package validation

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

// Decode public key struct from PEM string
func ExportPEMStrToPubKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		fmt.Println("failed to decode PEM block containing public key")
	}
	key, err2 := x509.ParsePKCS1PublicKey(block.Bytes)
	if err2 != nil {
		fmt.Println("Error on Parsing:", err2)
	}
	return key
}

// Read data from file
func ReadKeyFromFile(filename string) []byte {
	key, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Error on Reading file:", err)
	}
	return key
}
