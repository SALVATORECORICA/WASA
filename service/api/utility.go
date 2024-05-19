package api

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Utility function
func isValidID(nickname string) bool {
	if len(nickname) >= 3 && len(nickname) <= 16 {
		return true
	} else {
		return false
	}
}

// function to create a new folder

func createFolders(id int) error {
	fmt.Println(id)
	// Obtain the complete path
	path, err := os.Executable()
	if err != nil {
		return err
	}
	// Obtain the directory of the father

	fatherDir := filepath.Dir(path)

	// Convert id from int to string
	idString := strconv.Itoa(id)

	//Create the folder
	completePath := filepath.Join(fatherDir, idString)
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
	fmt.Println("IL path e: ", subfolderPath)
	return nil
}

// function to extract the bearer
func extractBearer(authorization string) string {
	var tokens = strings.Split(authorization, " ")
	if len(tokens) == 2 {
		return strings.Trim(tokens[1], " ")
	}
	return ""
}

// function to check of the comment is valid

func validComment(comment string) bool {
	if len(comment) == 0 || len(comment) > 30 {
		return false
	}
	return true
}

// Funzione per determinare il tipo di immagine basato sui primi byte
func detectImageType(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	if bytes.Equal(data[0:2], []byte{0xFF, 0xD8}) {
		return true
	} else if bytes.Equal(data[0:4], []byte{0x89, 0x50, 0x4E, 0x47}) {
		return true
	} else {
		return false
	}
}
