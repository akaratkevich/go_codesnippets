package getUsernamePassword

import (
	"bufio"
	"github.com/pterm/pterm"
	"log"
	"os"
	"strings"
)

func GetUsernamePassword() (string, string) {

	//Setup a buffered reader
	reader := bufio.NewReader(os.Stdin)

	// Get username and password from user input
	pterm.Info.Println("Enter Username:")
	username, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	pterm.Info.Println("Enter Password:")
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	// Remove the trailing newline or carriage return characters
	username = strings.TrimRight(username, "\r\n")
	password = strings.TrimRight(password, "\r\n")

	return username, password
}
