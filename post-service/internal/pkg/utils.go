package pkg

import (
	"bytes"
	"encoding/json"
)

func Copy(dst, src any) error {
	data, err := json.Marshal(src)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		if _, ok := err.(*json.UnmarshalTypeError); ok || err.Error() == "json: unknown field \"created_at\"" {
			return nil
		}
		return err
	}
	return nil
}

func ConvertListIntToInt64(nums []int) []int64 {
	var resp []int64
	for _, num := range nums {
		resp = append(resp, int64(num))
	}
	return resp
}

func ConvertListUintToInt64(nums []uint64) []int64 {
	var resp []int64
	for _, num := range nums {
		resp = append(resp, int64(num))
	}
	return resp
}
