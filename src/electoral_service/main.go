package main

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	dependencyinjection "electoral_service/service/dependencies_injection"
	"fmt"
)

func main() {
	electoral_service := dependencyinjection.Injection()
	electoral_service.DropDataBases()
	electoral_service.GetElectionSettings()
}

func DecryptVote(publicKeyVoter *rsa.PublicKey, signature []byte) {
	msgHash := sha256.New()
	msgHash.Write([]byte("1"))
	msgHashSum := msgHash.Sum(nil)
	err := rsa.VerifyPSS(publicKeyVoter, crypto.SHA256, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("Verification failed: ", err)
	} else {
		fmt.Println("Message verified.")
	}
}
