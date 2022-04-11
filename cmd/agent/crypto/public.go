package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func Encrypt(publicKeyPath string, data []byte) ([]byte, error) {
	publicKey, err := newPublicKey(publicKeyPath)
	if err != nil {
		return nil, err
	}

	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}

	return cipherToPem(cipher), nil
}

func newPublicKey(filePath string) (*rsa.PublicKey, error) {
	var publicKey *rsa.PublicKey

	keyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)

	publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return publicKey, nil
}

func cipherToPem(cipher []byte) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "MESSAGE",
			Bytes: cipher,
		},
	)
}
