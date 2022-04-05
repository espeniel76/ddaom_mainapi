package define

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	HTTP_SERVER   = "http://192.168.1.20:81"
	HTTP_PORT     = "3011"
	HTTP_PORT_SSL = "3012"
	SOCKET_SERVER = "http://221.146.220.5:5150"
	DSN_MASTER    = "espeniel:anjgkrp@tcp(221.146.220.5)/ddaom?charset=utf8mb4&parseTime=True&loc=Local"
	DSN_SLAVE1    = "espeniel:anjgkrp@tcp(221.146.220.5)/ddaom?charset=utf8mb4&parseTime=True&loc=Local"
	DSN_LOG1      = "espeniel:anjgkrp@tcp(221.146.220.5)/ddaom_user1?charset=utf8mb4&parseTime=True&loc=Local"
	DSN_LOG2      = "espeniel:anjgkrp@tcp(221.146.220.5)/ddaom_user2?charset=utf8mb4&parseTime=True&loc=Local"

	DSN_REDIS = "redis://anjgkrp@221.146.220.5:6379"

	JWT_ACCESS_SECRET  = "a969cefccb3c64dbdc8f19f36f4400b5b35216be30667fc2671036f0bc48040cfb11ff4ad6d3384b7de874d09b03bd4be27d2fae226ff807e3982211f8db5517"
	JWT_REFRESH_SECRET = "d83ae8d11d2cda872b0c2d510a504c65d98eccf6f9d73df46413622042d5cc6356ada1bef8becc76adadde8e600227a0f0518c631ae3c0180ab2fef86a46b659"

	FILE_UPLOAD_PATH = "/home/samba/espeniel/www_ddaom/upload/"
	REPLACE_PATH     = "/home/samba/espeniel/www_ddaom"
	DEFAULT_PROFILE  = "/default.png"
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
