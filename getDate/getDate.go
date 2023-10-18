package getDate

import (
	"fmt"
	"github.com/pterm/pterm"
	"strings"
	"time"
)

func GetDate() time.Time {
	// Get the date input
	pterm.BgBlue.Print("Enter date (format: 02/01/06): ")
	pterm.Println()
	dateStr, err := pterm.DefaultInteractiveTextInput.WithMultiLine(false).Show()
	if err != nil {
		fmt.Println("Error reading input:", err)
		return time.Time{}
	}

	// Remove the trailing newline or carriage return characters
	dateStr = strings.TrimRight(dateStr, "\r\n")

	// Parse the date string into a time.Time value
	date, err := time.Parse("02/01/06", dateStr)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return time.Time{}
	}
	pterm.Info.Println("You have entered the following date:", dateStr)

	return date
}
