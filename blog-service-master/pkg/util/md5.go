package util

//一个上传文件的工具库，功能是针对上传文件时的一些相关处理。
import (
	"crypto/md5"
	"encoding/hex"
)

// EncodeMD5 用于针对上传后的文件名格式化，
// 将文件名 MD5 后再进行写入，防止暴露原始名称
func EncodeMD5(value string) string {
	//新建一个md存储文件名
	m := md5.New()
	m.Write([]byte(value))
	//将md内格式化后的文件名以字符串形式返回
	return hex.EncodeToString(m.Sum(nil))
}
