package rest

import (
	"ddaom/define"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func checkToken(token []string) (bool, string) {
	_token := ""
	_isBool := false

	if len(token) > 0 {
		_token = token[0]
		_isBool = true
	}

	return _isBool, _token
}

func checkRequireMethod(r *http.Request, requireParameters map[string]interface{}) (bool, map[string]interface{}) {
	for kUri, vUri := range requireParameters {

		requestSlice := strings.Split(r.RequestURI, "?")
		requireSlice := strings.Split(kUri, "/:")

		m := vUri.(map[string]interface{})

		if len(requireSlice) > 1 {
			if len(requestSlice) > 1 {
				if requireSlice[0] == requestSlice[0] {
					if m["allow_method"] == r.Method {
						return true, m
					}
				}
			} else {

				matched, _ := regexp.MatchString(requireSlice[0]+"/[0-9]+", requestSlice[0])
				if matched {
					if m["allow_method"] == r.Method {
						return true, m
					}
				}

			}
		} else {
			if requireSlice[0] == requestSlice[0] {
				if m["allow_method"] == r.Method {
					return true, m
				}
			}
		}
	}
	return false, nil
}

func checkSetFormDataParameters(r *http.Request, selectFormat map[string]interface{}) (bool, string) {

	isExistKey := true
	noExistParam := define.NONE
	for _, v := range selectFormat["require_parameters"].([]interface{}) {
		isExistKey = false
		require := v.(map[string]interface{})
		if require["required"].(bool) {
			switch require["type"].(string) {
			case "blob":
				_, _, err := r.FormFile(require["name"].(string))
				if err == nil {
					isExistKey = true
				}
			case "string":
				val := r.FormValue(require["name"].(string))
				if val != "" {
					isExistKey = true
				}
			case "int":
				val, _ := strconv.Atoi(r.FormValue(require["name"].(string)))
				if val != 0 {
					isExistKey = true
				}
			}
			if !isExistKey {
				noExistParam = fmt.Sprintf("%v (%v)", require["name"], require["type"])
				break
			}
		} else {
			isExistKey = true
			break
		}
	}
	return isExistKey, noExistParam
}

func checkRequireParameters(requestParamameter map[string]interface{}, selectFormat map[string]interface{}) (bool, string) {

	isExistKey := true
	noExistParam := define.NONE
	fmt.Println("##### Check parameters #####")
	for _, v := range selectFormat["require_parameters"].([]interface{}) {
		isExistKey = false
		require := v.(map[string]interface{})
		if require["required"].(bool) {
			for request, _ := range requestParamameter {
				if require["name"] == request {
					isExistKey = true
					break
				}
			}
		} else {
			isExistKey = true
			break
		}

		if !isExistKey {
			noExistParam = fmt.Sprintf("%v (%v)", require["name"], require["type"])
			break
		}
	}

	return isExistKey, noExistParam
}

func checkParameterType(requestParamameter map[string]interface{}, selectFormat map[string]interface{}) (bool, string) {
	isMatchType := true
	noMatchParam := define.NONE

	for _, v := range selectFormat["require_parameters"].([]interface{}) {
		isMatchType = false
		require := v.(map[string]interface{})
		if require["required"].(bool) {
			for kRequest, vRequest := range requestParamameter {
				if require["name"] == kRequest {
					convertedRequireValue := fmt.Sprintf("%v", require["type"])
					convertedRequestValue := fmt.Sprintf("%v", reflect.TypeOf(vRequest))
					if convertedRequireValue == convertedRequestValue {
						isMatchType = true
						break
					}
				}
			}
		} else {
			isMatchType = true
			break
		}

		if !isMatchType {
			noMatchParam = fmt.Sprintf("%s (%v)", require["name"], require["type"])
			break
		}
	}

	return isMatchType, noMatchParam
}
