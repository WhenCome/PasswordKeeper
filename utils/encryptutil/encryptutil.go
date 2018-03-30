package encryptutil

import (
	"crypto/rsa"
	"crypto/x509"
	"crypto/rand"
	"encoding/pem"
	"os"
	"encoding/base64"
	"../fileutil"
	"../../config"
)

// 生成RSA证书信息，并返回证书配置内容
func GenRsaKey(bits int, storePath string) (*config.CertConfig, error) {
	cfg := &config.CertConfig{}
	cfg.PrivateKeyFile = fileutil.BuildPath(storePath, "private.pem")
	cfg.PublicKeyFile = fileutil.BuildPath(storePath, "public.pem")
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create(cfg.PrivateKeyFile)
	if err != nil {
		return nil,err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return nil,err
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil,err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create(cfg.PublicKeyFile)
	if err != nil {
		return nil,err
	}
	err = pem.Encode(file, block)
	if err != nil {
		return nil,err
	}
	return cfg,nil
}

// 加密数据
func EncryptData(data string) (string, error) {
	encBytes, err := rsaEncrypt([]byte(data))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(encBytes), nil
}

// 解密密数据
func DecryptData(data string) (string, error) {
	byteData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	decBytes, err := rsaDecrypt(byteData)
	if err != nil {
		return "", err
	}
	return string(decBytes), nil
}