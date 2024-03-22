package upload

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/go-programming-tour-book/blog-service/global"
	"github.com/go-programming-tour-book/blog-service/pkg/util"
)

// FileType 作为类别标识的基础类型
type FileType int

// TypeImage 为了后续有其它的需求，能标准化的进行处理
const TypeImage FileType = iota + 1

// GetFileName 获取文件名称
func GetFileName(name string) string {
	//先通过获取文件名后缀
	ext := GetFileExt(name)
	//TrimSuffix（）减去后缀
	fileName := strings.TrimSuffix(name, ext)
	//进行DM5加密（encode）后返回加密后端文件名
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

// GetFileExt 获取文件后缀
func GetFileExt(name string) string {
	return path.Ext(name)
}

// GetSavePath 获取文件保存地址
func GetSavePath() string {
	return global.AppSetting.UploadSavePath
}

func GetServerUrl() string {
	return global.AppSetting.UploadServerUrl
}

// CheckSavePath 检查保存目录是否存在
func CheckSavePath(dst string) bool {
	_, err := os.Stat(dst)
	//IsNotExist判断错误是否存在
	return os.IsNotExist(err)
}

// CheckContainExt 检查文件后缀是否包含在约定的后缀配置项中
func CheckContainExt(t FileType, name string) bool {
	ext := GetFileExt(name)
	//统一格式
	ext = strings.ToUpper(ext)
	switch t {
	case TypeImage:
		for _, allowExt := range global.AppSetting.UploadImageAllowExts {
			if strings.ToUpper(allowExt) == ext {
				return true
			}
		}

	}

	return false
}

// CheckMaxSize 检查文件大小是否超出最大大小限制。
func CheckMaxSize(t FileType, f multipart.File) bool {
	content, _ := io.ReadAll(f)
	size := len(content)
	switch t {
	case TypeImage:
		if size >= global.AppSetting.UploadImageMaxSize*1024*1024 {
			return true
		}
	}

	return false
}

// CheckPermission 检查文件权限是否足够
func CheckPermission(dst string) bool {
	_, err := os.Stat(dst)

	return os.IsPermission(err)
}

func CreateSavePath(dst string, perm os.FileMode) error {
	//os.MkdirAll将会以传入的 os.FileMode 权限位去递归创建所需的目录结构，
	//若涉及的目录均已存在，则不会进行任何操作，直接返回 nil。
	err := os.MkdirAll(dst, perm)
	if err != nil {
		return err
	}

	return nil
}

func SaveFile(file *multipart.FileHeader, dst string) error {
	//拿到文件描述符
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	//打开目录
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	//将文件写入目录
	_, err = io.Copy(out, src)
	return err
}
