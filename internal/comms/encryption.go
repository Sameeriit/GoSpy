package comms

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"golang.org/x/crypto/scrypt"
)

// byteEncryptor controls the encryption and decryption of bytes.
type byteEncryptor struct {
	passwordBytes []byte
}

// newByteEncryptor instantiates a new byteEncryptor.
func newByteEncryptor(password string) byteEncryptor {
	return byteEncryptor{[]byte(password)}
}

// Encrypt takes a slice of bytes and encrypts is using the given password.
func (e byteEncryptor) Encrypt(data []byte) (ciphertext []byte, err error) {
	key, salt, err := e.deriveKey(nil)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	// The nonce is prepended to the ciphertext and the salt is appended.
	ciphertext = gcm.Seal(nonce, nonce, data, nil)
	ciphertext = append(ciphertext, salt...)
	return ciphertext, nil
}

// Decrypt takes a slice of bytes and decrypts is using the given password.
func (e byteEncryptor) Decrypt(data []byte) (plaintext []byte, err error) {
	salt, data := data[len(data)-32:], data[:len(data)-32]

	key, _, err := e.deriveKey(salt)
	if err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err = gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// deriveKey takes a salt (can be nil), and returns a 32 byte secure key along with the salt for the given password.
func (e byteEncryptor) deriveKey(salt []byte) (key []byte, generatedSalt []byte, err error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err = scrypt.Key(e.passwordBytes, salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
