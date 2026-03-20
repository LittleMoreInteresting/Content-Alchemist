package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"os"
)

// Crypto 加密工具
type Crypto struct {
	key []byte
}

// NewCrypto 创建加密工具，使用机器唯一标识生成密钥
func NewCrypto() (*Crypto, error) {
	// 使用机器标识生成密钥
	machineID := getMachineID()
	hash := sha256.Sum256([]byte(machineID))
	return &Crypto{key: hash[:]}, nil
}

// getMachineID 获取机器标识
func getMachineID() string {
	// 尝试从环境变量获取唯一标识
	if id := os.Getenv("COMPUTERNAME"); id != "" {
		return id + "_content_alchemist_salt_v1"
	}
	if id := os.Getenv("HOSTNAME"); id != "" {
		return id + "_content_alchemist_salt_v1"
	}
	return "default_machine_content_alchemist_salt_v1"
}

// Encrypt 加密字符串
func (c *Crypto) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt 解密字符串
func (c *Crypto) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}
	
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
