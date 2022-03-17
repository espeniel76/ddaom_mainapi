package tools

import "math"

func GetLimitStart(page int64, sizePerPage int64) int64 {
	return (page - 1) * sizePerPage
}

func GetTotalPage(totalData int64, sizePerPage int64) int {
	return int(math.Ceil(float64(totalData) / float64(sizePerPage)))
}
