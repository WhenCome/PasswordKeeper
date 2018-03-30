package encryptutil

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"errors"
	"../../config"
)

// 加密
func rsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(config.GetPublicKey())
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func rsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(config.GetPrivateKey())
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}