package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/joho/godotenv"
)

func generateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	// This method requires a random number of bits.
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// The public key is part of the PrivateKey struct
	return privateKey, &privateKey.PublicKey
}

func exportPubKeyAsPEMStr(pubkey *rsa.PublicKey) string {
	pubKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pubkey),
		},
	))
	return pubKeyPem
}

// Export private key as a string in PEM format
func exportPrivKeyAsPEMStr(privkey *rsa.PrivateKey) string {
	privKeyPem := string(pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(privkey),
		},
	))
	return privKeyPem
}

// Save string to a file
func saveKeyToFile(keyPem, filename string) {
	pemBytes := []byte(keyPem)
	ioutil.WriteFile(filename, pemBytes, 0400)
}

func main() {
	// Generate a 2048-bits key
	//privateKey, publicKey := generateKeyPair(2048)
	//publicKey = publicKey
	// Create PEM string
	//privKeyStr := exportPrivKeyAsPEMStr(privateKey)
	//pubKeyStr := exportPubKeyAsPEMStr(publicKey)

	//saveKeyToFile(privKeyStr, "privkey.pem")
	//saveKeyToFile(pubKeyStr, "pubkey.pem")
	godotenv.Load()
	CreateVotes()
	fmt.Println("Press Enter to exit")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

func getPrivateKey() *rsa.PrivateKey {
	privateKeyPEM := ReadKeyFromFile("./privkey.pem")
	privateKey := ExportPEMStrToPrivKey(privateKeyPEM)
	return privateKey
}

func ExportPEMStrToPrivKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)
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

func encryptVote(vote *VoteModel) {
	publicKey := getPublicKey()

	vote.IdVoter = encryptText(vote.IdVoter, publicKey)
	vote.IdCandidate = encryptText(vote.IdCandidate, publicKey)
	vote.IdElection = encryptText(vote.IdElection, publicKey)
}

func getPublicKey() *rsa.PublicKey {
	publicKeyPEM := ReadKeyFromFile("./pubkey_appEV.pem")
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		fmt.Println("failed to decode PEM block containing public key")
	}
	publicKey, err2 := x509.ParsePKCS1PublicKey(block.Bytes)
	if err2 != nil {
		fmt.Println("Error on Parsing:", err2)
	}

	return publicKey
}

func encryptText(text string, publicKey *rsa.PublicKey) string {

	secretMessage := []byte(text)
	rng := rand.Reader

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, publicKey, secretMessage, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		panic(err)
	}
	ciphertextBase64 := b64.StdEncoding.EncodeToString(ciphertext)
	return ciphertextBase64
}
