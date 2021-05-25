package slice

// EntryExists checks if an entry exists in a slice of strings
// The function returns a boolean value:
//		true if the entry exists or
// 		false if the entry does not exist
func EntryExists(slice []string, entry string) bool {
	for i := 0; i < len(slice); i++ {
		if slice[i] == entry {
			return true
		}
	}
	return false
}
