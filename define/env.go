package define

import (
	"ddaom/domain"
	"encoding/json"
	"io/ioutil"
	"os"
)

var Mconfig map[string]interface{}

func SetDefineApiParse() {
	jsonFile, _ := os.Open("db/initial/define_api_v2.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	Mconfig = make(map[string]interface{})
	_ = json.Unmarshal([]byte(byteValue), &Mconfig)
}

var Mconn domain.ConnectionInfos

func SetConnectionInfosParse() {
	jsonFile, _ := os.Open("/etc/ddaom/define_conn_v0.1.json")
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	Mconn = domain.ConnectionInfos{}
	_ = json.Unmarshal([]byte(byteValue), &Mconn)
}
