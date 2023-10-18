package compareNeighborsBGP

import (
	"bufio"
	"github.com/pterm/pterm"
	"log"
	"os"
	"regexp"
	"strings"
)

// Neighbor representation to include the command, IP, and State
type Neighbor struct {
	Command string
	IP      string
	State   string
}

func CompareNeighbors(pairs [][]string) {
	for _, pair := range pairs {
		// Extract node name from the pair (e.g. D:\\vdcMigrations\\ld5-1a_route-checks-pre)
		path := pair[0]
		pathComponents := strings.Split(path, "\\")
		fileName := pathComponents[len(pathComponents)-1]
		fileNameComponents := strings.Split(fileName, "_")
		nodeName := fileNameComponents[0]

		pterm.DefaultSection.Println("Section/BGP Neighbor comparison:", nodeName)
		pterm.Info.Printf("Comparing BGP Neighbors: %s and %s\n", pair[1], pair[0])

		// Read file 1
		file1Bytes, err := os.ReadFile(pair[1])
		if err != nil {
			log.Fatal(err)
		}
		file1Content := string(file1Bytes)
		file1Neighbors := extractNeighbors(file1Content)

		// Read file 2
		file2Bytes, err := os.ReadFile(pair[0])
		if err != nil {
			log.Fatal(err)
		}
		file2Content := string(file2Bytes)
		file2Neighbors := extractNeighbors(file2Content)

		panel1 := "(*) Existing BGP Neighbors/State:\n"
		panel1 += "==================================\n"

		// Group neighbors by command for file 1
		file1NeighborsByCommand := groupNeighborsByCommand(file1Neighbors)
		for command, neighbors := range file1NeighborsByCommand {
			panel1 += pterm.FgLightYellow.Sprintf("BGP Summary: %s\n", command)
			for _, neighbor := range neighbors {
				panel1 += pterm.FgLightYellow.Sprintf(" * %s: %s\n", neighbor.IP, neighbor.State)
			}
		}

		// Panel 2 preparation
		panel2 := "(-) Missing BGP Neighbors/State Change:\n"
		panel2 += "==================================\n"
		// Group neighbors by command for file 2
		file2NeighborsByCommand := groupNeighborsByCommand(file2Neighbors)

		// Loop to go through file 1 neighbors and compare with file 2
		for command, neighbors := range file1NeighborsByCommand {
			file2Neighbors, found := file2NeighborsByCommand[command]
			if !found {
				panel2 += pterm.FgLightRed.Sprintf("%s\n", command)
				for _, neighbor := range neighbors {
					panel2 += pterm.FgLightRed.Sprintf("- %s\n", neighbor.IP)
				}
			} else {
				// Compare neighbors within the same command
				for _, neighbor := range neighbors {
					file2Neighbor := findNeighborByIP(file2Neighbors, neighbor.IP)
					if file2Neighbor != nil {
						if file2Neighbor.State != neighbor.State {
							panel2 += pterm.FgLightRed.Sprintf("[%s] State change %s: %s -> %s\n", command, neighbor.IP, neighbor.State, file2Neighbor.State)
						}
					}
				}
			}
		}

		panel3 := "(+) Extra BGP Neighbors:\n"
		panel3 += "==================================\n"

		for command, neighbors := range file2NeighborsByCommand {
			file1Neighbors, found := file1NeighborsByCommand[command]
			if !found {
				panel3 += pterm.FgLightGreen.Sprintf("%s\n", command)
				for _, neighbor := range neighbors {
					panel3 += pterm.FgLightGreen.Sprintf("+ %s\n", neighbor.IP)
				}
			} else {
				// Compare neighbors within the same command
				for _, neighbor := range neighbors {
					found := false
					for _, file1Neighbor := range file1Neighbors {
						if file1Neighbor.IP == neighbor.IP {
							found = true
							break
						}
					}
					if !found {
						panel3 += pterm.FgLightGreen.Sprintf("+ [%s] %s\n", command, neighbor.IP)
					}
				}
			}
		}

		panels := pterm.Panels{
			{{Data: pterm.Sprintf(panel1)}, {Data: pterm.Sprintf(panel2)}, {Data: pterm.Sprintf(panel3)}},
		}
		pterm.DefaultPanel.WithPanels(panels).Render()
	}
}

func extractNeighbors(content string) []Neighbor {
	neighbors := []Neighbor{}
	scanner := bufio.NewScanner(strings.NewReader(content))

	// Reg ex to extract neighbors and the commands from the file
	reCommand := regexp.MustCompile(`show bgp summary group (\S+)`)
	reNeighbor := regexp.MustCompile(`^(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s+\d+\s+\d+\s+\d+\s+\d+\s+\d+\s+\S+\s+\S+\s+(\S+)`)

	var currentCommand string

	for scanner.Scan() {
		line := scanner.Text()

		commandMatch := reCommand.FindStringSubmatch(line)

		if len(commandMatch) > 0 {
			currentCommand = commandMatch[1]
			continue
		}

		neighborMatch := reNeighbor.FindStringSubmatch(line)

		if len(neighborMatch) > 0 {
			neighbor := Neighbor{
				Command: currentCommand,
				IP:      neighborMatch[1],
				State:   neighborMatch[2],
			}
			neighbors = append(neighbors, neighbor)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return neighbors
}

// Group neighbors by command
func groupNeighborsByCommand(neighbors []Neighbor) map[string][]Neighbor {
	neighborsByCommand := make(map[string][]Neighbor)

	for _, neighbor := range neighbors {
		neighborsByCommand[neighbor.Command] = append(neighborsByCommand[neighbor.Command], neighbor)
	}

	return neighborsByCommand
}

// Find neighbor by IP
func findNeighborByIP(neighbors []Neighbor, ip string) *Neighbor {
	for _, neighbor := range neighbors {
		if neighbor.IP == ip {
			return &neighbor
		}
	}
	return nil
}
