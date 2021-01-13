package auth_crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func Decrypt(msg string) string {
	//e, cipherKey := getEncKey()
	////log.Info(string(cipherKey))
	//if e != nil {
	//	panic(e)
	//	return e.Error() + "@@" + msg
	//}
	//m, err := base64.StdEncoding.DecodeString(msg)
	//if err != nil {
	//	panic(e)
	//}
	//err, ret := decrypt(cipherKey, m)
	//if err != nil {
	//	return err.Error() + "@@" + msg
	//}
	//return ret
	return msg
}

func decrypt(key []byte, ciphertext []byte) (error, string) {
	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return err, ""
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	bm := cipher.NewCBCDecrypter(block, iv)
	bm.CryptBlocks(ciphertext, ciphertext)
	ciphertext, err = pkcs7Unpad(ciphertext, aes.BlockSize)
	if err != nil {
		return err, ""
	}
	return nil, string(ciphertext)
}

var (
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

// pkcs7Unpad validates and unpads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7Unpad(b []byte, blocksize int) ([]byte, error) {
	if blocksize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if b == nil || len(b) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	if len(b)%blocksize != 0 {
		return nil, ErrInvalidPKCS7Padding
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, ErrInvalidPKCS7Padding
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, ErrInvalidPKCS7Padding
		}
	}
	return b[:len(b)-n], nil
}
