package domain

import (
	"mime/multipart"
	"time"
)

type RequestParameter struct {
	AllowMethod       string                 `json:"allow_method"`
	DataType          string                 `json:"data_type"`
	RequireToken      bool                   `json:"require_token"`
	RequireParameters map[string]interface{} `json:"require_parameters" default:"[]"`
}

type FileStructure struct {
	File        multipart.File
	FileName    string
	ContentType string
	Size        int64
}

type UserToken struct {
	Authorized bool
	SeqMember  int64
	Email      string
	UserLevel  int
	Exp        time.Time
	Allocated  int8
}

type RequireParameterItem struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}
