package bgpSummary

import (
	"encoding/xml"
	"fmt"
	"github.com/pterm/pterm"
	"io"
	"log"
	"os"
	"strings"
)

type BGP struct {
	XMLName        xml.Name `xml:"rpc-reply"`
	Text           string   `xml:",chardata"`
	Junos          string   `xml:"junos,attr"`
	BgpInformation struct {
		Text                     string `xml:",chardata"`
		Xmlns                    string `xml:"xmlns,attr"`
		BgpThreadMode            string `xml:"bgp-thread-mode"`
		ThreadState              string `xml:"thread-state"`
		DefaultEbgpAdvertiseMode string `xml:"default-ebgp-advertise-mode"`
		DefaultEbgpReceiveMode   string `xml:"default-ebgp-receive-mode"`
		GroupCount               string `xml:"group-count"`
		PeerCount                string `xml:"peer-count"`
		DownPeerCount            string `xml:"down-peer-count"`
		BgpRib                   []struct {
			Text                          string `xml:",chardata"`
			Style                         string `xml:"style,attr"`
			Name                          string `xml:"name"`
			TotalPrefixCount              string `xml:"total-prefix-count"`
			ReceivedPrefixCount           string `xml:"received-prefix-count"`
			AcceptedPrefixCount           string `xml:"accepted-prefix-count"`
			ActivePrefixCount             string `xml:"active-prefix-count"`
			SuppressedPrefixCount         string `xml:"suppressed-prefix-count"`
			HistoryPrefixCount            string `xml:"history-prefix-count"`
			DampedPrefixCount             string `xml:"damped-prefix-count"`
			TotalExternalPrefixCount      string `xml:"total-external-prefix-count"`
			ActiveExternalPrefixCount     string `xml:"active-external-prefix-count"`
			AcceptedExternalPrefixCount   string `xml:"accepted-external-prefix-count"`
			SuppressedExternalPrefixCount string `xml:"suppressed-external-prefix-count"`
			TotalInternalPrefixCount      string `xml:"total-internal-prefix-count"`
			ActiveInternalPrefixCount     string `xml:"active-internal-prefix-count"`
			AcceptedInternalPrefixCount   string `xml:"accepted-internal-prefix-count"`
			SuppressedInternalPrefixCount string `xml:"suppressed-internal-prefix-count"`
			PendingPrefixCount            string `xml:"pending-prefix-count"`
			BgpRibState                   string `xml:"bgp-rib-state"`
			VpnRibState                   string `xml:"vpn-rib-state"`
		} `xml:"bgp-rib"`
		BgpPeer []struct {
			Text            string `xml:",chardata"`
			Style           string `xml:"style,attr"`
			Heading         string `xml:"heading,attr"`
			PeerAddress     string `xml:"peer-address"`
			PeerAs          string `xml:"peer-as"`
			InputMessages   string `xml:"input-messages"`
			OutputMessages  string `xml:"output-messages"`
			RouteQueueCount string `xml:"route-queue-count"`
			FlapCount       string `xml:"flap-count"`
			ElapsedTime     struct {
				Text    string `xml:",chardata"`
				Seconds string `xml:"seconds,attr"`
			} `xml:"elapsed-time"`
			Description string `xml:"description"`
			PeerState   struct {
				Text   string `xml:",chardata"`
				Format string `xml:"format,attr"`
			} `xml:"peer-state"`
			BgpRib []struct {
				Text                  string `xml:",chardata"`
				Style                 string `xml:"style,attr"`
				Name                  string `xml:"name"`
				ActivePrefixCount     string `xml:"active-prefix-count"`
				ReceivedPrefixCount   string `xml:"received-prefix-count"`
				AcceptedPrefixCount   string `xml:"accepted-prefix-count"`
				SuppressedPrefixCount string `xml:"suppressed-prefix-count"`
			} `xml:"bgp-rib"`
		} `xml:"bgp-peer"`
	} `xml:"bgp-information"`
	Cli struct {
		Text   string `xml:",chardata"`
		Banner string `xml:"banner"`
	} `xml:"cli"`
}

