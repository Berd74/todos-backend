package utils

import (
	"strconv"
	"strings"
)

func StringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func StringToBoolOrNil(s string) *bool {
	if s == "" {
		return nil
	}

	done, errParse := strconv.ParseBool(s)
	if errParse != nil {
		return nil
	}
	return &done
}

func SplitOrNil(s *string) *[]string {
	if s == nil {
		return nil
	}
	arr := strings.Split(*s, ",")
	return &arr
}
