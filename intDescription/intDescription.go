package intDescription

import (
	"encoding/xml"
	"github.com/pterm/pterm"
	"io"
	"log"
	"os"
	"strings"
)

type IntDesc struct {
	XMLName xml.Name `xml:"rpc-reply"`
	Junos   string   `xml:"xmlns:junos,attr"`
	Info    struct {
		XMLName   xml.Name `xml:"interface-information"`
		Junos     string   `xml:"junos:style,attr"`
		PhysIntfs []struct {
			XMLName     xml.Name `xml:"physical-interface"`
			Name        string   `xml:"name"`
			AdminStatus string   `xml:"admin-status"`
			OperStatus  string   `xml:"oper-status"`
			Description string   `xml:"description"`
		} `xml:"physical-interface"`
		LogIntfs []struct {
			XMLName     xml.Name `xml:"logical-interface"`
			Name        string   `xml:"name"`
			AdminStatus string   `xml:"admin-status"`
			OperStatus  string   `xml:"oper-status"`
			Description string   `xml:"description"`
		} `xml:"logical-interface"`
	} `xml:"interface-information"`
	Cli struct {
		Banner string `xml:"banner"`
	} `xml:"cli"`
}

func IntDescription(pairs [][]string) {
	for _, pair := range pairs {
		// extract node name from the pair
		path := pair[0] // extract the first element of the pair
		// split the path by backslashes and get the last component (file name)
		pathComponents := strings.Split(path, "\\")
		fileName := pathComponents[len(pathComponents)-1]
		// split the file name by hyphens and get the second component (node name)
		fileNameComponents := strings.Split(fileName, "_")
		nodeName := fileNameComponents[0]

		pterm.DefaultSection.Println("Section/interface descriptions:", nodeName)
		pterm.Info.Printf("Comparing interface descriptions between %s and %s\n", pair[1], pair[0])
		file1Content, err := os.Open(pair[1])
		if err != nil {
			// Handle error
			continue
		}
		defer file1Content.Close()
		//Read XML Data from file
		xmlData1, err := io.ReadAll(file1Content)
		if err != nil {
			log.Fatal(err)
		}

		file2Content, err := os.Open(pair[0])
		if err != nil {
			// Handle error
			continue
		}
		defer file2Content.Close()
		//Read XML Data from file
		xmlData2, err := io.ReadAll(file2Content)
		if err != nil {
			log.Fatal(err)
		}

		// Parse the XML data into a Go structs
		var intInfo1 IntDesc
		err = xml.Unmarshal(xmlData1, &intInfo1)
		if err != nil {
			log.Fatal(err)
		}

		var intInfo2 IntDesc
		err = xml.Unmarshal(xmlData2, &intInfo2)
		if err != nil {
			log.Fatal(err)
		}
		var differences []string

		for _, intf1 := range intInfo1.Info.PhysIntfs {
			for _, intf2 := range intInfo2.Info.PhysIntfs {
				if intf1.Name == intf2.Name && (intf1.Description != intf2.Description) {
					differences = append(differences, pterm.FgRed.Sprintf("Interface Description Changed %s: was %s -> now %s", intf1.Name, intf1.Description, intf2.Description))
				}
				if intf1.Name == intf2.Name && (intf1.AdminStatus != intf2.AdminStatus) {
					differences = append(differences, pterm.FgRed.Sprintf("Interface Admin Status Changed %s: was %s -> now %s", intf1.Name, intf1.AdminStatus, intf2.AdminStatus))
				}
				if intf1.Name == intf2.Name && (intf1.OperStatus != intf2.OperStatus) {
					differences = append(differences, pterm.FgRed.Sprintf("Interface Operational Status Changed %s: was %s -> now %s", intf1.Name, intf1.OperStatus, intf2.OperStatus))
				}

			}
		}
		// Loop through logical interfaces for the current physical interface
		for _, logintf1 := range intInfo1.Info.LogIntfs {
			for _, logintf2 := range intInfo2.Info.LogIntfs {
				if logintf1.Name == logintf2.Name && (logintf1.Description != logintf2.Description) {
					differences = append(differences, pterm.FgRed.Sprintf("Interface Description Changed: %s was %s -> now %s", logintf1.Name, logintf1.Description, logintf2.Description))
				}
				if logintf1.Name == logintf2.Name && (logintf1.AdminStatus != logintf2.AdminStatus) {
					differences = append(differences, pterm.FgRed.Sprintf("Interface Admin Status Changed %s: was %s -> now %s", logintf1.Name, logintf1.AdminStatus, logintf2.AdminStatus))
				}
				if logintf1.Name == logintf2.Name && (logintf1.OperStatus != logintf2.OperStatus) {
					differences = append(differences, pterm.FgRed.Sprintf("Interface Operational Status Changed %s: was %s -> now %s", logintf1.Name, logintf1.OperStatus, logintf2.OperStatus))
				}
			}
		}

		if len(differences) == 0 {
			pterm.Success.Println("No differences found.\n")
		} else {
			pterm.Warning.Println("Differences found. Please check the following interfaces:\n")
			for _, diff := range differences {
				pterm.Println(diff)
			}
		}
	}

}
