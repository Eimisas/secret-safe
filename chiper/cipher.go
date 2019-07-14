package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

// Encrypt takes key and plaintext and returns an encrypted (hex) version of it
func Encrypt(key, plaintext string) (string, error) {

	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream, err := stream(key, iv)
	if err != nil {
		return "", err
	}
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plaintext))

	return fmt.Sprintf("%x", cipherText), nil
}

func decryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := createCipherBlock(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBDecrypter(block, iv), nil

}

// Decrypt takes key and cipherHex and decrypts it, returning a string
func Decrypt(key, cipherHex string) (string, error) {

	cipherText, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("Decrypt: cipherText is too short, no IV!!")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream, err := decryptStream(key, iv)
	if err != nil {
		return "", err
	}
	// Ciphertex replace inplace
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil

}

// DecryptReader will return a reader that decodes the data from provided reader
func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if n < len(iv) || err != nil {
		errors.New("DecryptWriter: unable to read the full iv")
	}
	stream, err := decryptStream(key, iv)
	if err != nil {
		return nil, err
	}
	return &cipher.StreamReader{
		S: stream, R: r,
	}, nil

}

func stream(key string, iv []byte) (cipher.Stream, error) {
	block, err := createCipherBlock(key)
	if err != nil {
		return nil, err
	}

	return cipher.NewCFBEncrypter(block, iv), nil

}

// EncryptWriter will return a writer that will write encrypted data to original writer
func EncryptWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream, err := stream(key, iv)
	if err != nil {
		return nil, err
	}
	n, err := w.Write(iv)
	if n != len(iv) || err != nil {
		return nil, errors.New("EncryptWriter: unable to write full iv to writer")
	}
	return &cipher.StreamWriter{S: stream, W: w}, nil
}

func createCipherBlock(key string) (cipher.Block, error) {
	// Hasher hash the input key in order to normalize it's size to 16bit
	// as NewCipher function takes specific length of AES key
	hasher := md5.New()
	_, err := io.WriteString(hasher, key)
	if err != nil {
		return nil, err
	}
	chiperKey := hasher.Sum(nil)

	block, err := aes.NewCipher(chiperKey)
	if err != nil {
		return nil, err
	}
	return block, nil
}
