package evpnBrief

import (
	"encoding/xml"
	"github.com/pterm/pterm"
	"io"
	"log"
	"os"
	"strings"
)

type RpcReply struct {
	XMLName          xml.Name         `xml:"rpc-reply"`
	Junos            string           `xml:"xmlns:junos,attr"`
	EvpnInstanceInfo EvpnInstanceInfo `xml:"evpn-instance-information"`
}

type EvpnInstanceInfo struct {
	XMLName       xml.Name       `xml:"evpn-instance-information"`
	EvpnInstances []EvpnInstance `xml:"evpn-instance"`
}

type EvpnInstance struct {
	XMLName              xml.Name `xml:"evpn-instance"`
	EvpnInstanceName     string   `xml:"evpn-instance-name"`
	NumLocalMacs         int      `xml:"num-local-macs"`
	NumRemoteMacs        int      `xml:"num-remote-macs"`
	LocalInterfaces      int      `xml:"local-interfaces"`
	LocalInterfacesUp    int      `xml:"local-interfaces-up"`
	IrbInterfaces        int      `xml:"irb-interfaces"`
	IrbInterfacesUp      int      `xml:"irb-interfaces-up"`
	NumProtectInterfaces int      `xml:"num-protect-interfaces"`
	EvpnNumNeighbors     int      `xml:"evpn-num-neighbors"`
	EvpnNumEsi           int      `xml:"evpn-num-esi"`
}

func EvpnBrief(pairs [][]string) {
	for _, pair := range pairs {

		// extract node name from the pair
		path := pair[0] // extract the first element of the pair
		// split the path by backslashes and get the last component (file name)
		pathComponents := strings.Split(path, "\\")
		fileName := pathComponents[len(pathComponents)-1]
		// split the file name by hyphens and get the second component (node name)
		fileNameComponents := strings.Split(fileName, "_")
		nodeName := fileNameComponents[0]

		pterm.DefaultSection.Println("Section/evpn brief:", nodeName)
		pterm.Info.Printf("Comparing EVPN Brief between %s and %s\n", pair[1], pair[0])
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
		var instanceInfo1 RpcReply
		err = xml.Unmarshal(xmlData1, &instanceInfo1)
		if err != nil {
			log.Fatal(err)
		}

		var instanceInfo2 RpcReply
		err = xml.Unmarshal(xmlData2, &instanceInfo2)
		if err != nil {
			log.Fatal(err)
		}

		// Compare the contents of the two files
		var differences []string
		for i, instance1 := range instanceInfo1.EvpnInstanceInfo.EvpnInstances {
			instance2 := instanceInfo2.EvpnInstanceInfo.EvpnInstances[i]
			if instance1.EvpnInstanceName == instance2.EvpnInstanceName && instance1.NumLocalMacs != instance2.NumLocalMacs || instance1.NumRemoteMacs != instance2.NumRemoteMacs || instance1.LocalInterfaces != instance2.LocalInterfaces || instance1.LocalInterfacesUp != instance2.LocalInterfacesUp || instance1.IrbInterfaces != instance2.IrbInterfaces || instance1.IrbInterfacesUp != instance2.IrbInterfaces || instance1.NumProtectInterfaces != instance2.NumProtectInterfaces || instance1.EvpnNumNeighbors != instance2.EvpnNumNeighbors || instance1.EvpnNumEsi != instance2.EvpnNumEsi {
				differences = append(differences, pterm.FgLightWhite.Sprintf("EVPN Instance Name: %s", instance1.EvpnInstanceName))
				if instance1.NumLocalMacs != instance2.NumLocalMacs {
					percentageDiff := ((float64(instance2.NumLocalMacs) - float64(instance1.NumLocalMacs)) / float64(instance1.NumLocalMacs)) * 100
					if percentageDiff > 0 {
						differences = append(differences, pterm.FgLightGreen.Sprintf("Local MAC Count: was %d -> now %d (percentage difference: %.2f%%)", instance1.NumLocalMacs, instance2.NumLocalMacs, percentageDiff))
					} else {
						differences = append(differences, pterm.FgLightMagenta.Sprintf("Local MAC Count: was %d -> now %d (percentage difference: %.2f%%)", instance1.NumLocalMacs, instance2.NumLocalMacs, percentageDiff))
					}
				}
				if instance1.NumRemoteMacs != instance2.NumRemoteMacs {
					percentageDiff := ((float64(instance2.NumRemoteMacs) - float64(instance1.NumRemoteMacs)) / float64(instance1.NumRemoteMacs)) * 100
					if percentageDiff > 0 {
						differences = append(differences, pterm.FgLightGreen.Sprintf("Remote MAC Count: was %d -> now %d (percentage difference: %.2f%%)", instance1.NumRemoteMacs, instance2.NumRemoteMacs, percentageDiff))
					} else {
						differences = append(differences, pterm.FgLightMagenta.Sprintf("Remote MAC Count: was %d -> now %d (percentage difference: %.2f%%)", instance1.NumRemoteMacs, instance2.NumRemoteMacs, percentageDiff))
					}
				}
				if instance1.LocalInterfaces != instance2.LocalInterfaces {
					differences = append(differences, pterm.FgRed.Sprintf("Local Interfaces: was %d -> now %d", instance1.LocalInterfaces, instance2.LocalInterfaces))
				}
				if instance1.LocalInterfacesUp != instance2.LocalInterfacesUp {
					differences = append(differences, pterm.FgRed.Sprintf("Local Interfaces UP: was %d -> now %d", instance1.LocalInterfacesUp, instance2.LocalInterfacesUp))
				}
				if instance1.IrbInterfaces != instance2.IrbInterfaces {
					differences = append(differences, pterm.FgRed.Sprintf("IRB Interfaces: was %d -> now %d", instance1.IrbInterfaces, instance2.IrbInterfaces))
				}
				if instance1.IrbInterfacesUp != instance2.IrbInterfacesUp {
					differences = append(differences, pterm.FgRed.Sprintf("IRB Interfaces UP: was %d -> now %d", instance1.IrbInterfacesUp, instance2.IrbInterfacesUp))
				}
				if instance1.NumProtectInterfaces != instance2.NumProtectInterfaces {
					differences = append(differences, pterm.FgRed.Sprintf("Protected Interfaces: was %d -> now %d", instance1.NumProtectInterfaces, instance2.NumProtectInterfaces))
				}
				if instance1.EvpnNumNeighbors != instance2.EvpnNumNeighbors {
					differences = append(differences, pterm.FgRed.Sprintf("EVPN Neighbors: was %d -> now %d", instance1.EvpnNumNeighbors, instance2.EvpnNumNeighbors))
				}
				if instance1.EvpnNumEsi != instance2.EvpnNumEsi {
					differences = append(differences, pterm.FgRed.Sprintf("EVPN ESI: was %d -> now %d", instance1.EvpnNumEsi, instance2.EvpnNumEsi))
				}
			}
		}

		if len(differences) == 0 {
			pterm.Success.Println("No differences found.\n")
		} else {
			pterm.Warning.Println("Differences found. Please check the following Instances:\n")
			for _, diff := range differences {
				pterm.Println(diff)
			}
		}
	}

}
