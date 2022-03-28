package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func Encrypt(publicKeyPath string, data []byte) ([]byte, error) {
	publicKey, err := getPublicKey(publicKeyPath)
	if err != nil {
		return nil, err
	}

	cipher, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, data)
	if err != nil {
		return nil, err
	}

	return cipherToPem(cipher), nil
}

func getPublicKey(file_path string) (*rsa.PublicKey, error) {
	var public_key *rsa.PublicKey

	keyBytes, err := ioutil.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)

	public_key, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return public_key, nil
}

func cipherToPem(cipher []byte) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "MESSAGE",
			Bytes: cipher,
		},
	)
}
