package getFiles

import (
	"bufio"
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"strings"
)

// Declare a struct for the input files
type captureFile struct {
	Path string
	Name string
}

func GetFiles() ([]captureFile, error) {

	var FileList []captureFile

	for {
		pterm.Info.Print("Enter the file name paths for the required files (format: D:\\user\\anton\\config.txt). Enter DONE to finish: ")
		file := bufio.NewScanner(os.Stdin)
		file.Scan()
		filePath := strings.TrimSpace(file.Text())

		if filePath == "DONE" {
			break
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			pterm.Error.Printf("Error: %s\n", err.Error())
			continue
		}

		if fileInfo.IsDir() {
			pterm.Error.Println("Error: Please provide a path to a file, not a directory")
			continue
		}

		//extension := filepath.Ext(filePath)
		//if extension != ".txt" {
		//	pterm.Error.Println("Error: Only .txt files are allowed")
		//	continue
		//}

		FileList = append(FileList, captureFile{
			Path: filePath,
			Name: fileInfo.Name(),
		})
	}

	if len(FileList) == 0 {
		return nil, fmt.Errorf("No files provided")
	}

	return FileList, nil
}
