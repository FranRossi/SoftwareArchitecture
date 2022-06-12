package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
)

func exportPEMStrToPrivKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
	return key
}

// Decode public key struct from PEM string
func exportPEMStrToPubKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	key, _ := x509.ParsePKCS1PublicKey(block.Bytes)
	return key
}

// Read data from file
func readKeyFromFile(filename string) []byte {
	key, _ := ioutil.ReadFile(filename)
	return key
}

func main() {
	privKeyPEM := readKeyFromFile("privkey.pem")
	privKeyFile := exportPEMStrToPrivKey(privKeyPEM)
	fmt.Printf("Private key: %v\n", privKeyFile)
}
