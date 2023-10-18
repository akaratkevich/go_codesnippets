package getNodes

import (
	"bufio"
	"fmt"
	"github.com/pterm/pterm"
	"os"
	"strings"
)

// Declare a struct for the input files
type captureNode struct {
	Node string
}

func GetNodes() ([]captureNode, error) {

	var NodeList []captureNode

	for {
		pterm.Info.Print("Enter the nodes. Enter 'DONE' to finish: ")
		list := bufio.NewScanner(os.Stdin)
		list.Scan()
		fullList := strings.TrimSpace(list.Text())

		if fullList == "DONE" {
			break
		}

		NodeList = append(NodeList, captureNode{
			Node: fullList,
		})
	}
	if len(NodeList) == 0 {
		return nil, fmt.Errorf("No files provided")
	}

	return NodeList, nil

}
