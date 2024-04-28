package api

// Utility function
func isValidID(id string) bool {
	if len(id) >= 3 && len(id) <= 16 {
		return true
	} else {
		return false
	}
}
