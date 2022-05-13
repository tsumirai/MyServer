package util

import (
	"MyServer/middleware/logger"
	"context"
	"strconv"
)

// ConvertIntSliceToStringSlice 将int格式的切片转为string格式的切片
func ConvertIntSliceToStringSlice(param []int) []string {
	result := make([]string, 0, len(param))
	for _, v := range param {
		result = append(result, strconv.FormatInt(int64(v), 10))
	}
	return result
}

// ConvertInt64SliceToStringSlice 将int64格式的切片转为string格式的切片
func ConvertInt64SliceToStringSlice(param []int64) []string {
	result := make([]string, 0, len(param))
	for _, v := range param {
		result = append(result, strconv.FormatInt(v, 10))
	}
	return result
}

// ConvertStringSliceToInt64Slice 将string格式的切片转为int64格式的切片
func ConvertStringSliceToInt64Slice(ctx context.Context, param []string) ([]int64, error) {
	result := make([]int64, 0, len(param))
	for _, v := range param {
		temp, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			logger.Error(ctx, "ConvertStringSliceToInt64Slice", logger.LogArgs{"err": err, "param": v})
			return result, err
		}
		result = append(result, temp)
	}
	return result, nil
}

// ConvertStringSliceToIntSlice 将string格式的切片转为int格式的切片
func ConvertStringSliceToIntSlice(ctx context.Context, param []string) ([]int, error) {
	result := make([]int, 0, len(param))
	for _, v := range param {
		temp, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			logger.Error(ctx, "ConvertStringSliceToIntSlice", logger.LogArgs{"err": err, "param": v})
			return result, err
		}
		result = append(result, int(temp))
	}
	return result, nil
}
