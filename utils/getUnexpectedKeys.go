package utils

func GetUnexpectedKeys(myMap *map[string]interface{}, allowedKeys []string) []string {
	var unexpectedKeys []string

	allowed := make(map[string]bool)
	for _, key := range allowedKeys {
		allowed[key] = true
	}
	for key := range *myMap {
		if !allowed[key] {
			if unexpectedKeys == nil {
				unexpectedKeys = make([]string, 0)
			}
			unexpectedKeys = append(unexpectedKeys, key)
		}
	}
	return unexpectedKeys
}
