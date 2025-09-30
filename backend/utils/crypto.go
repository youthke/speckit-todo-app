package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// CryptoService handles encryption and decryption of sensitive data
type CryptoService struct {
	key []byte
}

// NewCryptoService creates a new crypto service from environment
func NewCryptoService() (*CryptoService, error) {
	key := os.Getenv("JWT_SECRET")
	if key == "" {
		return nil, errors.New("JWT_SECRET environment variable is not set")
	}

	// Use first 32 bytes of JWT_SECRET as encryption key
	// In production, use a separate ENCRYPTION_KEY
	keyBytes := []byte(key)
	if len(keyBytes) < 32 {
		// Pad key if too short
		padded := make([]byte, 32)
		copy(padded, keyBytes)
		keyBytes = padded
	} else if len(keyBytes) > 32 {
		// Truncate if too long
		keyBytes = keyBytes[:32]
	}

	return &CryptoService{
		key: keyBytes,
	}, nil
}

// Encrypt encrypts a plaintext string using AES-GCM
func (s *CryptoService) Encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a ciphertext string using AES-GCM
func (s *CryptoService) Decrypt(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptToken encrypts an OAuth token for storage
func (s *CryptoService) EncryptToken(token string) (string, error) {
	return s.Encrypt(token)
}

// DecryptToken decrypts an OAuth token from storage
func (s *CryptoService) DecryptToken(encryptedToken string) (string, error) {
	return s.Decrypt(encryptedToken)
}

// HashPassword hashes a password using bcrypt
// Note: This is a placeholder - actual password hashing should use bcrypt library
func HashPassword(password string) (string, error) {
	// In production, use golang.org/x/crypto/bcrypt
	// For now, this is a placeholder
	return password, errors.New("password hashing not implemented - use bcrypt")
}

// VerifyPassword verifies a password against a hash
func VerifyPassword(password, hash string) bool {
	// In production, use bcrypt.CompareHashAndPassword
	// For now, this is a placeholder
	return false
}

// GenerateSecureToken generates a cryptographically secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// GenerateSecureRandomBytes generates cryptographically secure random bytes
func GenerateSecureRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	return bytes, nil
}

// EncryptOAuthTokens encrypts OAuth access and refresh tokens
func EncryptOAuthTokens(accessToken, refreshToken string) (string, string, error) {
	crypto, err := NewCryptoService()
	if err != nil {
		return "", "", err
	}

	encryptedAccess, err := crypto.EncryptToken(accessToken)
	if err != nil {
		return "", "", err
	}

	encryptedRefresh, err := crypto.EncryptToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	return encryptedAccess, encryptedRefresh, nil
}

// DecryptOAuthTokens decrypts OAuth access and refresh tokens
func DecryptOAuthTokens(encryptedAccess, encryptedRefresh string) (string, string, error) {
	crypto, err := NewCryptoService()
	if err != nil {
		return "", "", err
	}

	accessToken, err := crypto.DecryptToken(encryptedAccess)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := crypto.DecryptToken(encryptedRefresh)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}