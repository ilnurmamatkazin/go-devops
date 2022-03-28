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

func getPrivateKey(file_path string) (*rsa.PrivateKey, error) {
	var private_key *rsa.PrivateKey

	keyBytes, err := ioutil.ReadFile(file_path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)

	private_key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return private_key, nil
}

func pemToCipher(encryptedMessage []byte) []byte {
	b, _ := pem.Decode(encryptedMessage)

	return b.Bytes
}
