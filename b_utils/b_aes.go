package b_utils

import (
	"bytes"
	"crypto/aes"
)

type AESTool struct {
	Key []byte
	IV  []byte
}

func paddingKey(key []byte) []byte {

	if len(key) == 16 || len(key) == 24 || len(key) == 32 {
		return key
	} else if len(key) < 16 {

		return append(key, bytes.Repeat([]byte{byte(0)}, 16-len(key))...)
	} else if len(key) < 24 {
		return append(key, bytes.Repeat([]byte{byte(0)}, 24-len(key))...)
	} else if len(key) < 32 {
		return append(key, bytes.Repeat([]byte{byte(0)}, 32-len(key))...)
	} else {
		return key[0:32]
	}

}

func NewAESTool(key []byte, iv []byte) *AESTool {
	key = paddingKey(key)
	return &AESTool{Key: key, IV: iv}
}

func (tool *AESTool) ZeroPadding(src []byte) []byte {

	paddingCount := aes.BlockSize - len(src)%aes.BlockSize

	if paddingCount == 0 {
		return src
	} else {
		return append(src, bytes.Repeat([]byte{byte(0)}, paddingCount)...)
	}
}

func (tool *AESTool) PKCS5Padding(src []byte) []byte {

	padding := aes.BlockSize - len(src)%aes.BlockSize
	return append(src, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func (tool *AESTool) PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:length-unPadding]
}

func (tool *AESTool) ZeroUnPadding(src []byte) []byte {

	for i := len(src) - 1; ; i-- {
		if src[i] != 0 {
			return src[:i+1]
		}
	}
	return nil
}

func (tool *AESTool) ECBEncrypt(src []byte) ([]byte, error) {

	block, err := aes.NewCipher(tool.Key)

	if err != nil {
		return nil, err
	}

	src = tool.PKCS5Padding(src)
	//返回加密的结果
	encryptData := make([]byte, len(src))
	tmpData := make([]byte, aes.BlockSize)

	for index := 0; index < len(src); index += aes.BlockSize {

		block.Encrypt(tmpData, src[index:index+aes.BlockSize])
		copy(encryptData[index:], tmpData)
	}

	return encryptData, nil
}

func (tool *AESTool) ECBDecrypt(src []byte) ([]byte, error) {
	block, err := aes.NewCipher(tool.Key)

	if err != nil {
		return nil, err
	}

	decryptData := make([]byte, len(src))
	tmpData := make([]byte, aes.BlockSize)

	for index := 0; index < len(src); index += aes.BlockSize {
		block.Decrypt(tmpData, src[index:index+aes.BlockSize])
		copy(decryptData[index:], tmpData)
	}
	return tool.PKCS5UnPadding(decryptData), nil
}
