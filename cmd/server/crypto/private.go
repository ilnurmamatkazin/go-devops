package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
)

func Decrypt(privateKeyPath string, data []byte) ([]byte, error) {
	privateKey, err := getPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	cipher, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, pemToCipher(data))
	if err != nil {
		return nil, err
	}

	return cipher, nil
}

func getPrivateKey(filePath string) (*rsa.PrivateKey, error) {
	var privateKey *rsa.PrivateKey

	keyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func pemToCipher(encryptedMessage []byte) []byte {
	b, _ := pem.Decode(encryptedMessage)

	return b.Bytes
}
