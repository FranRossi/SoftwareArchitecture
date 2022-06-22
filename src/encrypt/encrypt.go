package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"
	"os"
	l "own_logger"
)

func DecryptText(encryptedText string) string {
	privKey := GetInstancePrivateKey()
	rng := rand.Reader

	bytesToDecrypt, _ := b64.StdEncoding.DecodeString(encryptedText)
	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rng, privKey, bytesToDecrypt, nil)
	if err != nil {
		l.LogError("Error from decryption: " + err.Error())
	}

	plaintext := string(decryptedBytes)
	return plaintext
}

func EncryptText(text string) string {

	publicKeyPEM := ReadKeyFromFile("./../encrypt/pubkey_appEV.pem")
	publicKey := ExportPEMStrToPubKey(publicKeyPEM)

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

var privateKeyInstance = readPrivateKey()

func GetInstancePrivateKey() *rsa.PrivateKey {
	if privateKeyInstance == nil {
		privateKeyInstance = readPrivateKey()
	}
	return privateKeyInstance
}

func readPrivateKey() *rsa.PrivateKey {
	privateKeyPEM := ReadKeyFromFile("./../encrypt/privkey_appEV.pem")
	privKey := ExportPEMStrToPrivKey(privateKeyPEM)
	return privKey
}
