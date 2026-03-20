package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/zalando/go-keyring"
)

const (
	keyringService = "Content-Alchemist"
	keyringUser    = "encryption-key"
)

// Crypto 加密工具
type Crypto struct {
	key []byte
}

// NewCrypto 创建加密工具，从系统钥匙串获取或生成密钥
func NewCrypto() (*Crypto, error) {
	// 尝试从系统钥匙串获取密钥
	keyStr, err := keyring.Get(keyringService, keyringUser)
	if err == nil && keyStr != "" {
		// 解码已存储的密钥
		key, err := base64.StdEncoding.DecodeString(keyStr)
		if err == nil && len(key) == 32 {
			return &Crypto{key: key}, nil
		}
	}

	// 钥匙串中没有密钥，需要生成新密钥
	// 注意：调用方需要处理这种情况，引导用户完成初始化
	return nil, fmt.Errorf("encryption key not found in keyring: %w", ErrKeyNotFound)
}

// ErrKeyNotFound 密钥未找到错误
var ErrKeyNotFound = fmt.Errorf("encryption key not found")

// GenerateAndStoreKey 生成新密钥并存储到系统钥匙串
func GenerateAndStoreKey() error {
	// 生成32字节随机密钥
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return fmt.Errorf("generate random key failed: %w", err)
	}

	// Base64 编码后存储
	keyStr := base64.StdEncoding.EncodeToString(key)
	
	// 存储到系统钥匙串
	if err := keyring.Set(keyringService, keyringUser, keyStr); err != nil {
		return fmt.Errorf("store key to keyring failed: %w", err)
	}

	return nil
}

// HasKey 检查系统中是否存在加密密钥
func HasKey() bool {
	_, err := keyring.Get(keyringService, keyringUser)
	return err == nil
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

// DeleteKey 从系统钥匙串删除密钥（谨慎使用）
func DeleteKey() error {
	return keyring.Delete(keyringService, keyringUser)
}
