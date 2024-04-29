package api

// Utility function
func isValidID(nickname string) bool {
	if len(nickname) >= 3 && len(nickname) <= 16 {
		return true
	} else {
		return false
	}
}

// function to create a new folder

func createFolders(nickname string) error {

	// Obtain the complete path
	path, err := os.Executable()
	if err != nil {
		return err
	}
	// Obtain the directory of the father

	fatherDir := filepath.Dir(path)

	//Create the folder
	completePath := filepath.Join(fatherDir, nickname)
	err = os.Mkdir(completePath, 0777)
	if err != nil {
		return err
	}

	// Create the sub-folder
	subfolderPath := filepath.Join(completePath, "photos")
	err = os.Mkdir(subfolderPath, 0777)
	if err != nil {
		return err
	}
	return nil
}