func BGPSummary(pairs [][]string) {

	for _, pair := range pairs {
		// Create a panel for bgp summary
		panel1 := ""

		// extract node name from the pair
		path := pair[0] // extract the first element of the pair
		// split the path by backslashes and get the last component (file name)
		pathComponents := strings.Split(path, "\\")
		fileName := pathComponents[len(pathComponents)-1]
		// split the file name by hyphens and get the second component (node name)
		fileNameComponents := strings.Split(fileName, "_")
		nodeName := fileNameComponents[0]

		pterm.DefaultSection.Println("Section/bgp summary:", nodeName)
		pterm.Info.Printf("Comparing BGP summary between %s and %s\n", pair[1], pair[0])
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
		var instanceInfo1 BGP
		err = xml.Unmarshal(xmlData1, &instanceInfo1)
		if err != nil {
			log.Fatal(err)
		}

		var instanceInfo2 BGP
		err = xml.Unmarshal(xmlData2, &instanceInfo2)
		if err != nil {
			log.Fatal(err)
		}

		// Var for differences
		var differences []string

		// compare the peer count in the summary
		peers1 := instanceInfo1.BgpInformation.PeerCount
		peers2 := instanceInfo2.BgpInformation.PeerCount

		if peers1 == peers2 {
			differences = append(differences, pterm.FgLightBlue.Sprintf("\n*Number of peers is the same: %s\n", instanceInfo1.BgpInformation.PeerCount))
		} else if peers1 != peers2 {
			differences = append(differences, pterm.FgLightRed.Sprintf("\n*Number of peers has changed from %s -> %s\n", instanceInfo1.BgpInformation.PeerCount, instanceInfo2.BgpInformation.PeerCount))
		}

		minLen := min(len(instanceInfo1.BgpInformation.BgpRib), len(instanceInfo2.BgpInformation.BgpRib))
		for i := 0; i < minLen; i++ {
			rib1 := instanceInfo1.BgpInformation.BgpRib[i]
			rib2 := instanceInfo2.BgpInformation.BgpRib[i]
			if rib1.Name == rib2.Name {
				differences = append(differences, pterm.DefaultSection.WithLevel(2).Sprintf(rib1.Name))
				if rib1.TotalPrefixCount == rib2.TotalPrefixCount {
					differences = append(differences, pterm.FgLightBlue.Sprintf("TOTAL: OK |"))
				} else if rib1.TotalPrefixCount != rib2.TotalPrefixCount {
					differences = append(differences, pterm.FgLightRed.Sprintf("TOTAL: Changed from %s -> %s |", rib1.TotalPrefixCount, rib2.TotalPrefixCount))
				}
				if rib1.ReceivedPrefixCount == rib2.ReceivedPrefixCount {
					differences = append(differences, pterm.FgLightBlue.Sprintf("RECEIVED: OK |"))
				} else if rib1.ReceivedPrefixCount != rib2.ReceivedPrefixCount {
					differences = append(differences, pterm.FgLightRed.Sprintf("RECEIVED: Changed from %s -> %s |", rib1.ReceivedPrefixCount, rib2.ReceivedPrefixCount))
				}
				if rib1.AcceptedPrefixCount == rib2.AcceptedPrefixCount {
					differences = append(differences, pterm.FgLightBlue.Sprintf("ACCEPTED: OK |"))
				} else if rib1.AcceptedPrefixCount != rib2.AcceptedPrefixCount {
					differences = append(differences, pterm.FgLightRed.Sprintf("ACCEPTED: Changed from %s -> %s |", rib1.AcceptedPrefixCount, rib2.AcceptedPrefixCount))
				}
				if rib1.ActivePrefixCount == rib2.ActivePrefixCount {
					differences = append(differences, pterm.FgLightBlue.Sprintf("ACTIVE: OK"))
				} else if rib1.ActivePrefixCount != rib2.ActivePrefixCount {
					differences = append(differences, pterm.FgLightRed.Sprintf("ACTIVE: Changed from %s -> %s", rib1.ActivePrefixCount, rib2.ActivePrefixCount))
				}
			} else {
				differences = append(differences, pterm.FgLightRed.Sprintf("\n*Out of bounds error detected for %s\n", rib1.Name))
			}
		}

		// loop through the RIs and find differences

		minLen2 := min(len(instanceInfo1.BgpInformation.BgpPeer), len(instanceInfo2.BgpInformation.BgpPeer))
		for i := 0; i < minLen2; i++ {
			bgpPeer1 := instanceInfo1.BgpInformation.BgpPeer[i]
			bgpPeer2 := instanceInfo2.BgpInformation.BgpPeer[i]

			if bgpPeer1.PeerAddress == bgpPeer2.PeerAddress {
				differences = append(differences, pterm.Sprintf("\n"))
				differences = append(differences, pterm.DefaultSection.WithLevel(3).Sprintf("Peer address: %s, Peer AS: %s, Description: %s", bgpPeer1.PeerAddress, bgpPeer1.PeerAs, bgpPeer1.Description))
				minLen3 := min(len(bgpPeer1.BgpRib), len(bgpPeer2.BgpRib))
				for j := 0; j < minLen3; j++ {
					bgpRib1 := bgpPeer1.BgpRib[j]
					bgpRib2 := bgpPeer2.BgpRib[j]
					if bgpRib1.Name == bgpRib2.Name {
						differences = append(differences, pterm.DefaultSection.WithLevel(4).Sprintf(bgpRib1.Name))
						if bgpRib1.AcceptedPrefixCount != bgpRib2.AcceptedPrefixCount {
							differences = append(differences, pterm.FgLightRed.Sprintf("ACCEPTED: %s -> %s ", bgpRib1.AcceptedPrefixCount, bgpRib2.AcceptedPrefixCount))
						} else {
							differences = append(differences, pterm.FgLightGreen.Sprintf("ACCEPTED: OK "))
						}
						if bgpRib1.ReceivedPrefixCount != bgpRib2.ReceivedPrefixCount {
							differences = append(differences, pterm.FgLightRed.Sprintf("| RECEIVED: %s -> %s ", bgpRib1.ReceivedPrefixCount, bgpRib2.ReceivedPrefixCount))
						} else {
							differences = append(differences, pterm.FgLightGreen.Sprintf("| RECEIVED: OK "))
						}
						if bgpRib1.ActivePrefixCount != bgpRib2.ActivePrefixCount {
							differences = append(differences, pterm.FgLightRed.Sprintf("| ACTIVE: %s -> %s ", bgpRib1.ActivePrefixCount, bgpRib2.ActivePrefixCount))
						} else {
							differences = append(differences, pterm.FgGreen.Sprintf("| ACTIVE: OK "))
						}
						if bgpRib1.SuppressedPrefixCount != bgpRib2.SuppressedPrefixCount {
							differences = append(differences, pterm.FgLightRed.Sprintf("| SUPPRESSED: %s -> %s", bgpRib1.SuppressedPrefixCount, bgpRib2.SuppressedPrefixCount))
						} else {
							differences = append(differences, pterm.FgLightGreen.Sprintf("| SUPPRESSED: OK"))
						}
					} else {
						differences = append(differences, pterm.FgLightRed.Sprintf("\n*Out of bounds error detected for %s in bgpPeer1 and bgpPeer2\n", bgpRib1.Name))
					}
				}
			} else {
				differences = append(differences, pterm.FgLightRed.Sprintf("\n*Out of bounds error detected for %s in bgpPeer1 and bgpPeer2\n", bgpPeer1.PeerAddress))
			}
		}

		if len(differences) == 0 {
			pterm.Success.Println("No differences found.\n")
		} else {
			pterm.Warning.Println("Differences found. Please check the following RIBs:\n")
			for _, diff := range differences {
				panel1 += fmt.Sprintf(diff)
			}
		}
		// Arrange panels
		panels := pterm.Panels{
			{{Data: pterm.Sprintf(panel1)}},
			//{{Data: pterm.Sprintf(panel3)}},
			//{{Data: pterm.Sprintf(panel4)}},
		}
		// Print panels.
		pterm.DefaultPanel.WithPanels(panels).Render()
		// End panels

	}

}

// Helper function to check the length
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
