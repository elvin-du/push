package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"io"
	"log"
)

var (
	AES_ENCRYPTED_DATA_INVALID = errors.New("Data invalid")
)

func AesEncrypt(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	//padding
	remain := block.BlockSize() - buf.Len()%block.BlockSize()
	for i := 0; i < remain; i++ {
		buf.WriteByte(byte(i))
	}

	src_len := buf.Len()

	iv := make([]byte, block.BlockSize())
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return nil, err
	}

	buf.Write(iv)
	bin := buf.Bytes()

	cipher_text := bin[:src_len]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipher_text, cipher_text)

	return bin, nil
}

func AesEncryptToHex(key, data []byte) (string, error) {
	bin, err := AesEncrypt(key, data)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bin), nil
}

func AesDecrypt(key, bin []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(bin) < block.BlockSize()*2 || len(bin)%block.BlockSize() != 0 {
		return nil, AES_ENCRYPTED_DATA_INVALID
	}

	src_len := len(bin) - block.BlockSize()
	iv := bin[src_len:]

	mode := cipher.NewCBCDecrypter(block, iv)
	src := bin[:src_len]
	mode.CryptBlocks(src, src)
	padding_size := int(src[src_len-1])
	if padding_size >= block.BlockSize() {
		return nil, AES_ENCRYPTED_DATA_INVALID
	}

	return src[:src_len-padding_size-1], nil
}

func AesDecryptFromHex(key []byte, in string) ([]byte, error) {
	bin, err := hex.DecodeString(in)
	if err != nil {
		return nil, err
	}

	return AesDecrypt(key, bin)
}

func HMacHash(key, message []byte) []byte {
	hash := hmac.New(sha1.New, key)
	hash.Write(message)
	return hash.Sum(nil)
}

func HMacHashToHex(key, message []byte) string {
	return hex.EncodeToString(HMacHash(key, message))
}

func HMacEqual(mac1, mac2 []byte) bool {
	return hmac.Equal(mac1, mac2)
}

func HMacEqualFromHex(mac1, mac2 string) bool {
	mac_bin1, err := hex.DecodeString(mac1)
	if err != nil {
		log.Println(err)
		return false
	}

	mac_bin2, err := hex.DecodeString(mac2)
	if err != nil {
		log.Println(err)
		return false
	}

	return HMacEqual(mac_bin1, mac_bin2)
}
