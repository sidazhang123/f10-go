package auth_crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
)

func getEncKey() (error, []byte) {
	if encKey == nil {
		return fmt.Errorf("failed to get encKey"), nil
	}
	return nil, encKey
}

func Encrypt(msg string) string {
	//e, cipherKey := getEncKey()
	////log.Info(string(cipherKey))
	//if e != nil {
	//	log.Error(e.Error())
	//	return ""
	//}
	//ret := encrypt(cipherKey, []byte(msg))
	//
	//return ret
	return msg
}

func encrypt(key []byte, plaintext []byte) string {

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	plaintext, _ = pkcs7Pad(plaintext, block.BlockSize())
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, block.BlockSize()+len(plaintext))
	for i := 0; i < aes.BlockSize; i++ {
		ciphertext[i] = byte(1 + i)
	}
	iv := ciphertext[:aes.BlockSize]
	//log.Info(fmt.Sprintf("iv=%+v", iv))
	//log.Info(fmt.Sprintf("ciphertext=%+v", ciphertext))

	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	//log.Info(fmt.Sprintf("%+v", ciphertext))
	return base64.StdEncoding.EncodeToString(ciphertext)
}

var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")
)

func pkcs7Pad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	n := blocksize - (len(b) % blocksize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}
