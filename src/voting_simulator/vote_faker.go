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
	"google.golang.org/grpc/credentials/insecure"
	"io/ioutil"
	"log"
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
	//privateKey, publicKey := generateKeyPair(2048)
	//publicKey = publicKey
	// Create PEM string
	//privKeyStr := exportPrivKeyAsPEMStr(privateKey)
	//pubKeyStr := exportPubKeyAsPEMStr(publicKey)

	//saveKeyToFile(privKeyStr, "privkey.pem")
	//saveKeyToFile(pubKeyStr, "pubkey.pem")
	privateKeyPEM := ReadKeyFromFile("./privkey.pem")
	privateKey := ExportPEMStrToPrivKey(privateKeyPEM)
	SignVote(privateKey)
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

func SignVote(privateKey *rsa.PrivateKey) {
	vote := VoteModel{
		IdElection:  "1",
		IdVoter:     "10000000",
		IdCandidate: "2",
		Circuit:     "5",
	}
	voter := []byte(vote.IdVoter)
	msgHash := sha256.New()
	msgHash.Write(voter)
	msgHashSBytes := msgHash.Sum(nil)
	signature, _ := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSBytes, nil)
	vote.Signature = signature
	Vote(vote)

}

type VoteModel struct {
	IdElection  string
	IdVoter     string
	Circuit     string
	IdCandidate string
	Signature   []byte
}

const addr = "localhost:50004"

func Vote(vote VoteModel) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("cannot connect: %s", err)
	}
	defer conn.Close()

	client := proto.NewVoteServiceClient(conn)
	request := &proto.VoteRequest{
		IdElection:  vote.IdElection,
		IdVoter:     vote.IdVoter,
		Circuit:     vote.Circuit,
		IdCandidate: vote.IdCandidate,
		Signature:   vote.Signature,
	}
	response, err2 := client.Vote(context.Background(), request)
	if err2 != nil {
		log.Fatalf("could not vote: %v", err2)
	}
	fmt.Printf("Vote: %s\n", response.Message)
}
