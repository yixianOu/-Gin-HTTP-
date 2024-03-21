package service

import (
	"errors"
	"mime/multipart"
	"os"

	"github.com/go-programming-tour-book/blog-service/global"

	"github.com/go-programming-tour-book/blog-service/pkg/upload"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

// UploadFile 将上传文件工具库与我们具体的业务接口结合起来
func (svc *Service) UploadFile(fileType upload.FileType, file multipart.File, fileHeader *multipart.FileHeader) (*FileInfo, error) {
	//获取文件所需的基本信息
	fileName := upload.GetFileName(fileHeader.Filename)
	//文件检查(文件大小，后缀）
	if !upload.CheckContainExt(fileType, fileName) {
		return nil, errors.New("file suffix is not supported.")
	}
	if upload.CheckMaxSize(fileType, file) {
		return nil, errors.New("exceeded maximum file limit.")
	}

	uploadSavePath := upload.GetSavePath()
	//写入文件前，判断是否具备写入条件（目的地是否存在，有无权限）
	if upload.CheckSavePath(uploadSavePath) {
		if err := upload.CreateSavePath(uploadSavePath, os.ModePerm); err != nil {
			return nil, errors.New("failed to create save directory.")
		}
	}
	if upload.CheckPermission(uploadSavePath) {
		return nil, errors.New("insufficient file permissions.")
	}

	dst := uploadSavePath + "/" + fileName
	//写入文件操作
	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
