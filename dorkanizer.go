// Dorkanizer v2.0
// Author: Pratyaksh Choudhary

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func init() {
	fmt.Println("\nDorkanizer [v1.0]")
	fmt.Printf("------------------------------------------------------\n\n")
}

func main() {

	// Checking platform
	if runtime.GOOS != "windows" && runtime.GOOS != "linux" {
		fmt.Println("This program has only been tested on Windows. Please don't run this on your system")
		os.Exit(0)
	}

	userProfile, err := user.Current() // getting userProfile of the current user
	exitIfError(err)
	downloadFolderPath := filepath.Join(userProfile.HomeDir, "Downloads") // getting path of the downloads folder

	dir, err := ioutil.ReadDir(downloadFolderPath) // getting details of files in the directory
	exitIfError(err)

	err = makeHistoryFile(downloadFolderPath, dir)
	exitIfError(err)

	// getting a map from JSON or using a pre-defined map if file doesn't exist
	extensions, err := getCategoryMap("./extensions.json")
	if err != nil {
		fmt.Println("No valid JSON file found, now using pre-defined categories...")
	}

	err = moveFiles(downloadFolderPath, dir, extensions)
	exitIfError(err)

	fmt.Println("All operations executed successfully...")
	exitMessage()
}

// getCategoryMap returns a map consisting of the categories of files and their extension names
// it calls the readFromJSON function and if there is not JSON file, returns a pre-defined map
func getCategoryMap(path string) (map[string][]string, error) {
	result, err := readFromJSON(path)

	if err != nil {
		result = map[string][]string{
			"Videos":            {".mp4", ".mkv", ".flv", ".avi", ".webm", ".mov", ".wmv"},
			"Audio":             {".mp3", ".m4a", ".mpeg", ".wav", ".aac", ".wma", ".flac"},
			"Applications":      {".exe", ".msi"},
			"Compressed":        {".zip", ".rar", ".7z", ".tar", ".gz"},
			"Not Miscellaneous": {".iso", ".ini"},
			"Miscellaneous":     {".torrent"},
			"Documents":         {".pdf", ".docx", ".doc", ".pptx", ".ppt", ".xlsx", ".xls", ".csv", ".tsv", ".txt"},
			"Images":            {".jpg", ".jpeg", ".png", ".gif", ".bmp"},
			"Programming":       {".py", ".java", ".class", ".c", ".cpp", ".cs", ".go", ".mod", ".html", ".css", ".js", ".ts", ".php", ".json", ".r", ".kt", ".md", ".ipynb", ".sh", ".xml"},
		}
	}

	return result, err
}

// readFromJSON reads a JSON file and returns it's values in a map
func readFromJSON(path string) (map[string][]string, error) {
	var result map[string][]string     // creating a map
	data, err := ioutil.ReadFile(path) // reading JSON in bytecote

	if err != nil {
		return result, err
	}

	json.Unmarshal(data, &result) // Unmarshalling the JSON into the map

	return result, nil
}

// makeDirIfNotExists makes a directory only if it doesn't exist already
// path defines where to create the directory
// dirName defines the name of the directory to be created
func makeDirIfNotExists(path string, dirName string) error {
	path = filepath.Join(path, dirName) // creating path for folder

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm) // making directory
		return nil
	}
	return fmt.Errorf("FileExistsError")
}

// makeHistoryFile creates a history.txt file in a given path
// path defines the path where the file is to be created
func makeHistoryFile(path string, dir []os.FileInfo) error {
	path = filepath.Join(path, "history.txt") // absolute path of history file

	// getting current date and time
	dateTime := time.Now()
	currDate := dateTime.Format("02/01/2006")
	timezone, _ := dateTime.Zone()
	currTime := fmt.Sprintf("%s (%s)", dateTime.Format("15:04:05"), timezone) // appending timezone to current time

	// opening the file or creating if it doesn't exist
	histFile, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	defer histFile.Close() // defer closing the file
	if err != nil {
		return err
	}

	// writing date and time to the file
	_, err = histFile.WriteString("\n    Date: ")
	if err != nil {
		return err
	}
	histFile.WriteString(currDate)
	histFile.WriteString("\n    Time: ")
	histFile.WriteString(currTime)
	histFile.WriteString("\n\n\n")

	// writing filenames to the file
	for _, file := range dir {
		histFile.WriteString(file.Name()) // writing the name

		// checking if a folder
		if file.IsDir() {
			histFile.WriteString("   [*Folder*]\n")
		} else {
			histFile.WriteString("\n")
		}
	}

	// adding a separator
	histFile.WriteString("\n")
	histFile.WriteString(strings.Repeat("-", 100))
	histFile.WriteString("\n")

	return nil
}

// moveFiles moves all files in a path according to the categories defined in the extensions map
func moveFiles(path string, dir []os.FileInfo, extensions map[string][]string) error {

	var source, dest string                  // to store source and destination
	folderDirName := "Uncategorized Folders" // folder name which will have only the folders
	categories := getKeysFromMap(extensions) // slice of all keys from extensions map
	otherExceptions := []string{
		"history.txt",
		"src",
		"dorkanizer.go",
		"extensions.json",
		"dorkanizer.exe",
		"dorkanizer.zip",
		"dorkanizer",
		"dorkanizer.bin",
	} // other exceptions to be ignored

	// looping though each file/folder in the directory
	for _, file := range dir {

		filename := file.Name()                // filename being iterated
		source = filepath.Join(path, filename) // absolute path of the file

		// taking care of folders
		if existsInSlice(filename, categories) || filename == folderDirName || existsInSlice(filename, otherExceptions) {
			continue
		} else if file.IsDir() {
			makeDirIfNotExists(path, folderDirName) // making directory for folders
			dest = filepath.Join(path, folderDirName, filename)
			err := os.Rename(source, dest) // moving the folder

			if err != nil {
				return err
			}
			continue
		}

		ext := strings.ToLower(filepath.Ext(source)) // getting extension of file

		isMoved := false // flag to check if file has been moved
		for key, values := range extensions {
			for _, extension := range values {
				if ext == extension {
					isMoved = true                            // setting true as file is being moved
					makeDirIfNotExists(path, key)             // making directory of category
					dest = filepath.Join(path, key, filename) // assigning file the destination
					err := os.Rename(source, dest)            // moving the file

					if err != nil {
						return err
					}
				}
			}
		}

		// moving file to Miscellaneous sections if not moved in the above block
		if !isMoved {
			makeDirIfNotExists(path, "Miscellaneous")             // making directory of category
			dest = filepath.Join(path, "Miscellaneous", filename) // assigning file the destination
			err := os.Rename(source, dest)                        // moving the file

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// getKeysFromMap returns a slice of keys from a map
func getKeysFromMap(hashmap map[string][]string) []string {
	var keySlice []string // empty slice

	for key := range hashmap {
		keySlice = append(keySlice, key) // appending each key to slice
	}

	return keySlice
}

// existsInSlice returns a bool true if a string exists in a slice
// false if doesn't exist in slice
func existsInSlice(str string, slice []string) bool {
	for _, value := range slice {
		if value == str {
			return true
		}
	}

	return false
}

func exitMessage() {
	fmt.Printf("\nPress ENTER to exit the program...")
	fmt.Scanln()
}

func exitIfError(err error) {
	if err != nil {
		fmt.Println(err)
		exitMessage()
		os.Exit(0)
	}
}
