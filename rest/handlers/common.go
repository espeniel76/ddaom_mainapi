package handlers

import (
	"ddaom/define"
	"ddaom/domain"
	"ddaom/tools"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Cp(s map[string]interface{}, t string) string {
	return fmt.Sprintf("%v", s[t])
}
func CpInt64(s map[string]interface{}, t string) int64 {
	n, _ := strconv.ParseInt(fmt.Sprintf("%v", s[t]), 10, 64)
	return n
}
func CpBool(s map[string]interface{}, t string) bool {
	n, _ := strconv.ParseBool(fmt.Sprintf("%v", s[t]))
	return n
}

func Cf(s map[string]interface{}, t string, r *http.Request) (*domain.FileStructure, error) {

	_, _, err := r.FormFile(t)
	if err != nil {
		return nil, err
	}

	item := s[t+"_header"].(*multipart.FileHeader)

	o := domain.FileStructure{
		File:        s[t].(multipart.File),
		FileName:    item.Filename,
		ContentType: item.Header.Get("Content-Type"),
		Size:        item.Size,
	}

	return &o, nil
}

func SaveFile(_path string, oFile *domain.FileStructure) (string, error) {

	var err error

	defer oFile.File.Close()

	now := time.Now()
	custom := now.Format("200601")

	path := define.Mconn.FileUploadPath + _path + "/" + custom + "/"
	_, err = os.ReadDir(path)
	if err != nil {
		err = os.MkdirAll(path, 0755)
	}
	_id := tools.MakeShortId()
	var ext string
	switch oFile.ContentType {
	case "image/png":
		ext = "png"
	case "image/jpeg":
		ext = "jpg"
	case "image/gif":
		ext = "gif"
	default:
		alist := strings.Split(oFile.FileName, ".")
		_ext := alist[len(alist)-1]
		if _ext == "png" || _ext == "jpg" || _ext == "gif" {
			ext = _ext
		} else {
			err = errors.New("not allowed image format")
			return "", err
		}
	}

	saveFileName := _id + "." + strings.ToLower(ext)
	fullPath := path + saveFileName
	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	_, err = io.Copy(f, oFile.File)
	if err != err {
		fmt.Println(err)
	}

	fullPath = strings.Replace(fullPath, define.Mconn.ReplacePath, "", -1)

	return fullPath, nil
}

func SaveFileS3(_path string, oFile *domain.FileStructure) (string, error) {

	s3 := tools.S3Info{
		AwsProfileName: "ddaom",
		AwsS3Region:    define.Mconn.AwsS3Region,
		AwsSecretKey:   define.Mconn.AwsSecretKey,
		AwsAccessKey:   define.Mconn.AwsAccessKey,
		BucketName:     define.Mconn.AwsBucketName,
	}

	err := s3.SetS3ConfigByKey()
	if err != nil {
		return "", err
	}
	if oFile.ContentType == "" {
		return "", err
	}
	defer oFile.File.Close()

	now := time.Now()
	custom := now.Format("200601")
	path := define.Mconn.FileUploadPath + _path + "/" + custom + "/"
	_id := tools.MakeShortId()
	var ext string
	switch oFile.ContentType {
	case "image/png":
		ext = "png"
	case "image/jpeg":
		ext = "jpg"
	case "image/gif":
		ext = "gif"
	case "application/octet-stream":
		alist := strings.Split(oFile.FileName, ".")
		_ext := alist[len(alist)-1]
		if _ext == "png" || _ext == "jpg" || _ext == "gif" {
			ext = _ext
		} else {
			err = errors.New("not allowed image format")
			return "", err
		}
	}

	saveFileName := _id + "." + strings.ToLower(ext)
	fullPath := path + saveFileName

	s3.UploadFile(oFile.File, fullPath, oFile.ContentType)

	return "/" + fullPath, nil
}
