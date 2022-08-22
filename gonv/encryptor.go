package gonv

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"os"

	"github.com/hay-kot/yal"
)

func EncryptFile(path, passphrase string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	encrypted, err := encryptBytes(contents, passphrase)
	if err != nil {
		return err
	}

	return os.WriteFile(path, toB64(encrypted), 0644)
}

func DecryptFile(path, passphrase string) error {
	contents, err := os.ReadFile(path)
	if err != nil {
		yal.Debugf("os.ReadFile(path=%s) failed with error: %s", path, err)
		return err
	}
	contents, err = fromB64(contents)
	if err != nil {
		yal.Debugf("fromB64(contents) failed with error: %s", err)
		return err
	}

	decrypted, err := decryptBytes(contents, passphrase)
	if err != nil {
		yal.Debugf("decryptBytes failed with error: %s", err)
		return err
	}
	return os.WriteFile(path, decrypted, 0644)
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encryptBytes(b []byte, passphrase string) ([]byte, error) {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, b, nil)
	return ciphertext, nil
}

func decryptBytes(b []byte, passphrase string) ([]byte, error) {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		yal.Debugf("aes.NewCipher(key) failed with error: %s", err)
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		yal.Debugf("cipher.NewGCM(block) failed with error: %s", err)
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := b[:nonceSize], b[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		yal.Debugf("gcm.Open(nil, nonce, ciphertext, nil) failed with error: %s", err)
		return nil, err
	}
	return plaintext, nil
}

func toB64(b []byte) []byte {
	return []byte(base64.StdEncoding.EncodeToString(b))
}

func fromB64(b []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(b))
}
