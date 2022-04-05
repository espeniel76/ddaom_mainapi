package handlers

import (
	"ddaom/db"
	"ddaom/define"
	"ddaom/domain"
	"ddaom/tools"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
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

	fmt.Println(item.Header.Get("Content-Type"))
	o := domain.FileStructure{
		File:        s[t].(multipart.File),
		FileName:    item.Filename,
		ContentType: item.Header.Get("Content-Type"),
		Size:        item.Size,
	}

	return &o, nil
}

func GetMyLogDb(allocated int8) *gorm.DB {

	var myLogDb *gorm.DB
	switch allocated {
	case 1:
		myLogDb = db.List[define.DSN_LOG1]
	case 2:
		myLogDb = db.List[define.DSN_LOG2]
	}
	return myLogDb
}

func SaveFile(_path string, oFile *domain.FileStructure) (string, error) {

	var err error
	// if oFile.ContentType == "" {
	// 	return "", err
	// }

	defer oFile.File.Close()

	now := time.Now()
	custom := now.Format("200601")

	path := define.FILE_UPLOAD_PATH + _path + "/" + custom + "/"
	_, err = os.ReadDir(path)
	if err != nil {
		err = os.MkdirAll(path, 0755)
	}
	uuid := tools.MakeUUID()
	var ext string
	switch oFile.ContentType {
	case "image/png":
		ext = "png"
	case "image/jpeg":
		ext = "jpg"
	case "image/gif":
		ext = "gif"
	default:
		tmp := strings.Split(oFile.FileName, ".")
		ext = tmp[len(tmp)-1]
	}

	saveFileName := uuid + "." + ext
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

	fullPath = strings.Replace(fullPath, define.REPLACE_PATH, "", -1)

	return fullPath, nil
}

func corm(o *gorm.DB, res *domain.CommonResponse) bool {
	isError := false
	if o.Error != nil {
		res.ResultCode = define.DB_ERROR_ORM
		res.ErrorDesc = o.Error.Error()
		isError = true
	}
	return isError
}
