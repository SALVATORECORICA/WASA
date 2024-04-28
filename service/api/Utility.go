package api

// Utility function
func isValidID(nickname string) bool {
	if len(nickname) >= 3 && len(nickname) <= 16 {
		return true
	} else {
		return false
	}
}
