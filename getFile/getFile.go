package getFile

import (
	"fmt"
	"github.com/pterm/pterm"
	"strings"
)

func GetFile() string {

	pterm.BgBlue.Print("Enter the file path (format: D:\\user\\anton\\document.doc): ")
	pterm.Println()
	fileName, err := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
	if err != nil {
		fmt.Println("Error reading input:", err)
		return ""
	}

	// Remove the trailing newline or carriage return characters
	fileName = strings.TrimRight(fileName, "\r\n")
	pterm.Info.Println("You have entered the following file path:", fileName)

	return fileName
}
