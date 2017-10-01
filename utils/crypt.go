package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"io"
	"crypto/des"
	"crypto/cipher"
	"bytes"
)

func Md5(in string) (m5 string) {
	hash := md5.New()
	io.WriteString(hash, in)
	bytes := hash.Sum(nil)
	m5 = hex.EncodeToString(bytes)
	return
}

func Sha256(in string) (s256 string) {
	hash := sha256.New()
	io.WriteString(hash, in)
	bytes := hash.Sum(nil)
	s256 = hex.EncodeToString(bytes)
	return
}

func Sha512(in string) (s512 string) {
	hash := sha512.New()
	io.WriteString(hash, in)
	bytes := hash.Sum(nil)
	s512 = hex.EncodeToString(bytes)
	return
}

/**
 * BKDR Hash Function
 * 把字符串hash到64位整数
 */
func BKDRHash(str string) uint64 {
	var seed uint64 = 131 // 31 131 1313 13131 131313 etc..
	var hash uint64 = 0

	for i := 0; i < len(str); i++ {
		var ui = uint64(str[i])
		hash = hash*seed + ui
	}
	return hash
}

// 3DES加密
func TripleDesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:8])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
// 3DES解密
func TripleDesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key[:8])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}