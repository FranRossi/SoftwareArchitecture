package main

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"
	"voting_simulator/proto"
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
	privateKey, publicKey := generateKeyPair(2048)

	fmt.Printf("Private key: %v\n", privateKey)
	fmt.Printf("Public Key: %v", publicKey)

	// Create PEM string
	privKeyStr := exportPrivKeyAsPEMStr(privateKey)
	pubKeyStr := exportPubKeyAsPEMStr(publicKey)

	fmt.Println(privKeyStr)
	fmt.Println(pubKeyStr)

	saveKeyToFile(privKeyStr, "privkey.pem")
	saveKeyToFile(pubKeyStr, "pubkey.pem")
	SignVote(privateKey)
}

func SignVote(privateKey *rsa.PrivateKey) {
	vote := VoteModel{
		IdElection:  "1",
		IdVoter:     "10000000",
		IdCandidate: "1",
		Circuit:     "5",
	}
	candidate := []byte(vote.IdCandidate)
	msgHash := sha256.New()
	msgHash.Write(candidate)
	msgHashSBytes := msgHash.Sum(nil)
	signature, _ := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSBytes, nil)
	vote.Signature = string(signature)
	Vote(vote)

}

type VoteModel struct {
	IdElection  string
	IdVoter     string
	Circuit     string
	IdCandidate string
	Signature   string
}

const addr = "localhost:50004"

func Vote(vote VoteModel) {
	var opts []grpc.DialOption
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	conn.Connect()
	client := proto.NewVoteServiceClient(conn)
	request := &proto.VoteRequest{
		IdElection:  vote.IdElection,
		IdVoter:     vote.IdVoter,
		Circuit:     vote.Circuit,
		IdCandidate: vote.IdCandidate,
		Signature:   vote.Signature,
	}
	client.Vote(context.Background(), request)
}
