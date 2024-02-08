package hash

import (
	"crypto/md5"
	"fmt"
	"github.com/spaolacci/murmur3"
)

func Hash(data []byte) uint64 {
	return murmur3.Sum64(data)
}

// Md5 对传入的字节数据进行哈希计算，并返回一个 16 字节的 MD5 哈希值
func Md5(data []byte) []byte {
	digest := md5.New()
	digest.Write(data)
	return digest.Sum(nil)
}

// Md5Hex 无论输入多长的数据，最后返回的字符串长度固定为 32 个字符。
func Md5Hex(data []byte) string {
	return fmt.Sprintf("%x", Md5(data))
}
