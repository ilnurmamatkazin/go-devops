package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	// http://www.inanzzz.com/index.php/post/1nw3/data-encryption-and-decryption-with-x509-public-and-private-key-example-in-golang
	// https://developpaper.com/question/how-to-use-go-to-realize-rsa-public-key-encryption-and-private-key-decryption/
	// https://gist.github.com/miguelmota/3ea9286bd1d3c2a985b67cac4ba2130a
	// https://gist.github.com/jshap70/259a87a7146393aab5819873a193b88c
	// https://russianblogs.com/article/46489585/
	// https://russianblogs.com/article/6656426244/#_52

	key, err := New()
	if err != nil {
		log.Fatalln(err)
	}

	if err = key.PublicKeyToFile(); err != nil {
		log.Fatalln(err)
	}

	if err = key.PrivateKeyToFile(); err != nil {
		log.Fatalln(err)
	}
}

type Key struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func New() (Key, error) {
	var k Key

	privateKey, err := rsa.GenerateKey(rand.Reader, 8192)
	if err != nil {
		return k, err
	}

	k.publicKey = &privateKey.PublicKey
	k.privateKey = privateKey

	return k, nil
}

func (k Key) PublicKeyToFile() error {
	file, err := os.Create("public.pem")
	if err != nil {
		return err
	}

	pem.Encode(
		file,
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(k.publicKey),
		},
	)

	return nil
}

func (k Key) PrivateKeyToFile() error {
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}

	pem.Encode(
		file,
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(k.privateKey),
		},
	)

	return nil
}
