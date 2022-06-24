package define

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	JWT_ACCESS_SECRET  = "a969cefccb3c64dbdc8f19f36f4400b5b35216be30667fc2671036f0bc48040cfb11ff4ad6d3384b7de874d09b03bd4be27d2fae226ff807e3982211f8db5517"
	JWT_REFRESH_SECRET = "d83ae8d11d2cda872b0c2d510a504c65d98eccf6f9d73df46413622042d5cc6356ada1bef8becc76adadde8e600227a0f0518c631ae3c0180ab2fef86a46b659"
	PUSH_SERVER_KEY    = "AAAAs8DEFV4:APA91bHjJF63wpyefl-6IBMhJ0PVb0VPePwirNxes3PzRgMxg7wb1Q8ykTyzxnTrCVVMX8cE5ROxvjWJLLZ9cRw8pt5daXUsd-mxiK4jqgdVkR_XWaUW1snEXBSFFnebSR_D2L-Pn-wY"

	// real
	// HTTP_SERVER      = "hhttps://ddaom.s3.ap-northeast-2.amazonaws.com/default.png"
	// HTTP_PORT        = "80"
	// HTTP_PORT_SSL    = "443"
	// DSN_MASTER       = "espeniel:anjgkrp@tcp(172.31.33.10)/ddaom?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_SLAVE        = "espeniel:anjgkrp@tcp(localhost)/ddaom?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_LOG1_MASTER  = "espeniel:anjgkrp@tcp(172.31.33.10)/ddaom_user1?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_LOG1_SLAVE   = "espeniel:anjgkrp@tcp(localhost)/ddaom_user1?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_LOG2_MASTER  = "espeniel:anjgkrp@tcp(172.31.33.10)/ddaom_user2?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_LOG2_SLAVE   = "espeniel:anjgkrp@tcp(localhost)/ddaom_user2?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_REDIS_MASTER = "redis://anjgkrp@172.31.33.10:6379"
	// DSN_REDIS_SLAVE  = "redis://anjgkrp@localhost:6379"
	// DSN_MONGODB      = "mongodb://172.31.33.10:27017"

	// home
	// HTTP_SERVER      = "http://192.168.1.20:81"
	// HTTP_PORT        = "3011"
	// HTTP_PORT_SSL    = "3012"
	// DSN_MASTER       = "espeniel:anjgkrp@tcp(localhost:3307)/ddaom?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_SLAVE        = "espeniel:anjgkrp@tcp(localhost:3307)/ddaom?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_LOG1_MASTER  = "espeniel:anjgkrp@tcp(localhost:3307)/ddaom_user1?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_LOG1_SLAVE   = "espeniel:anjgkrp@tcp(localhost:3307)/ddaom_user1?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_LOG2_MASTER  = "espeniel:anjgkrp@tcp(localhost:3307)/ddaom_user2?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_LOG2_SLAVE   = "espeniel:anjgkrp@tcp(localhost:3307)/ddaom_user2?charset=utf8mb4&parseTime=True&loc=Local"
	// DSN_REDIS_MASTER = "redis://anjgkrp@localhost:6379"
	// DSN_REDIS_SLAVE  = "redis://anjgkrp@localhost:6379"
	// DSN_MONGODB      = "mongodb://localhost:27017"

	// office
	HTTP_SERVER      = "http://192.168.1.20:81"
	HTTP_PORT        = "3011"
	HTTP_PORT_SSL    = "3012"
	DSN_MASTER       = "espeniel:anjgkrp@tcp(192.168.1.20:3306)/ddaom?charset=utf8mb4&parseTime=True&loc=Local"
	DSN_SLAVE        = "espeniel:anjgkrp@tcp(192.168.1.20:3306)/ddaom?charset=utf8mb4&parseTime=True&loc=Local"
	DSN_LOG1_MASTER  = "espeniel:anjgkrp@tcp(192.168.1.20:3306)/ddaom_user1?charset=utf8mb4&parseTime=True&loc=Local"
	DSN_LOG1_SLAVE   = "espeniel:anjgkrp@tcp(192.168.1.20:3306)/ddaom_user1?charset=utf8mb4&parseTime=True&loc=Local"
	DSN_LOG2_MASTER  = "espeniel:anjgkrp@tcp(192.168.1.20:3306)/ddaom_user2?charset=utf8mb4&parseTime=True&loc=Local"
	DSN_LOG2_SLAVE   = "espeniel:anjgkrp@tcp(192.168.1.20:3306)/ddaom_user2?charset=utf8mb4&parseTime=True&loc=Local"
	DSN_REDIS_MASTER = "redis://anjgkrp@192.168.1.20:6379"
	DSN_REDIS_SLAVE  = "redis://anjgkrp@192.168.1.20:6379"
	DSN_MONGODB      = "mongodb://192.168.1.20:27017"

	// FILE_UPLOAD_PATH = "/home/samba/espeniel/www_ddaom/upload/"
	// REPLACE_PATH     = "/home/samba/espeniel/www_ddaom"

	// S3 로 변경
	FILE_UPLOAD_PATH = "upload/"
	REPLACE_PATH     = "www"
	DEFAULT_PROFILE  = "/default.png"
	AWS_S3_REGION    = "ap-northeast-2"
	AWS_SECRET_KEY   = "tmdrd9InqdB6zxz0qWhFgMi2yHDSAycTDAiwFXYf"
	AWS_ACCESS_KEY   = "AKIA225UPJLVEA7EMINP"
	AWS_BUCKET_NAME  = "ddaom"
	AWS_S3_SERVER    = "https://ddaom.s3.ap-northeast-2.amazonaws.com"
)

func GetDefineApiParse() (*map[string]interface{}, error) {
	jsonFile, err := os.Open("db/initial/define_api_v2.json")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	m := make(map[string]interface{})
	err = json.Unmarshal([]byte(byteValue), &m)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

var Mconfig map[string]interface{}

func SetDefineApiParse() {
	jsonFile, _ := os.Open("db/initial/define_api_v2.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	Mconfig = make(map[string]interface{})
	_ = json.Unmarshal([]byte(byteValue), &Mconfig)
}
