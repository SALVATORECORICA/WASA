package api

func isValidID(string id) bool {
	if len(id) >= 3 && len(id) <= 16 {
		return true
	} else {
		return false
	}
}
