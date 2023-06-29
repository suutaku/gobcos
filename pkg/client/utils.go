package client

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/pbkdf2"
)

func encryptAES(plainText string, key string) (string, error) {
	dk := pbkdf2.Key([]byte(key), []byte("dapp.blockchain.webank.com"), 65536, 16, sha256.New)
	iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

	block, err := aes.NewCipher(dk)
	if err != nil {
		fmt.Println("key error1", err)
	}
	if plainText == "" {
		fmt.Println("plain content empty")
	}

	ecb := cipher.NewCBCEncrypter(block, iv)
	content := []byte(plainText)
	content = pKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	// return hex string
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func pKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (c *Client) publicHeader() PublicRequestHeader {
	return PublicRequestHeader{
		Token:    c.dapp.Token,
		BizSeqNo: c.bizSeqNo,
	}
}

func (c *Client) invoke(method Method, payload interface{}, headrs interface{}) (*CommonResponse, error) {
	buf := bytes.Buffer{}
	if err := json.NewEncoder(&buf).Encode(payload); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", method(), &buf)
	if err != nil {
		return nil, err
	}
	headMap := make(map[string]string)
	err = mapstructure.Decode(headrs, &headMap)
	if err != nil {
		return nil, err
	}
	for k, v := range headMap {
		req.Header.Add(k, v)
	}
	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	return NewCommonResponse(resp.Body)
}
