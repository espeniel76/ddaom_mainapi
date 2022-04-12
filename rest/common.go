package rest

import (
	"ddaom/define"
	"ddaom/domain"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func common(f func(*domain.CommonRequest) domain.CommonResponse) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// s, err := ioutil.ReadAll(r.Body)
		// fmt.Println(err)
		// fmt.Println(string(s))

		var req = domain.CommonRequest{}
		var res = domain.CommonResponse{}
		isCheck := true
		var checkRequireParametersValue string
		var requestParameters *map[string]interface{}
		var contentType string

		authorization := r.Header["Authorization"]
		isToken, token := checkToken(authorization)

		if len(r.Header["Content-Type"]) < 1 {
			res.ResultCode = define.NO_EXIST_CONTENT_TYPE
			isCheck = false
		}

		if isCheck {
			contentType = r.Header["Content-Type"][0]
			if contentType != "application/json" {
				contentType = strings.Split(contentType, ";")[0]
			}

			fmt.Println("Content-Type: ", contentType)
		}

		var selectFormat map[string]interface{}
		if isCheck {
			isCheck, selectFormat = checkRequireMethod(r, define.Mconfig)
			if !isCheck {
				res.ResultCode = define.INCORRECT_HTTP_METHOD
				isCheck = false
			}
		}

		if isCheck {
			if selectFormat["require_token"].(bool) {
				if !isToken {
					res.ResultCode = define.NO_TOKEN
					res.ErrorDesc = "There are no tokens. You must have an appropriate authentication token."
					isCheck = false
				} else {
					verifyResult, err := define.VerifyToken(token, define.JWT_ACCESS_SECRET)
					if err != nil {
						res.ResultCode = verifyResult
						res.ErrorDesc = err.Error()
						isCheck = false
					}
				}
			}
		}

		selectFormatDateType := fmt.Sprintf("%v", selectFormat["data_type"])

		if isCheck {
			if contentType != selectFormatDateType {
				res.ResultCode = define.CONTENT_TYPE_MISMATCH
				res.ErrorDesc = fmt.Sprintf("Require: %s, Now: %s", selectFormatDateType, contentType)
				isCheck = false
			}
		}

		if isCheck {
			switch contentType {
			case "application/json":
				d := json.NewDecoder(r.Body)
				d.UseNumber()
				if err := d.Decode(&requestParameters); err != nil {
					res.ResultCode = define.JSON_SYNTAX_ERROR
					res.ErrorDesc = err.Error()
					isCheck = false
				}
			case "multipart/form-data":
				isCheck, checkRequireParametersValue = checkSetFormDataParameters(r, selectFormat)
				if !isCheck {
					res.ResultCode = define.NO_PARAMETER
					res.ErrorDesc = checkRequireParametersValue
				}
			case "text/plain":
			}
		}

		if isCheck && contentType == "application/json" {
			isCheck, checkRequireParametersValue = checkRequireParameters(*requestParameters, selectFormat)
			if !isCheck {
				res.ResultCode = define.NO_PARAMETER
				res.ErrorDesc = checkRequireParametersValue
			}
		}

		if isCheck {
			req.HttpRquest = r
			req.JWToken = token
			if contentType == "multipart/form-data" {
				res = f(&req)
			} else {
				switch r.Method {
				case http.MethodGet:
					req.Vars = mux.Vars(r)
					res = f(&req)
				case http.MethodPut:
					fallthrough
				case http.MethodPatch:
					fallthrough
				case http.MethodDelete:
					fallthrough
				case http.MethodPost:
					req.Parameters = *requestParameters
					req.Vars = mux.Vars(r)
					res = f(&req)
				}
			}
			if res.ResultCode == "" {
				res.ResultCode = define.OK
			}
		}

		data, _ := json.Marshal(res)
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, string(data))
	}
}
