package main

// insertAt inserts substring sub into string str at the specified position pos
func insertAt(str string, sub string, pos int) string {
	return str[:pos] + sub + str[pos:]
}

// googleIdToUuid formats a given string by inserting hyphens at specific positions and appending a key
func googleIdToUuid(id string) string {
	key := "todo" // must have 4 characters
	id = insertAt(id, "-", 20)
	id = insertAt(id, "-", 16)
	id = insertAt(id, "-", 12)
	id = insertAt(id, "-", 8)
	return id + key
}
