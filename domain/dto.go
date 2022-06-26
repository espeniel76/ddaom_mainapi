package domain

import "net/http"

type CommonRequest struct {
	Vars       map[string]string
	Parameters map[string]interface{}
	HttpRquest *http.Request
	JWToken    string
}

type CommonResponse struct {
	ResultCode string      `json:"result_code"`
	ErrorDesc  string      `json:"error_desc"`
	Data       interface{} `json:"data" default:"[]"`
}

type ConnectionInfos struct {
	JwtAccessSecret  string `json:"jwt_access_secret"`
	JwtRefreshSecret string `json:"jwt_refresh_secret"`
	PushServerKey    string `json:"push_server_key"`
	DefaultProfile   string `json:"default_profile"`
	HTTPServer       string `json:"http_server"`
	HTTPPort         string `json:"http_port"`
	HTTPPortSsl      string `json:"http_port_ssl"`
	DsnMaster        string `json:"dsn_master"`
	DsnSlave         string `json:"dsn_slave"`
	DsnLog1Master    string `json:"dsn_log1_master"`
	DsnLog1Slave     string `json:"dsn_log1_slave"`
	DsnLog2Master    string `json:"dsn_log2_master"`
	DsnLog2Slave     string `json:"dsn_log2_slave"`
	DsnRedisMaster   string `json:"dsn_redis_master"`
	DsnRedisSlave    string `json:"dsn_redis_slave"`
	DsnMongodb       string `json:"dsn_mongodb"`
	FileUploadPath   string `json:"file_upload_path"`
	ReplacePath      string `json:"replace_path"`
	AwsS3Region      string `json:"aws_s3_region"`
	AwsSecretKey     string `json:"aws_secret_key"`
	AwsAccessKey     string `json:"aws_access_key"`
	AwsBucketName    string `json:"aws_bucket_name"`
	AwsS3Server      string `json:"aws_s3_server"`
}
