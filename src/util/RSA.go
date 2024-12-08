package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"pm-ssl-management/src/log"
)

func GenerateRSAKeyPair() (private string, public string) {
	var privateKey *rsa.PrivateKey
	var err error
	// 生成 RSA 私钥
	privateKey, err = rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		log.ComLoggerClient.Error("生成 RSA 私钥失败: ", err.Error())
		return
	}
	// 将私钥转换为 PEM 格式
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyPEMBytes := pem.EncodeToMemory(privateKeyPEM)

	var publicKeyASN1 []byte
	// 提取公钥
	publicKey := &privateKey.PublicKey
	// 将公钥转换为 ASN.1 DER 格式
	publicKeyASN1, err = x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		log.ComLoggerClient.Error("生成 RSA 公钥失败: ", err.Error())
		return
	}
	// 将公钥的 ASN.1 DER 格式转换为 PEM 格式
	publicKeyPEM := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyASN1,
	}
	publicKeyPEMBytes := pem.EncodeToMemory(publicKeyPEM)

	private = string(privateKeyPEMBytes)
	public = string(publicKeyPEMBytes)
	return
}

func Encrypt(plainText string, private string) (cipherText string) {
	privateKeyPEM := []byte(private)
	// 解析私钥
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		log.ComLoggerClient.Error("failed to decode private key")
		return
	}
	var privateKey *rsa.PrivateKey
	var err error
	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.ComLoggerClient.Error("failed to decode private key: ", err.Error())
		return
	}
	// 使用私钥解密数据
	plainTextByte, _ := base64.StdEncoding.DecodeString(plainText)
	cipherTextByte, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, plainTextByte)
	if err != nil {
		log.ComLoggerClient.Error("failed to decrypt data: ", err.Error())
		return
	}
	cipherText = string(cipherTextByte)
	return
}
