package encryptutil

import (
	"crypto/md5"
	"io"
	"encoding/hex"
)

/**
 * md5加密
 * @param $prestr 需要签名的字符串
 * @param $key 私钥
 * return 签名结果
 */
func Md5(str, salt string) string {
	h := md5.New()
	io.WriteString(h, str)
	io.WriteString(h, salt)
	cipherStr := h.Sum(nil)
	encryptedData := hex.EncodeToString(cipherStr)
	return encryptedData
}