package util

import (
	"fmt"
	"github.com/wumansgy/goEncrypt/aes"
)

const AesSecretKey = "7GGTPKEDeoE0M4Co"
const CommonPrefix = "//+/"

func Encrypt(msg string) string {
	if msg == "" {
		return ""
	}
	length := len(CommonPrefix)
	if len(msg) > length && msg[:length] == CommonPrefix {
		fmt.Println("已加密过")
		return msg
	}
	text, err := aes.AesCbcEncryptHex([]byte(msg), []byte(AesSecretKey), nil)
	if err != nil {
		fmt.Println(err)
		return msg
	}
	return CommonPrefix + text
}

func Decrypt(msg string) string {
	if msg == "" {
		return ""
	}
	length := len(CommonPrefix)
	if len(msg) <= length || msg[:length] != CommonPrefix {
		return msg
	}
	msg = msg[length:]
	plaintext, err := aes.AesCbcDecryptByHex(msg, []byte(AesSecretKey), nil)
	if err != nil {
		fmt.Println(err)
		return msg
	}
	return string(plaintext)
}
