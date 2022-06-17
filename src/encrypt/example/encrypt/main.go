package main

import (
	"crypto/rand"
	"crypto/rsa"
	"electoral_service/models"
	"encrypt"
	"fmt"
	"io/ioutil"
)

func main() {

	//	encryptedText := "BDq9eDQUapkK/93zLbGGCxgqCNNHb1eJmIQQOXubvbGRR2KQs8MQc5QrzHnN9cQv3i3ifxoZuZpzQYKFPFsAIjmdRlCOYFZ7EdXuh8EHvhQ4vxVrSdfIge1aoUfg5ZlhYEaucUhA/H8EGDRbHcb5iHsIdmDQKKl/wECHjD9LGp23EzYpeC0Bsdj7Y+F4leqHadvCxEWaRRIW5bcHaUXftnvU44Cx7kt3OIiE6Hp+HhtmohTaz43tTE/D3USM0RzXc/U4PnVra7PRNFd/QP178uI4cxm3ug4zkysEU7JTBeruPTJrbu/abIxpyfgpA7m/z76K36+wwF0iaTOKSijKQA=="
	encryptedText := encrypt.EncryptText("hola")
	decryptedText := encrypt.DecryptText(encryptedText)
	fmt.Println(encryptedText)
	fmt.Println(decryptedText)
	//Generate a 2048-bits key
	// privateKey, publicKey := generateKeyPair(2048)
	// //Create PEM string
	// privKeyStr := encrypt.ExportPrivKeyAsPEMStr(privateKey)
	// pubKeyStr := encrypt.ExportPubKeyAsPEMStr(publicKey)

	// saveKeyToFile(privKeyStr, "privkey.pem")
	// saveKeyToFile(pubKeyStr, "pubkey.pem")

	voter := models.VoterModel{
		Id:        "123",
		FullName:  "John Doe",
		Sex:       "Male",
		BirthDate: "1/1/2000",
		Voted:     2,
		Email:     "example@gmail.com",
		OtherFields: map[string]any{
			"Test":       "cxvcv",
			"Numero":     234,
			"Bool":       true,
			"otro valor": []int{2, 3, 2},
		},
	}
	encrypt.EncryptVoter(&voter)
	fmt.Println(voter)
	encrypt.DecryptVoter(&voter)
	fmt.Println(voter)
}

func generateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	// This method requires a random number of bits.
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// The public key is part of the PrivateKey struct
	return privateKey, &privateKey.PublicKey
}

// Save string to a file
func saveKeyToFile(keyPem, filename string) {
	pemBytes := []byte(keyPem)
	ioutil.WriteFile(filename, pemBytes, 0400)
}
