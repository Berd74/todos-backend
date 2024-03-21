package utils

func StringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
