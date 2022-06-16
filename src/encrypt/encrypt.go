package encrypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	b64 "encoding/base64"
	"fmt"
	"os"
)

func DecryptText(encryptedText string) string {
	privateKeyPEM := ReadKeyFromFile("./../encrypt/privkey.pem")
	privKey := ExportPEMStrToPrivKey(privateKeyPEM)
	rng := rand.Reader

	bytesToDecrypt, _ := b64.StdEncoding.DecodeString(encryptedText)
	decryptedBytes, err := rsa.DecryptOAEP(sha256.New(), rng, privKey, bytesToDecrypt, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)
	}

	textInBase64 := string(decryptedBytes)
	textInBytes, _ := b64.StdEncoding.DecodeString(textInBase64)
	plaintext := string(textInBytes)
	return plaintext
}

func EncryptText(text string) string {

	publicKeyPEM := ReadKeyFromFile("./../encrypt/pubkey.pem")
	publicKey := ExportPEMStrToPubKey(publicKeyPEM)

	textInBase64 := b64.StdEncoding.EncodeToString([]byte(text))
	secretMessage := []byte(textInBase64)
	rng := rand.Reader

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, publicKey, secretMessage, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		panic(err)
	}
	ciphertextBase64 := b64.StdEncoding.EncodeToString(ciphertext)
	return ciphertextBase64
}
